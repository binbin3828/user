package service

import (
	"encoding/json"
	"fmt"
	"user/constant"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (s *Service) ForgotPassword(c *gin.Context) {
	reqBody, _ := c.GetRawData()

	var req forgotPasswordReq
	if err := json.Unmarshal(reqBody, &req); err != nil {
		s.returnError(c, constant.ERROR_PARAM_ERR, err.Error())
		return
	}
	if err := validateReq(req); err != nil {
		s.returnError(c, constant.ERROR_PARAM_ERR, err.Error())
		return
	}

	ctx := c.Request.Context()
	user, err := s.UserDao.FindUserByEmail(ctx, req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.returnSuccess(c, gin.H{"message": "if the email is registered, a reset link has been sent"})
			return
		}
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}

	resetToken, err := s.PasswordResetDao.CreateToken(ctx, user.Id)
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}

	resetLink := fmt.Sprintf("https://localhost/reset-password?token=%s", resetToken.Token)
	body := fmt.Sprintf("Hello %s,\n\nUse this link to reset your password (valid for 15 minutes):\n%s", user.Name, resetLink)
	if err := s.Mailer.Send(user.Email, "Password Reset", body); err != nil {
		s.Logger.Errorf("failed to send reset email to %s: %v", user.Email, err)
	}

	tokenInResponse := resetToken.Token
	if tokenInResponse == "" {
		tokenInResponse = "sent"
	}
	s.returnSuccess(c, gin.H{
		"message": "if the email is registered, a reset link has been sent",
		"token":   tokenInResponse,
	})
}

func (s *Service) ResetPassword(c *gin.Context) {
	reqBody, _ := c.GetRawData()

	var req resetPasswordReq
	if err := json.Unmarshal(reqBody, &req); err != nil {
		s.returnError(c, constant.ERROR_PARAM_ERR, err.Error())
		return
	}
	if err := validateReq(req); err != nil {
		s.returnError(c, constant.ERROR_PARAM_ERR, err.Error())
		return
	}

	ctx := c.Request.Context()
	t, err := s.PasswordResetDao.FindValidToken(ctx, req.Token)
	if err != nil {
		s.returnError(c, constant.ERROR_PARAM_ERR, err.Error())
		return
	}

	hashedPwd, err := hashPassword(req.NewPassword)
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, "internal error")
		return
	}

	err = s.UserDao.UpdateUser(ctx, t.UID, map[string]interface{}{"password": hashedPwd})
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}

	err = s.PasswordResetDao.MarkTokenUsed(ctx, t.Id)
	if err != nil {
		s.returnError(c, constant.ERROR_MYSQL_ERR, sanitizeErr(err).Error())
		return
	}

	s.returnSuccess(c, gin.H{"message": "password has been reset successfully"})
}
