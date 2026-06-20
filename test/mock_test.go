package test

import (
	"context"
	"errors"
	"fmt"
	"time"
	"user/model"
	"user/pkg/logger"
	"user/pkg/util"

	"gorm.io/gorm"

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

	CreateUserErr  error
	FindUserErr    error
	FindByNameErr  error
	FindByEmailErr error
	DeleteUserErr  error
	UpdateUserErr  error
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
	user.Email = user.Email
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

func (m *MockUserDao) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	if m.FindByEmailErr != nil {
		return nil, m.FindByEmailErr
	}
	for _, u := range m.Users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
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
	if v, ok := modifyArr["password"]; ok {
		user.Password = v.(string)
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
	NearbyStrangersData []*model.RetNearbyFriendsList
	NearbyStrangersErr  error
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

func (m *MockFriendsDao) GetNearbyStrangers(ctx context.Context, uid int, subStr string, limit, offset int) ([]*model.RetNearbyFriendsList, error) {
	if m.NearbyStrangersErr != nil {
		return nil, m.NearbyStrangersErr
	}
	if m.NearbyStrangersData != nil {
		if offset >= len(m.NearbyStrangersData) {
			return []*model.RetNearbyFriendsList{}, nil
		}
		end := offset + limit
		if end > len(m.NearbyStrangersData) {
			end = len(m.NearbyStrangersData)
		}
		return m.NearbyStrangersData[offset:end], nil
	}
	return []*model.RetNearbyFriendsList{}, nil
}

func (m *MockFriendsDao) CountNearbyStrangers(ctx context.Context, uid int, subStr string) (int64, error) {
	if m.NearbyStrangersErr != nil {
		return 0, m.NearbyStrangersErr
	}
	if m.NearbyStrangersData != nil {
		return int64(len(m.NearbyStrangersData)), nil
	}
	return 0, nil
}

type MockFriendRequestDao struct {
	Requests []*model.FriendRequest
	nextID   int

	CreateErr            error
	GetIncomingErr       error
	GetOutgoingErr       error
	GetByIDErr           error
	UpdateStatusErr      error
	HasPendingErr        error
	AreAlreadyFriendsErr error

	AlreadyFriends bool
	HasPending     bool
}

func NewMockFriendRequestDao() *MockFriendRequestDao {
	return &MockFriendRequestDao{nextID: 1}
}

func (m *MockFriendRequestDao) CreateRequest(ctx context.Context, fromUID, toUID int) (*model.FriendRequest, error) {
	if m.CreateErr != nil {
		return nil, m.CreateErr
	}
	now := util.JsonTime(time.Now())
	req := &model.FriendRequest{
		Id:        m.nextID,
		FromUID:   fromUID,
		ToUID:     toUID,
		Status:    "pending",
		CreatedAt: now,
		UpdatedAt: now,
	}
	m.nextID++
	m.Requests = append(m.Requests, req)
	return req, nil
}

func (m *MockFriendRequestDao) GetIncomingRequests(ctx context.Context, toUID int, status string, limit, offset int) ([]*model.FriendRequest, int64, error) {
	if m.GetIncomingErr != nil {
		return nil, 0, m.GetIncomingErr
	}
	var list []*model.FriendRequest
	for _, r := range m.Requests {
		if r.ToUID == toUID {
			if status == "" || r.Status == status {
				list = append(list, r)
			}
		}
	}
	total := int64(len(list))
	if list == nil {
		list = []*model.FriendRequest{}
		return list, 0, nil
	}
	if offset >= len(list) {
		return []*model.FriendRequest{}, total, nil
	}
	end := offset + limit
	if end > len(list) {
		end = len(list)
	}
	return list[offset:end], total, nil
}

func (m *MockFriendRequestDao) GetOutgoingRequests(ctx context.Context, fromUID int, status string, limit, offset int) ([]*model.FriendRequest, int64, error) {
	if m.GetOutgoingErr != nil {
		return nil, 0, m.GetOutgoingErr
	}
	var list []*model.FriendRequest
	for _, r := range m.Requests {
		if r.FromUID == fromUID {
			if status == "" || r.Status == status {
				list = append(list, r)
			}
		}
	}
	total := int64(len(list))
	if list == nil {
		list = []*model.FriendRequest{}
		return list, 0, nil
	}
	if offset >= len(list) {
		return []*model.FriendRequest{}, total, nil
	}
	end := offset + limit
	if end > len(list) {
		end = len(list)
	}
	return list[offset:end], total, nil
}

func (m *MockFriendRequestDao) GetRequestByID(ctx context.Context, id int) (*model.FriendRequest, error) {
	if m.GetByIDErr != nil {
		return nil, m.GetByIDErr
	}
	for _, r := range m.Requests {
		if r.Id == id {
			return r, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockFriendRequestDao) UpdateRequestStatus(ctx context.Context, id int, status string) error {
	if m.UpdateStatusErr != nil {
		return m.UpdateStatusErr
	}
	for _, r := range m.Requests {
		if r.Id == id {
			r.Status = status
			r.UpdatedAt = util.JsonTime(time.Now())
			return nil
		}
	}
	return errors.New("record not found")
}

func (m *MockFriendRequestDao) HasPendingRequest(ctx context.Context, fromUID, toUID int) (bool, error) {
	if m.HasPendingErr != nil {
		return false, m.HasPendingErr
	}
	if m.HasPending {
		return true, nil
	}
	for _, r := range m.Requests {
		if r.FromUID == fromUID && r.ToUID == toUID && r.Status == "pending" {
			return true, nil
		}
	}
	return false, nil
}

func (m *MockFriendRequestDao) AreAlreadyFriends(ctx context.Context, uid1, uid2 int) (bool, error) {
	if m.AreAlreadyFriendsErr != nil {
		return false, m.AreAlreadyFriendsErr
	}
	return m.AlreadyFriends, nil
}

type MockBlacklistDao struct {
	Entries []*model.Blacklist

	BlockErr     error
	UnblockErr   error
	IsBlockedErr error
	GetListErr   error

	IsBlockedResult bool
}

func NewMockBlacklistDao() *MockBlacklistDao {
	return &MockBlacklistDao{}
}

func (m *MockBlacklistDao) Block(ctx context.Context, uid, blockedUID int) error {
	if m.BlockErr != nil {
		return m.BlockErr
	}
	now := util.JsonTime(time.Now())
	m.Entries = append(m.Entries, &model.Blacklist{Uid: uid, BlockedUID: blockedUID, CreatedAt: now})
	return nil
}

func (m *MockBlacklistDao) Unblock(ctx context.Context, uid, blockedUID int) error {
	if m.UnblockErr != nil {
		return m.UnblockErr
	}
	for i, e := range m.Entries {
		if e.Uid == uid && e.BlockedUID == blockedUID {
			m.Entries = append(m.Entries[:i], m.Entries[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *MockBlacklistDao) IsBlocked(ctx context.Context, uid, targetUID int) (bool, error) {
	if m.IsBlockedErr != nil {
		return false, m.IsBlockedErr
	}
	if m.IsBlockedResult {
		return true, nil
	}
	for _, e := range m.Entries {
		if e.Uid == uid && e.BlockedUID == targetUID {
			return true, nil
		}
	}
	return false, nil
}

func (m *MockBlacklistDao) GetBlockedList(ctx context.Context, uid int, limit, offset int) ([]*model.Blacklist, int64, error) {
	if m.GetListErr != nil {
		return nil, 0, m.GetListErr
	}
	var list []*model.Blacklist
	for _, e := range m.Entries {
		if e.Uid == uid {
			list = append(list, e)
		}
	}
	total := int64(len(list))
	if list == nil {
		list = []*model.Blacklist{}
		return list, 0, nil
	}
	if offset >= len(list) {
		return []*model.Blacklist{}, total, nil
	}
	end := offset + limit
	if end > len(list) {
		end = len(list)
	}
	return list[offset:end], total, nil
}

type MockPasswordResetDao struct {
	Tokens []*model.PasswordResetToken
	nextID int

	CreateErr    error
	FindValidErr error
	MarkUsedErr  error
}

func NewMockPasswordResetDao() *MockPasswordResetDao {
	return &MockPasswordResetDao{nextID: 1}
}

func (m *MockPasswordResetDao) CreateToken(ctx context.Context, uid int) (*model.PasswordResetToken, error) {
	if m.CreateErr != nil {
		return nil, m.CreateErr
	}
	now := util.JsonTime(time.Now())
	t := &model.PasswordResetToken{
		Id:        m.nextID,
		UID:       uid,
		Token:     "reset-token-" + fmt.Sprint(m.nextID),
		ExpiresAt: util.JsonTime(time.Now().Add(15 * time.Minute)),
		Used:      false,
		CreatedAt: now,
	}
	m.nextID++
	m.Tokens = append(m.Tokens, t)
	return t, nil
}

func (m *MockPasswordResetDao) FindValidToken(ctx context.Context, token string) (*model.PasswordResetToken, error) {
	if m.FindValidErr != nil {
		return nil, m.FindValidErr
	}
	for _, t := range m.Tokens {
		if t.Token == token && !t.Used && time.Time(t.ExpiresAt).After(time.Now()) {
			return t, nil
		}
	}
	return nil, errors.New("invalid or expired token")
}

func (m *MockPasswordResetDao) MarkTokenUsed(ctx context.Context, id int) error {
	if m.MarkUsedErr != nil {
		return m.MarkUsedErr
	}
	for _, t := range m.Tokens {
		if t.Id == id {
			t.Used = true
			return nil
		}
	}
	return errors.New("token not found")
}
