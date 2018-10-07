package redis

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var startTime time.Time
var serverRedisConfig string
var passwordRedisConfig string
var clientRedis *redis.Client
var remainDB int

func uptime() time.Duration {
	return time.Since(startTime)
}

func loadConf() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("../../")
	viper.AddConfigPath("../")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("\n Fatal error config file: %s ", err))
	}
	serverRedisConfig = viper.GetString("server_redis")
	passwordRedisConfig = viper.GetString("password_redis")
	remainDB = viper.GetInt("remain_DB")
}

func init() {
	startTime = time.Now()
	loadConf()
	clientRedis, _ = GetConnection(remainDB)
}
func main() {
	clientRedis, _ = GetConnection(remainDB)
}
func GetConnection(redisDB int) (*redis.Client, error) {
	if clientRedis == nil {
		return ConnectRedisServer(redisDB)
	}
	pong, err := clientRedis.Ping().Result()
	if pong == "PONG" && err == nil {
		//log.Println("GetConnection redis: ", pong)
		return clientRedis, nil
	} else {
		//log.Println("GetConnection redis: ", err)
		time.Sleep(2 * time.Second)
		return ConnectRedisServer(redisDB)
	}
}
func ConnectRedisServer(redisDB int) (*redis.Client, error) {
	serverRedises := serverRedisConfig
	client := redis.NewClient(&redis.Options{
		Addr:     serverRedises,
		Password: passwordRedisConfig,
		DB:       redisDB,
	})
	pong, err := client.Ping().Result()
	if err == nil {
		log.Println("Connection redis: ", pong)
	} else {
		log.Println("Connection redis: ", err)
		time.Sleep(2 * time.Second)
		return ConnectRedisServer(redisDB)
	}
	return client, err
}
