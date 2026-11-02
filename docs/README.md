# db — MongoDB / Redis 连接池

> 仓库：https://github.com/open4go/db ｜ 模块：`github.com/open4go/db` ｜ Go 1.25.8 ｜ 最近提交：2026-08-21（最新 tag v0.0.15）

## 项目简介
从 viper 配置（`db.mongo`、`db.redis`）读取多个数据源，启动时建立 MongoDB 客户端和 Redis 客户端，放进进程级全局池 `DBPool`，按名称取 handler。

## 目录结构
```
db/
├── mongo/  conf.go(MongoClientConf) handler.go(DataBasePool) init.go(Init) README.md
├── redis/  conf.go(RedisClientConf) handler.go(DataBasePool) init.go(Init) README.md
└── .github/workflows/go.yml
```

## 对外 API
**mongo**
- `var DBPool *DataBasePool`
- `Init(ctx) error`：读 `db.mongo`（`[]{host,name}`），逐个 `GetClient`，失败时返回 error
- `NewDataBasePool(ctx)`、`(*DataBasePool).GetClient(ctx, host, name)`（连接 + Ping，按 `host+name` 缓存）
- `GetHandler(name) (*mongo.Database, error)`、`CloseAll(ctx)`

**redis**
- `var DBPool *DataBasePool`
- `Init(ctx)`：读 `db.redis`（`[]{host,name,db}`），**失败时直接 `Fatal`**
- `NewDataBasePool()`、`GetClient(ctx, host, name, db)`（PoolSize 100、MaxRetries 3，按 `name` 缓存）
- `GetHandler(name)`、`CloseAll()`

## 配置示例（mongo/README.md）
```yaml
db:
  mongo:
    - host: mongodb://localhost:27077
      name: sys_auth
  redis:
    - host: localhost:6379
      name: cache
      db: 0
```

## 使用示例（真实调用方）
```go
// middle/db.go、auth/db.go、r2id/db.go、data/db/db.go、util9s/db/redis.go、p7/*/db.go 都有类似写法
handler, err := r2redis.DBPool.GetHandler("cache")
if err != nil { log.Log(ctx).Fatal(err) }
// auth/app.go
handler, err := r2mongo.DBPool.GetHandler("sys_auth")
```

## 依赖
`go.mongodb.org/mongo-driver v1.17.7`、`github.com/redis/go-redis/v9 v9.7.3`、`github.com/spf13/viper`、`github.com/open4go/log v0.0.12`

## 构建
`go build ./... && go vet ./...` 通过，没有测试。

## 与其他仓库的关系
被 `auth`、`middle`、`r2id`、`data`、`util9s`、`p7`、`req5rsp`（间接）依赖。依赖 `log`（旧版 v0.0.12）。
