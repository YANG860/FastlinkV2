package db

import (
	"context"
	"fastlink/src/config"
	
	"log/slog"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var RedisClient *redis.Client
var MySQLClient *gorm.DB
var Ctx = context.Background()

func connectRedis() error {

	redisConfig := config.Redis()
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})
	return RedisClient.Ping(context.Background()).Err()
}

func connectMySQL() error {
	mysqlConfig := config.MySQL()

	var err error
	MySQLClient, err = gorm.Open(mysql.Open(mysqlConfig.DSN), &gorm.Config{
	})
	if err == nil {
		err = MySQLClient.AutoMigrate(&User{}, &Link{}, &Admin{})
	}
	return err
}

func init() {
	err := connectRedis()
	if err != nil {
		slog.Error("Redis connection failed", "error", err)
		panic(err)
	}
	err = connectMySQL()
	if err != nil {
		slog.Error("Database connection failed", "error", err)
		panic(err)
	}

	slog.Info("Database connected successfully")
}
