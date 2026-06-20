package test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"user/model"

	"github.com/gin-gonic/gin"
)

func TestGetNearbyFriend_Success(t *testing.T) {
	svc, userDao, friendsDao := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice", LocGeohash: "wx4g0" + "xxxxx"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob", LocGeohash: "wx4g0" + "yyyyy"}
	friendsDao.AddFriend(context.Background(), 1, 2)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/nearbyfriends/1", nil)
	c.Params = gin.Params{{Key: "uid", Value: "1"}}
	authContextSet(c, 1)

	svc.GetNearbyFriend(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
	data, ok := resp["data"].([]interface{})
	if !ok || len(data) == 0 {
		t.Fatal("expected at least 1 nearby friend")
	}
	first := data[0].(map[string]interface{})
	if int(first["fri_uid"].(float64)) != 2 {
		t.Errorf("expected fri_uid=2, got %v", first["fri_uid"])
	}
	p := resp["pagination"].(map[string]interface{})
	if int(p["total"].(float64)) != 1 {
		t.Errorf("expected total=1, got %v", p["total"])
	}
}

func TestGetNearbyFriend_MissingUID(t *testing.T) {
	svc, _, _ := newTestService()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/nearbyfriends/", nil)

	svc.GetNearbyFriend(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) == 0 {
		t.Fatal("expected error, got success")
	}
}

func TestGetNearbyFriend_UserNotFound(t *testing.T) {
	svc, _, _ := newTestService()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/nearbyfriends/999", nil)
	c.Params = gin.Params{{Key: "uid", Value: "999"}}
	authContextSet(c, 999)

	svc.GetNearbyFriend(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) == 0 {
		t.Fatal("expected error, got success")
	}
}

func TestGetNearbyFriend_EmptyList(t *testing.T) {
	svc, userDao, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice", LocGeohash: "wx4g0xxxxx"}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/nearbyfriends/1", nil)
	c.Params = gin.Params{{Key: "uid", Value: "1"}}
	authContextSet(c, 1)

	svc.GetNearbyFriend(c)

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

func TestGetNearbyFriend_WithPrecision(t *testing.T) {
	svc, userDao, friendsDao := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice", LocGeohash: "wx4g0xxxxx"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob", LocGeohash: "wx4g0yyyyy"}
	friendsDao.AddFriend(context.Background(), 1, 2)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/nearbyfriends/1?precision=4", nil)
	c.Params = gin.Params{{Key: "uid", Value: "1"}}
	authContextSet(c, 1)

	svc.GetNearbyFriend(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
	data, ok := resp["data"].([]interface{})
	if !ok || len(data) == 0 {
		t.Fatal("expected at least 1 nearby friend with precision=4")
	}
	p := resp["pagination"].(map[string]interface{})
	if int(p["total"].(float64)) != 1 {
		t.Errorf("expected total=1, got %v", p["total"])
	}
}

func TestGetNearbyFriend_InvalidPrecision(t *testing.T) {
	svc, userDao, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice", LocGeohash: "wx4g0xxxxx"}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/nearbyfriends/1?precision=abc", nil)
	c.Params = gin.Params{{Key: "uid", Value: "1"}}
	authContextSet(c, 1)

	svc.GetNearbyFriend(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) == 0 {
		t.Fatal("expected error for invalid precision, got success")
	}
}
