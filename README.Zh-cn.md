# user
A simple restful API framework,  a mysql CURD example for  user table

一个自己实现的简单的RESTFULL API框架，实现了对USER表的CRUD操作

### 特性

- 遵循 RESTful API 设计规范

- 遵循 MVC 接口代码规范

- 基于 GORM 的数据库存储，可扩展多种类型数据库
 
- 基于 GORM 数据库连接池

- 配置文件采用 .yaml

- 支持滚动日志

- 支持 makefile

- 支持 test 单元测试



### 获取代码


```bash
# 获取接口代码
git clone https://github.com/binbin3828/user.git

```

### 代码结构目录





### 启动说明(linux环境)

```bash
# 进入 user 后端项目
cd ./user

# 编译项目, 生成 UserServer 二进制执行文件
cd ./cmd
make 

# 修改配置,修改数据库连接信息  
# 文件路径  ./user/pkg/config/config.yaml 
vi ./pkg/config/config.yaml 

#启动
cd ./cmd
./UserServer   #前台启动
./UserServer & #后台启动

```


### 交叉编译

```bash
# windows
env GOOS=windows GOARCH=amd64 go build main.go

# or
# linux
env GOOS=linux GOARCH=amd64 go build main.go
```



### 在线体验
在线接口地址

http://21.196.204.236:8080/user


###  License

Copyright (c) 2022 guobin