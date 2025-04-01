package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func reverseProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.Errorf("Trying to redirect")
		remote, err := url.Parse(target)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid target URL"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		c.Request.URL.Host = remote.Host
		c.Request.URL.Scheme = remote.Scheme
		c.Request.Host = remote.Host

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func runWWW(bind string) {
	var www *gin.Engine

	gin.SetMode(gin.ReleaseMode)

	www = gin.New()
	www.Use(gin.Recovery())

	www.RedirectTrailingSlash = true
	www.RedirectFixedPath = true

	api := www.Group("/api/")

	{
		api.POST("/account/register", reverseProxy("http://"+config.AccountsServiceHost))
		api.POST("/account/login", reverseProxy("http://"+config.AccountsServiceHost))
		api.GET("/account/:account_id/profile", reverseProxy("http://"+config.AccountsServiceHost))
		api.PATCH("/account/profile", reverseProxy("http://"+config.AccountsServiceHost))
	}

	logrus.Infof("Application starting on addres: %s", bind)
	logrus.Fatalln(www.Run(bind))
}

func main() {

	loadConfig()
	go runWWW(config.Bind)

	select {}
}
