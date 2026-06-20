package dao

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
	"user/model"
	"user/pkg/logger"
	"user/pkg/util"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

var daoResetTokenTracer = otel.Tracer("dao.password_reset")

type IPasswordResetDao interface {
	CreateToken(ctx context.Context, uid int) (*model.PasswordResetToken, error)
	FindValidToken(ctx context.Context, token string) (*model.PasswordResetToken, error)
	MarkTokenUsed(ctx context.Context, id int) error
}

var _ IPasswordResetDao = (*PasswordResetDao)(nil)

type PasswordResetDao struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewPasswordResetDao(db *gorm.DB, log logger.Logger) *PasswordResetDao {
	return &PasswordResetDao{db: db, logger: log}
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (T *PasswordResetDao) CreateToken(ctx context.Context, uid int) (*model.PasswordResetToken, error) {
	_, span := daoResetTokenTracer.Start(ctx, "PasswordResetDao.CreateToken",
		trace.WithAttributes(attribute.Int("uid", uid)),
	)
	defer span.End()

	tok, err := generateToken()
	if err != nil {
		return nil, err
	}

	t := &model.PasswordResetToken{
		UID:       uid,
		Token:     tok,
		ExpiresAt: util.JsonTime(time.Now().Add(15 * time.Minute)),
		Used:      false,
		CreatedAt: util.JsonTime(time.Now()),
	}

	err = T.db.WithContext(ctx).Table("password_reset_tokens").Create(t).Error
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (T *PasswordResetDao) FindValidToken(ctx context.Context, token string) (*model.PasswordResetToken, error) {
	_, span := daoResetTokenTracer.Start(ctx, "PasswordResetDao.FindValidToken")
	defer span.End()

	var t model.PasswordResetToken
	err := T.db.WithContext(ctx).Table("password_reset_tokens").
		Where("token = ? AND used = ? AND expires_at > ?", token, false, time.Now()).
		First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid or expired token")
		}
		return nil, err
	}
	return &t, nil
}

func (T *PasswordResetDao) MarkTokenUsed(ctx context.Context, id int) error {
	_, span := daoResetTokenTracer.Start(ctx, "PasswordResetDao.MarkTokenUsed")
	defer span.End()

	return T.db.WithContext(ctx).Table("password_reset_tokens").
		Where("id = ?", id).
		Update("used", true).Error
}
