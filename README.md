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

### 实现功能

- 对 User表的 增删改查

- zap实现的滚动日志，gorm日志输出（高亮，打印sql执行时间）

- 添加好友，好友列表

- geohash算法实现附近的人



### 获取代码


```bash
# 获取接口代码
git clone https://github.com/AmberGroup-WhaleFin/go-backend-test-guobin.git

# 切换到 guobin-user分支
git checkout -b guobin-user origin/guobin-user

```

### 代码结构目录

    ├── cmd     ---------------------------------------------------main文件入口，编译目录
    │   ├── main.go               // main文件
    │   ├── makefile              // makefile
    │   └── UserServer            // make 后生成的二进制执行文件
    ├── constant  ------------------------------------------------- 定义的常量/错误码目录
    │   └── ErrorCode.go          // 错误码定义
    ├── dao       ------------------- ----------------------------- dao层(数据的访问和操作)
    │   ├── FriendsDao.go         // 好友dao
    │   └── UserDao.go            // 用户dao
    ├── go.mod                    // go mod包依赖管理
    ├── go.sum                    // go mod包依赖管理
    ├── log       -------------------------------------------------  日志目录
    │   └── user.log              // 日志文件滚动
    ├── model     -------------------------------------------------- 数据模型
    │   ├── Friends.go            //好友 model
    │   └── User.go               //用户 model
    ├── pkg       -------------------------------------------------- 通用工具包
    │   ├── config       ------------------------------------------- 配置文件相关            
    │   │   ├── config.go         //配置文件管理包
    │   │   └── config.yaml       //配置文件.yaml
    │   ├── dbconn                ----------------------------------- 数据库连接管理  
    │   │   └── mysqlconn.go      //数据库连接
    │   ├── logger                -----------------------------------日志管理
    │   │   ├── GormLogger.go     // goorm日志
    │   │   └── Logger.go         // zip日志管理
    │   └── util                  ------------------------------------其他工具包
    │       ├── JsonTime.go       // json转换时间工具
    │       └── net.go            // 网络消息返回定义
    ├── README.md                 ------------------------------------- 说明文档
    ├── service                   ------------------------------------- 接口服务
    │   ├── FriendsServer.go      // 好友相关接口文件
    │   ├── Routes.go             // 接口路由
    │   └── UserService.go        // 用户相关接口文件
    ├── test                      --------------------------------------单元测试
    │   ├── AddFriends_test.go          //添加好友
    │   ├── CreateUser_test.go          //创建用户
    │   ├── DeleteUser_test.go          //删除用户
    │   ├── GetFriendsList_test.go      //好友列表
    │   ├── GetNearbyFriendList_test.go //附近的好友
    │   ├── GetUser_test.go             //获取用户信息
    │   └── ModifyUser_test.go          //修改用户信息
    └── UserAPI接口协议文档.md     //接口协议文档


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

### 在线体验
在线接口地址

http://121.196.204.236:8080/


###  License

Copyright (c) 2022 guobin