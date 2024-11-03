// @Author Chen_dark
// @Date 2024/10/31 11:16:00
// @Desc
package redis

type RedisConfig struct {
	Network  string `mapstructure:"network"`  // 网络类型 TCP
	Address  string `mapstructure:"address"`  // redis地址 ip:port
	Password string `mapstructure:"password"` // redis密码
	Topic    string `mapstructure:"topic"`    // topic名称
}
type Consumer struct {
	consumerGroup string //消费者组名称
	consumerID    string //消费者名称
}
