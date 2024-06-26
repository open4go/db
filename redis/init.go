package redis

import (
	"context"
	"github.com/open4go/log"
	"github.com/spf13/viper"
)

// Init 快速执行初始化
func Init(ctx context.Context) {
	// 初始化 MongoDB 连接池
	NewDataBasePool()

	// 解析 YAML 数据到结构体数组
	var services []RedisClientConf
	err := viper.UnmarshalKey("db.redis", &services)
	if err != nil {
		log.Log(ctx).WithField("err", "Error unmarshaling services").
			Fatal(err)
		return
	}

	// 初始化所有数据库
	for _, i := range services {
		_, err := DBPool.GetClient(ctx, i.Host, i.Name, i.DB)
		if err != nil {
			log.Log(ctx).Fatal(err)
		}
	}
}
