package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"soa-hw-ilyaleshchyk/internal/tools"

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
		c.Request.Header.Set("accountID", c.Writer.Header().Get("accountID"))

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

	jwtManager := &tools.JWTManager{}
	jwtManager.InitDB(config.DB, config.DBDebug)
	jwtManager.InitJWT(config.PrivateSecret, config.PublicSecret)

	api := www.Group("/api/")

	apiAuth := api.Group("", jwtManager.CheckJwt)

	{
		// account service
		api.POST("/account/register", reverseProxy("http://"+config.AccountsServiceHost))
		api.POST("/account/login", reverseProxy("http://"+config.AccountsServiceHost))
		apiAuth.GET("/account/:account_id/profile", reverseProxy("http://"+config.AccountsServiceHost))
		apiAuth.PATCH("/account/profile", reverseProxy("http://"+config.AccountsServiceHost))

		// posts service
		apiAuth.POST("/posts/create", reverseProxy("http://"+config.PostsServiceHost))
		apiAuth.DELETE("/posts/:post_id", reverseProxy("http://"+config.PostsServiceHost))
		apiAuth.PATCH("/posts/:post_id", reverseProxy("http://"+config.PostsServiceHost))
		apiAuth.GET("/posts/:post_id", reverseProxy("http://"+config.PostsServiceHost))
		apiAuth.GET("/account/posts/:owner_id", reverseProxy("http://"+config.PostsServiceHost))
	}

	logrus.Infof("Application starting on addres: %s", bind)
	logrus.Fatalln(www.Run(bind))
}

func main() {

	loadConfig()
	go runWWW(config.Bind)

	select {}
}
