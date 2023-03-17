package main

import (
	"kvsql/apiserver"
	"kvsql/common"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/rs/zerolog/log"
)

func main() {
	err := common.Client.Init("postgres", "root", "kvstore", "localhost")
	if err != nil {
		log.Error().Err(err)
		return
	}
	err = common.Client.Test()
	if err != nil {
		log.Error().Err(err)
		return
	}
	err = common.Client.CreateTable()
	if err != nil {
		log.Error().Err(err)
		return
	}
	wg := sync.WaitGroup{}
	defer wg.Wait()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Debug().Msg("running clean job")
			err := common.Client.BatchDelete()
			if err != nil {
				log.Fatal().Err(err)
			}
			time.Sleep(1 * time.Minute)
		}
	}()
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", apiserver.Test)
	e.POST("/putKey", apiserver.PutKey)
	e.GET("/getKey", apiserver.GetKey)
	e.DELETE("/deleteKey", apiserver.DeleteKey)
	// Start server
	e.Logger.Fatal(e.Start(":1323"))

}
