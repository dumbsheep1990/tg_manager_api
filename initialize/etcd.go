package initialize

import (
	"context"
	"time"
	"tg_manager_api/global"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

// InitEtcd 初始化etcd客户端
func InitEtcd() {
	etcdConfig := global.Config.Etcd
	
	// 创建etcd客户端配置
	config := clientv3.Config{
		Endpoints:   etcdConfig.Endpoints,
		DialTimeout: time.Duration(etcdConfig.DialTimeout) * time.Second,
	}
	
	// 如果配置了用户名和密码，则添加认证信息
	if etcdConfig.Username != "" && etcdConfig.Password != "" {
		config.Username = etcdConfig.Username
		config.Password = etcdConfig.Password
	}
	
	// 创建etcd客户端
	client, err := clientv3.New(config)
	if err != nil {
		global.Logger.Error("etcd连接失败", zap.Error(err))
		return
	}
	
	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	_, err = client.Status(ctx, etcdConfig.Endpoints[0])
	if err != nil {
		global.Logger.Error("etcd服务不可用", zap.Error(err))
		return
	}
	
	global.EtcdClient = client
	global.Logger.Info("etcd初始化成功", zap.Strings("endpoints", etcdConfig.Endpoints))
}

// GetTaskKey 获取任务状态的key
func GetTaskKey(taskID string) string {
	return global.Config.Etcd.TaskPrefix + taskID
}

// GetLockKey 获取分布式锁的key
func GetLockKey(resourceID string) string {
	return global.Config.Etcd.LockPrefix + resourceID
}

// SaveTaskStatus 保存任务状态到etcd
func SaveTaskStatus(taskID string, status string) error {
	if global.EtcdClient == nil {
		return nil
	}
	
	key := GetTaskKey(taskID)
	
	// 设置带TTL的任务状态
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// 设置任务状态并设置TTL
	lease, err := global.EtcdClient.Grant(ctx, int64(global.Config.Etcd.TaskTTL))
	if err != nil {
		global.Logger.Error("创建etcd租约失败", zap.Error(err))
		return err
	}
	
	_, err = global.EtcdClient.Put(ctx, key, status, clientv3.WithLease(lease.ID))
	if err != nil {
		global.Logger.Error("保存任务状态失败", zap.String("taskID", taskID), zap.Error(err))
		return err
	}
	
	global.Logger.Debug("保存任务状态成功", zap.String("taskID", taskID), zap.String("status", status))
	return nil
}

// GetTaskStatus 从etcd获取任务状态
func GetTaskStatus(taskID string) (string, error) {
	if global.EtcdClient == nil {
		return "", nil
	}
	
	key := GetTaskKey(taskID)
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	resp, err := global.EtcdClient.Get(ctx, key)
	if err != nil {
		global.Logger.Error("获取任务状态失败", zap.String("taskID", taskID), zap.Error(err))
		return "", err
	}
	
	if len(resp.Kvs) == 0 {
		return "", nil
	}
	
	return string(resp.Kvs[0].Value), nil
}

// WatchTaskStatus 监控任务状态变化
func WatchTaskStatus(taskID string, callback func(status string)) {
	if global.EtcdClient == nil {
		return
	}
	
	key := GetTaskKey(taskID)
	
	// 创建观察通道
	watchChan := global.EtcdClient.Watch(context.Background(), key)
	
	// 处理状态变更事件
	go func() {
		for watchResp := range watchChan {
			for _, event := range watchResp.Events {
				status := string(event.Kv.Value)
				callback(status)
			}
		}
	}()
}
