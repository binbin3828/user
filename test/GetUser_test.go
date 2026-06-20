package test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"user/constant"
	"user/model"
	"user/pkg/mailer"
	"user/pkg/ratelimit"
	"user/service"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestService() (*service.Service, *MockUserDao, *MockFriendsDao, *MockFriendRequestDao, *MockBlacklistDao, *MockPasswordResetDao) {
	log := &MockLogger{}
	userDao := NewMockUserDao()
	friendsDao := NewMockFriendsDao()
	friendReqDao := NewMockFriendRequestDao()
	blacklistDao := NewMockBlacklistDao()
	passwordResetDao := NewMockPasswordResetDao()
	mailerInst := &mailer.DevMailer{}
	rl := ratelimit.NewMemoryLimiter(10, time.Minute)
	svc := service.NewService(log, userDao, friendsDao, friendReqDao, blacklistDao, passwordResetDao, mailerInst, rl)
	return svc, userDao, friendsDao, friendReqDao, blacklistDao, passwordResetDao
}

func TestGetUser_Success(t *testing.T) {
	svc, userDao, _, _, _, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "bobby", Address: "shenzhen"}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/user/1", nil)
	c.Params = gin.Params{{Key: "uid", Value: "1"}}
	authContextSet(c, 1)

	svc.GetUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
	data := resp["data"].(map[string]interface{})
	if int(data["id"].(float64)) != 1 || data["name"] != "bobby" {
		t.Errorf("got %+v, want id=1 name=bobby", data)
	}
}

func TestGetUser_MissingUID(t *testing.T) {
	svc, _, _, _, _, _ := newTestService()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/user/", nil)

	svc.GetUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -1 {
		t.Errorf("expected code=-1, got code=%v", resp["code"])
	}
}

func TestGetUser_InvalidUID(t *testing.T) {
	svc, _, _, _, _, _ := newTestService()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/user/abc", nil)
	c.Params = gin.Params{{Key: "uid", Value: "abc"}}

	svc.GetUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) == 0 {
		t.Fatal("expected error, got success")
	}
}

func TestGetUser_NotFound(t *testing.T) {
	svc, _, _, _, _, _ := newTestService()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/user/999", nil)
	c.Params = gin.Params{{Key: "uid", Value: "999"}}

	svc.GetUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != constant.ERROR_PERMISSION_DENIED {
		t.Errorf("expected code=%d, got %d", constant.ERROR_PERMISSION_DENIED, code)
	}
	if resp["msg"] != "user not found" {
		t.Errorf("expected 'user not found', got '%v'", resp["msg"])
	}
}
