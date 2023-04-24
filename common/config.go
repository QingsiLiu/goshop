package common

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func GetConsulConfig(url, fileKey string) (*viper.Viper, error) {
	conf := viper.New()
	err := conf.AddRemoteProvider("consul", url, fileKey)
	if err != nil {
		return nil, errors.Wrap(err, "get consul remote config error")
	}
	conf.SetConfigType("json")
	err = conf.ReadRemoteConfig()
	if err != nil {
		return nil, errors.Wrap(err, "viper conf error")
	}
	return conf, nil
}

func GetMysqlFromConsul(vip *viper.Viper) (db *gorm.DB, err error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			Colorful:      true,
			LogLevel:      logger.Info,
		},
	)

	userName := vip.GetString("user")
	pwd := vip.GetString("pwd")
	host := vip.GetString("host")
	port := vip.GetString("port")
	database := vip.GetString("database")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local", userName, pwd, host, port, database)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, errors.Wrap(err, "init mysql db error")
	}

	return db, nil
}

// 获取redis配置
func GetRedisFromConsul(vip *viper.Viper) (red *redis.Client, err error) {
	red = redis.NewClient(&redis.Options{
		Addr:         vip.GetString("addr"),
		Password:     vip.GetString("password"),
		DB:           vip.GetInt("DB"),
		PoolSize:     vip.GetInt("poolSize"),
		MinIdleConns: vip.GetInt("minIdleConn"),
	})

	// 集群client
	cluterClients := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{},
	})
	fmt.Println(cluterClients)

	return
}

// 设置用户token
func SetUserToken(red *redis.Client, key string, val []byte, timeTTL time.Duration) {
	red.Set(context.Background(), key, val, timeTTL)
}

// 获取用户token
func GetUserToken(red *redis.Client, key string) (string, error) {
	res, err := red.Get(context.Background(), key).Result()
	if err != nil {
		log.Printf("get user token from redis error: %v", err)
	}
	return res, nil
}
