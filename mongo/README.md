# 如何使用数据库连接池

### config.yaml 配置

```shell
db:
  mongo:
    - host: mongodb://localhost:27077
      name: sys_auth
    - host: mongodb://localhost:27077
      name: m3s
    - host: mongodb://localhost:27077
      name: r2client
```

### 在main.go 中初始化

```shell
	// 初始化 MongoDB 连接池
	db.NewMongoDBPool()

	// 解析 YAML 数据到结构体数组
	var services []db.MongoClientConf
	err := viper.UnmarshalKey("db.mongo", &services)
	if err != nil {
		fmt.Printf("Error unmarshaling services: %s\n", err)
		return
	}

	// 初始化所有数据库
	for _, i := range services {
		_, err := db.GlobalMongoDBPool.GetClient(context.TODO(), i.Host, i.Name)
		if err != nil {
			log.Log().Fatal(err)
		}
	}
```

## 优化说明
>之前MongoDB连接和管理库的代码结构总体上设计得不错，但有一些潜在的问题和改进的地方。下面是我对您代码的分析和提出的改进建议：

### 1. 代码改进点

#### a) 避免全局状态
您使用全局的 `DBPool` 变量来管理数据库连接，这在并发环境下可能会导致问题。建议通过函数参数或者使用依赖注入的方式来传递这些依赖，避免全局状态。

#### b) 错误处理
代码中使用 `log.Fatal` 在遇到错误时终止程序。这对于库函数来说通常不是最佳实践，因为它会导致调用者应用程序的意外退出。最好是返回错误给调用者，让调用者决定如何处理错误。

#### c) 线程安全与连接复用
- 在 `GetClient` 函数中，`ctx` 被用作全局变量，这在多线程环境下可能会引起问题。每个请求应该有自己的 context。
- 数据库连接池应该支持连接复用，而不是每次都建立新连接。`mongo.Connect` 应该在确认没有现有连接的情况下才执行。

#### d) 斜线处理
将斜线加到 host 的末尾可能不是必要的，因为 MongoDB 的 URI 标准不要求在 host 后面加斜线。

### 2. 其他考虑
- **连接字符串**: 确保连接字符串遵循MongoDB的标准格式。
- **配置与环境管理**: 使用 `viper` 管理配置是个好做法，确保配置格式和使用符合预期。
- **安全性**: 确保敏感信息如数据库密码等不要硬编码在代码中，使用环境变量或安全的配置文件。