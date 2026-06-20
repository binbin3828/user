package test

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"user/model"

	"github.com/gin-gonic/gin"
)

func TestBlockUser_Success(t *testing.T) {
	svc, userDao, _, _, blacklistDao, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob"}

	body := `{"blocked_uid":2}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/blacklist", strings.NewReader(body))
	authContextSet(c, 1)

	svc.BlockUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
	if len(blacklistDao.Entries) != 1 {
		t.Errorf("expected 1 blacklist entry, got %d", len(blacklistDao.Entries))
	}
}

func TestBlockUser_Self(t *testing.T) {
	svc, userDao, _, _, _, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}

	body := `{"blocked_uid":1}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/blacklist", strings.NewReader(body))
	authContextSet(c, 1)

	svc.BlockUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -1 || resp["msg"] != "cannot block yourself" {
		t.Errorf("got code=%v msg=%v, want code=-1 msg='cannot block yourself'", resp["code"], resp["msg"])
	}
}

func TestBlockUser_TargetNotFound(t *testing.T) {
	svc, _, _, _, _, _ := newTestService()

	body := `{"blocked_uid":999}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/blacklist", strings.NewReader(body))
	authContextSet(c, 1)

	svc.BlockUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -1 || resp["msg"] != "target user not found" {
		t.Errorf("got code=%v msg=%v, want code=-1 msg='target user not found'", resp["code"], resp["msg"])
	}
}

func TestBlockUser_AlreadyBlocked(t *testing.T) {
	svc, userDao, _, _, blacklistDao, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob"}
	blacklistDao.Entries = append(blacklistDao.Entries, &model.Blacklist{Uid: 1, BlockedUID: 2})

	body := `{"blocked_uid":2}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/blacklist", strings.NewReader(body))
	authContextSet(c, 1)

	svc.BlockUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -1 || resp["msg"] != "user already blocked" {
		t.Errorf("got code=%v msg=%v, want code=-1 msg='user already blocked'", resp["code"], resp["msg"])
	}
}

func TestBlockUser_MissingParam(t *testing.T) {
	svc, _, _, _, _, _ := newTestService()

	body := `{}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/blacklist", strings.NewReader(body))

	svc.BlockUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -1 || resp["msg"] != "param blockeduid not set" {
		t.Errorf("got code=%v msg=%v, want code=-1 msg='param blockeduid not set'", resp["code"], resp["msg"])
	}
}

func TestBlockUser_DAOError(t *testing.T) {
	svc, userDao, _, _, blacklistDao, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob"}
	blacklistDao.BlockErr = errors.New("db error")

	body := `{"blocked_uid":2}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/blacklist", strings.NewReader(body))
	authContextSet(c, 1)

	svc.BlockUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) == 0 {
		t.Fatal("expected error, got success")
	}
	if resp["msg"] != "internal error" {
		t.Errorf("expected 'internal error', got '%v'", resp["msg"])
	}
}

func TestUnblockUser_Success(t *testing.T) {
	svc, _, _, _, blacklistDao, _ := newTestService()
	blacklistDao.Entries = append(blacklistDao.Entries, &model.Blacklist{Uid: 1, BlockedUID: 2})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/blacklist/2", nil)
	c.Params = gin.Params{{Key: "uid", Value: "2"}}
	authContextSet(c, 1)

	svc.UnblockUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
	if len(blacklistDao.Entries) != 0 {
		t.Errorf("expected 0 entries after unblock, got %d", len(blacklistDao.Entries))
	}
}

func TestUnblockUser_MissingUID(t *testing.T) {
	svc, _, _, _, _, _ := newTestService()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/blacklist/", nil)

	svc.UnblockUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -1 || resp["msg"] != "param uid not set" {
		t.Errorf("got code=%v msg=%v, want code=-1 msg='param uid not set'", resp["code"], resp["msg"])
	}
}

func TestGetBlockedList_Success(t *testing.T) {
	svc, _, _, _, blacklistDao, _ := newTestService()
	blacklistDao.Entries = append(blacklistDao.Entries,
		&model.Blacklist{Uid: 1, BlockedUID: 2},
		&model.Blacklist{Uid: 1, BlockedUID: 3},
	)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/blacklist", nil)
	authContextSet(c, 1)

	svc.GetBlockedList(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d", code)
	}
	pagination := resp["pagination"].(map[string]interface{})
	if int(pagination["total"].(float64)) != 2 {
		t.Errorf("expected total=2, got %v", pagination["total"])
	}
}

func TestGetBlockedList_Empty(t *testing.T) {
	svc, _, _, _, _, _ := newTestService()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/blacklist", nil)
	authContextSet(c, 1)

	svc.GetBlockedList(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d", code)
	}
}

func TestGetUser_Blocked(t *testing.T) {
	svc, userDao, _, _, blacklistDao, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob"}
	blacklistDao.Entries = append(blacklistDao.Entries, &model.Blacklist{Uid: 1, BlockedUID: 2})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/user/1", nil)
	c.Params = gin.Params{{Key: "uid", Value: "1"}}
	authContextSet(c, 2)

	svc.GetUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != -9 || resp["msg"] != "you have been blocked by this user" {
		t.Errorf("got code=%d msg=%v, want code=-9 msg='you have been blocked by this user'", code, resp["msg"])
	}
}

func TestSendFriendRequest_BlockedByTarget(t *testing.T) {
	svc, userDao, _, _, blacklistDao, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob"}
	blacklistDao.Entries = append(blacklistDao.Entries, &model.Blacklist{Uid: 2, BlockedUID: 1})

	body := `{"to_uid":2}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/friend-requests", strings.NewReader(body))
	authContextSet(c, 1)

	svc.SendFriendRequest(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != -9 || resp["msg"] != "you have been blocked by this user" {
		t.Errorf("got code=%d msg=%v, want code=-9 msg='you have been blocked by this user'", code, resp["msg"])
	}
}
