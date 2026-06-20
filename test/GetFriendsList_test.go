package test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"user/constant"
	"user/model"

	"github.com/gin-gonic/gin"
)

func TestGetFriendsList_Success(t *testing.T) {
	svc, userDao, friendsDao := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob"}
	friendsDao.AddFriend(context.Background(), 1, 2)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/friends/1", nil)
	c.Params = gin.Params{{Key: "uid", Value: "1"}}
	authContextSet(c, 1)

	svc.GetFriendsList(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
	data := resp["data"].([]interface{})
	if len(data) != 1 {
		t.Fatalf("expected 1 friend, got %d", len(data))
	}
	first := data[0].(map[string]interface{})
	if int(first["fri_uid"].(float64)) != 2 {
		t.Errorf("expected fri_uid=2, got %v", first["fri_uid"])
	}
	p := resp["pagination"].(map[string]interface{})
	if int(p["total"].(float64)) != 1 {
		t.Errorf("expected total=1, got %v", p["total"])
	}
	if int(p["page"].(float64)) != 1 {
		t.Errorf("expected page=1, got %v", p["page"])
	}
	if int(p["page_size"].(float64)) != 20 {
		t.Errorf("expected page_size=20, got %v", p["page_size"])
	}
}

func TestGetFriendsList_MissingUID(t *testing.T) {
	svc, _, _ := newTestService()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/friends/", nil)

	svc.GetFriendsList(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) == 0 {
		t.Fatal("expected error, got success")
	}
}

func TestGetFriendsList_UserNotFound(t *testing.T) {
	svc, _, _ := newTestService()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/friends/999", nil)
	c.Params = gin.Params{{Key: "uid", Value: "999"}}
	authContextSet(c, 999)

	svc.GetFriendsList(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code == 0 {
		t.Fatal("expected error, got success")
	}
	if code != constant.ERROR_MYSQL_ERR {
		t.Errorf("expected code=%d, got %d", constant.ERROR_MYSQL_ERR, code)
	}
}

func TestGetFriendsList_EmptyList(t *testing.T) {
	svc, userDao, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/friends/1", nil)
	c.Params = gin.Params{{Key: "uid", Value: "1"}}
	authContextSet(c, 1)

	svc.GetFriendsList(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
	data, _ := resp["data"].([]interface{})
	if len(data) != 0 {
		t.Errorf("expected empty list, got %d items", len(data))
	}
	p := resp["pagination"].(map[string]interface{})
	if int(p["total"].(float64)) != 0 {
		t.Errorf("expected total=0, got %v", p["total"])
	}
}
