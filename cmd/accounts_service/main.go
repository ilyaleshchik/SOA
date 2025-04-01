package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type wwwHandlerT func(*gin.Context) error

func wwwHandler(f wwwHandlerT) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := f(c); err != nil {
			logrus.Errorf("WWW %s error: %v", c.FullPath(), err.Error())
		}
	}
}

func main() {
	// ctx := context.Background()

	server := NewServer()
	loadConfig()
	server.InitDB()
	server.InitJWTManager()

	go server.runWWW()

	select {}
}
