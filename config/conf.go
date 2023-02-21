package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
)

type Mysql struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

var dbConf Mysql

func init() {
	if _, err := toml.DecodeFile("D:\\mydy\\miniDy\\config\\config.toml", &dbConf); err != nil {
		panic(err)
	}
}

func DSNString() string {
	arg := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Database)
	log.Println(arg)

	return arg
}
