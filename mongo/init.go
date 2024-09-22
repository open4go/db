package mongo

import (
	"context"
	"github.com/spf13/viper"
)

// Init 快速执行初始化
func Init(ctx context.Context) error {
	DBPool = NewDataBasePool(ctx)

	var services []MongoClientConf
	if err := viper.UnmarshalKey("db.mongo", &services); err != nil {
		return err // 处理配置解析错误
	}

	for _, service := range services {
		if _, err := DBPool.GetClient(ctx, service.Host, service.Name); err != nil {
			return err // 返回错误，由调用者决定如何处理
		}
	}

	return nil // 成功初始化
}
