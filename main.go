package main

import (
	"app/shared/database/mysql"
	"app/shared/database/redis"
	"app/shared/jsonconfig"
	"app/shared/server"
	"app/controller"
	"app/model/news"
	"app/model/user"
	"encoding/json"
	"time"
	"os"
)

func init() {
	if os.Getenv("ENV") == "production" {
		jsonconfig.Load("config"+string(os.PathSeparator)+"config.json", config)
	} else {
		jsonconfig.Load("config"+string(os.PathSeparator)+"config-dev.json", config)
	}
	mysql.Connect(config.Database)
	mysql.ConnectDB2(config.Database2)
	go func() {
		for {
			news.LoadNewsSimilarityToCache()
			time.Sleep(time.Duration(1) * time.Minute)
		}
	}()

	go func() {
		for {
			news.LoadTrendingNews()
			time.Sleep(time.Duration(5) * time.Minute)
		}
	}()

	go func() {
		for {
			news.LoadProbabilityCategoryEqualCi()
			time.Sleep(time.Duration(15) * time.Minute)
		}
	}()

	go func() {
		for {
			user.LoadUserProfileToCache()
			time.Sleep(time.Duration(5) * time.Minute)
		}
	}()
}

func main() {
	controller.Load()
}

var config = &configuration{}

// configuration contains the application settings
type configuration struct {
	ServerInfo server.ServerInfo `json:"Server"`
	Database   mysql.MySQLInfo   `json:"MYSQL"`
	Database2   mysql.MySQLInfo   `json:"MYSQL2"`
	RedisInfo  redis.RedisInfo   `json:"Redis"`
}

// ParseJSON unmarshals bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}
