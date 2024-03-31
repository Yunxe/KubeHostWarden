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

type MysqlClient struct {
	client *gorm.DB
}

func SetupMysql(database *gorm.DB) *MysqlClient {
	if database != nil {
		return &MysqlClient{
			client: database,
		}
	}

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
		panic("failed to connect database")
	}

	return &MysqlClient{
		client: client,
	}
}
