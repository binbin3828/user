package test

import (
	"context"
	"errors"
	"fmt"
	"time"
	"user/model"
	"user/pkg/logger"
	"user/pkg/util"
)

type MockLogger struct{}

func (m *MockLogger) Info(args ...interface{})                  {}
func (m *MockLogger) Infof(format string, args ...interface{})  {}
func (m *MockLogger) Error(args ...interface{})                 {}
func (m *MockLogger) Errorf(format string, args ...interface{}) {}
func (m *MockLogger) Debug(args ...interface{})                 {}
func (m *MockLogger) Debugf(format string, args ...interface{}) {}
func (m *MockLogger) Warn(args ...interface{})                  {}
func (m *MockLogger) Warnf(format string, args ...interface{})  {}

var _ logger.Logger = (*MockLogger)(nil)

type MockUserDao struct {
	Users  map[int]*model.User
	nextID int

	CreateUserErr error
	FindUserErr   error
	DeleteUserErr error
	UpdateUserErr error
}

func NewMockUserDao() *MockUserDao {
	return &MockUserDao{
		Users:  make(map[int]*model.User),
		nextID: 1,
	}
}

func (m *MockUserDao) CreateUser(ctx context.Context, user *model.User) error {
	if m.CreateUserErr != nil {
		return m.CreateUserErr
	}
	user.Id = m.nextID
	m.nextID++
	user.CreateAt = util.JsonTime(time.Now())
	m.Users[user.Id] = user
	return nil
}

func (m *MockUserDao) FindUser(ctx context.Context, id int) (*model.User, error) {
	if m.FindUserErr != nil {
		return nil, m.FindUserErr
	}
	user, ok := m.Users[id]
	if !ok {
		return nil, errors.New("record not found")
	}
	return user, nil
}

func (m *MockUserDao) DeleteUser(ctx context.Context, uid int) error {
	if m.DeleteUserErr != nil {
		return m.DeleteUserErr
	}
	if _, ok := m.Users[uid]; !ok {
		return errors.New("record not found")
	}
	delete(m.Users, uid)
	return nil
}

func (m *MockUserDao) UpdateUser(ctx context.Context, uid int, modifyArr map[string]interface{}) error {
	if m.UpdateUserErr != nil {
		return m.UpdateUserErr
	}
	user, ok := m.Users[uid]
	if !ok {
		return errors.New("record not found")
	}
	if v, ok := modifyArr["name"]; ok {
		user.Name = v.(string)
	}
	if v, ok := modifyArr["dob"]; ok {
		user.Dob = v.(string)
	}
	if v, ok := modifyArr["address"]; ok {
		user.Address = v.(string)
	}
	if v, ok := modifyArr["description"]; ok {
		user.Description = v.(string)
	}
	if v, ok := modifyArr["latitude"]; ok {
		user.Latitude = v.(float64)
	}
	if v, ok := modifyArr["longitude"]; ok {
		user.Longitude = v.(float64)
	}
	if v, ok := modifyArr["loc_geohash"]; ok {
		user.LocGeohash = v.(string)
	}
	return nil
}

type MockFriendsDao struct {
	Friends []*model.Friends

	AddFriendErr        error
	GetFriendsListErr   error
	GetNearbyFriendErr  error
}

func NewMockFriendsDao() *MockFriendsDao {
	return &MockFriendsDao{}
}

func (m *MockFriendsDao) AddFriend(ctx context.Context, uid, friendID int) error {
	if m.AddFriendErr != nil {
		return m.AddFriendErr
	}
	now := util.JsonTime(time.Now())
	m.Friends = append(m.Friends, &model.Friends{Uid: uid, FriendID: friendID, CreateTime: now})
	m.Friends = append(m.Friends, &model.Friends{Uid: friendID, FriendID: uid, CreateTime: now})
	return nil
}

func (m *MockFriendsDao) GetFriendsList(ctx context.Context, uid int) ([]*model.RetListFriends, error) {
	if m.GetFriendsListErr != nil {
		return nil, m.GetFriendsListErr
	}
	now := util.JsonTime(time.Now())
	var list []*model.RetListFriends
	for _, f := range m.Friends {
		if f.Uid == uid {
			list = append(list, &model.RetListFriends{
				FriUid:     f.FriendID,
				FriName:    fmt.Sprintf("user_%d", f.FriendID),
				CreateTime: now,
			})
		}
	}
	if list == nil {
		// return empty slice, not nil, so JSON serializes as [] not null
		return []*model.RetListFriends{}, nil
	}
	return list, nil
}

func (m *MockFriendsDao) GetNearbyFriend(ctx context.Context, uid int, subStr string) ([]*model.RetNearbyFriendsList, error) {
	if m.GetNearbyFriendErr != nil {
		return nil, m.GetNearbyFriendErr
	}
	now := util.JsonTime(time.Now())
	var list []*model.RetNearbyFriendsList
	for _, f := range m.Friends {
		if f.Uid == uid {
			list = append(list, &model.RetNearbyFriendsList{
				FriUid:     f.FriendID,
				FriName:    fmt.Sprintf("user_%d", f.FriendID),
				CreateTime: now,
				Latitude:   39.91,
				Longitude:  116.41,
				LocGeohash: subStr + "xxx",
			})
		}
	}
	if list == nil {
		return []*model.RetNearbyFriendsList{}, nil
	}
	return list, nil
}
