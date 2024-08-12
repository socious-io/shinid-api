package main

import (
	"fmt"
	"shinid/src/app"
	"shinid/src/config"
	"shinid/src/database"
	"time"
)

func main() {
	config.Init("config.yml")
	fmt.Println(config.Config)
	database.Connect(&database.ConnectOption{
		URL:         config.Config.Database.URL,
		SqlDir:      config.Config.Database.SqlDir,
		MaxRequests: 5,
		Interval:    30 * time.Second,
		Timeout:     5 * time.Second,
	})

	app.Serve()
}
