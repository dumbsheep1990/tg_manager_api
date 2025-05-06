# TG营销系统 - 后端API服务

## 项目简介

TG营销系统后端API服务是基于Golang开发的Telegram营销管理平台服务端，提供电报账号管理、任务调度、消息发送等功能的API接口。系统采用模块化架构设计，支持分布式部署，通过etcd管理任务状态，使用nacos实现服务注册发现，通过RabbitMQ与Python Worker协作处理Telegram API交互。

## 技术栈

- **核心框架**：Golang + Gin
- **数据库**：MySQL
- **缓存**：Redis
- **消息队列**：RabbitMQ
- **分布式协调**：etcd
- **服务注册发现**：nacos
- **日志系统**：zap

## 项目结构

```
tg_manager_api/
├── config/                # 配置文件目录
│   └── config.go          # 配置结构定义
├── core/                  # 核心组件
│   └── ...
├── global/                # 全局变量
│   └── ...
├── initialize/            # 初始化组件目录
│   ├── etcd.go            # etcd初始化与操作
│   ├── nacos.go           # nacos服务注册与发现
│   ├── db.go              # 数据库连接初始化
│   ├── redis.go           # Redis连接初始化
│   ├── rabbitmq.go        # RabbitMQ连接初始化
│   └── ...
├── model/                 # 数据模型定义
│   ├── account.go         # 账号模型
│   ├── account_group.go   # 账号分组模型
│   ├── task.go            # 任务模型
│   └── ...
├── services/              # 服务模块目录
│   ├── core/              # 服务集成和管理
│   │   └── service.go     # 服务组织和初始化
│   ├── account/           # 账号管理服务
│   │   └── service/       # 服务实现
│   ├── task/              # 任务管理服务
│   │   └── service/       # 服务实现
│   ├── tdata/             # TData处理服务
│   │   └── service/       # 服务实现
│   ├── message/           # 消息管理服务
│   │   └── service/       # 服务实现
│   └── dashboard/         # 仪表盘服务
│       └── service/       # 服务实现
├── test/                  # 单元测试目录
│   ├── account/           # 账号服务测试
│   ├── task/              # 任务服务测试
│   ├── tdata/             # TData服务测试
│   ├── message/           # 消息服务测试
│   ├── dashboard/         # 仪表盘服务测试
│   └── initialize/        # 初始化组件测试
├── config.toml            # 配置文件
├── go.mod                 # Go模块定义
├── go.sum                 # Go依赖版本锁定
└── main.go                # 程序入口文件
```

## 核心服务模块

### 账号服务 (account)
管理Telegram账号和账号分组，支持TData文件导入账号

**主要功能**：
- 账号CRUD操作
- 账号分组管理
- 账号状态检测
- TData账号导入

### 任务服务 (task)
管理各类营销任务的创建、执行、状态更新和终止

**主要功能**：
- 任务创建和配置
- 任务状态管理
- 任务进度跟踪
- 任务执行控制(启动/暂停/停止)

### TData服务 (tdata)
管理TData文件的上传、保存和导入处理

**主要功能**：
- TData文件上传
- TData文件存储管理
- TData导入状态跟踪
- 账号提取和处理

### 消息服务 (message)
管理消息模板和消息发送功能

**主要功能**：
- 消息模板创建和管理
- 消息发送请求处理
- 消息发送状态追踪
- 发送历史记录管理

### 仪表盘服务 (dashboard)
提供系统数据统计和概览

**主要功能**：
- 系统运行状态监控
- 账号统计数据
- 任务执行统计
- 消息发送统计

## 核心组件

### etcd
用于分布式任务状态管理和锁机制，确保任务在分布式环境下的一致性

**主要功能**：
- 保存和获取任务状态
- 分布式锁实现
- 服务健康检查

### nacos
实现服务注册发现，支持系统的横向扩展和负载均衡

**主要功能**：
- 服务注册
- 服务发现
- 配置管理

### RabbitMQ
作为Go后端与Python Worker的通信桥梁，实现任务分发

**主要功能**：
- 任务队列管理
- 消息发送请求分发
- Worker执行结果回调

## 安装与配置

### 环境要求
- Go 1.18+
- MySQL 5.7+
- Redis 6.0+
- RabbitMQ 3.8+
- etcd 3.5+
- nacos 2.0+

### 安装步骤

1. 克隆仓库:
```bash
git clone <repository-url>
cd tg_manager_api
```

2. 安装依赖:
```bash
go mod tidy
```

3. 配置环境:
编辑 `config.toml` 配置数据库、Redis、RabbitMQ、etcd和nacos等连接信息:

```toml
# 数据库配置
[mysql]
host = "127.0.0.1"
port = 3306
database = "tg_manager"
username = "root"
password = "password"

# Redis配置
[redis]
host = "127.0.0.1"
port = 6379
password = ""
db = 0

# 其他配置...
```

4. 启动服务:
```bash
go run main.go
```

## 接口文档

API接口文档通过Swagger自动生成，启动服务后访问:
```
http://localhost:8888/swagger/index.html
```

## 单元测试

运行所有测试:
```bash
go test ./test/...
```

运行特定模块测试:
```bash
go test ./test/account/...
go test ./test/task/...
```

## 部署

### 标准部署
```bash
go build -o tg_manager_api
./tg_manager_api
```

### Docker部署
```bash
docker build -t tg-manager-api .
docker run -p 8888:8888 tg-manager-api
```

## 架构图

```
                   ┌──────────────┐
                   │   前端应用    │
                   └───────┬──────┘
                           │
                           ▼
┌───────────────────────────────────────────────┐
│                  API服务层                     │
├───────────┬────────────┬───────────┬──────────┤
│ 账号服务   │  任务服务   │  消息服务  │ TData服务 │
└───────────┴────────────┴───────────┴──────────┘
                           │
                ┌──────────┴───────────┐
                │                      │
                ▼                      ▼
        ┌───────────────┐      ┌───────────────┐
        │    RabbitMQ   │      │     etcd      │
        └───────┬───────┘      └───────┬───────┘
                │                      │
                ▼                      │
        ┌───────────────┐              │
        │ Python Worker │              │
        └───────────────┘              │
                                       ▼
                                ┌───────────────┐
                                │     nacos     │
                                └───────────────┘
```
