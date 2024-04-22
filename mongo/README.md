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