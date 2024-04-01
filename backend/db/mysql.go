package mysql

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type mysqlConfig struct {
	user     string
	password string
	address  string
	port     string
	database string
}

var mc *MysqlClient

type MysqlClient struct {
	Client *gorm.DB
}

func GetClient() *MysqlClient {
	return mc
}

func SetupMysql() error {
	config := &mysqlConfig{
		user:     os.Getenv("MYSQL_USER"),
		password: os.Getenv("MYSQL_PASSWORD"),
		address:  os.Getenv("MYSQL_ADDRESS"),
		port:     os.Getenv("MYSQL_PORT"),
		database: os.Getenv("MYSQL_DATABASE"),
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.user, config.password, config.address, config.port, config.database)
	client, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}
	mc = &MysqlClient{Client: client}
	return nil
}
