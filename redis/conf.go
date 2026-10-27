package redis

type RedisClientConf struct {
	Host string `mapstructure:"host"`
	Name string `mapstructure:"name"`
	DB   int    `mapstructure:"db"`
	// Password is the Redis AUTH password (optional; empty = no auth)
	Password string `mapstructure:"password"`
}
