# [API 接口文档]()

## 基础信息

| 项目     | 内容                   |
| -------- | ---------------------- |
| 基础URL  | `http://localhost/api` |
| 数据格式 | JSON                   |
| 字符编码 | UTF-8                  |

## 认证方式

除注册和登录外，所有接口需要在请求头中携带 JWT Token：

Token 通过登录或注册接口获取，有效期 24 小时。

## 状态码说明

| 状态码 | 含义             |
| ------ | ---------------- |
| 200    | 成功             |
| 400    | 参数错误         |
| 401    | 未授权/Token无效 |
| 404    | 资源不存在       |
| 409    | 资源冲突         |
| 500    | 服务器错误       |

---

## 1. 用户注册

| 项目 | 内容            |
| ---- | --------------- |
| 路径 | `/api/register` |
| 方法 | POST            |
| 鉴权 | 无              |

**请求参数**

| 参数     | 类型   | 必填 | 说明             |
| -------- | ------ | ---- | ---------------- |
| username | string | 是   | 用户名，3-50字符 |
| password | string | 是   | 密码，至少6字符  |

**请求示例**
```json
{
    "username": "zhangsan",
    "password": "123456"
}
```

**成功响应（200）**

```
{
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
        "id": 1,
        "username": "zhangsan",
        "created_at": "2026-05-16T03:00:00Z"
    }
}
```

**错误响应**

| 状态码 | 错误信息     | 说明                   |
| :----- | :----------- | :--------------------- |
| 400    | 参数错误     | 用户名或密码格式不正确 |
| 400    | 用户名已存在 | 用户名已被注册         |

## 2. 用户登录

| 项目 | 内容         |
| :--- | :----------- |
| 路径 | `/api/login` |
| 方法 | POST         |
| 鉴权 | 无           |

**请求参数**

| 参数     | 类型   | 必填 | 说明   |
| :------- | :----- | :--- | :----- |
| username | string | 是   | 用户名 |
| password | string | 是   | 密码   |

**请求示例**

```
{
    "username": "zhangsan",
    "password": "123456"
}
```

**成功响应（200）**

```
{
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
        "id": 1,
        "username": "zhangsan",
        "created_at": "2026-05-16T03:00:00Z"
    }
}
```

**错误响应**

| 状态码 | 错误信息   | 说明         |
| :----- | :--------- | :----------- |
| 401    | 用户不存在 | 用户名未注册 |
| 401    | 密码错误   | 密码不正确   |

------

## 3. 查询单词

| 项目 | 内容              |
| :--- | :---------------- |
| 路径 | `/api/word/query` |
| 方法 | GET               |
| 鉴权 | 需要 Bearer Token |

**请求参数**

| 参数        | 类型   | 必填 | 说明             |
| :---------- | :----- | :--- | :--------------- |
| word        | string | 是   | 要查询的单词     |
| ai_provider | string | 是   | AI提供商：`qwen` |

**请求示例**

```
GET /api/word/query?word=hello&ai_provider=qwen
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**成功响应（200）- 从 AI 查询**

```
{
    "from_cache": false,
    "word": "hello",
    "definition": "你好；喂",
    "sentences": [
        "Hello, how are you today?",
        "She said hello to everyone.",
        "Hello is often the first word people learn."
    ],
    "ai_provider": "qwen"
}
```

**成功响应（200）- 从缓存返回**

```
{
    "from_cache": true,
    "word": "hello",
    "definition": "你好；喂",
    "sentences": [
        "Hello, how are you today?",
        "She said hello to everyone.",
        "Hello is often the first word people learn."
    ],
    "ai_provider": "qwen"
}
```

**错误响应**

| 状态码 | 错误信息    | 说明                     |
| :----- | :---------- | :----------------------- |
| 401    | 未提供token | 请求头缺少 Authorization |
| 401    | token无效   | Token过期或签名错误      |
| 500    | AI服务错误  | 调用 AI API 失败         |

------

## 4. 保存单词

| 项目 | 内容              |
| :--- | :---------------- |
| 路径 | `/api/word/save`  |
| 方法 | POST              |
| 鉴权 | 需要 Bearer Token |

**请求参数**

| 参数        | 类型     | 必填 | 说明           |
| :---------- | :------- | :--- | :------------- |
| word        | string   | 是   | 单词           |
| definition  | string   | 是   | 中文释义       |
| sentences   | []string | 是   | 3条例句        |
| ai_provider | string   | 是   | AI来源（qwen） |

**请求示例**

```
{
    "word": "hello",
    "definition": "你好；喂",
    "sentences": [
        "Hello, how are you?",
        "Say hello to your mother.",
        "Hello world!"
    ],
    "ai_provider": "qwen"
}
```

**成功响应（200）**

```
{
    "message": "保存成功",
    "id": 5
}
```

**错误响应**

| 状态码 | 错误信息    | 说明               |
| :----- | :---------- | :----------------- |
| 401    | 未提供token | 未认证             |
| 409    | 单词已保存  | 该单词已在单词本中 |

------

## 5. 获取单词列表

| 项目 | 内容              |
| :--- | :---------------- |
| 路径 | `/api/words`      |
| 方法 | GET               |
| 鉴权 | 需要 Bearer Token |

**请求参数**

| 参数      | 类型 | 必填 | 默认值 | 说明              |
| :-------- | :--- | :--- | :----- | :---------------- |
| page      | int  | 否   | 1      | 页码，从1开始     |
| page_size | int  | 否   | 10     | 每页数量，最大100 |

**请求示例**

```
GET /api/words?page=2&page_size=10
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**成功响应（200）**

```
{
    "total": 25,
    "page": 2,
    "page_size": 10,
    "words": [
        {
            "id": 10,
            "word": "apple",
            "definition": "苹果",
            "sentences": [
                "I eat an apple every day.",
                "Apple is a fruit.",
                "The apple tree is blooming."
            ],
            "ai_provider": "qwen",
            "created_at": "2026-05-15"
        }
    ]
}
```

**错误响应**

| 状态码 | 错误信息    | 说明   |
| :----- | :---------- | :----- |
| 401    | 未提供token | 未认证 |

------

## 6. 删除单词

| 项目 | 内容              |
| :--- | :---------------- |
| 路径 | `/api/word/:id`   |
| 方法 | DELETE            |
| 鉴权 | 需要 Bearer Token |

**请求参数**

| 参数 | 类型 | 位置     | 说明    |
| :--- | :--- | :------- | :------ |
| id   | int  | 路径参数 | 单词 ID |

**请求示例**

```
DELETE /api/word/10
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**成功响应（200）**

```
{
    "message": "删除成功"
}
```

**错误响应**

| 状态码 | 错误信息    | 说明                           |
| :----- | :---------- | :----------------------------- |
| 401    | 未提供token | 未认证                         |
| 404    | 单词不存在  | 单词 ID 不存在或不属于当前用户 |

**补充说明**：采用软删除机制，删除后设置 `deleted_at` 字段，查询时不会返回。