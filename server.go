package main

import (
	"fmt"
	"godrink/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

var config util.ConfigData
var errChan chan error

func runServer() {
	engine := gin.Default()
	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{"pong": true})
	})
	err := engine.Run(fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port))
	if err != nil {
		errChan <- err
	}
}

func main() {
	config = util.GetConfig()
	errChan = make(chan error)
	go runServer()
	fmt.Printf("Server running on %s:%s", config.Server.Host, config.Server.Port)

	if err := <-errChan; err != nil {
		panic(fmt.Errorf("error during server execution: %w", err))
	}
}
