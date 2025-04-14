# qqlx

## 简介

GO 编写的 WEB 框架，此框架基于流行的`Gin`框架构建。
前端地址：[https://github.com/qinquanliuxiang/go-web-react.git](https://github.com/qinquanliuxiang/go-web-react.git)

## 核心特性

1. 用户管理：用户注册、登录、资料更新等功能。
2. 集成 `ldap`, 用户同步 `ldap` 并且支持用户组。
3. 角色管理：支持多种角色定义，便于组织结构化权限分配。
4. 权限控制：利用`Casbin`实现精确到资源级别的访问控制。

## 技术栈

1. `Gin`：轻量级的Go web框架，提供路由、中间件支持等功能。
2. `GORM`：强大的对象关系映射（ORM）工具，简化数据库操作。
3. `Casbin`：功能丰富的访问控制库，实现`RBAC`访问模型。
4. `MySQL`：可靠的关系型数据库管理系统，适用于大规模数据存储需求。
5. `Redis`：高性能键值存储系统，用于加速角色信息检索过程。
6. `Wire`：依赖注入

## 接口文档

<https://apifox.com/apidoc/shared-8db2216b-8451-4ead-b0fd-019ce8676f1e>

## 初始化数据

```bash
go build -o .

# env
export CONFIG_PATH=configPath
export CASBIN_MODE_PATH=modelPath

# wire install
go install github.com/google/wire/cmd/wire@latest
# wire build
wire ./cmd/wire.go

# edit config
cp example.yaml config.yaml

# init data
./qqlx init
# init data with options
./qqlx init -C config.yaml -M model.conf
```

## **启动服务**

### Docker 启动

```bash
# 拷贝修改配置文件
cp example.yaml config.yaml

# 启动
docker compose up -d
```
