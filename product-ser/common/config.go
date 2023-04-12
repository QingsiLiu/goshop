package common

import (
	"fmt"
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
