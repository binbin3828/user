package test

import (
	"context"
	"errors"
	"fmt"
	"time"
	"user/model"
	"user/pkg/logger"
	"user/pkg/util"

	"github.com/gin-gonic/gin"
)

func authContextSet(c *gin.Context, userID int) {
	c.Set("user_id", userID)
}

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

	CreateUserErr   error
	FindUserErr     error
	FindByNameErr   error
	DeleteUserErr   error
	UpdateUserErr   error
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

func (m *MockUserDao) FindUserByName(ctx context.Context, name string) (*model.User, error) {
	if m.FindByNameErr != nil {
		return nil, m.FindByNameErr
	}
	for _, u := range m.Users {
		if u.Name == name {
			return u, nil
		}
	}
	return nil, errors.New("record not found")
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

func (m *MockFriendsDao) GetFriendsList(ctx context.Context, uid int, limit, offset int) ([]*model.RetListFriends, error) {
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
		return []*model.RetListFriends{}, nil
	}
	return list, nil
}

func (m *MockFriendsDao) CountFriendsList(ctx context.Context, uid int) (int64, error) {
	if m.GetFriendsListErr != nil {
		return 0, m.GetFriendsListErr
	}
	var count int64
	for _, f := range m.Friends {
		if f.Uid == uid {
			count++
		}
	}
	return count, nil
}

func (m *MockFriendsDao) GetNearbyFriend(ctx context.Context, uid int, subStr string, limit, offset int) ([]*model.RetNearbyFriendsList, error) {
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

func (m *MockFriendsDao) CountNearbyFriend(ctx context.Context, uid int, subStr string) (int64, error) {
	if m.GetNearbyFriendErr != nil {
		return 0, m.GetNearbyFriendErr
	}
	var count int64
	for _, f := range m.Friends {
		if f.Uid == uid {
			count++
		}
	}
	return count, nil
}
