package main

import (
	"shin/src/app"
	"shin/src/app/services"
	"shin/src/config"
	"shin/src/database"
	"time"
)

func main() {
	config.Init("config.yml")
	database.Connect(&database.ConnectOption{
		URL:         config.Config.Database.URL,
		SqlDir:      config.Config.Database.SqlDir,
		MaxRequests: 5,
		Interval:    30 * time.Second,
		Timeout:     5 * time.Second,
	})
	services.InitSendGridService()

	app.Serve()
}
