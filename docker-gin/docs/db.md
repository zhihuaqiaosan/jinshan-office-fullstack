# 数据库设计文档

## 基础信息

| 项目       | 内容               |
| ---------- | ------------------ |
| 数据库类型 | MySQL 8.0          |
| 字符集     | utf8mb4            |
| 排序规则   | utf8mb4_unicode_ci |
| 数据库名   | wordbook           |

## 一、users 表（用户表）

存储用户账号信息。

### 字段设计

| 字段       | 数据类型        | 约束                         | 说明                |
| ---------- | --------------- | ---------------------------- | ------------------- |
| id         | bigint unsigned | PRIMARY KEY, AUTO_INCREMENT  | 用户唯一标识        |
| username   | varchar(50)     | NOT NULL, UNIQUE             | 用户名，登录凭证    |
| password   | varchar(255)    | NOT NULL                     | bcrypt 加密后的密码 |
| created_at | datetime(3)     | DEFAULT CURRENT_TIMESTAMP(3) | 注册时间            |

### 索引设计

| 索引名             | 字段     | 类型     | 说明                         |
| ------------------ | -------- | -------- | ---------------------------- |
| PRIMARY            | id       | 主键索引 | 聚簇索引，快速定位记录       |
| idx_users_username | username | 唯一索引 | 快速查询用户，保证用户名唯一 |

### 建表语句

```sql
CREATE TABLE IF NOT EXISTS `users` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `username` varchar(50) NOT NULL,
    `password` varchar(255) NOT NULL,
    `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_users_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

## 二、words 表（单词本表）

存储用户保存的单词记录。

### 字段设计

| 字段        | 数据类型        | 约束                         | 说明                        |
| :---------- | :-------------- | :--------------------------- | :-------------------------- |
| id          | bigint unsigned | PRIMARY KEY, AUTO_INCREMENT  | 单词记录唯一标识            |
| user_id     | bigint unsigned | NOT NULL, INDEX, FOREIGN KEY | 所属用户 ID，关联 users.id  |
| word        | varchar(100)    | NOT NULL, INDEX              | 英文单词                    |
| definition  | text            | NOT NULL                     | 中文释义                    |
| sentences   | text            | NOT NULL                     | JSON 格式，存储 3 条例句    |
| ai_provider | varchar(50)     | DEFAULT NULL                 | AI 来源（qwen）             |
| created_at  | datetime(3)     | DEFAULT CURRENT_TIMESTAMP(3) | 保存时间                    |
| deleted_at  | datetime(3)     | DEFAULT NULL, INDEX          | 软删除时间，NULL 表示未删除 |

### 索引设计

| 索引名               | 字段       | 类型     | 说明                     |
| :------------------- | :--------- | :------- | :----------------------- |
| PRIMARY              | id         | 主键索引 | 聚簇索引，快速定位记录   |
| idx_words_user_id    | user_id    | 普通索引 | 快速查询某用户的所有单词 |
| idx_words_word       | word       | 普通索引 | 快速按单词查询           |
| idx_words_deleted_at | deleted_at | 普通索引 | 快速过滤已删除记录       |

### 外键约束

```
CONSTRAINT `fk_words_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
```

- `user_id` 引用 `users.id`
- `ON DELETE CASCADE`：用户删除时，级联删除该用户的所有单词记录

### 建表语句

```
CREATE TABLE IF NOT EXISTS `words` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_id` bigint unsigned NOT NULL,
    `word` varchar(100) NOT NULL,
    `definition` text NOT NULL,
    `sentences` text NOT NULL,
    `ai_provider` varchar(50) DEFAULT NULL,
    `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
    `deleted_at` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_words_user_id` (`user_id`),
    KEY `idx_words_word` (`word`),
    KEY `idx_words_deleted_at` (`deleted_at`),
    CONSTRAINT `fk_words_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

## 三、表关联关系说明

| 项目     | 说明                                             |
| :------- | :----------------------------------------------- |
| 关系类型 | 一对多（One-to-Many）                            |
| 关联字段 | `words.user_id` → `users.id`                     |
| 业务含义 | 一个用户可以拥有多个单词，每个单词只属于一个用户 |
| 级联删除 | 用户删除时，自动删除该用户的所有单词记录         |

### 关联查询示例

**查询用户的所有单词（分页）**：

```
SELECT id, word, definition, sentences, ai_provider, created_at
FROM words
WHERE user_id = ? AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT ? OFFSET ?;
```

**查询用户是否已保存某个单词**：

```
SELECT COUNT(*) FROM words
WHERE user_id = ? AND word = ? AND deleted_at IS NULL;
```

**软删除单词**：

```
UPDATE words
SET deleted_at = NOW()
WHERE id = ? AND user_id = ? AND deleted_at IS NULL;
```

## 四、初始化脚本

完整建表脚本见 `docs/init.sql`，Docker 启动时自动执行。