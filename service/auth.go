package service

import (
	"encoding/json"
	"time"
	"user/constant"
	"user/pkg/config"
	"user/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func jwtSecret() []byte {
	s, ok := config.Get("config.jwt.secret").(string)
	if !ok || s == "" {
		panic("FATAL: jwt.secret is not configured. Set it in config.yaml or via JWT_SECRET env var.")
	}
	return []byte(s)
}

func (s *Service) generateToken(userID int) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret())
}

func (s *Service) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret(), nil
	})
	if err != nil {
		return nil, util.NewCodeError(constant.ERROR_AUTH_FAIL, "invalid token")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, util.NewCodeError(constant.ERROR_AUTH_FAIL, "invalid token")
	}
	return claims, nil
}

type loginReq struct {
	Name     string `json:"name"     validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (s *Service) Login(c *gin.Context) {
	reqBody, _ := c.GetRawData()

	var req loginReq
	if err := json.Unmarshal(reqBody, &req); err != nil {
		s.returnError(c, constant.ERROR_PARAM_ERR, err.Error())
		return
	}
	if err := validateReq(req); err != nil {
		s.returnError(c, constant.ERROR_PARAM_ERR, err.Error())
		return
	}

	user, err := s.UserDao.FindUserByName(c.Request.Context(), req.Name)
	if err != nil {
		s.returnError(c, constant.ERROR_AUTH_FAIL, "invalid credentials")
		return
	}

	if !checkPassword(req.Password, user.Password) {
		loginAttempts.WithLabelValues("fail").Inc()
		s.returnError(c, constant.ERROR_AUTH_FAIL, "invalid credentials")
		return
	}

	loginAttempts.WithLabelValues("success").Inc()
	token, err := s.generateToken(user.Id)
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, "internal error")
		return
	}

	s.returnSuccess(c, gin.H{"token": token, "user_id": user.Id})
}
