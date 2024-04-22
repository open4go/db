package mongo

type MongoClientConf struct {
	Host string `mapstructure:"host"`
	Name string `mapstructure:"name"`
}
