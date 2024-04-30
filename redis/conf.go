package redis

type RedisClientConf struct {
	Host string `mapstructure:"host"`
	Name string `mapstructure:"name"`
	DB   int    `mapstructure:"db"`
}
