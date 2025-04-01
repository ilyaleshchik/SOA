package main

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"os"
	acc "soa-hw-ilyaleshchyk/internal/account"
	"soa-hw-ilyaleshchyk/internal/post"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// @title posts API
// @version 1.0
// @description API for maniging posts
// @BasePath /api
// @securityDefinitions.apikey Bearer
// @in header
// @name Auth
// @description Enter the token with the `Bearer: ` prefix.

type Server struct {
	db         *gorm.DB
	jwtPrivate *rsa.PrivateKey
	jwtPublic  *rsa.PublicKey
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) runWWW() {
	var www *gin.Engine

	gin.SetMode(gin.ReleaseMode)

	www = gin.New()
	www.Use(gin.Recovery())

	www.RedirectTrailingSlash = true
	www.RedirectFixedPath = true

	api := www.Group("/api/")

	{
		api.POST("/posts/create", wwwHandler(s.createPost))
		api.DELETE("/posts/:post_id", wwwHandler(s.deletePost))
		api.PATCH("/posts/:post_id", wwwHandler(s.updatePost))
		api.GET("/posts/:post_id", wwwHandler(s.getPost))
		api.GET("/account/posts/:owner_id", wwwHandler(s.getPosts))
	}

	logrus.Infof("Application starting on addres: %s", config.Bind)
	logrus.Fatalln(www.Run(config.Bind))
	logrus.Infof("Application started on addres: %s", config.Bind)
}

func (s *Server) InitDB() {
	var err error

	cfg := &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}

	if config.DBDebug {
		cfg.Logger = logger.Default.LogMode(logger.Info)
	}

	s.db, err = gorm.Open(postgres.Open(config.DB), cfg)
	if err != nil {
		logrus.Panic("failed to connect database:", err)
	}

	sqlDB, err := s.db.DB()
	if err != nil {
		logrus.Panic("failed to get sql database:", err)
	}

	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)

	if err := s.db.AutoMigrate(&post.Post{}); err != nil {
		logrus.Panic("can't automig posts")
	}

}

func (s *Server) InitJWT() {
	private, err := os.ReadFile(config.PrivateSecret)
	if err != nil {
		logrus.Panic(err.Error() + config.PrivateSecret + "--------")
	}

	public, err := os.ReadFile(config.PublicSecret)
	if err != nil {
		logrus.Panic(err.Error() + config.PublicSecret)
	}

	s.jwtPrivate, err = jwt.ParseRSAPrivateKeyFromPEM(private)
	if err != nil {
		logrus.Panic(err)
	}

	s.jwtPublic, err = jwt.ParseRSAPublicKeyFromPEM(public)
	if err != nil {
		logrus.Panic(err)
	}
}

func (s *Server) genJWT(username string) (string, error) {

	claims := jwt.MapClaims{
		"username": username,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour * 12).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(s.jwtPrivate)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *Server) checkJwt(c *gin.Context) {

	jwtSession := c.GetHeader("Auth")

	token, err := jwt.Parse(jwtSession, func(token *jwt.Token) (interface{}, error) {
		if alg, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		} else if alg != jwt.SigningMethodRS256 {
			return nil, fmt.Errorf("signing method does not match: %v", token.Header["alg"])
		}

		return s.jwtPublic, nil
	})

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, claimsOk := token.Claims.(jwt.MapClaims)
	if !claimsOk || !token.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	username, ok := claims["username"].(string)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	account, err := acc.GetAccountByUsername(s.db, username)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if account == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("account", account)
}
