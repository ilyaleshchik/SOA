package main

import (
	"net/http"
	"soa-hw-ilyaleshchyk/api"
	acc "soa-hw-ilyaleshchyk/internal/account"
	"soa-hw-ilyaleshchyk/internal/tools"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// registerAccount registers a new account
// @Summary Register an account
// @Description Creates a new account with username, email, and password
// @Tags account
// @Accept  json
// @Produce  json
// @Param request body api.AccountRegisterRequest true "Account Register Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func (s *Server) registerAccount(c *gin.Context) error {

	req := api.AccountRegisterRequest{}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong format of request",
		})
		return nil
	}

	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "bad email",
		})
		return nil
	}

	if req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "bad username",
		})
		return nil
	}

	if req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "bad password",
		})
		return nil
	}

	logrus.Errorf("Req.Password: %s", req.Password)
	hashedPassword, err := tools.HashPassword(req.Password)
	if err != nil {
		return err
	}
	logrus.Errorf("Hashed.Password: %s", string(hashedPassword))

	account, err := acc.GetAccountByUsername(s.db, req.Username)
	if err != nil {
		return err
	}

	if account != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "username already exists",
		})
		return nil
	}

	newAcc := &acc.Account{
		ID:       tools.GenerateRandomAccountID(),
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := acc.SaveAccount(s.db, newAcc); err != nil {
		return err
	}

	session, err := s.jwtManager.GenJWT(newAcc.Username)

	if err != nil {
		return err
	}

	c.Header("Authorization", session)
	c.JSON(http.StatusOK, gin.H{"jwt": session})
	return nil
}

// login logs in an existing user
// @Summary Login to account
// @Description Authenticates user and returns JWT session token
// @Tags account
// @Accept  json
// @Produce  json
// @Param request body api.LoginRequest true "Login Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func (s *Server) login(c *gin.Context) error {
	var req api.LoginRequest

	if err := c.BindJSON(&req); err != nil {
		return err
	}

	if req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "empty username",
		})
		return nil
	}

	if req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "empty password",
		})
		return nil
	}

	acc, err := acc.GetAccountByUsername(s.db, req.Username)
	if err != nil {
		return err
	}

	if acc == nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return nil
	}

	logrus.Errorf("Req.Password: %s", req.Password)
	hashedPassword, err := tools.HashPassword(req.Password)
	if err != nil {
		return err
	}
	logrus.Errorf("Hashed.Password: %s", string(hashedPassword))

	if err := tools.IsSamePass([]byte(req.Password), []byte(acc.Password)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"status": "wrong password",
		})
		return nil
	}

	session, err := s.jwtManager.GenJWT(acc.Username)

	if err != nil {
		return err
	}

	c.Header("Authorization", session)
	c.JSON(http.StatusOK, gin.H{"jwt": session})
	return nil
}

// getAccountProfile retrieves the profile of an account
// @Summary Get account profile
// @Description Fetches account profile details
// @Tags account
// @Accept  json
// @Produce  json
// @Param account_id path int true "Account ID"
// @Success 200 {object} api.AccountProfile
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /account/{account_id} [get]
func (s *Server) getAccountProfile(c *gin.Context) error {

	reqID, err := strconv.ParseInt(c.Param("account_id"), 10, 64)
	if err != nil {
		return err
	}

	accountID := acc.AccountID(reqID)

	account, err := acc.GetAccountByID(s.db, accountID)
	if err != nil {
		return err
	}

	accProfile, err := acc.GetOrCreateAccountProfile(s.db, accountID)
	if err != nil {
		return err
	}

	resp := &api.AccountProfile{
		Username:  account.Username,
		Email:     account.Email,
		BirthDate: time.Unix(account.Birthday.Unix(), 0).Format("2006-01-02"),
		Bio:       accProfile.Bio,
		LastName:  accProfile.LastName,
		Name:      accProfile.Name,
		Phone:     account.Phone,
	}

	c.JSON(http.StatusOK, resp)
	return nil
}

// updateAccountProfile updates the profile of the logged-in user
// @Summary Update account profile
// @Description Modifies the account profile details
// @Tags account
// @Accept  json
// @Produce  json
// @Param request body api.UpdateAccountProfileRequest true "Update Account Profile Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /account/profile [patch]
func (s *Server) updateAccountProfile(c *gin.Context) error {

	accountID, err := strconv.ParseInt(c.GetHeader("accountID"), 10, 64)
	if err != nil {
		return err
	}

	account, err := acc.GetAccountByID(s.db, acc.AccountID(accountID))
	if err != nil {
		return err
	}

	accProfile, err := acc.GetOrCreateAccountProfile(s.db, account.ID)

	if err != nil {
		return err
	}

	var req api.UpdateAccountProfileRequest

	if err := c.BindJSON(&req); err != nil {
		return err
	}

	accountToUpdate := make(map[string]interface{})
	profileToUpdate := make(map[string]interface{})

	if req.Name != nil {
		accProfile.Name = *req.Name
		profileToUpdate["name"] = *req.Name
	}

	if req.LastName != nil {
		accProfile.LastName = *req.LastName
		profileToUpdate["last_name"] = *req.LastName
	}

	if req.Phone != nil {
		accountToUpdate["phone"] = *req.Phone
		account.Phone = *req.Phone
	}

	if req.Email != nil {
		accountToUpdate["email"] = *req.Email
		account.Email = *req.Email
	}

	if req.Birthday != nil {
		var birth time.Time
		if len(*req.Birthday) > 10 {
			birth, err = time.Parse(time.RFC3339, *req.Birthday)
		} else {
			birth, err = time.Parse("2006-01-02", *req.Birthday)
		}

		account.Birthday = birth
		accountToUpdate["birthday"] = birth
	}

	err = s.db.Model(&acc.Account{}).Where("id = ?", account.ID).Updates(accountToUpdate).Error
	if err != nil {
		return err
	}

	err = s.db.Model(&acc.AccountProfile{}).Where("account_id = ?", account.ID).Updates(profileToUpdate).Error

	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{})
	return nil
}

//curl -v -X POST http://127.0.0.1:8081/api/account/login -H "Content-Type: application/json" -d '{"username": "lolfucjj", "password":"au228"}'
