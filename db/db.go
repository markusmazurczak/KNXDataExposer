package db

import (
	"fmt"

	"gitea.realmottek.duckdns.org/nexus/KNXDataExposer/handler"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(logger *zap.Logger) *gorm.DB {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable&TimeZone=%s",
		viper.GetString("postgres.username"),
		viper.GetString("postgres.password"),
		viper.GetString("postgres.server"),
		viper.GetInt("postgres.port"),
		viper.GetString("postgres.dbName"),
		viper.GetString("postgres.timeZone"))

	logger.Info("Postgres DB connection",
		zap.String("host", viper.GetString("postgres.server")),
		zap.String("port", viper.GetString("postgres.port")),
		zap.String("user", viper.GetString("postgres.username")),
		zap.String("database", viper.GetString("postgres.dbName")),
		zap.String("timezone", viper.GetString("postgres.timeZone")))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		logger.Sugar().Fatalf(err.Error())
	}

	db.AutoMigrate(&handler.Dataset{})

	return db
}
