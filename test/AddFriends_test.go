package test

import (
	"errors"
	"strings"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"user/model"

	"github.com/gin-gonic/gin"
)

func TestAddFriend_Success(t *testing.T) {
	svc, userDao, friendsDao := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob"}

	body := `{"uid":1,"fri":2}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/friends", strings.NewReader(body))
	authContextSet(c, 1)

	svc.AddFriend(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
	data := resp["data"].(map[string]interface{})
	if int(data["uid"].(float64)) != 1 || int(data["friend_id"].(float64)) != 2 {
		t.Errorf("got uid=%v friend_id=%v, want uid=1 friend_id=2", data["uid"], data["friend_id"])
	}
	found := 0
	for _, f := range friendsDao.Friends {
		if f.Uid == 1 && f.FriendID == 2 {
			found++
		}
		if f.Uid == 2 && f.FriendID == 1 {
			found++
		}
	}
	if found != 2 {
		t.Errorf("expected 2 friend records, found %d", found)
	}
}

func TestAddFriend_MissingUID(t *testing.T) {
	svc, _, _ := newTestService()
	body := `{"fri":2}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/friends", strings.NewReader(body))

	svc.AddFriend(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -1 || resp["msg"] != "param uid not set" {
		t.Errorf("got code=%v msg=%v, want code=-1 msg='param uid not set'", resp["code"], resp["msg"])
	}
}

func TestAddFriend_MissingFri(t *testing.T) {
	svc, userDao, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}

	body := `{"uid":1}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/friends", strings.NewReader(body))

	svc.AddFriend(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -1 || resp["msg"] != "param fri not set" {
		t.Errorf("got code=%v msg=%v, want code=-1 msg='param fri not set'", resp["code"], resp["msg"])
	}
}

func TestAddFriend_UserNotFound(t *testing.T) {
	svc, _, _ := newTestService()
	body := `{"uid":999,"fri":2}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/friends", strings.NewReader(body))
	authContextSet(c, 999)

	svc.AddFriend(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) == 0 {
		t.Fatal("expected error, got success")
	}
	// should not reveal which user doesn't exist
	if resp["msg"] != "invalid friend request" {
		t.Errorf("expected 'invalid friend request', got '%v'", resp["msg"])
	}
}

func TestAddFriend_FriendNotFound(t *testing.T) {
	svc, userDao, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}

	body := `{"uid":1,"fri":999}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/friends", strings.NewReader(body))
	authContextSet(c, 1)

	svc.AddFriend(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) == 0 {
		t.Fatal("expected error, got success")
	}
	if resp["msg"] != "invalid friend request" {
		t.Errorf("expected 'invalid friend request', got '%v'", resp["msg"])
	}
}

func TestAddFriend_DAOAddError(t *testing.T) {
	svc, userDao, friendsDao := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob"}
	friendsDao.AddFriendErr = errors.New("db error")

	body := `{"uid":1,"fri":2}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/friends", strings.NewReader(body))
	authContextSet(c, 1)

	svc.AddFriend(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) == 0 {
		t.Fatal("expected error, got success")
	}
	// DAO error sanitized to generic message
	if resp["msg"] != "internal error" {
		t.Errorf("expected 'internal error', got '%v'", resp["msg"])
	}
}
