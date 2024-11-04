// @Author Chen_dark
// @Date 2024/10/31 11:16:00
// @Desc
package redis

type Config struct {
	Network  string `mapstructure:"network"`  // 网络类型 TCP
	Address  string `mapstructure:"address"`  // redis地址 ip:port
	Password string `mapstructure:"password"` // redis密码
	//Topic    string `mapstructure:"topic"`    // topic名称
}
type ConsumerConfig struct {
	RedisConfig Config
	GroupID     string //消费者组名称
	ConsumerID  string //消费者名称
	Topic       string
}

type ProducerConfig struct {
	RedisConfig Config
	Topic       string `mapstructure:"topic"`
}
