package dao

import (
	"context"
	"time"
	"user/model"
	"user/pkg/logger"
	"user/pkg/util"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

var daoBlacklistTracer = otel.Tracer("dao.blacklist")

type IBlacklistDao interface {
	Block(ctx context.Context, uid, blockedUID int) error
	Unblock(ctx context.Context, uid, blockedUID int) error
	IsBlocked(ctx context.Context, uid, targetUID int) (bool, error)
	GetBlockedList(ctx context.Context, uid int, limit, offset int) ([]*model.Blacklist, int64, error)
}

var _ IBlacklistDao = (*BlacklistDao)(nil)

type BlacklistDao struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewBlacklistDao(db *gorm.DB, log logger.Logger) *BlacklistDao {
	return &BlacklistDao{db: db, logger: log}
}

func (T *BlacklistDao) Block(ctx context.Context, uid, blockedUID int) error {
	_, span := daoBlacklistTracer.Start(ctx, "BlacklistDao.Block",
		trace.WithAttributes(
			attribute.Int("uid", uid),
			attribute.Int("blocked_uid", blockedUID),
		),
	)
	defer span.End()

	entry := model.Blacklist{
		Uid:        uid,
		BlockedUID: blockedUID,
		CreatedAt:  util.JsonTime(time.Now()),
	}
	return T.db.WithContext(ctx).Table("blacklist").Create(&entry).Error
}

func (T *BlacklistDao) Unblock(ctx context.Context, uid, blockedUID int) error {
	_, span := daoBlacklistTracer.Start(ctx, "BlacklistDao.Unblock",
		trace.WithAttributes(
			attribute.Int("uid", uid),
			attribute.Int("blocked_uid", blockedUID),
		),
	)
	defer span.End()

	return T.db.WithContext(ctx).Table("blacklist").
		Where("uid = ? AND blocked_uid = ?", uid, blockedUID).
		Delete(&model.Blacklist{}).Error
}

func (T *BlacklistDao) IsBlocked(ctx context.Context, uid, targetUID int) (bool, error) {
	_, span := daoBlacklistTracer.Start(ctx, "BlacklistDao.IsBlocked")
	defer span.End()

	var count int64
	err := T.db.WithContext(ctx).Table("blacklist").
		Where("uid = ? AND blocked_uid = ?", uid, targetUID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (T *BlacklistDao) GetBlockedList(ctx context.Context, uid int, limit, offset int) ([]*model.Blacklist, int64, error) {
	_, span := daoBlacklistTracer.Start(ctx, "BlacklistDao.GetBlockedList",
		trace.WithAttributes(attribute.Int("uid", uid)),
	)
	defer span.End()

	var list []*model.Blacklist
	var total int64

	query := T.db.WithContext(ctx).Table("blacklist").Where("uid = ?", uid)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	if list == nil {
		list = []*model.Blacklist{}
	}
	return list, total, nil
}
