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

var daoFriendReqTracer = otel.Tracer("dao.friend_request")

type IFriendRequestDao interface {
	CreateRequest(ctx context.Context, fromUID, toUID int) (*model.FriendRequest, error)
	GetIncomingRequests(ctx context.Context, toUID int, status string, limit, offset int) ([]*model.FriendRequest, int64, error)
	GetOutgoingRequests(ctx context.Context, fromUID int, status string, limit, offset int) ([]*model.FriendRequest, int64, error)
	GetRequestByID(ctx context.Context, id int) (*model.FriendRequest, error)
	UpdateRequestStatus(ctx context.Context, id int, status string) error
	HasPendingRequest(ctx context.Context, fromUID, toUID int) (bool, error)
	AreAlreadyFriends(ctx context.Context, uid1, uid2 int) (bool, error)
}

var _ IFriendRequestDao = (*FriendRequestDao)(nil)

type FriendRequestDao struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewFriendRequestDao(db *gorm.DB, log logger.Logger) *FriendRequestDao {
	return &FriendRequestDao{db: db, logger: log}
}

func (T *FriendRequestDao) CreateRequest(ctx context.Context, fromUID, toUID int) (*model.FriendRequest, error) {
	ctx, span := daoFriendReqTracer.Start(ctx, "FriendRequestDao.CreateRequest",
		trace.WithAttributes(
			attribute.Int("from_uid", fromUID),
			attribute.Int("to_uid", toUID),
		),
	)
	defer span.End()

	req := &model.FriendRequest{
		FromUID:   fromUID,
		ToUID:     toUID,
		Status:    "pending",
		CreatedAt: util.JsonTime(time.Now()),
		UpdatedAt: util.JsonTime(time.Now()),
	}
	err := T.db.WithContext(ctx).Table("friend_requests").Create(req).Error
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (T *FriendRequestDao) GetIncomingRequests(ctx context.Context, toUID int, status string, limit, offset int) ([]*model.FriendRequest, int64, error) {
	ctx, span := daoFriendReqTracer.Start(ctx, "FriendRequestDao.GetIncomingRequests",
		trace.WithAttributes(attribute.Int("to_uid", toUID)),
	)
	defer span.End()

	var list []*model.FriendRequest
	var total int64

	query := T.db.WithContext(ctx).Table("friend_requests").Where("to_uid = ?", toUID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	if list == nil {
		list = []*model.FriendRequest{}
	}
	return list, total, nil
}

func (T *FriendRequestDao) GetOutgoingRequests(ctx context.Context, fromUID int, status string, limit, offset int) ([]*model.FriendRequest, int64, error) {
	ctx, span := daoFriendReqTracer.Start(ctx, "FriendRequestDao.GetOutgoingRequests",
		trace.WithAttributes(attribute.Int("from_uid", fromUID)),
	)
	defer span.End()

	var list []*model.FriendRequest
	var total int64

	query := T.db.WithContext(ctx).Table("friend_requests").Where("from_uid = ?", fromUID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	if list == nil {
		list = []*model.FriendRequest{}
	}
	return list, total, nil
}

func (T *FriendRequestDao) GetRequestByID(ctx context.Context, id int) (*model.FriendRequest, error) {
	ctx, span := daoFriendReqTracer.Start(ctx, "FriendRequestDao.GetRequestByID",
		trace.WithAttributes(attribute.Int("id", id)),
	)
	defer span.End()

	var req model.FriendRequest
	err := T.db.WithContext(ctx).Table("friend_requests").Where("id = ?", id).First(&req).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (T *FriendRequestDao) UpdateRequestStatus(ctx context.Context, id int, status string) error {
	ctx, span := daoFriendReqTracer.Start(ctx, "FriendRequestDao.UpdateRequestStatus",
		trace.WithAttributes(
			attribute.Int("id", id),
			attribute.String("status", status),
		),
	)
	defer span.End()

	return T.db.WithContext(ctx).Table("friend_requests").
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": util.JsonTime(time.Now()),
		}).Error
}

func (T *FriendRequestDao) HasPendingRequest(ctx context.Context, fromUID, toUID int) (bool, error) {
	ctx, span := daoFriendReqTracer.Start(ctx, "FriendRequestDao.HasPendingRequest")
	defer span.End()

	var count int64
	err := T.db.WithContext(ctx).Table("friend_requests").
		Where("from_uid = ? AND to_uid = ? AND status = ?", fromUID, toUID, "pending").
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (T *FriendRequestDao) AreAlreadyFriends(ctx context.Context, uid1, uid2 int) (bool, error) {
	ctx, span := daoFriendReqTracer.Start(ctx, "FriendRequestDao.AreAlreadyFriends")
	defer span.End()

	var count int64
	err := T.db.WithContext(ctx).Table("friends").
		Where("(uid = ? AND fri = ?) OR (uid = ? AND fri = ?)", uid1, uid2, uid2, uid1).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
