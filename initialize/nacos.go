package initialize

import (
	"fmt"
	"net"
	"strconv"
	"tg_manager_api/global"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go.uber.org/zap"
)

// InitNacos 初始化nacos客户端并注册服务
func InitNacos() {
	nacosConfig := global.Config.Nacos
	
	// 创建ServerConfig
	serverConfig := []constant.ServerConfig{
		{
			IpAddr: nacosConfig.IPAddr,
			Port:   nacosConfig.Port,
		},
	}
	
	// 创建ClientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         nacosConfig.NamespaceID,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              nacosConfig.LogDir,
		CacheDir:            nacosConfig.CacheDir,
		LogLevel:            "info",
	}
	
	// 如果配置了用户名和密码
	if nacosConfig.Username != "" && nacosConfig.Password != "" {
		clientConfig.Username = nacosConfig.Username
		clientConfig.Password = nacosConfig.Password
	}
	
	// 创建nacos客户端
	namingClient, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfig,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		global.Logger.Error("创建nacos客户端失败", zap.Error(err))
		return
	}
	
	global.NacosClient = namingClient
	global.Logger.Info("nacos客户端初始化成功")
	
	// 注册服务
	RegisterService()
}

// RegisterService 向nacos注册服务
func RegisterService() {
	nacosConfig := global.Config.Nacos
	
	// 获取本机IP
	ip, err := getLocalIP()
	if err != nil {
		global.Logger.Error("获取本机IP失败", zap.Error(err))
		// 使用127.0.0.1作为回退IP
		ip = "127.0.0.1"
	}
	
	port := global.Config.System.Port
	
	// 服务注册参数
	params := vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        uint64(port),
		ServiceName: nacosConfig.ServiceName,
		Weight:      nacosConfig.Weight,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"version": "1.0"},
		ClusterName: nacosConfig.ClusterName,
		GroupName:   nacosConfig.GroupName,
	}
	
	// 注册服务实例
	success, err := global.NacosClient.RegisterInstance(params)
	if err != nil {
		global.Logger.Error("服务注册失败", zap.Error(err))
		return
	}
	
	if success {
		global.Logger.Info("服务注册成功",
			zap.String("ip", ip),
			zap.Int("port", port),
			zap.String("serviceName", nacosConfig.ServiceName))
	} else {
		global.Logger.Error("服务注册失败")
	}
}

// DeregisterService 从nacos注销服务
func DeregisterService() {
	if global.NacosClient == nil {
		return
	}
	
	nacosConfig := global.Config.Nacos
	
	// 获取本机IP
	ip, err := getLocalIP()
	if err != nil {
		global.Logger.Error("获取本机IP失败", zap.Error(err))
		ip = "127.0.0.1"
	}
	
	port := global.Config.System.Port
	
	// 服务注销参数
	params := vo.DeregisterInstanceParam{
		Ip:          ip,
		Port:        uint64(port),
		ServiceName: nacosConfig.ServiceName,
		GroupName:   nacosConfig.GroupName,
		Cluster:     nacosConfig.ClusterName,
		Ephemeral:   true,
	}
	
	// 注销服务实例
	success, err := global.NacosClient.DeregisterInstance(params)
	if err != nil {
		global.Logger.Error("服务注销失败", zap.Error(err))
		return
	}
	
	if success {
		global.Logger.Info("服务注销成功")
	} else {
		global.Logger.Error("服务注销失败")
	}
}

// DiscoverService 发现服务
func DiscoverService(serviceName string) ([]vo.Instance, error) {
	if global.NacosClient == nil {
		return nil, fmt.Errorf("nacos客户端未初始化")
	}
	
	nacosConfig := global.Config.Nacos
	
	// 服务发现参数
	params := vo.SelectInstancesParam{
		ServiceName: serviceName,
		GroupName:   nacosConfig.GroupName,
		HealthyOnly: true,
	}
	
	// 获取服务实例列表
	instances, err := global.NacosClient.SelectInstances(params)
	if err != nil {
		global.Logger.Error("服务发现失败", 
			zap.String("serviceName", serviceName), 
			zap.Error(err))
		return nil, err
	}
	
	return instances, nil
}

// GetServiceURL 获取服务URL
func GetServiceURL(serviceName string) (string, error) {
	instances, err := DiscoverService(serviceName)
	if err != nil {
		return "", err
	}
	
	if len(instances) == 0 {
		return "", fmt.Errorf("没有可用的服务实例")
	}
	
	// 简单的负载均衡策略（选择第一个实例）
	instance := instances[0]
	return "http://" + instance.Ip + ":" + strconv.FormatUint(instance.Port, 10), nil
}

// getLocalIP 获取本机IP
func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	
	return "", fmt.Errorf("无法获取本机IP")
}
