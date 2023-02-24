package main

import (
	"onlinestatusindicator/redis"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var redisClient redis.RClient = redis.RClient{}

func main() {
	// Echo instance
	e := echo.New()

	redisClient.Init("localhost:6379")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", getUsers)
	e.GET("/onlineUser/:usrName", userOnline)
	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
