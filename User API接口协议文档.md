# 1. User API 接口文档

## 1.1. API V1 接口说明

- 接口本地测试基准地址：`http://127.0.0.1:8080`

- 线上运行基准地址  `http://121.196.204.236:8080`

- 数据返回格式统一使用 JSON

### 1.1.1. 支持的请求方法

- GET（SELECT）：从服务器取出资源（一项或多项）。
- POST（CREATE）：在服务器新建一个资源。
- PUT（UPDATE）：在服务器更新资源（客户端提供改变后的完整资源）。
- DELETE（DELETE）：从服务器删除资源。


### 1.1.2. 通用返回状态说明

| *状态码* | *含义*                | *说明*                                              |
| -------- | --------------------- | --------------------------------------------------- |
| 200      | OK                    | 请求成功                                            |


### 1.1.3. 通用成功消息
```json
{
    "code": 0,
    "data": {
        "key1": "value1",
        "key2": "value2"
         ...
    }
}
```

### 1.1.4. 通用错误返回消息
```json
{
    "code": -1,
    "desc": "erros messge..."
}
```

### 1.1.5 code错误说明
code=0 表示成功  
code<0 表示失败

code错误码 | 说明
---|---
0 | 成功消息
-1 | param error
-2 | mysql/internal error
-3 | auth fail
-4 | permission denied
-5 | already friends
-6 | friend request already exists
-7 | friend request not found
-8 | friend request no longer pending
... | ...



------

## 1.2. 创建用户


- 请求路径：http://127.0.0.1:8080/user
- 请求方法：post
- 请求参数

| 参数名   | 参数说明 | 是否必传     | 字段类型 |
| -------- | -------- | -------- |------|
| name | 用户名   | 是 | string|
| dob | 生日     | 否 |string|
| address | 地址     | 否 |string|
| description | 描述     | 否 |string|

- 请求示例：
```json
{
    "name":"bobby哥",
    "dob":"1990-01-10",
    "address":"sz",
    "description":"i am a boy, a coder"
}
```

- 响应参数

| 参数名   | 参数说明    | 字段类型     |
| -------- | ----------- | ------------ |
| id       | 用户 ID     |        int64 |
| dob      | 生日        |        string|
| address  | 地址        |        string|
| description   | 描述   |        string|
| name     | 姓名        |        string|
| create_at| 创建时间    |        string|

- 响应成功数据

```json
{
    "code": 0,
    "data": {
        "id": 12,
        "name": "bobby哥",
        "dob": "1990-01-10",
        "address": "sz",
        "description": "i am a boy, a coder",
        "create_at": "2022-06-09 05:03:53"
    }
}
```

- 响应失败示例
```json
{
    "code":-1
    "desc": "param name not set"
}
```



## 1.3. 查询用户


- 请求路径：http://127.0.0.1:8080/user/{uid}
- 请求方法：get
- 请求参数

| 参数名   | 参数说明     | 是否必传     |
| -------- | ------------ | -------- |
| uid    | 玩家uid        | 是       |

- 请求示例

http://127.0.0.1:8080/user/12

- 响应参数

| 参数名   | 参数说明    | 字段类型     |
| -------- | ----------- | ------------ |
| id       | 用户 ID     |        int64 |
| dob      | 生日        |        string|
| address  | 地址        |        string|
| description   | 描述   |        string|
| name     | 姓名        |        string|
| create_at| 创建时间    |        string|
- 响应数据

```json
{
    "code": 0,
    "data": {
        "id": 12,
        "name": "bobby哥",
        "dob": "1990-01-10",
        "address": "sz",
        "description": "i am a boy, a coder",
        "create_at": "2022-06-09 05:03:53"
    }
}
```

## 1.4. 修改用户信息


- 请求路径：http://127.0.0.1:8080/user
- 请求方法：put
- 请求参数

| 参数名   | 参数说明    | 是否必传     | 字段类型 |
| -------- | ----------- | ------------ |--------|
| id       | 用户 ID     |        是 | int64|
| dob      | 生日        |        否| string|
| address  | 地址        |        否|string|
| description   | 描述   |        否|string|
| name     | 姓名        |        否|string|

请求示例:
```json
{
    "id":12,
    "name": "guobin"
}
```

- 响应参数

| 参数名   | 参数说明    | 字段类型     |
| -------- | ----------- | ------------ |
| id       | 用户 ID     |        int64 |
| dob      | 生日        |        string|
| address  | 地址        |        string|
| description   | 描述   |        string|
| name     | 姓名        |        string|
| create_at| 创建时间    |        string|
- 响应数据

```json
{
    "code": 0,
    "data": {
        "id": 12,
        "name": "guobin",
        "dob": "1990-01-10",
        "address": "sz",
        "description": "i am a boy, a coder",
        "create_at": "2022-06-09 05:03:53"
    }
}
```
- 错误示例
```json
{
    "code": -1,
    "msg": "user id is must param"
}
```

## 1.5. 删除用户


- 请求路径：http://127.0.0.1:8080/user/{uid}
- 请求方法：delete
- 请求参数

| 参数名   | 参数说明    | 是否必传     | 字段类型 |
| -------- | ----------- | ------------ |--------|
| id       | 用户 ID     |        是    | int64|


请求示例:

http://127.0.0.1:8080/user/12

- 响应参数

| 参数名   | 参数说明    | 字段类型     |
| -------- | ------------ | ------------ |
| data       | 成功消息说明|        string |

- 响应成功数据

```json
{
    "code": 0,
    "data": "delete succ"
}
```
- 错误示例
```json
{
    "code": -2,
    "msg": "delete user error in db"
}
```


## 1.6. 发起好友请求

- 请求路径：http://127.0.0.1:8080/v1/friend-requests
- 请求方法：POST
- 认证：需要 Bearer Token
- 请求参数

| 参数名   | 参数说明     | 是否必传     | 字段类型 |
| -------- | ------------ | ------------ | -------- |
| to_uid   | 目标用户 uid | 是           | int      |

- 请求示例：
```json
{
    "to_uid": 11
}
```

- 响应成功数据
```json
{
    "code": 0,
    "data": {
        "id": 1,
        "from_uid": 10,
        "to_uid": 11,
        "status": "pending",
        "message": "friend request sent"
    }
}
```

- 错误示例

| 错误码 | 说明 |
|--------|------|
| -1 | 参数错误（不能添加自己为好友） |
| -5 | 已经是好友 |
| -6 | 已发送过待处理的好友请求 |
| 0 | 若对方也已向你发起请求，则自动成为好友（status: accepted） |


## 1.7. 查询收到的好友请求

- 请求路径：http://127.0.0.1:8080/v1/friend-requests/incoming
- 请求方法：GET
- 认证：需要 Bearer Token
- 请求参数

| 参数名    | 参数说明     | 是否必传     | 默认值   |
| --------- | ------------ | ------------ | -------- |
| status    | 状态筛选     | 否           | pending  |
| page      | 页码         | 否           | 1        |
| page_size | 每页条数     | 否           | 20       |

- 请求示例：
```
http://127.0.0.1:8080/v1/friend-requests/incoming?status=pending&page=1&page_size=20
```

- 响应成功数据
```json
{
    "code": 0,
    "data": [
        {
            "id": 1,
            "from_uid": 9,
            "to_uid": 10,
            "status": "pending",
            "created_at": "2024-01-15 10:30:00",
            "updated_at": "2024-01-15 10:30:00"
        }
    ],
    "pagination": {
        "total": 1,
        "page": 1,
        "page_size": 20,
        "total_pages": 1
    }
}
```


## 1.8. 查询发出的好友请求

- 请求路径：http://127.0.0.1:8080/v1/friend-requests/outgoing
- 请求方法：GET
- 认证：需要 Bearer Token
- 请求参数

| 参数名    | 参数说明     | 是否必传     | 默认值   |
| --------- | ------------ | ------------ | -------- |
| status    | 状态筛选     | 否           | pending  |
| page      | 页码         | 否           | 1        |
| page_size | 每页条数     | 否           | 20       |

- 响应格式同收到的好友请求列表


## 1.9. 同意好友请求

- 请求路径：http://127.0.0.1:8080/v1/friend-requests/{id}/accept
- 请求方法：PUT
- 认证：需要 Bearer Token（仅请求接收者可操作）
- 请求参数

| 参数名 | 参数说明         | 是否必传     | 字段类型 |
| ------ | ---------------- | ------------ | -------- |
| id     | 好友请求 ID      | 是           | int      |

- 请求示例：
```
PUT http://127.0.0.1:8080/v1/friend-requests/1/accept
```

- 响应成功数据
```json
{
    "code": 0,
    "data": {
        "id": 1,
        "from_uid": 9,
        "to_uid": 10,
        "status": "accepted"
    }
}
```

- 错误示例

| 错误码 | 说明 |
|--------|------|
| -4 | 无权限（非请求接收者） |
| -7 | 好友请求不存在 |
| -8 | 好友请求已处理（非 pending 状态） |


## 1.10. 拒绝好友请求

- 请求路径：http://127.0.0.1:8080/v1/friend-requests/{id}/reject
- 请求方法：PUT
- 认证：需要 Bearer Token（仅请求接收者可操作）
- 请求参数

| 参数名 | 参数说明         | 是否必传     | 字段类型 |
| ------ | ---------------- | ------------ | -------- |
| id     | 好友请求 ID      | 是           | int      |

- 响应成功数据
```json
{
    "code": 0,
    "data": {
        "id": 1,
        "from_uid": 9,
        "to_uid": 10,
        "status": "rejected"
    }
}
```


## 1.11. 互相添加好友（直接）


- 请求路径：http://127.0.0.1:8080/friends
- 请求方法：post
- 请求参数

| 参数名   | 参数说明 | 是否必传     | 字段类型 |
| -------- | -----------  | -------- |------|
| uid      | 我的uid      | 是       | int64|
| fri      | 添加的好友uid| 否       | int64|

- 请求示例：
```json
{
    "uid":10,
    "fri":11,
}
```

- 响应参数

| 参数名   | 参数说明    | 字段类型     |
| -------- | ----------- | ------------ |
| uid      | 我的uid      |        int64 |
| fri      | 添加的好友uid |        int64 |


- 响应成功数据

```json
{
    "code": 0,
    "data": {
        "uid":10,
        "fri":11,
    }
}
```
- 错误示例
```json
{
    "code": -2,
    "msg": "create error in db"
}
```

## 1.12. 查询好友列表


- 请求路径：http://127.0.0.1:8080/friends/{uid}
- 请求方法：get
- 请求参数

| 参数名   | 参数说明     | 是否必传     |
| -------- | ------------ | -------- |
| uid    | 玩家uid        | 是       |

- 请求示例

http://127.0.0.1:8080/friends/10

- 响应参数

| 参数名   | 参数说明    | 字段类型     |
| -------- | ----------- | ------------ |
| list     | 好友数组        |         arr |
| fri_uid  | 好友uid         |        int64|
| fri_name | 好友名字        |        string|
| create_at| 好友关系建立时间 |        string|
- 响应数据

```json
{
    "code": 0,
    "data": {
        "list": [
            {
                "fri_uid": 9,
                "fri_name": "bobby9",
                "create_at": "2022-06-09 16:42:29"
            },
            {
                "fri_uid": 11,
                "fri_name": "",
                "create_at": "2022-06-09 16:32:06"
            }
        ],
        "uid": 10
    }
}
```

## 1.13. 查询附近的好友列表


- 请求路径：http://127.0.0.1:8080/nearbyfriends/{uid}
- 请求方法：get
- 请求参数

| 参数名   | 参数说明     | 是否必传     |
| -------- | ------------ | -------- |
| uid    | 玩家uid        | 是       |

- 请求示例

http://127.0.0.1:8080/nearbyfriends/10

- 响应参数

| 参数名    | 参数说明    | 字段类型     |
| ---------| -----------| ------------ |
| list         | 好友数组     |      arr     |
| fri_uid     | 好友uid      |      int64   |
| fri_name    | 好友名字     |      string  |
| latitude    | 纬度         |      float64 |
| longitude   | 经度         |      float64 |
| loc_geohash | geo哈希算法字符串 |  string  |
| create_at   | 好友关系建立时间  |  string  |
- 响应数据

```json
{
	"code": 0,
	"data": {
		"list": [{
			"fri_uid": 9,
			"fri_name": "bobby9",
			"create_at": "2022-06-09 16:42:29",
			"latitude": 39.910935,
			"longitude": 116.4133,
			"loc_geohash": "wx4g119d"
		}, {
			"fri_uid": 11,
			"fri_name": "",
			"create_at": "2022-06-09 16:32:06",
			"latitude": 39.910987,
			"longitude": 116.413311,
			"loc_geohash": "wx4g119d"
		}],
		"uid": 10
	}
}
```

## 1.14. 发现附近陌生人

- 请求路径：http://127.0.0.1:8080/v1/nearby-users/{uid}
- 请求方法：GET
- 认证：需要 Bearer Token
- 说明：基于当前用户的地理位置（Geohash），推荐附近非好友、无待处理请求的用户。已排除自己、好友、pending 请求的用户。
- 请求参数

| 参数名    | 参数说明             | 是否必传     | 默认值   |
| --------- | -------------------- | ------------ | -------- |
| uid       | 用户 UID（路径参数） | 是           | -        |
| precision | Geohash 精度（1-12）  | 否           | 6        |
| page      | 页码                 | 否           | 1        |
| page_size | 每页条数             | 否           | 20       |

- 请求示例：

http://127.0.0.1:8080/v1/nearby-users/10?precision=6&page=1&page_size=20

- 响应成功数据
```json
{
    "code": 0,
    "data": [
        {
            "fri_uid": 9,
            "fri_name": "bobby9",
            "latitude": 39.910935,
            "longitude": 116.4133,
            "loc_geohash": "wx4g119d"
        }
    ],
    "pagination": {
        "total": 1,
        "page": 1,
        "page_size": 20,
        "total_pages": 1
    }
}
```

- 空结果（用户未设置位置）
```json
{
    "code": 0,
    "data": []
}
```


## 1.15. 拉黑用户

- 请求路径：http://127.0.0.1:8080/v1/blacklist
- 请求方法：POST
- 认证：需要 Bearer Token
- 说明：拉黑后，对方无法查看你的资料、向你发起好友请求，也不会出现在附近陌生人推荐中。
- 请求参数

| 参数名      | 参数说明         | 是否必传 | 字段类型 |
| ----------- | ---------------- | -------- | -------- |
| blocked_uid | 要拉黑的用户 UID  | 是       | int      |

- 请求示例：
```json
{
    "blocked_uid": 11
}
```

- 响应成功数据
```json
{
    "code": 0,
    "data": {
        "uid": 10,
        "blocked_uid": 11,
        "status": "blocked"
    }
}
```

- 错误示例

| 错误码 | 说明 |
|--------|------|
| -1 | 参数错误（不能拉黑自己、用户已拉黑） |


## 1.16. 取消拉黑

- 请求路径：http://127.0.0.1:8080/v1/blacklist/{uid}
- 请求方法：DELETE
- 认证：需要 Bearer Token
- 请求参数

| 参数名 | 参数说明           | 是否必传 | 字段类型 |
| ------ | ------------------ | -------- | -------- |
| uid    | 要取消拉黑的用户 UID | 是       | int      |

- 请求示例：
```
DELETE http://127.0.0.1:8080/v1/blacklist/11
```

- 响应成功数据
```json
{
    "code": 0,
    "data": {
        "uid": 10,
        "blocked_uid": 11,
        "status": "unblocked"
    }
}
```


## 1.17. 查看黑名单

- 请求路径：http://127.0.0.1:8080/v1/blacklist
- 请求方法：GET
- 认证：需要 Bearer Token
- 请求参数

| 参数名    | 参数说明 | 是否必传 | 默认值 |
| --------- | -------- | -------- | ------ |
| page      | 页码     | 否       | 1      |
| page_size | 每页条数 | 否       | 20     |

- 请求示例：
```
http://127.0.0.1:8080/v1/blacklist?page=1&page_size=20
```

- 响应成功数据
```json
{
    "code": 0,
    "data": [
        {
            "uid": 10,
            "blocked_uid": 11,
            "created_at": "2024-01-15 10:30:00"
        }
    ],
    "pagination": {
        "total": 1,
        "page": 1,
        "page_size": 20,
        "total_pages": 1
    }
}
```


## 1.18. 登录

- 请求路径：http://127.0.0.1:8080/v1/auth/login
- 请求方法：POST
- 认证：否
- 限流：同一 IP 每分钟限 10 次
- 请求参数

| 参数名   | 参数说明 | 是否必传     | 字段类型 |
| -------- | -------- | ------------ | -------- |
| name     | 用户名   | 是           | string   |
| password | 密码     | 是           | string   |

- 请求示例：
```json
{
    "name": "bobby哥",
    "password": "mypassword"
}
```

- 响应成功数据
```json
{
    "code": 0,
    "data": {
        "token": "eyJhbGciOi...",
        "user_id": 12
    }
}
```


## 1.19. 忘记密码

- 请求路径：http://127.0.0.1:8080/v1/auth/forgot-password
- 请求方法：POST
- 认证：否
- 说明：输入注册邮箱，系统生成有时限的重置令牌。开发环境直接返回 token 到响应中，生产环境通过 SMTP 发送邮件。
- 请求参数

| 参数名 | 参数说明 | 是否必传     | 字段类型 |
| ------ | -------- | ------------ | -------- |
| email  | 注册邮箱 | 是           | string   |

- 请求示例：
```json
{
    "email": "user@example.com"
}
```

- 响应成功数据
```json
{
    "code": 0,
    "data": {
        "message": "if the email is registered, a reset link has been sent",
        "token": "a1b2c3d4e5f6..."
    }
}
```

- 注意：无论邮箱是否已注册，都返回相同消息，防止用户枚举。


## 1.20. 重置密码

- 请求路径：http://127.0.0.1:8080/v1/auth/reset-password
- 请求方法：POST
- 认证：否
- 请求参数

| 参数名       | 参数说明                   | 是否必传     | 字段类型 |
| ------------ | -------------------------- | ------------ | -------- |
| token        | 忘记密码接口返回的重置令牌 | 是           | string   |
| new_password | 新密码（最少 8 位）        | 是           | string   |

- 请求示例：
```json
{
    "token": "a1b2c3d4e5f6...",
    "new_password": "newpassword123"
}
```

- 响应成功数据
```json
{
    "code": 0,
    "data": {
        "message": "password has been reset successfully"
    }
}
```

- 错误示例

| 错误码 | 说明 |
|--------|------|
| -1 | 令牌无效或已过期 |


## 1.21. 批量查询在线状态

- 请求路径：http://127.0.0.1:8080/v1/users/online
- 请求方法：GET
- 认证：需要 Bearer Token
- 说明：基于 Redis 心跳检测，认证后的每个请求自动续期 5 分钟。未配置 Redis 时全部返回 false。
- 请求参数

| 参数名 | 参数说明           | 是否必传 | 字段类型 |
| ------ | ------------------ | -------- | -------- |
| uids   | 用户 UID 数组       | 是       | []int    |

- 请求示例：
```
http://127.0.0.1:8080/v1/users/online?uids=9&uids=10&uids=11
```

- 响应成功数据
```json
{
    "code": 0,
    "data": [
        { "uid": 9,  "is_online": true },
        { "uid": 10, "is_online": false },
        { "uid": 11, "is_online": true }
    ]
}
```