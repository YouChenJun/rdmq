// @Author Chen_dark
// @Date 2024/10/31 11:16:00
// @Desc
package redmq

type RedisConfig struct {
	Network  string `mapstructure:"network"`
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	Topic    string `mapstructure:"topic"`
}
