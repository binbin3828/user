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
-2 | mysql error
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


## 1.6. 互相添加好友


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

## 1.7. 查询好友列表


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

## 1.8. 查询附近的好友列表


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