package tools

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	acc "soa-hw-ilyaleshchyk/internal/account"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type JWTManager struct {
	jwtPrivate *rsa.PrivateKey
	jwtPublic  *rsa.PublicKey
	db         *gorm.DB
}

func (j *JWTManager) InitJWT(privSecret, publicSecret string) {
	private, err := os.ReadFile(privSecret)
	if err != nil {
		logrus.Panic(err.Error() + privSecret + "--------")
	}

	public, err := os.ReadFile(publicSecret)
	if err != nil {
		logrus.Panic(err.Error() + publicSecret)
	}

	j.jwtPrivate, err = jwt.ParseRSAPrivateKeyFromPEM(private)
	if err != nil {
		logrus.Panic(err)
	}

	j.jwtPublic, err = jwt.ParseRSAPublicKeyFromPEM(public)
	if err != nil {
		logrus.Panic(err)
	}
}

func (s *JWTManager) CheckJwt(c *gin.Context) {

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

	c.Header("accountID", strconv.FormatInt(int64(account.ID), 10))
}

func (s *JWTManager) InitDB(db string, dbDebug bool) {
	var err error

	cfg := &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}

	if dbDebug {
		cfg.Logger = logger.Default.LogMode(logger.Info)
	}

	s.db, err = gorm.Open(postgres.Open(db), cfg)
	if err != nil {
		logrus.Panic("failed to connect database:", err)
	}

	sqlDB, err := s.db.DB()
	if err != nil {
		logrus.Panic("failed to get sql database:", err)
	}

	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)
}

func (j *JWTManager) GenJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour * 12).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(j.jwtPrivate)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
