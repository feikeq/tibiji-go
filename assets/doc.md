# 全局公共参数

**全局Header参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| TIBIJI_SERV_ENV | development | String | 是 | 环境变量判断生产和测试环境 |
| Authorization | aAo5aky0AitM4izCyqh4JgW+YUpztJfNNi6xbHpxNS9IcODUbA+l8c7WPf5wAj9+TYJMJvKFFDyG3/qMKzJWXrJ0gcyVlsda3II= | String | 是 | 自定义密钥认证方式 |

**全局Query参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**全局Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**全局认证方式**

> Bearer Token

> 在Header添加参数 Authorization，其值为在Bearer之后拼接空格和访问令牌

> Authorization: Bearer your_access_token

# 状态码说明

| 状态码 | 中文描述 |
| --- | ---- |
| 200 (http) | 请求成功 |
| 401 (http) | 无法访问 |
| 10000 | 未知错误 |
| 10001 | 错误的请求 |
| 10004 | 未授权访问 |
| 10007 | 内部服务器错误 |
| 10010 | 请求的资源不存在 |
| 10020 | 传入的参数为空或者不合法 |
| 10030 | 请求过于频繁请稍后再试 |
| 20000 | 数据库错误 |
| 20010 | 没有找到指定类型的记录 |
| 20020 | 资源已存在 |
| 20022 | 闲置空资源待认领 |
| 20040 | 帐号已被禁用 |
| 20050 | 帐号无权限操作 |
| 20060 | 帐号尚未激活 |
| 20100 | 用户名或密码错误 |
| 20104 | 授权标头为空 |
| 20200 | Token错误或过期 |
| 20301 | 验证码错误 |
| 20302 | 验证码已过期 |
| 20400 | 格式不正确 |
| 20500 | 您不是VIP或已过期 |

# 用户中心

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-05-24 11:18:19

> 更新时间: 2023-05-24 12:06:15

```text
暂无描述
```

**目录Header参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**目录Query参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**目录Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**目录认证信息**

> 继承父级

## 用户登录

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-05-22 09:42:21

> 更新时间: 2024-09-02 16:23:16

**支持 用户名 、手机号 、邮箱地址 、身份证号 多种方式进行登录。,注意：当开启 SERV_REGISTER_CHECK是否开启验证 时 ，且登录方式是用邮箱/手机时 code 和 ticket 才为必填字段，当为 用户名/身份证 登录时可不用验证码就可以登录。**

**接口状态**

> 已完成

**接口URL**

> /user/login/

**请求方式**

> POST

**Content-Type**

> urlencoded

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| name | ccavgogogo | String | 是 | 用户名/邮箱/手机/身份证 |
| pwd | liudehuadhah | String | 是 | 密码 |
| code | 374627 | String | 否 | 验证码（开启验证时必填） |
| ticket | J8fm9IISDIMmyr8pxrwLDWZEgY0AJHPpnQl6DNI33X6gisepEXreVsSHn+oMItgktr0BDqxiyB5FUE8= | String | 否 | 验证令牌（开启验证时必填） |

**认证方式**

> 继承父级

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": {
		"address": "",
		"balance": 88.88,
		"bankcard": "",
		"birthday": "0001-01-01",
		"cell": "1-5****8888",
		"cid": 0,
		"city": "",
		"company": "",
		"country": "",
		"email": "lsd****qdi@ccav.tv",
		"exptime": "2024-12-12 12:12:12",
		"fname": "",
		"grouptag": "",
		"headimg": "",
		"identity_card": "",
		"intime": "2012-08-08 08:08",
		"inviter": 0,
		"manage": 2,
		"material_object": "{123}",
		"material_remark": "用户充值",
		"material_uptime": "2024-08-07 16:00:28",
		"nickname": "",
		"object": "hahahahah",
		"password": true,
		"province": "",
		"referer": "",
		"regip": "192.168.129.18",
		"remark": "",
		"sex": 0,
		"state": 1,
		"tag": "010201020",
		"token": "YV9qOhvhAigcty3AyPp7dw7uYRUg55mYNyviPn56P3sUdOHSbgumpcGBOhZQ/aAZkFSLDH6yysJS934=",
		"uid": 12938,
		"uptime": "2023-06-08 13:54",
		"userclan": "",
		"username": "oggoccav",
		"vip": 3
	},
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": "err debug",
	"data": {
		"name": "admin",
		"pwd": "admin2"
	},
	"msg": "密码错误"
}
```

* 验证失败(401)

```javascript
{
	"code": 20108,
	"msg": "验证码错误"
}
```

* 过期(401)

```javascript
{
	"code": 20109,
	"msg": "验证码已过期"
}
```

* 禁用(200)

```javascript
{
	"code": 20010,
	"msg": "当前登录账户已被禁用，请联系管理员"
}
```

* 激活(200)

```javascript
{
	"code": 20013,
	"msg": "尚未激活"
}
```

## 新增用户

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-05-10 13:56:16

> 更新时间: 2024-08-07 16:19:26

**注册用户也就是新增添加用户的基本数据项。,注意：当开启 SERV_REGISTER_CHECK是否开启验证 时 ，code和ticket 为必填字段，且username只支持手机和邮箱注册才能成功发送验证码，否则使用 用户名/身份证 是拿不到验证码的。,开启验证后添加成功后会返回新增用户的登录数据。,管理员可直接添加用户跳过验证且不会返回用户登录数据**

**接口状态**

> 已完成

**接口URL**

> /user/

**请求方式**

> POST

**Content-Type**

> urlencoded

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| ticket | 8/2mW9KRcMAOWf0sWxvWq3s+BoAtJQnwyTYlc6dJZb2o7AGne0aUmk3WtAtrVkOIpSmbafanrJrZqvE= | String | 否 | 验证令牌（开启验证时必填） |
| code | 372827 | String | 否 | 验证码（开启验证时必填） |
| username | test00233 | String | 是 | 用户名 |
| ciphers | 123456789 | String | 否 | 密码 |
| email | test@163.com | String | 否 | 邮箱 |
| cell | 13838389438 | String | 否 | 电话 |
| nickname | 测试号2 | String | 否 | 昵称 |
| headimg | - | String | 否 | 头像 |
| sex | 2 | Integer | 否 | 性别(0未知 1男 2女) |
| birthday | 1999-09-09 | String | 否 | 出生日期 |
| company | 清北 | String | 否 | 公司 |
| address | 东三环中路 | String | 否 | 地址 |
| city | 长沙 | String | 否 | 城市 |
| province | 湖南 | String | 否 | 省份 |
| country | 中国 | String | 否 | 国家 |
| referer | API测试 | String | 否 | 用户来源 |
| inviter | 1684899785994911 | Number | 否 | 邀请者UID |
| fname | 张三 | String | 否 | 真实姓名 |
| bankcard | - | String | 否 | 银行卡 |
| identity_card | - | String | 否 | 身份证 |
| grouptag | tag | String | 否 | 用户组 |
| remark | .... | String | 否 | 备注 |
| object | {} | String | 否 | 预留字段 |

**认证方式**

> Bearer Token

> 在Header添加参数 Authorization，其值为在Bearer之后拼接空格和访问令牌

> Authorization: Bearer your_access_token

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": 1719306990345093,
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": 20000,
	"msg": "数据库错误"
}
```

* 验证码错误(401)

```javascript
{
	"code": 20108,
	"msg": "验证码错误"
}
```

* 过期(401)

```javascript
{
	"code": 20109,
	"msg": "验证码已过期"
}
```

* 注册并登录(200)

```javascript
{
	"code": 0,
	"data": {
		"address": "东三环中路",
		"bankcard": "",
		"birthday": "1999-09-09",
		"cell": "",
		"city": "长沙",
		"company": "清北",
		"country": "中国",
		"email": "",
		"fname": "张*",
		"grouptag": "tag",
		"headimg": "",
		"identity_card": "",
		"intime": "2024-07-01 11:57",
		"inviter": 1684899785994911,
		"nickname": "测试号2",
		"object": "{}",
		"password": true,
		"province": "湖南",
		"referer": "API测试",
		"regip": "192.168.172.88",
		"remark": "....",
		"sex": 2,
		"state": 1,
		"token": "aAo6ZkC3AilJ7y3DzaF/JlW5ZxYvts7ON362Oip3PilAduSBOA2io5rRPq8nBG95TYJPKv2FGDqH3sNBFE4lQAKpWDHvPbo4ODo=",
		"uid": 1719806279610972,
		"uptime": "2024-07-01 11:57",
		"userclan": "1684899785994911",
		"username": "test001"
	},
	"msg": ""
}
```

## 用户信息

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-05-10 13:51:51

> 更新时间: 2024-08-07 15:41:53

**获取指定用户UID的基本资料，非管理员只能获取自己的用户详情**

**接口状态**

> 已完成

**接口URL**

> /user/{uid}/

**请求方式**

> GET

**Content-Type**

> none

**认证方式**

> Bearer Token

> 在Header添加参数 Authorization，其值为在Bearer之后拼接空格和访问令牌

> Authorization: Bearer your_access_token

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": {
		"address": "东三环中路",
		"bankcard": "",
		"birthday": "1999-09-09",
		"cell": "138****9438",
		"city": "长沙",
		"company": "清北",
		"country": "中国",
		"email": "t**t@163.com",
		"fname": "张*",
		"grouptag": "tag",
		"headimg": "",
		"identity_card": "",
		"intime": "2023-05-24 11:44",
		"inviter": 1684899785994911,
		"nickname": "测试号",
		"object": "{}",
		"password": true,
		"province": "湖南",
		"referer": "API测试",
		"regip": "127.0.0.1",
		"remark": "....",
		"sex": 2,
		"state": 1,
		"uid": 1684899854801307,
		"userclan": "1684899785994911,1684899408028392,",
		"username": "test"
	},
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": 20105,
	"msg": "Token错误或过期"
}
```

## 修改用户

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-05-19 14:42:15

> 更新时间: 2024-08-07 14:30:34

**修改指定用户UID的基本资料（因条件会带当前用户ID）**

**接口状态**

> 已完成

**接口URL**

> /user/{uid}/

**请求方式**

> PUT

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| nickname | 哈哈哈 | String | 否 | 昵称 |
| headimg | - | String | 否 | 头像 |
| sex | 1 | Integer | 否 | 性别(0未知 1男 2女) |
| birthday | 1990-01-01 | String | 否 | 出生日期 |
| company | - | String | 否 | 公司 |
| address | - | String | 否 | 地址 |
| city | - | String | 否 | 城市 |
| province | - | String | 否 | 省份 |
| country | - | String | 否 | 国家 |
| referer | - | String | 否 | 用户来源  ( *管理员才能修改 ) |
| inviter | 333 | String | 否 | 邀请者UID  ( *管理员才能修改 ) |
| fname | - | String | 否 | 真实姓名 |
| bankcard | - | String | 否 | 银行卡 |
| grouptag | - | String | 否 | 用户组  ( *管理员才能修改 ) |
| remark | - | String | 否 | 备注 |
| object | - | String | 否 | 预留字段 |
| state | 2 | String | 否 | 状态 ( *管理员才能修改 ) |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": 1,
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": 0,
	"data": 0,
	"msg": ""
}
```

## 更新用户

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-05-25 09:35:16

> 更新时间: 2024-09-10 17:34:17

**只有自己或管理员才能操作,更新用户的密码、邮箱、电话、用户名、身份证 ，每次从只能更新其中一项数据！,注意：当开启 SERV_REGISTER_CHECK是否开启验证 时 ，且需要更新的是邮箱或手机时 code 和 ticket 才为必填字段，当更新密码时要填写pwe原来的旧密码正确才能成功（否则只能使用忘记密码找回密码来重置密码），当更新 用户名或身份证 可不用验证码就可以直接更新。**

**接口状态**

> 已完成

**接口URL**

> /user/{uid}/

**请求方式**

> PATCH

**Content-Type**

> form-data

**路径变量**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| uid | 1725436024707002 | String | 是 | - |

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| ciphers | 111 | String | 否 | 新密码 |
| pwd | 111 | String | 是 | 原密码（当更新密码时必填） |
| email | 22121@1212.com | String | 否 | 邮箱 |
| cell | 13838389438 | String | 否 | 电话 |
| identity_card | 2222 | String | 否 | 身份证(当实名认证时可传fname) |
| username | kkkkkkk | String | 否 | 帐号 |
| code | 219936 | String | 是 | 验证码（开启验证时必填） |
| ticket | MVn9A+zMAqf/iIqoGPnHPi58MfXJU1BjZOMI2EAXbblEw1zcEyv3DRAB71ajI8Lyf4sWcTs4u4H85Ck= | String | 是 | 验证令牌（开启验证时必填） |
| fname | - | String | 否 | 真实姓名(当修改身份证时才能传它) |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": "email",
	"msg": "ok"
}
```

* 失败(200)

```javascript
{
	"code": 10005,
	"msg": "传入的参数为空或者不合法"
}
```

* 原密码错误(200)

```javascript
{
	"code": 20400,
	"data": "ciphers",
	"msg": "格式不正确"
}
```

## 删除用户

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-05-26 10:12:23

> 更新时间: 2024-08-08 11:36:56

**通常是管理操作，用户只能删除自己也就是注销用户了。,data返回当前用户ID则是逻辑删除、返回1则是物理删除（判断当前用户状态为停用时自动进行物理删除），* 物理删除时用户附属资料表数据虽没有删除但无关联就相当于此条数据已废弃！**

**接口状态**

> 已完成

**接口URL**

> /user/{uid}/

**请求方式**

> DELETE

**Content-Type**

> none

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": 1684999387650990,
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": "err debug",
	"data": 0,
	"msg": "sql: no rows in result set"
}
```

## 忘记密码

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-05-31 10:55:24

> 更新时间: 2024-09-09 10:53:07

**限制访问频率
接口会返回验证令牌在data数据里，同时将通过邮件和短信发送6位数字的验证码给用户,尊敬的提笔记用户您好：
您的用户名：tbj
您的邮箱：fw@126.com
您的电话：+868943
请务必在0.5小时内通过下面这个地址修改您的密码，此链接将在2023-06-01 14:35:04后失效！,您的的验证码: 504578,提笔记安全中心 2023-06-01 14:05:04**

**接口状态**

> 已完成

**接口URL**

> /user/password/

**请求方式**

> GET

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| name | tibiji_test | String | 是 | 用户名/邮箱/手机/身份证 |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
    "code": 0,
    "data": "MqapC0UeruLH9rBq4MggXc8Di6GlceEve0jtrEyZK9OKzzECsh4SLEl3hnbVrFuEcAKuFyu9N93IIe0=",
    "msg": "c**v@ccav.tv 199****9999"
}
```

* 失败(200)

```javascript
{
	"code": 20001,
	"msg": "没有找到指定类型的记录"
}
```

## 找回密码

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-05-31 13:43:39

> 更新时间: 2024-09-09 13:44:11

**已限制访问频率**

**接口状态**

> 已完成

**接口URL**

> /user/password/

**请求方式**

> POST

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| pwd | 13838389438 | String | 是 | 用户名/邮箱/手机/身份证 |
| ticket | 6Jalnc0uDimvhLA2zKFQO5x8MIO/o5d7bqOe3rSS0/u+65TugF6nJuM/VME5uRh88t3bUY4oItBmRrk= | String | 是 | 验证令牌 |
| code | 504572 | String | 是 | 验证码 |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": {
		"cell": "+86****8943",
		"email": "f****w@126.com",
		"identity_card": "",
		"username": "tbj"
	},
	"msg": "密码重置成功"
}
```

* 失败(200)

```javascript
{
	"code": 20108,
	"msg": "证码验证失败"
}
```

## 获取验证码

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2024-06-25 18:01:14

> 更新时间: 2024-08-07 16:14:15

**已限制访问频率,接口会返回验证令牌在data数据里，同时将通过邮件或手机发送6位数字的验证码给用户。,当开启 SERV_OPEN_CHECK开启验证 时，某些接口需要调用该接口，如注册新用户或使用邮箱或手机登录时。,本接口有流量限制不能过度频繁的调用，否则会提示请求过于频繁请稍后再试。**

**接口状态**

> 已完成

**接口URL**

> /user/captcha/

**请求方式**

> PATCH

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| name | test11 | String | 是 | 用户名/邮箱/手机/身份证 |

**认证方式**

> 继承父级

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": "GKREWQBIZfgujcuaDIC/uWm6FimLM5TV/hE0imPxl/Cc2ziEj/WIUiRqQXgZD7HH1XtHCYx2cnUV+fE=",
	"msg": "cell"
}
```

* 失败(200)

```javascript
{
	"code": 20001,
	"msg": "没有找到指定类型的记录"
}
```

* 异常(429)

```javascript
{
	"code": 10030,
	"msg": "请求过于频繁请稍后再试"
}
```

## 用户日志

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-05-25 11:54:02

> 更新时间: 2024-09-04 12:39:24

**非管理员时只会查询自己时uid的记录，管理员可指定uid或不传uid查全部,条件过滤字段筛选数据样本
"sex":"!=::0","headimg":"!=::","username":"LIKE::fk%","state":1,"intime":"<::2023-05-23 07:00:00","country":"NOT LIKE::%美国%"
生成的sql条件是
AND headimg!= '' ANDsex!= 0 ANDcountryNOT LIKE '%美国%' ANDstate= 1 ANDintime< '2023-05-23 07:00:00' ANDusername LIKE 'fk%'
注意：因为MAP是无序的所以生成的筛选条件顺序也完全是随机的**

```
undefined
```

**接口状态**

> 已完成

**接口URL**

> /user/logs/

**请求方式**

> GET

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| pageSize | 2 | String | 否 | 每页条数 (默认 1) |
| pageNumber | - | String | 否 | 分页页码 (默认 20) |
| pageOrder | - | String | 否 | 排序方向(desc,asc) |
| pageField | - | String | 否 | 排序字段 |
| uid | 2 | String | 否 | 用户ID（查询自己时必填） |
| action | login | String | 否 | 操作标识 |
| note | - | String | 否 | 操作说明(微信登录、修改密码、充值等) |
| actip | - | String | 否 | 操作IP地址 |
| ua | - | String | 否 | 设备信息 |
| intime | - | String | 否 | 记录时间 |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": {
		"total": 11,
		"pageNumber": 1,
		"pageSize": 2,
		"pageOrder": "desc",
		"pageField": "id",
		"list": [
			{
				"id": 19,
				"uid": 168233853578546,
				"action": "login",
				"note": "username",
				"actip": "127.0.0.1",
				"ua": "PostmanRuntime-ApipostRuntime/1.1.0",
				"intime": "2023-05-25 12:10:37"
			},
			{
				"id": 17,
				"uid": 168233853578546,
				"action": "login",
				"note": "username",
				"actip": "127.0.0.1",
				"ua": "PostmanRuntime-ApipostRuntime/1.1.0",
				"intime": "2023-05-25 11:11:46"
			}
		]
	},
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": 0,
	"data": {
		"total": 0,
		"pageNumber": 1,
		"pageSize": 2,
		"pageOrder": "desc",
		"pageField": "id",
		"list": null
	},
	"msg": ""
}
```

## 用户列表

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-05-23 12:26:20

> 更新时间: 2024-09-03 15:24:44

**本接口只有管理员才能访问！,条件过滤字段筛选数据样本
"sex":"!=::0","headimg":"!=::","username":"LIKE::fk%","state":1,"intime":"<::2023-05-23 07:00:00","country":"NOT LIKE::%美国%"
生成的sql条件是
AND headimg!= '' ANDsex!= 0 ANDcountryNOT LIKE '%美国%' ANDstate= 1 ANDintime< '2023-05-23 07:00:00' ANDusername LIKE 'fk%'
注意：因为MAP是无序的所以生成的筛选条件顺序也完全是随机的**

```
undefined
```

**searchOR多字段搜索包括 uid、username、email、cell、nickname、fname、regip**

**接口状态**

> 已完成

**接口URL**

> /user/

**请求方式**

> GET

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| pageSize | 10 | String | 否 | 每页条数 (默认 1) |
| pageNumber | 1 | String | 否 | 分页页码 (默认 20) |
| pageOrder | desc | String | 否 | 排序方向(desc,asc) |
| pageField | uid | String | 否 | 排序字段 |
| username | LIKE::test021% | String | 否 | 帐号 ( AND username LIKE ’fk%‘ ) |
| headimg | !=:: | String | 是 | 头像 ( AND headimg != '' ) |
| state | 1 | String | 是 | 状态 (  AND state = 1 ) |
| intime | <::2023-05-23 07:00 | String | 是 | 注册时间 ( AND intime< '2023-05-23 07:00' ) |
| country | NOT LIKE::%美国% | String | 是 | 国家 ( AND country NOT LIKE '%美国%' ) |
| searchOR | - | String | 是 | 多字段搜索 |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
    "code": 0,
    "data": {
        "total": 211232,
        "pageNumber": 1,
        "pageSize": 20,
        "pageOrder": "desc",
        "pageField": "uid",
        "list": [
            {
                "uid": 1724998239665380,
                "username": "liudehua",
                "ciphers": "5276330af90f4989ef73c0d0cec81d22",
                "email": "",
                "cell": "",
                "nickname": "",
                "headimg": "",
                "sex": 0,
                "birthday": "0001-01-01 00:00:00",
                "company": "",
                "address": "",
                "city": "",
                "province": "",
                "country": "",
                "regip": "192.168.172.88",
                "referer": "",
                "inviter": 0,
                "userclan": "",
                "fname": "",
                "bankcard": "",
                "identity_card": "",
                "grouptag": "",
                "remark": "",
                "object": "",
                "state": 1,
                "intime": "2024-08-30 14:10:39",
                "uptime": "2024-08-30 14:10:39",
                "cid": 0,
                "balance": 0,
                "vip": 0,
                "exptime": "0001-01-01 00:00:00",
                "manage": 0,
                "tag": "",
                "material_remark": "",
                "material_object": "",
                "material_uptime": "2024-08-30 14:10:39"
            },
            {
                "uid": 1724834347400113,
                "username": "taquna",
                "ciphers": "a334d6a5040d3d292f300f1912f0960c",
                "email": "",
                "cell": "",
                "nickname": "",
                "headimg": "",
                "sex": 0,
                "birthday": "0001-01-01 00:00:00",
                "company": "",
                "address": "",
                "city": "",
                "province": "",
                "country": "",
                "regip": "192.168.172.88",
                "referer": "",
                "inviter": 0,
                "userclan": "",
                "fname": "",
                "bankcard": "",
                "identity_card": "",
                "grouptag": "",
                "remark": "",
                "object": "",
                "state": 1,
                "intime": "2024-08-28 16:39:07",
                "uptime": "2024-08-28 16:39:07",
                "cid": 0,
                "balance": 0,
                "vip": 0,
                "exptime": "0001-01-01 00:00:00",
                "manage": 0,
                "tag": "",
                "material_remark": "",
                "material_object": "",
                "material_uptime": "2024-08-28 16:39:07"
            }
        ]
    },
    "msg": ""
}
```

* 失败(200)

```javascript
{
	"code": 0,
	"data": {
		"total": 0,
		"pageNumber": 1,
		"pageSize": 10,
		"pageOrder": "desc",
		"pageField": "uid",
		"list": null
	},
	"msg": ""
}
```

## 获取附属资料

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-05-25 08:54:42

> 更新时间: 2024-09-19 11:04:29

**只能查看自己或管理员才能查看，如果不是管理员会对手机号、邮箱、银行卡、身份证、真实姓名进行脱敏处理**

**接口状态**

> 已完成

**接口URL**

> /user/{uid}/material/

**请求方式**

> GET

**Content-Type**

> none

**路径变量**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| uid | 1988218 | String | 是 | - |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": {
		"uid": 1988218,
		"username": "gofalse",
		"ciphers": "1",
		"email": "123123@163.com",
		"cell": "123456789",
		"nickname": "",
		"headimg": "",
		"sex": 0,
		"birthday": "0001-01-01 00:00:00",
		"company": "",
		"address": "",
		"city": "",
		"province": "",
		"country": "",
		"regip": "218.241.129.18",
		"referer": "",
		"inviter": 0,
		"userclan": "",
		"fname": "",
		"bankcard": "",
		"identity_card": "",
		"grouptag": "",
		"remark": "",
		"object": "",
		"state": 1,
		"intime": "2012-08-08 08:08",
		"uptime": "2023-06-08 13:54",
		"cid": 0,
		"balance": 88.88,
		"vip": 3,
		"exptime": "2024-12-12 12:13:14",
		"manage": 2,
		"tag": "010201020",
		"material_remark": "用户充值",
		"material_object": "{123}",
		"material_uptime": "2024-08-07 16:00"
	},
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": "err debug",
	"data": 16847180152984132,
	"msg": "sql: no rows in result set"
}
```

## 修改附属资料

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2024-09-04 13:59:15

> 更新时间: 2024-09-04 16:10:42

**只有管理员可操作 - 修改用户资料和附属资料同时进行,其中以下字段为附属资料的修改
cid 绑定的联系人ID
balance	 余额（元）
vip 会员等级
exptime 会员到期时间
manage 管理权限(0普通用户 1管理员 2超管)
tag 权限标识
material_remark 附属备注
material_object 附属预留值**

**接口状态**

> 已完成

**接口URL**

> /user/{uid}/material/

**请求方式**

> PUT

**Content-Type**

> form-data

**路径变量**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| uid | - | String | 是 | - |

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| username | - | String | 是 | 帐号 |
| nickname | 哈哈哈 | String | 否 | 昵称 |
| headimg | - | String | 否 | 头像 |
| sex | 1 | Integer | 否 | 性别(0未知 1男 2女) |
| birthday | 1990-01-01 | String | 否 | 出生日期 |
| company | - | String | 否 | 公司 |
| address | - | String | 否 | 地址 |
| city | - | String | 否 | 城市 |
| province | - | String | 否 | 省份 |
| country | - | String | 否 | 国家 |
| referer | - | String | 否 | 用户来源  ( *管理员才能修改 ) |
| inviter | 333 | String | 否 | 邀请者UID  ( *管理员才能修改 ) |
| fname | - | String | 否 | 真实姓名 |
| bankcard | - | String | 否 | 银行卡 |
| grouptag | - | String | 否 | 用户组  ( *管理员才能修改 ) |
| remark | - | String | 否 | 备注 |
| object | - | String | 否 | 预留字段 |
| state | 2 | String | 否 | 状态 ( *管理员才能修改 ) |
| cid | 0 | String | 否 | 绑定的联系人ID |
| balance | 88.88 | String | 否 | 余额（元） |
| vip | 3 | String | 否 | 会员等级 |
| exptime | - | Date | 否 | 会员到期时间 |
| manage | 2 | String | 否 | 管理权限(0普通用户 1管理员 2超管) |
| tag | 010201020 | String | 否 | 权限标识 |
| material_remark | - | String | 否 | 附属备注 |
| material_object | - | String | 否 | 附属预留值 |
| ciphers | - | String | 否 | 密码 |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": 1,
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": 0,
	"data": 0,
	"msg": ""
}
```

# 平台接入

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-05-29 13:29:03

> 更新时间: 2023-05-29 13:29:03

```text
暂无描述
```

**目录Header参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**目录Query参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**目录Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**目录认证信息**

> 继承父级

## 手动接入

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-05-18 15:41:37

> 更新时间: 2024-09-10 13:31:23

**场景情况：,一、openid存在的直接登录成功(返回登录信息uid和token等),二、openid不存在则查询同一平台的unionid如果存在则入库并登录成功(返回登录信息uid和token等),三、同一平台的unionid不存在则存入库并返回开放平台资料(需绑定:返回oid且uid为0的第三方平台信息),四、如已登录则直接入库绑定，如已绑定别的则是登录别的帐号(返回绑定成功或登录信息),五、如果是未登录的情况下再使用“手动绑定”接口进行绑定**

**接口状态**

> 已完成

**接口URL**

> /oauth/

**请求方式**

> POST

**Content-Type**

> urlencoded

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| platfrom | hhhh | String | 是 | 外接平台名(微信 支付宝..) |
| openid | hhhhhhhh | String | 是 | 外接平台身份ID |
| unionid | 222212222 | String | 否 | 外接唯一标识[英文帐号] |
| headimg | - | String | 否 | 头像 |
| nickname | - | String | 否 | 昵称 |
| sex | - | Integer | 否 | 性别(0未知 1男 2女) |
| city | - | String | 否 | 城市 |
| province | - | String | 否 | 省份 |
| country | - | String | 否 | 国家 |
| language | - | String | 否 | 语言（简中zh_CN、EN） |
| privilege | - | String | 否 | 用户特权信息(json数组)[关注数粉丝数等] |
| token | - | String | 否 | 平台接口授权凭证 |
| expires | - | String | 否 | 平台授权失效时间 |
| refresh | - | String | 否 | 平台刷新token |
| scope | - | String | 否 | 用户授权的作用域(使用逗号分隔) |
| subscribe | - | Integer | 否 | 是否关注[或是follow数] |
| subscribetime | - | String | 否 | 关注时间(时间戳) |
| grouptag | - | String | 否 | 分组 |
| tidings | - | String | 否 | 用户动态(json数组)[或最近一条社交信息] |
| remark | eee | String | 否 | 备注 |
| object | - | String | 否 | 预留字段 |

**认证方式**

> 继承父级

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": {
		"address": "",
		"bankcard": "",
		"birthday": "1988-01-01",
		"cell": "135****5555",
		"city": "",
		"company": "",
		"country": "",
		"email": "c****av@aav.com",
		"fname": "",
		"grouptag": "",
		"headimg": "",
		"identity_card": "",
		"intime": "2012-08-05 05:05",
		"inviter": 0,
		"nickname": "",
		"object": "ooo121",
		"password": true,
		"province": "",
		"referer": "",
		"regip": "218.241.129.18",
		"remark": "",
		"sex": 0,
		"state": 1,
		"token": "YV5tO0uwDS0b5HiXxKtwIwG4MUsusZvNOyjmMX1yMS0VdOHRYgmhosuFNBPRKHKMhsHTr+qCceEcWaQ=",
		"uid": 18888288,
		"uptime": "2023-06-08 13:54",
		"userclan": "",
		"username": "ccavtest"
	},
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": 10005,
	"msg": "传入的参数为空或者不合法"
}
```

* 需绑定(200)

```javascript
{
	"code": 0,
	"data": {
		"oid": 17194668571127180,
		"uid": 0,
		"platfrom": "test",
		"openid": "123",
		"headimg": "",
		"unionid": "456",
		"nickname": "",
		"sex": 0,
		"city": "",
		"province": "",
		"country": "",
		"language": "",
		"privilege": "",
		"token": "",
		"expires": "",
		"refresh": "",
		"scope": "",
		"subscribe": 0,
		"subscribetime": "",
		"grouptag": "",
		"tidings": "",
		"remark": "eee",
		"object": "",
		"intime": "",
		"uptime": ""
	},
	"msg": ""
}
```

## 手动绑定

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-07-10 14:45:03

> 更新时间: 2024-08-08 15:33:38

**在已登录的情况下使用应接口对oid进行绑定到此登录用户**

**接口状态**

> 已完成

**接口URL**

> /oauth/

**请求方式**

> PUT

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| oid | 41000 | Number | 是 | 开放平台oid |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": 1,
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": 10002,
	"msg": "未授权访问"
}
```

## 手动解绑

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-07-14 11:53:17

> 更新时间: 2024-09-04 11:13:00

**只能解绑自己的第三方登录**

**接口状态**

> 已完成

**接口URL**

> /oauth/{oid}/

**请求方式**

> DELETE

**Content-Type**

> form-data

**路径变量**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| oid | 41000 | String | 是 | 开放平台oid |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": 1,
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": 10002,
	"msg": "未授权访问"
}
```

## 绑定查询

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2024-09-04 10:40:45

> 更新时间: 2024-09-04 10:44:44

**查询当前用户所有绑定的第三方登录**

**接口状态**

> 开发中

**接口URL**

> /oauth/

**请求方式**

> GET

**Content-Type**

> none

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": [
		{
			"oid": 3214097,
			"uid": 21182,
			"platfrom": "新浪微博",
			"openid": "1032582180",
			"headimg": "",
			"unionid": "",
			"nickname": "haha",
			"sex": 0,
			"city": "",
			"province": "",
			"country": "",
			"language": "",
			"privilege": "",
			"token": "",
			"expires": "",
			"refresh": "",
			"scope": "",
			"subscribe": 0,
			"subscribetime": "",
			"grouptag": "",
			"tidings": "",
			"remark": "",
			"object": "",
			"intime": "2012-10-15 15:20:05",
			"uptime": "2023-06-08 13:54:07"
		},
		{
			"oid": 3435626,
			"uid": 21182,
			"platfrom": "腾讯qq",
			"openid": "F7112E56BA62ED0A44EE1BD8CAAB9F88",
			"headimg": "",
			"unionid": "",
			"nickname": "haha",
			"sex": 0,
			"city": "",
			"province": "",
			"country": "",
			"language": "",
			"privilege": "",
			"token": "",
			"expires": "",
			"refresh": "",
			"scope": "",
			"subscribe": 0,
			"subscribetime": "",
			"grouptag": "",
			"tidings": "",
			"remark": "",
			"object": "",
			"intime": "2012-10-15 15:20:05",
			"uptime": "2023-06-08 13:54:07"
		}
	],
	"msg": ""
}
```

* 失败(404)

```javascript
暂无数据
```

## 接入查询

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-07-07 10:31:47

> 更新时间: 2024-09-04 11:14:13

**只要知道外接平台身份ID或外接唯一标识 + 平台名称 就能查询任何人的第三方登录信息**

**接口状态**

> 已完成

**接口URL**

> /oauth/{openid_unionid}/

**请求方式**

> GET

**Content-Type**

> form-data

**路径变量**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| openid_unionid | 22222222 | String | 是 | 外接平台身份ID或外接唯一标识 |

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| platfrom | 测试平台 | String | 是 | 外接平台名(微信 支付宝..) |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": {
		"id": 16887147793343852,
		"uid": 0,
		"platfrom": "wx",
		"openid": "111111111",
		"headimg": "",
		"unionid": "22222222",
		"nickname": "",
		"sex": 0,
		"city": "",
		"province": "",
		"country": "",
		"language": "",
		"privilege": "",
		"token": "",
		"expires": "",
		"refresh": "",
		"scope": "",
		"subscribe": 0,
		"subscribetime": "",
		"grouptag": "",
		"tidings": "",
		"remark": "eee",
		"object": "",
		"intime": "",
		"uptime": ""
	},
	"msg": ""
}
```

* 失败(404)

```javascript
{
	"code": 20105,
	"msg": "Token错误或过期"
}
```

## 微信授权

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-07-10 09:37:27

> 更新时间: 2023-07-18 12:40:38

```text
暂无描述
```

**接口状态**

> 已完成

**接口URL**

> /oauth/wx/{scope}

**请求方式**

> GET

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| redirect_uri | aaa | String | 是 | 授权后重定向的回调链接地址 |
| state | 1234 | String | 否 | 重定向后会带上state参数，企业可以填写a-zA-Z0-9的参数值，长度不可超过128个字节 |
| appid | - | String | 否 | 公众号的唯一标识 |
| appsecret | - | String | 否 | 公众号的appsecret |
| agentid | 888 | String | 否 | 企业微信 应用agentid |

**认证方式**

> 继承父级

**响应示例**

* 成功(200)

```javascript
暂无数据
```

* 失败(404)

```javascript
暂无数据
```

## 微信接入

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-05-26 09:37:05

> 更新时间: 2024-09-10 17:24:54

**获取第三方授权跳转地址 与  拿TOKEN和用户信息 是母子关系接口，先调用get授权获取code再调用post进行下一步接入操作,get 方法： 用户同意授权获取code
post 方法： 用code换取access_token再拿用户信息,// 微信接入 - 网页授权任意请求类型访问 POST:/oauth/wx/{code}/
// 微信授权 - 网页授权任意请求类型访问 GET:/oauth/wx/{scope}/**

**接口状态**

> 已完成

**接口URL**

> /oauth/wx/{code}/

**请求方式**

> POST

**Content-Type**

> form-data

**路径变量**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| code | - | String | 是 | - |

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| platfrom | hahaha | String | 是 | 外接平台名(微信 支付宝..) |
| appid | - | String | 否 | 公众号的唯一标识 |
| appsecret | - | String | 否 | 公众号的appsecret |
| redirect_uri | - | String | 是 | 回调址址（默认为发起请求的地址） |
| state | - | String | 否 | 用于保持请求和回调的状态 |
| display | - | String | 是 | 用于展示的样式。不传则默认展示为PC下的样式。 如果传入“mobile”，则展示为mobile端下的样式 |

**认证方式**

> 继承父级

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": 16850655499153531,
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": "err debug",
	"data": {
		"object": "",
		"openid": "hhhhhhhh",
		"platfrom": "hhhh",
		"privilege": "",
		"remark": "",
		"tidings": "",
		"uid": "16850647765720025"
	},
	"msg": "Error 1062 (23000): Duplicate entry 'hhhh-hhhhhhhh' for key 'sys_oauth.platfrom_openid'"
}
```

## 微信签名票据

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-05-26 17:01:02

> 更新时间: 2024-07-01 13:44:39

```text
暂无描述
```

**接口状态**

> 已完成

**接口URL**

> /oauth/wx/jsticket/

**请求方式**

> GET

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| url | http://test.fk68.net/ | String | 是 | 当前网页的URL，不包含#及其后面部分 |
| appid | 1000003 | String | 否 | 公众号的唯一标识 |
| appsecret | WYPhNKDk0mBCQmY8BHkQn_-zEV2VlpGXAsmahuOL5Dw | String | 否 | 公众号的appsecret |

**认证方式**

> 继承父级

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": {
		"appId": "wx880af53101eeee3e",
		"timestamp": 1685096490,
		"nonceStr": "FK68TaJiuShiYYSD",
		"signature": "5f0610ae5439675a760210b099523c5655d878ba"
	},
	"msg": ""
}
```

* 失败(404)

```javascript
暂无数据
```

## 微博接入

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2024-09-10 17:18:41

> 更新时间: 2024-09-10 17:25:00

**获取第三方授权跳转地址 与  拿TOKEN和用户信息 是母子关系接口，先调用get授权获取code再调用post进行下一步接入操作,get 方法： 用户同意授权获取code
post 方法： 用code换取access_token再拿用户信息,// 微博接入 - 网页授权任意请求类型访问 POST:/oauth/wb/{code}/
// 微博授权 - 网页授权任意请求类型访问 GET:/oauth/wb/{scope}/**

**接口状态**

> 已完成

**接口URL**

> /oauth/wb/{code}/

**请求方式**

> POST

**Content-Type**

> form-data

**路径变量**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| code | - | String | 是 | - |

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| platfrom | hahaha | String | 是 | 外接平台名(微信 支付宝..) |
| appid | - | String | 否 | 公众号的唯一标识 |
| appsecret | - | String | 否 | 公众号的appsecret |
| redirect_uri | - | String | 是 | 回调址址（默认为发起请求的地址） |
| state | - | String | 否 | 用于保持请求和回调的状态 |
| display | - | String | 是 | 用于展示的样式。不传则默认展示为PC下的样式。 如果传入“mobile”，则展示为mobile端下的样式 |

**认证方式**

> 继承父级

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": 16850655499153531,
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": "err debug",
	"data": {
		"object": "",
		"openid": "hhhhhhhh",
		"platfrom": "hhhh",
		"privilege": "",
		"remark": "",
		"tidings": "",
		"uid": "16850647765720025"
	},
	"msg": "Error 1062 (23000): Duplicate entry 'hhhh-hhhhhhhh' for key 'sys_oauth.platfrom_openid'"
}
```

## QQ接入

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2024-09-10 17:19:02

> 更新时间: 2024-09-10 17:25:04

**获取第三方授权跳转地址 与  拿TOKEN和用户信息 是母子关系接口，先调用get授权获取code再调用post进行下一步接入操作,get 方法： 用户同意授权获取code
post 方法： 用code换取access_token再拿用户信息,// 腾讯QQ接入 - 网页授权任意请求类型访问 POST:/oauth/qq/{code}/
// 腾讯QQ授权 - 网页授权任意请求类型访问 GET:/oauth/qq/{scope}/**

**接口状态**

> 已完成

**接口URL**

> /oauth/qq/{code}/

**请求方式**

> POST

**Content-Type**

> form-data

**路径变量**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| code | - | String | 是 | - |

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| platfrom | hahaha | String | 是 | 外接平台名(微信 支付宝..) |
| appid | - | String | 否 | 公众号的唯一标识 |
| appsecret | - | String | 否 | 公众号的appsecret |
| redirect_uri | - | String | 是 | 回调址址（默认为发起请求的地址） |
| state | - | String | 否 | 用于保持请求和回调的状态 |
| display | - | String | 是 | 用于展示的样式。不传则默认展示为PC下的样式。 如果传入“mobile”，则展示为mobile端下的样式 |

**认证方式**

> 继承父级

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": 16850655499153531,
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": "err debug",
	"data": {
		"object": "",
		"openid": "hhhhhhhh",
		"platfrom": "hhhh",
		"privilege": "",
		"remark": "",
		"tidings": "",
		"uid": "16850647765720025"
	},
	"msg": "Error 1062 (23000): Duplicate entry 'hhhh-hhhhhhhh' for key 'sys_oauth.platfrom_openid'"
}
```

# 通用功能

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-05-30 10:20:52

> 更新时间: 2023-05-30 10:20:52

```text
暂无描述
```

**目录Header参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**目录Query参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**目录Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**目录认证信息**

> 继承父级

## 上传文件

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-05-29 16:29:55

> 更新时间: 2024-08-23 16:30:46

**上传文件后返回是拼接好的全量地址，入库时自动存半截地址**

**接口状态**

> 已完成

**接口URL**

> /common/upload/

**请求方式**

> POST

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| upfile | - | File | 是 | 文件字段名 |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": [
		"/uploads/20230602/168567691443933532.jpeg"
	],
	"msg": ""
}
```

* 失败(404)

```javascript
{
	"code": 10020,
	"msg": "请使用 name=upfile 做为file字段名"
}
```

## 测试通知

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-05-30 11:02:58

> 更新时间: 2024-09-23 13:34:42

**已限制访问频率**

**接口状态**

> 已完成

**接口URL**

> /common/testmsg/

**请求方式**

> GET

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| key | CCAV | String | 是 | API高级密钥 |
| email | test@126.cn | String | 是 | 邮箱 |
| phone | 1333333333 | String | 是 | 手机号 |
| types | 1 | String | 否 | 0其它 1生日 2纪念日 3闹铃 (只针对短信) |

**认证方式**

> 继承父级

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": 0,
	"msg": "Email or Phone Send OK"
}
```

* 失败(404)

```javascript
{
	"code": 10002,
	"msg": "未授权访问"
}
```

* 异常(200)

```javascript
{
	"code": 10000,
	"msg": "未知错误"
}
```

# 记帐理财

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-01 17:20:57

> 更新时间: 2024-09-03 15:32:15

```text
暂无描述
```

**目录Header参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**目录Query参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**目录Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**目录认证信息**

> Bearer Token

> 在Header添加参数 Authorization，其值为在Bearer之后拼接空格和访问令牌

> Authorization: Bearer your_access_token

## 添加帐单

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-02 09:10:06

> 更新时间: 2024-08-20 16:12:20

**item 默认值为 支出 ， accounts 默认值为 现金 ，btime 默认值为 当天,当传 cid 时请同时设置  object 减少一次查询**

**接口状态**

> 已完成

**接口URL**

> /account/

**请求方式**

> POST

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| item | 支出 | String | 是 | 操作项目 |
| class | 基本生活 | String | 是 | 主分类 |
| sort | 日常用品 | String | 是 | 子类别 |
| cid | 0 | String | 否 | 收支对象ID |
| object | 张三 | String | 否 | 收支对象 |
| accounts | 现金 | String | 是 | 操作账户 |
| money | 10.55 | Number | 是 | 金额 |
| note | 买菜 | String | 否 | 备注说明 |
| btime | 2023-06-04 18:18 | String | 是 | 帐单时间 |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": 16856698287568482,
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": 10002,
	"msg": "未授权访问"
}
```

## 帐单详情

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2024-08-20 16:26:05

> 更新时间: 2024-08-20 16:37:54

**根据帐单ID从数据库中获取单条帐单记录**

**接口状态**

> 已完成

**接口URL**

> /account/{aid}/

**请求方式**

> GET

**Content-Type**

> form-data

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": {
		"aid": 17241374483047958,
		"uid": 62345,
		"item": "支出",
		"class": "基本生活",
		"sort": "水果零食",
		"cid": 1024,
		"object": "张三",
		"accounts": "银行卡",
		"money": 22.22,
		"note": "现在正常了不",
		"btime": "2024-08-20T00:00:00+08:00",
		"intime": "2024-08-20T15:04:08+08:00",
		"uptime": "2024-08-20T15:04:08+08:00"
	},
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": 20000,
	"msg": "数据库错误"
}
```

* 异常(200)

```javascript
{
	"code": 10004,
	"msg": "未授权访问"
}
```

## 修改帐单

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-01 18:12:46

> 更新时间: 2024-08-20 16:26:55

```text
暂无描述
```

**接口状态**

> 已完成

**接口URL**

> /account/{aid}/

**请求方式**

> PUT

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| item | 收入 | String | 否 | 操作项目 |
| class | 人情往来 | String | 否 | 主分类 |
| sort | 节日礼金 | String | 否 | 子类别 |
| cid | 0 | String | 否 | 收支对象ID |
| object | 张三 | String | 否 | 收支对象 |
| accounts | 银行卡 | String | 否 | 操作账户 |
| money | 13.228 | Number | 否 | 金额 |
| note | 过生日呀 | String | 否 | 备注说明 |
| btime | 2012-12-12 | String | 否 | 帐单时间 |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": 1,
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": 10002,
	"msg": "未授权访问"
}
```

## 帐单日历

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-02 15:32:17

> 更新时间: 2024-08-08 14:25:58

**帐单日历：非连贯数据，某月没有帐单则不显示该月数据。,inc 收入 、out 支出 、 oth 其它,oth是item项目除了“收入”和“支出”以外的其它所有类型**

**接口状态**

> 已完成

**接口URL**

> /account/calendar/

**请求方式**

> GET

**Content-Type**

> none

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": [
		{
			"year": 2011,
			"month": 1,
			"inc": 0,
			"out": 300,
			"oth": 0
		},
		{
			"year": 2011,
			"month": 4,
			"inc": 0,
			"out": 4625.2,
			"oth": 0
		}
	],
	"msg": ""
}
```

* 失败(404)

```javascript
暂无数据
```

## 月份帐单

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-02 11:59:49

> 更新时间: 2024-08-13 11:59:52

**查询指定月份的下按天分组的帐单记录，主要用于记账首页日历下的历表呈现**

**接口状态**

> 已完成

**接口URL**

> /account/month/

**请求方式**

> GET

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| year | 2023 | Integer | 是 | 年份 |
| month | 6 | Integer | 是 | 月份 |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": [
		{
			"day": "2023-06-05",
			"moment": "今天 2023-06-05 星期一",
			"inc": 222,
			"out": 0,
			"oth": 0,
			"list": [
				{
					"aid": 16859372931904182,
					"uid": 6,
					"item": "收入",
					"class": "人情往来",
					"sort": "节日礼金",
					"cid": 0,
					"object": "张三",
					"accounts": "银行卡",
					"money": 222,
					"note": "过生日呀2",
					"btime": "2023-06-05 11:54",
					"intime": "2023-06-05 11:54",
					"uptime": "2023-06-05 11:54"
				}
			]
		},
		{
			"day": "2023-06-04",
			"moment": "昨天 2023-06-04 星期日",
			"inc": 0,
			"out": 29.54,
			"oth": 0,
			"list": [
				{
					"aid": 16859373895186803,
					"uid": 6,
					"item": "支出",
					"class": "基本生活",
					"sort": "日常用品",
					"cid": 0,
					"object": "张三",
					"accounts": "现金",
					"money": 10.55,
					"note": "买菜",
					"btime": "2023-06-04 18:18",
					"intime": "2023-06-05 11:56",
					"uptime": "2023-06-05 11:56"
				},
				{
					"aid": 16859374245411817,
					"uid": 6,
					"item": "支出",
					"class": "基本生活",
					"sort": "早餐晚餐",
					"cid": 0,
					"object": "",
					"accounts": "现金",
					"money": 18.99,
					"note": "吃粉",
					"btime": "2023-06-04 12:12",
					"intime": "2023-06-05 11:57",
					"uptime": "2023-06-05 11:57"
				}
			]
		}
	],
	"msg": ""
}
```

* 失败(404)

```javascript
暂无数据
```

## 帐单列表

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-02 10:52:30

> 更新时间: 2024-09-03 15:32:49

**条件过滤字段筛选数据样本 "sex":"!=::0","headimg":"!=::","username":"LIKE::fk%","state":1,"intime":"<::2023-05-23 07:00:00","country":"NOT LIKE::%美国%"
生成的sql条件是 
AND headimg!= '' ANDsex!= 0 ANDcountryNOT LIKE '%美国%' ANDstate= 1 ANDintime< '2023-05-23 07:00:00' ANDusername LIKE 'fk%' 
注意：因为MAP是无序的所以生成的筛选条件顺序也完全是随机的**

```
undefined
```

**多字段搜索searchOR包括：note、object、aid**

**接口状态**

> 已完成

**接口URL**

> /account/

**请求方式**

> GET

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| item | 收入 | String | 否 | 操作项目 |
| class | 人情往来 | String | 否 | 主分类 |
| sort | 节日礼金 | String | 否 | 子类别 |
| cid | 0 | String | 否 | 收支对象ID |
| object | 张三 | String | 否 | 收支对象 |
| accounts | 银行卡 | String | 否 | 操作账户 |
| money | 888.9999 | Number | 否 | 金额 |
| note | 过生日呀 | String | 否 | 备注说明 |
| btime | >::2023-01-02 | String | 否 | 帐单时间 |
| pageSize | 20 | String | 否 | 每页条数 (默认 1) |
| pageNumber | 1 | String | 否 | 分页页码 (默认 20) |
| pageOrder | desc | String | 否 | 排序方向(desc,asc) |
| pageField | aid | String | 否 | 排序字段 |
| searchOR | - | String | 是 | 多字段搜索 |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 2,
	"data": {
		"total": 2,
		"pageNumber": 1,
		"pageSize": 20,
		"pageOrder": "desc",
		"pageField": "aid",
		"list": [
			{
				"aid": 16856698287568482,
				"uid": 8,
				"item": "收入",
				"class": "人情往来",
				"sort": "节日礼金",
				"cid": 0,
				"object": "张三",
				"accounts": "银行卡",
				"money": 889,
				"note": "过生日呀",
				"btime": "2020-11-12 18:18",
				"intime": "2023-06-02 09:37",
				"uptime": "2023-06-02 09:37"
			},
			{
				"aid": 73732,
				"uid": 8,
				"item": "支出",
				"class": "",
				"sort": "",
				"cid": 0,
				"object": "",
				"accounts": "现金",
				"money": 13.23,
				"note": "",
				"btime": "2012-12-12 00:00",
				"intime": "2023-06-01 17:48",
				"uptime": "2023-06-02 09:00"
			}
		]
	},
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": 10002,
	"msg": "未授权访问"
}
```

## 帐单类目

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-02 11:21:59

> 更新时间: 2024-07-01 13:45:20

```text
暂无描述
```

**接口状态**

> 已完成

**接口URL**

> /account/type/

**请求方式**

> GET

**Content-Type**

> none

**认证方式**

> 继承父级

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": {
		"中奖收入": "其他收入",
		"书报影音": "文化进修",
		"人际往来": "人情往来",
		"休闲玩乐": "休闲娱乐",
		"公共交通": "交通通讯",
		"公务报销": "职业收入",
		"其他收益": "其他收入",
		"其他杂项": "其他杂项",
		"兼职收入": "职业收入",
		"利息收入": "其他收入",
		"加班收入": "职业收入",
		"化妆美容": "衣服饰品",
		"医药保健": "基本生活",
		"博彩彩票": "休闲娱乐",
		"奖金收入": "职业收入",
		"婚嫁礼金": "人情往来",
		"孝敬长辈": "人情往来",
		"宠物宝贝": "休闲娱乐",
		"家政服务": "其他杂项",
		"工资收入": "职业收入",
		"意外来财": "其他收入",
		"慈善捐助": "人情往来",
		"房屋房产": "其他杂项",
		"房租物业": "基本生活",
		"打车租车": "交通通讯",
		"投资收入": "职业收入",
		"投资理财": "其他杂项",
		"教育培训": "文化进修",
		"数码装备": "文化进修",
		"旅游度假": "休闲娱乐",
		"日常用品": "基本生活",
		"早午晚餐": "基本生活",
		"朋友聚会": "休闲娱乐",
		"服饰装扮": "衣服饰品",
		"柴米油盐": "基本生活",
		"水果零食": "基本生活",
		"水电煤气": "基本生活",
		"生日礼金": "人情往来",
		"电器家居": "其他杂项",
		"礼品礼金": "人情往来",
		"私车供养": "交通通讯",
		"经营收入": "职业收入",
		"节日礼金": "人情往来",
		"补助津贴": "职业收入",
		"话费网费": "交通通讯",
		"车产船产": "其他杂项",
		"运动健身": "休闲娱乐",
		"邮递快递": "交通通讯",
		"金融保险": "其他杂项",
		"鞋帽手套": "衣服饰品",
		"饰品包包": "衣服饰品"
	},
	"msg": ""
}
```

* 失败(404)

```javascript
暂无数据
```

## 收支对象

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-05 17:20:37

> 更新时间: 2024-08-15 15:01:00

```text
暂无描述
```

**接口状态**

> 已完成

**接口URL**

> /account/objects/

**请求方式**

> GET

**Content-Type**

> none

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": [
		{
			"key": 317666,
			"val": "刘德华",
			"str": "liu de hua "
		},
		{
			"key": 218621,
			"val": "张会妹",
			"str": "zhang hui mei "
		},
		{
			"key": 2312340,
			"val": "1空间",
			"str": "1 kong jian "
		},
	],
	"msg": ""
}
```

* 失败(404)

```javascript
暂无数据
```

## 删除帐单

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-05 12:25:45

> 更新时间: 2023-06-05 12:28:09

```text
暂无描述
```

**接口状态**

> 已完成

**接口URL**

> /account/{aid}

**请求方式**

> DELETE

**Content-Type**

> form-data

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": 1,
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": 0,
	"data": 0,
	"msg": ""
}
```

## 收支比例

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-05 13:19:37

> 更新时间: 2024-08-08 14:39:05

**查询指定年月的收支数据**

**接口状态**

> 已完成

**接口URL**

> /account/report/ratio/

**请求方式**

> GET

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| year | 2023 | Integer | 否 | 年份  (不填月分则是全部) |
| month | 6 | Integer | 否 | 月份  (不填月分则是全年) |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": {
		"year": 2023,
		"month": 6,
		"inc": 222,
		"out": 262.09,
		"oth": 0
	},
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": "err debug",
	"data": 6,
	"msg": "sql: no rows in result set"
}
```

## 近年统计

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-05 14:14:10

> 更新时间: 2024-08-08 14:39:24

**主要用于报表的展示**

**接口状态**

> 已完成

**接口URL**

> /account/report/ratios/

**请求方式**

> GET

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| limit | 6 | Integer | 否 | 获取几条 |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": [
		{
			"year": 2023,
			"month": 6,
			"inc": 222,
			"out": 262.09,
			"oth": 0
		},
		{
			"year": 2023,
			"month": 4,
			"inc": 600,
			"out": 0,
			"oth": 0
		},
		{
			"year": 2022,
			"month": 11,
			"inc": 500,
			"out": 0,
			"oth": 0
		},
		{
			"year": 2022,
			"month": 10,
			"inc": 0,
			"out": 2688,
			"oth": 0
		},
		{
			"year": 2022,
			"month": 7,
			"inc": 2800,
			"out": 0,
			"oth": 0
		},
		{
			"year": 2022,
			"month": 6,
			"inc": 0,
			"out": 400,
			"oth": 0
		}
	],
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": "err debug",
	"data": 6,
	"msg": "sql: no rows in result set"
}
```

## 收支明细

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-05 15:10:31

> 更新时间: 2024-08-08 15:06:40

**主要是统计某年或某月或所有时间内帐单的 sort子类别 的收入或支出总数（统计中不包含class 主分类），并以从大到下排列**

**接口状态**

> 已完成

**接口URL**

> /account/report/details/

**请求方式**

> GET

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| year | 2014 | Integer | 否 | 年份  (不填月分则是全部) |
| month | 6 | Integer | 否 | 月份 (不填月分则是全年) |
| item | 收入 | String | 是 | 操作项目(默认为：支出) |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": [
		{
			"sort": "日常用品",
			"total": 21.1
		},
		{
			"sort": "早餐晚餐",
			"total": 18.99
		}
	],
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": 0,
	"data": null,
	"msg": ""
}
```

## 流水账户

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-05 16:09:48

> 更新时间: 2024-08-08 15:28:25

**统计某一时间内或全部的 accounts 操作账户各分类的绝对值流水金额，也就是说不管是支出还是收入都会使用绝对值进行统计。**

**接口状态**

> 已完成

**接口URL**

> /account/report/accounts/

**请求方式**

> GET

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| year | 2022 | Integer | 否 | 年份  (不填月分则是全部) |
| month | 12 | Integer | 否 | 月份 (不填月分则是全年) |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": [
		{
			"sort": "现金",
			"total": 141416.82
		},
		{
			"sort": "银行卡",
			"total": 107674.35
		}
	],
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": 0,
	"data": null,
	"msg": ""
}
```

# 联系人

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-05 18:08:57

> 更新时间: 2023-06-05 18:08:57

```text
暂无描述
```

**目录Header参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**目录Query参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**目录Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**目录认证信息**

> 继承父级

## 添加联系人

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-05 18:09:52

> 更新时间: 2024-09-19 10:39:20

**phone 电话(手机,住宅,单位,自定义)格式 
[TEL::xxx||TEL_CELL_WORK::xxx||TEL_FAX::xxx||TEL_FAX_HOME::xxx],mail 邮箱(个人,单位,家用,自定义)格式 
[EMAIL::xxx||EMAIL_HOME::xxx||EMAIL_WORK::xxx],im 聊天(QQ,WeChat,Skype,Line,自定义)格式
[X-QQ::xxx||X-WECHAT::xxx],http 网址(个人,单位,自定义)格式 
[URL::xxx||URL_WORK::xxx],地址(家庭,单位,其他,自定义)格式
[ADR::xxx||ADR_HOME::xxx],remind 提醒方式(email,phone)格式 
{email::7,1,0||phone::7,1,0||..} 提前7天、1天、0当天,注意：一次性闹铃会在最后一天任务队列进行物理清理一次性闹铃记录**

**接口状态**

> 已完成

**接口URL**

> /contact/

**请求方式**

> POST

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| fullname | 中国人yysd | String | 是 | 姓名 |
| nickname | - | String | 否 | 昵称 |
| picture | - | String | 否 | 相片照片 |
| phone | - | String | 否 | 电话与传真(手机,住宅,单位,自定义) |
| mail | - | String | 否 | 邮箱(个人,单位,家用,自定义) |
| im | - | String | 否 | 聊天(QQ,WeChat,Skype,LINE,自定义) |
| http | - | String | 否 | 网址(个人,单位,自定义) |
| company | - | String | 否 | 公司部门 |
| position | - | String | 否 | 职位头衔 |
| address | - | String | 否 | 地址(家庭,单位,其他,自定义) |
| gender | - | Integer | 否 | 性别(0未知 1男 2女) |
| birthday | - | String | 否 | 生日时间(时间则为农历时辰) |
| lunar | - | Integer | 否 | 是否为农历 |
| grouptag | - | String | 否 | 分组 |
| remind | - | String | 否 | 提醒方式(email,phone) |
| relation | - | String | 否 | 关系 |
| family | - | String | 否 | 家庭户主 |
| note | - | String | 否 | 备注 |
| state | - | String | 否 | 状态(0删除 1正常 2纪念日 3闹铃) |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": 1685959897796961,
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": 10005,
	"msg": "传入的参数为空或者不合法"
}
```

## 联系人列表

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-06 13:47:42

> 更新时间: 2024-09-03 15:27:39

**获取自己的联系人列表,条件过滤字段筛选数据样本
"sex":"!=::0","headimg":"!=::","username":"LIKE::fk%","state":1,"intime":"<::2023-05-23 07:00:00","country":"NOT LIKE::%美国%"
生成的sql条件是
AND headimg!= '' ANDsex!= 0 ANDcountryNOT LIKE '%美国%' ANDstate= 1 ANDintime< '2023-05-23 07:00:00' ANDusername LIKE 'fk%'
注意：因为MAP是无序的所以生成的筛选条件顺序也完全是随机的**

```
undefined
```

**多字段搜索searchOR包括：fullname、phone、pinyin**

**接口状态**

> 已完成

**接口URL**

> /contact/

**请求方式**

> GET

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| fullname | 你好呀 | String | 是 | 姓名 |
| nickname | - | String | 否 | 昵称 |
| picture | - | String | 否 | 相片照片 |
| phone | - | String | 否 | 电话与传真(手机,住宅,单位,自定义) |
| mail | - | String | 否 | 邮箱(个人,单位,家用,自定义) |
| im | - | String | 否 | 聊天(QQ,WeChat,Skype,LINE,自定义) |
| http | - | String | 否 | 网址(个人,单位,自定义) |
| company | - | String | 否 | 公司部门 |
| position | - | String | 否 | 职位头衔 |
| address | - | String | 否 | 地址(家庭,单位,其他,自定义) |
| gender | - | String | 否 | 性别 |
| birthday | - | String | 否 | 生日时间 |
| lunar | - | Integer | 否 | 是否为农历 |
| grouptag | 家人 | String | 否 | 分组 |
| remind | - | String | 否 | 提醒方式(email,phone) |
| relation | - | String | 否 | 关系 |
| family | - | String | 否 | 家庭户主 |
| note | - | String | 否 | 备注 |
| state | - | String | 否 | 状态(0删除 1正常 2纪念日 3闹铃) |
| pageSize | 10 | String | 是 | 每页条数 (默认 1) |
| pageNumber | 1 | String | 是 | 分页页码 (默认 20) |
| pageOrder | DESC | String | 是 | 排序方向(desc,asc) |
| pageField | cid | String | 是 | 排序字段 |
| searchOR | 多字段搜索 | String | 是 | 多字段搜索 |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 1,
	"data": {
		"total": 1,
		"pageNumber": 1,
		"pageSize": 20,
		"pageOrder": "desc",
		"pageField": "uid",
		"list": [
			{
				"cid": 16860288931810703,
				"uid": 6,
				"fullname": "你好呀",
				"pinyin": "ni hao ya",
				"nickname": "我来了",
				"picture": "",
				"phone": "",
				"mail": "",
				"im": "",
				"http": "",
				"company": "",
				"position": "",
				"address": "",
				"gender": 1,
				"birthday": "0001-01-01 00:00",
				"lunar": 0,
				"grouptag": "",
				"remind": "",
				"relation": "",
				"family": "",
				"note": "",
				"state": 1,
				"intime": "2023-06-06 13:21",
				"uptime": "2023-06-06 13:57"
			}
		]
	},
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": 10005,
	"msg": "传入的参数为空或者不合法"
}
```

## 获取联系人

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2024-07-01 14:46:10

> 更新时间: 2024-08-08 10:17:19

**获取联系人功能必须是自己创建的才可以，管理员也无没权获取他人联系人**

**接口状态**

> 已完成

**接口URL**

> /contact/{cid}/

**请求方式**

> GET

**Content-Type**

> form-data

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": {
		"cid": 88271212,
		"uid": 121238,
		"fullname": "北京团队",
		"pinyin": "bei jing tuan dui ",
		"nickname": "GOGOGO",
		"picture": "./uploads/avatar/202202/20220218_170207_8.jpg",
		"phone": "TEL_FAX::1-301-100-5788||TEL::1-380-380-3838||",
		"mail": "EMAIL::ccav@tv.com||EMAIL_WORK::ccav@go.com||",
		"im": "X-QQ::12345678||X-SKYPE::xxxxxxxxx||",
		"http": "URL::http://www.fk68.net||",
		"company": "数字天才公司",
		"position": "高级工程师",
		"address": "ADR::100024 北京市 朝阳区 北京市朝阳区住总大楼9909||",
		"gender": 1,
		"birthday": "1984-08-08T05:30:00+08:00",
		"lunar": 1,
		"grouptag": "无",
		"remind": "",
		"relation": "我自己",
		"family": "刘德化",
		"note": "记录一下",
		"state": 1,
		"intime": "2015-10-15T15:21:15+08:00",
		"uptime": "2025-02-18T17:12:19+08:00"
	},
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": 0,
	"data": 0,
	"msg": ""
}
```

## 修改联系人

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-06 13:05:11

> 更新时间: 2024-08-26 10:54:18

**只能修改自己添加的联系人，因为更新条件会自带带当前用户ID**

**接口状态**

> 已完成

**接口URL**

> /contact/{cid}/

**请求方式**

> PUT

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| fullname | 你好呀 | String | 是 | 姓名 |
| nickname | 我来了 | String | 否 | 昵称 |
| picture | - | String | 否 | 相片照片 |
| phone | - | String | 否 | 电话与传真(手机,住宅,单位,自定义) |
| mail | - | String | 否 | 邮箱(个人,单位,家用,自定义) |
| im | - | String | 否 | 聊天(QQ,WeChat,Skype,LINE,自定义) |
| http | - | String | 否 | 网址(个人,单位,自定义) |
| company | - | String | 否 | 公司部门 |
| position | - | String | 否 | 职位头衔 |
| address | - | String | 否 | 地址(家庭,单位,其他,自定义) |
| gender | - | Integer | 否 | 性别(0未知 1男 2女) |
| birthday | - | String | 否 | 生日时间(时间则为农历时辰) |
| lunar | - | Integer | 否 | 是否为农历 |
| grouptag | - | String | 否 | 分组 |
| remind | - | String | 否 | 提醒方式(email,phone) |
| relation | - | String | 否 | 关系 |
| family | - | String | 否 | 家庭户主 |
| note | - | String | 否 | 备注 |
| state | - | String | 否 | 状态(0删除 1正常 2纪念日 3闹铃) |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": 1,
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": 0,
	"data": 0,
	"msg": ""
}
```

## 删除联系人

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-06 13:33:24

> 更新时间: 2024-08-08 10:17:31

**非物理删除，只能删除自己添加的联系人（因条件会带当前用户ID）**

**接口状态**

> 已完成

**接口URL**

> /contact/{cid}/

**请求方式**

> DELETE

**Content-Type**

> none

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": 16860288931810706,
	"msg": ""
}
```

* 失败(404)

```javascript
暂无数据
```

## 联系人分组

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-06 14:26:13

> 更新时间: 2024-08-07 13:50:20

**获取自己的联系人分组**

**接口状态**

> 已完成

**接口URL**

> /contact/groups/

**请求方式**

> GET

**Content-Type**

> none

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": [
		{
			"key": 0,
			"val": "基友",
			"str": "基友"
		},
		{
			"key": 0,
			"val": "艺人",
			"str": "艺人"
		}
	],
	"msg": ""
}
```

* 失败(404)

```javascript
暂无数据
```

## 导入VCF

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-25 15:31:23

> 更新时间: 2024-07-01 13:46:50

**导出 vCard,在 iCloud.com 上的“通讯录”中，点按以选择联系人列表中的所需联系人。如果你要导出多个联系人，请按住 Command 键（Mac 电脑上）或 Ctrl 键（Windows 电脑上），然后点按你要导出的每个联系人。
点按边栏中的 “显示操作菜单”弹出式按钮，然后选取“导出 vCard”。,如果你选择多个联系人，“通讯录”会导出一个包含所有联系人的 vCard。**

**接口状态**

> 已完成

**接口URL**

> /contact/vcards/

**请求方式**

> POST

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| upfile | - | File | 是 | - |

**认证方式**

> 继承父级

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": {
		"success": 192,
		"total": 203
	}
}
```

* 失败(404)

```javascript
暂无数据
```

## 关联联系人

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2024-07-01 13:43:13

> 更新时间: 2024-09-11 14:47:27

**一键关联用户资料到料到联系人，没有联系人时会自动创建联系人进行连接。,通常用于用户注册后进行初始化操作。如果已绑定则是更新联系人数据。**

**接口状态**

> 已完成

**接口URL**

> /contact/

**请求方式**

> PATCH

**Content-Type**

> none

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": 17230868771332108,
	"msg": ""
}
```

* 失败(404)

```javascript
{
	"code": 20000,
	"msg": "数据库错误"
}
```

# 生日提醒

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-07 11:55:15

> 更新时间: 2023-06-07 11:55:15

```text
暂无描述
```

**目录Header参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**目录Query参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**目录Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**目录认证信息**

> 继承父级

## 获取列表

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-05-19 16:21:46

> 更新时间: 2024-08-27 16:09:02

**无分页，,json:"remind_num" description:"剩余天数"
json:"remind_date" description:"生日的日期 YYYY-MM-DD"
json:"remind_cndate" description:"生日的日期 X月X日" // 跨年后会是 X月X日 XXXX年
json:"remind_year" description:"年龄或周年"
json:"remind_star" description:"星座"
json:"remind_zodiac" description:"生肖"
json:"birthday_cndate" description:"出生中文日期"
json:"remind_day" description:"今天明天后天N天后"
json:"remind_week" description:"星期几"**

```
undefined
```

**接口状态**

> 已完成

**接口URL**

> /remind/

**请求方式**

> GET

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| state | 2 | String | 否 | 状态类别(1正常 2纪念日 3闹铃) |
| lunar | 0 | String | 否 | 是否为农历 |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": [
		{
			"cid": 9995161,
			"uid": 8888,
			"fullname": "驾驶证积分清零",
			"pinyin": "jia shi zheng ji fen qing ling ",
			"nickname": "",
			"picture": "",
			"gender": 0,
			"birthday": "2004-02-24 00:00",
			"lunar": 0,
			"grouptag": "",
			"remind": "email::7,1,0||phone::7,1,0",
			"relation": "",
			"note": "",
			"state": 2,
			"remind_num": 256,
			"remind_date": "2024-02-24",
			"remind_year": 20,
			"remind_star": "双鱼",
			"remind_zodiac": "猴",
			"remind_cndate": "二〇〇四年二月初五",
			"remind_day": "256天后",
			"remind_week": "星期六"
		}
	],
	"msg": ""
}
```

* 失败(404)

```javascript
暂无数据
```

## 任务队列

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-08 14:50:36

> 更新时间: 2024-08-27 13:53:44

**已限制访问频率，每天在资源占用最小的时候生成队列，每天只执行一次。,成功后可直接通过api地址访问 /queue.json 文件**

**接口状态**

> 已完成

**接口URL**

> /remind/task/

**请求方式**

> GET

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| key | CCAV | String | 是 | API高级密钥 |

**认证方式**

> 继承父级

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": 168,
	"msg": "已成功生成队列任务"
}
```

* 失败(200)

```javascript
{
	"code": 10002,
	"msg": "未授权访问"
}
```

## 发送提醒

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-12 18:16:13

> 更新时间: 2024-08-30 11:48:41

**已限制访问频率，可以间隔十分钟执行一次，每天最多三至四次就差不多了。如果提醒队列中某一类型总条数大于0时执行发送任务。失败后会自动添加到新的队列中方便下次重试。,每次执行会返回操作的条数以及控制台会显示：,Email 队列总条数： 100,Phone 队列总条数： 88,Notice 队列总条数： 66,Message 队列总条数： 133,从队列中取出短信和邮件并发送**

**接口状态**

> 已完成

**接口URL**

> /remind/queue/

**请求方式**

> GET

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| key | CCAV | String | 是 | API高级密钥 |

**认证方式**

> 继承父级

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": 13,
	"msg": "已推送队列消息"
}
```

* 失败(200)

```javascript
{
	"code": 10002,
	"msg": "未授权访问"
}
```

# 人情往来

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-14 10:13:07

> 更新时间: 2023-06-14 10:13:07

```text
暂无描述
```

**目录Header参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**目录Query参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**目录Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**目录认证信息**

> 继承父级

## 人情列表

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-14 10:13:27

> 更新时间: 2024-09-03 15:09:02

**小提示：这里的人情账簿数据是从记账理财帐单中自动取得的，只要帐单中某笔记录有支出对象就会纳入这个人情列表里。 如果没有数据，你需要[记帐理材]下的[添加账单]在新增时添加一个支出对象（收支对象ID）即可。该接口不支持分页，且只会对item为"收入"或"支出"的进行统计，“其他”类形不参与统计。不支持分页操作,返回字段描述**


 - "inc": 对方回报 
 - "out": 我的付出
 - "inc_ratio": 收入占比
 - "out_ratio": 支出占比
 - "coef": 人情系数（可用于排序）
 - "last": 最后往来时间戳（默认排序）
 - "total": 往来总次数（可用于排序）
 - "inc_list": 回报帐单列表（list默认4条）
 - "out_list": 付出帐单列表（list默认4条）

**接口状态**

> 已完成

**接口URL**

> /humane/

**请求方式**

> GET

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| pageOrder | ASC | String | 否 | 排序方向(默认 desc 倒序、asc 正序) |
| pageField | inc | String | 否 | 排序字段(默认 last 往来时间、out 我的付出、inc 对方回报、total 往来次数 、coef 人情系数) |
| listSize | 2 | Number | 否 | 回报/付出帐单列表长度 |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": [
		{
			"cid": 1236028,
			"fullname": "刘德华",
			"pinyin": "liu de hua ",
			"picture": "",
			"inc": 22,
			"out": 22,
			"inc_ratio": 50,
			"out_ratio": 50,
			"coef": 9508.96,
			"last": 1724169600,
			"total": 1,
			"inc_total": 2,
			"out_total": 205,
			"inc_list": [
				{
					"item": "收入",
					"class": "其他收入",
					"sort": "中奖收入",
					"object": "刘德华",
					"accounts": "现金",
					"money": 22,
					"note": "中了一个亿",
					"btime": "2024-08-21",
					"uptime": "2024-08-21 09:32:33"
				}
			],
			"out_list": [
				{
					"item": "支出",
					"class": "人情往来",
					"sort": "人际往来",
					"object": "刘德华",
					"accounts": "现金",
					"money": 22,
					"note": "买花花",
					"btime": "2018-08-13",
					"uptime": "2023-07-10 16:17:59"
				}
			]
		}
	]
}
```

| 参数名 | 示例值 | 参数类型 | 参数描述 |
| --- | --- | ---- | ---- |
| code | 0 | Integer | 换取access_token的票据() |
| data | - | Array | 数据集 |
| data.cid | 1236028 | Integer | 收支对象ID |
| data.fullname | 刘德华 | String | 姓名 |
| data.pinyin | liu de hua | String | 拼音 |
| data.picture | - | String | 相片照片 |
| data.inc | 11 | Integer | 对方回报 |
| data.out | 22 | Integer | 我的付出 |
| data.inc_ratio | 100 | Integer | 收入占比 |
| data.out_ratio | 200 | Integer | 支出占比 |
| data.coef | 9508.96 | Number | 人情系数 |
| data.last | 1724169600 | Integer | 最后往来时间戳 |
| data.total | 1 | Integer | 往来总次数 |
| data.inc_total | 2 | Integer | 回报总次数 |
| data.out_total | 205 | Integer | 付出总次数 |
| data.inc_list | - | Array | 回报帐单列表 |
| data.inc_list.item | 收入 | String | 操作项目 |
| data.inc_list.class | 其他收入 | String | 主分类 |
| data.inc_list.sort | 中奖收入 | String | 子类别 |
| data.inc_list.object | 刘德华 | String | 预留字段 |
| data.inc_list.accounts | 现金 | String | 操作账户 |
| data.inc_list.money | 11 | Integer | 金额 |
| data.inc_list.note | 中了一个亿 | String | 备注说明 |
| data.inc_list.btime | 2024-08-21 | String | 帐单时间 |
| data.inc_list.uptime | 2024-08-21 09:32:33 | String | 更新时间 |
| data.out_list | - | Array | 付出帐单列表 |
| data.out_list.item | 支出 | String | 操作项目 |
| data.out_list.class | 人情往来 | String | 主分类 |
| data.out_list.sort | 人际往来 | String | 子类别 |
| data.out_list.object | 刘德华 | String | 预留字段 |
| data.out_list.accounts | 现金 | String | 操作账户 |
| data.out_list.money | 22 | Integer | 金额 |
| data.out_list.note | 买花花 | String | 备注说明 |
| data.out_list.btime | 2018-08-13 | String | 帐单时间 |
| data.out_list.uptime | 2023-07-10 16:17:59 | String | 更新时间 |

* 失败(404)

```javascript
暂无数据
```

## 人情详情

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-14 16:25:44

> 更新时间: 2024-08-22 14:37:15

**返回字段描述**


 - "inc": 对方回报 
 - "out": 我的付出
 - "inc_ratio": 收入占比
 - "out_ratio": 支出占比(通常使用inc_ratio就能确定inc_ratio)
 - "coef": 人情系数
 - "last": 最后往来时间戳
 - "total": 往来总次数
 - "inc_list": 回报帐单列表（全部）
 - "out_list": 付出帐单列表（全部）

**接口状态**

> 已完成

**接口URL**

> /humane/{cid}/

**请求方式**

> GET

**Content-Type**

> none

**路径变量**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| cid | - | String | 是 | - |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": {
		"cid": 1236010,
		"fullname": "张绘妹",
		"pinyin": "zhang hui mei ",
		"picture": "./uploads/avatar/201306/011850_6.jpg",
		"inc": 11,
		"out": 22,
		"inc_ratio": 30,
		"out_ratio": 70,
		"coef": 9508.96,
		"last": 1724169600,
		"total": 1,
		"inc_total": 2,
		"out_total": 205,
		"inc_list": [
			{
				"item": "收入",
				"class": "其他收入",
				"sort": "中奖收入",
				"object": "张绘妹",
				"accounts": "微信支付",
				"money": 11,
				"note": "中了一个亿",
				"btime": "2024-08-21",
				"uptime": "2024-08-21 09:32:33"
			}
		],
		"out_list": [
			{
				"item": "支出",
				"class": "人情往来",
				"sort": "人际往来",
				"object": "张绘妹",
				"accounts": "现金",
				"money": 22,
				"note": "买花花",
				"btime": "2018-08-13",
				"uptime": "2023-07-10 16:17:59"
			}
		]
	},
	"msg": ""
}
```

| 参数名 | 示例值 | 参数类型 | 参数描述 |
| --- | --- | ---- | ---- |
| code | 0 | Number | 换取access_token的票据() |
| data | - | Object | 返回数据 |
| data.cid | 1236010 | Number | 收支对象ID |
| data.fullname | 张绘妹 | String | 用户名/邮箱/手机/身份证 |
| data.pinyin | zhang hui mei | String | - |
| data.picture | ./uploads/avatar/201306/011850_6.jpg | String | 相片照片 |
| data.inc | 11 | Number | - |
| data.out | 22 | Number | - |
| data.inc_ratio | 100 | Number | - |
| data.out_ratio | 200 | Number | - |
| data.coef | 9508.96 | Number | - |
| data.last | 1724169600 | Number | - |
| data.total | 1 | Number | - |
| data.inc_total | 2 | Number | - |
| data.out_total | 205 | Number | - |
| data.inc_list | - | Array | - |
| data.inc_list.item | 收入 | String | 操作项目 |
| data.inc_list.class | 其他收入 | String | - |
| data.inc_list.sort | 中奖收入 | String | 子类别 |
| data.inc_list.object | 张绘妹 | String | 预留字段 |
| data.inc_list.accounts | 微信支付 | String | 操作账户 |
| data.inc_list.money | 11 | Number | 金额 |
| data.inc_list.note | 中了一个亿 | String | 操作说明(微信登录、修改密码、充值等) |
| data.inc_list.btime | 2024-08-21 | String | 帐单时间 |
| data.inc_list.uptime | 2024-08-21 09:32:33 | String | 更新时间 |
| data.out_list | - | Array | - |
| data.out_list.item | 支出 | String | 操作项目 |
| data.out_list.class | 人情往来 | String | - |
| data.out_list.sort | 人际往来 | String | 子类别 |
| data.out_list.object | 张绘妹 | String | 预留字段 |
| data.out_list.accounts | 现金 | String | 操作账户 |
| data.out_list.money | 22 | Number | 金额 |
| data.out_list.note | 买花花 | String | 操作说明(微信登录、修改密码、充值等) |
| data.out_list.btime | 2018-08-13 | String | 帐单时间 |
| data.out_list.uptime | 2023-07-10 16:17:59 | String | 更新时间 |
| msg | - | String | 返回文字描述 |

* 失败(404)

```javascript
暂无数据
```

## 统计比例

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-20 12:35:14

> 更新时间: 2024-07-01 13:47:33

```text
暂无描述
```

**接口状态**

> 已完成

**接口URL**

> /humane/report/ratio/

**请求方式**

> GET

**Content-Type**

> none

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": {
		"inc": 6814540.34,
		"out": 5564371.85
	},
	"msg": ""
}
```

* 失败(404)

```javascript
暂无数据
```

## 回报排名

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-20 13:09:01

> 更新时间: 2024-07-01 13:47:39

```text
暂无描述
```

**接口状态**

> 已完成

**接口URL**

> /humane/report/top/

**请求方式**

> GET

**Content-Type**

> none

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": [
		{
			"sort": "生日礼金",
			"money": 238492.82,
			"class": "礼金收入"
		}
	],
	"msg": ""
}
```

* 失败(404)

```javascript
暂无数据
```

# 云纸张

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-15 17:16:08

> 更新时间: 2023-06-15 17:16:08

```text
暂无描述
```

**目录Header参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**目录Query参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**目录Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| 暂无参数 |

**目录认证信息**

> 继承父级

## 纸张列表

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-16 11:47:28

> 更新时间: 2024-07-01 13:47:48

**获取当前用户所有纸张**

**接口状态**

> 已完成

**接口URL**

> /notepad/

**请求方式**

> GET

**Content-Type**

> none

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": [
		{
			"nid": 2131787,
			"uid": 222666,
			"url": "ccav",
			"share": "8JDcfaeFjf",
			"content": "",
			"pwd": "ccav",
			"caret": 0,
			"scroll": 0,
			"ip": "124.207.84.210",
			"referer": "csgo",
			"state": 1,
			"intime": "2012-06-19 16:44",
			"uptime": "2023-06-16 09:49"
		}
	],
	"msg": ""
}
```

* 失败(404)

```javascript
暂无数据
```

## 添加纸张

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-15 17:19:50

> 更新时间: 2024-07-01 13:47:53

**给当前添加纸张，会判查当前用户状态和当前用户纸张数据是否超过系统配置**

**接口状态**

> 已完成

**接口URL**

> /notepad/

**请求方式**

> POST

**Content-Type**

> form-data

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| url | ccavgo8 | String | 是 | 纸张访问名（云地址） |
| content | - | String | 否 | 内容 |
| pwd | hhhhh | String | 否 | 密码 |
| caret | - | String | 否 | 光标位置 |
| scroll | - | String | 否 | 滚动位置 |
| referer | QQ空是 | String | 否 | 纸张来源 |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": 16868208526358456,
	"msg": ""
}
```

* 失败(404)

```javascript
{
	"code": 20500,
	"msg": "您不是VIP或已过期"
}
```

* 异常(200)

```javascript
{
	"code": 20020,
	"msg": "资源已存在"
}
```

## 获取纸张

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-16 12:06:07

> 更新时间: 2024-09-12 13:06:16

**限制访问频率(2分钟内可访问3次每次间隔30秒),获取纸张全部数据不管当前纸张状态，且只有纸张所有人身份为VIP时才能获取到只读的分享地址share字段**

**接口状态**

> 已完成

**接口URL**

> /notepad/{url}/

**请求方式**

> GET

**Content-Type**

> form-data

**路径变量**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| url | - | String | 是 | - |

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| pwd | 443434332777 | String | 否 | 密码 |

**认证方式**

> 继承父级

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": {
		"nid": 8644,
		"uid": 0,
		"url": "ccav",
		"share": "iIFabJDcHg",
		"content": "hahahah",
		"pwd": "",
		"caret": 0,
		"scroll": 0,
		"ip": "124.126.135.10",
		"referer": "",
		"state": 1,
		"intime": "2013-09-11 17:40",
		"uptime": "2014-01-20 22:13"
	},
	"msg": ""
}
```

* 失败(404)

```javascript
{
	"code": 10004,
	"msg": "未授权访问"
}
```

* 异常(200)

```javascript
{
	"code": 20022,
	"msg": "闲置空资源待认领"
}
```

## 分享纸张

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-16 13:48:57

> 更新时间: 2024-09-12 18:11:56

**只获取纸张所主人的用户状态正常的内容**

**接口状态**

> 已完成

**接口URL**

> /notepad/share/{uuid}/

**请求方式**

> GET

**Content-Type**

> none

**路径变量**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| uuid | - | String | 是 | - |

**认证方式**

> 继承父级

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": "txt",
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": 10010,
	"msg": "请求的资源不存在"
}
```

## 更新纸张

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-16 10:31:31

> 更新时间: 2024-09-12 16:50:51

**更新自己或已知密码的他人的某一张纸，可以更新状态为锁定或修改url但不能修改共享地址。
当没有登录时要传key纸张的访问密码来进行操作
如果pwd为空字符串则重置密码为空下次他人访问无需密码**

**接口状态**

> 已完成

**接口URL**

> /notepad/{nid}/

**请求方式**

> PUT

**Content-Type**

> form-data

**路径变量**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| nid | - | String | 是 | - |

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| url | ccavgo8 | String | 是 | 纸张访问名（云地址） |
| content | 122212 | String | 否 | 内容 |
| pwd | - | String | 否 | 修改密码 |
| caret | - | String | 否 | 光标位置 |
| scroll | - | String | 否 | 滚动位置 |
| referer | - | String | 否 | 纸张来源 |
| key | - | String | 否 | 操作密码 |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": 1,
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": 20050,
	"msg": "帐号无权限操作"
}
```

## 删除纸张

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-06-16 10:46:14

> 更新时间: 2024-08-26 14:31:44

**物理删除自己的某一张纸,"data": 1时成功，0则失败**

**接口状态**

> 已完成

**接口URL**

> /notepad/{nid}/

**请求方式**

> DELETE

**Content-Type**

> form-data

**路径变量**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| nid | - | String | 是 | - |

**认证方式**

> 私密键值对

> 在Header添加参数

> key: 2jhskfdgjldfgjldf-9639-kiuwoiruk

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": 1,
	"msg": ""
}
```

* 失败(200)

```javascript
{
	"code": 0,
	"data": 0,
	"msg": ""
}
```

## 领取绑定

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2024-06-28 18:17:00

> 更新时间: 2024-08-26 14:31:52

**如果这页纸没有主人，则可以去认领绑定到自己名下，目的是为了兼容老用户。**

**接口状态**

> 开发中

**接口URL**

> /notepad/{url}/

**请求方式**

> PATCH

**Content-Type**

> form-data

**路径变量**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| url | - | String | 是 | - |

**请求Body参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| pwd | 123452 | String | 否 | 密码 |

**认证方式**

> 继承父级

**响应示例**

* 成功(200)

```javascript
{
	"code": 0,
	"data": 1,
	"msg": ""
}
```

* 失败(404)

```javascript
{
	"code": 20050,
	"msg": "帐号无权限操作"
}
```

# 文档

> 创建人: feikeq

> 更新人: feikeq

> 创建时间: 2023-05-29 16:28:11

> 更新时间: 2023-05-30 10:20:56

```text
暂无描述
```

**接口状态**

> 开发中

**接口URL**

> /%E6%8F%90%E7%AC%94%E8%AE%B0API%E6%8E%A5%E5%8F%A3%E6%96%87%E6%A1%A3%E6%8E%A5%E5%8F%A3%E6%96%87%E6%A1%A3.html?target_id=001

**请求方式**

> GET

**Content-Type**

> none

**请求Query参数**

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述 |
| --- | --- | ---- | ---- | ---- |
| target_id | 001 | String | 是 | - |

**认证方式**

> 继承父级

**响应示例**

* 成功(200)

```javascript
暂无数据
```

* 失败(404)

```javascript
暂无数据
```
