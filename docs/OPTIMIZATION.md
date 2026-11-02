# db 优化方案

| # | 优先级 | 问题 | 位置 | 建议 | 工作量 |
|---|---|---|---|---|---|
| 1 | 高 | 数据竞争：`mongo.GetHandler` 读 map 时没加锁，和 `GetClient` 的写操作并发；redis `GetHandler` 也一样 | `mongo/handler.go:71`、`redis/handler.go` GetHandler | 改用 `sync.RWMutex`，读加 RLock；或者初始化完成后把 map 冻结成只读 | S |
| 2 | 高 | mongo 的 handler 按 `name` 存，client 按 `host+name` 存：两个不同实例上的同名库会互相覆盖 | `mongo/handler.go:51` | key 统一用配置里的逻辑名（新增 `alias` 字段），重名时报错 | S |
| 3 | 中 | `CloseAll` 第一个 Disconnect 出错就 `return`，后面的连接不会关闭 | `mongo/handler.go:65` | 遍历全部，用 `errors.Join` 汇总错误并返回 | S |
| 4 | 中 | redis 的去斜杠逻辑写反了：注释说“没有斜线就加上”，代码却在没有 `/` 时 Trim，实际什么也没做 | `redis/handler.go:37` | 删掉这段（go-redis 的 Addr 不需要斜线） | S |
| 5 | 中 | 错误处理不一致：mongo.Init 返回 error，redis.Init 直接 `Fatal`；下游 8 个仓库各自写一遍 `GetRedisXxxHandler + Fatal` | `redis/init.go`；下游 `*/db.go` | redis.Init 改为返回 error；新增 `MustRedis(name)` / `Redis(name) (*Client, error)`，下游删除重复的 helper | M |
| 6 | 中 | 每次 `GetClient` 都用请求或启动时的 ctx 去 Connect/Ping，没有超时 | `GetClient` | 用 `context.WithTimeout(ctx, 5*time.Second)`；mongo 加 `SetServerSelectionTimeout` | S |
| 7 | 中 | Redis 只支持 Addr/DB，不支持密码、TLS、Sentinel/Cluster | `redis/conf.go` | 增加 `password`、`tls`、`url` 字段（`redis.ParseURL`） | S |
| 8 | 中 | README 示例过时（`db.NewMongoDBPool`、`GlobalMongoDBPool` 已不存在） | `mongo/README.md` | 按现有 API 重写 | S |
| 9 | 低 | 没有测试 | — | 用 miniredis + mongo testcontainers 覆盖 Init/GetHandler/并发 | M |
| 10 | 中 | govulncheck：`golang.org/x/text v0.17.0`（GO-2026-5970）；旧版 `log v0.0.12` 间接引入 docker（GO-2026-4887/4883） | `go.mod` | 升级 log 到 v0.0.23+ 和 x/text；修好后 CI 才能变绿 | S |
