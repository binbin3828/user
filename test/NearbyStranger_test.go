package test

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"user/model"

	"github.com/gin-gonic/gin"
)

func TestGetNearbyStranger_Success(t *testing.T) {
	svc, userDao, friendsDao, _, _, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice", LocGeohash: "wx4g119d"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob", LocGeohash: "wx4g119e"}
	userDao.Users[3] = &model.User{Id: 3, Name: "carol", LocGeohash: "wx4g119f"}

	friendsDao.NearbyStrangersData = []*model.RetNearbyFriendsList{
		{FriUid: 2, FriName: "bob", Latitude: 39.91, Longitude: 116.41, LocGeohash: "wx4g119e"},
		{FriUid: 3, FriName: "carol", Latitude: 39.92, Longitude: 116.42, LocGeohash: "wx4g119f"},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/nearby-users/1", nil)
	c.Params = gin.Params{{Key: "uid", Value: "1"}}
	authContextSet(c, 1)

	svc.GetNearbyStranger(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
	pagination := resp["pagination"].(map[string]interface{})
	if int(pagination["total"].(float64)) != 2 {
		t.Errorf("expected total=2, got %v", pagination["total"])
	}
}

func TestGetNearbyStranger_Empty(t *testing.T) {
	svc, userDao, _, _, _, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice", LocGeohash: "wx4g119d"}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/nearby-users/1", nil)
	c.Params = gin.Params{{Key: "uid", Value: "1"}}
	authContextSet(c, 1)

	svc.GetNearbyStranger(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
}

func TestGetNearbyStranger_NoGeohash(t *testing.T) {
	svc, userDao, _, _, _, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice", LocGeohash: ""}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/nearby-users/1", nil)
	c.Params = gin.Params{{Key: "uid", Value: "1"}}
	authContextSet(c, 1)

	svc.GetNearbyStranger(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
}

func TestGetNearbyStranger_MissingUID(t *testing.T) {
	svc, _, _, _, _, _ := newTestService()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/nearby-users/", nil)

	svc.GetNearbyStranger(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -1 || resp["msg"] != "param uid not set" {
		t.Errorf("got code=%v msg=%v, want code=-1 msg='param uid not set'", resp["code"], resp["msg"])
	}
}

func TestGetNearbyStranger_InvalidUID(t *testing.T) {
	svc, _, _, _, _, _ := newTestService()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/nearby-users/abc", nil)
	c.Params = gin.Params{{Key: "uid", Value: "abc"}}

	svc.GetNearbyStranger(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) == 0 {
		t.Fatal("expected error, got success")
	}
}

func TestGetNearbyStranger_UserNotFound(t *testing.T) {
	svc, _, _, _, _, _ := newTestService()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/nearby-users/999", nil)
	c.Params = gin.Params{{Key: "uid", Value: "999"}}
	authContextSet(c, 999)

	svc.GetNearbyStranger(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) == 0 {
		t.Fatal("expected error, got success")
	}
}

func TestGetNearbyStranger_InvalidPrecision(t *testing.T) {
	svc, userDao, _, _, _, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice", LocGeohash: "wx4g119d"}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/nearby-users/1?precision=13", nil)
	c.Params = gin.Params{{Key: "uid", Value: "1"}}
	authContextSet(c, 1)

	svc.GetNearbyStranger(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -1 || resp["msg"] != "param precision must be 1-12" {
		t.Errorf("got code=%v msg=%v, want code=-1 msg='param precision must be 1-12'", resp["code"], resp["msg"])
	}
}

func TestGetNearbyStranger_CustomPrecision(t *testing.T) {
	svc, userDao, friendsDao, _, _, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice", LocGeohash: "wx4g119d"}

	friendsDao.NearbyStrangersData = []*model.RetNearbyFriendsList{
		{FriUid: 2, FriName: "bob", Latitude: 39.91, Longitude: 116.41, LocGeohash: "wx4g119e"},
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/nearby-users/1?precision=5&page=1&page_size=10", nil)
	c.Params = gin.Params{{Key: "uid", Value: "1"}}
	authContextSet(c, 1)

	svc.GetNearbyStranger(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
	pagination := resp["pagination"].(map[string]interface{})
	if int(pagination["total"].(float64)) != 1 {
		t.Errorf("expected total=1, got %v", pagination["total"])
	}
}

func TestGetNearbyStranger_DAOError(t *testing.T) {
	svc, userDao, friendsDao, _, _, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice", LocGeohash: "wx4g119d"}
	friendsDao.NearbyStrangersErr = errors.New("db error")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/nearby-users/1", nil)
	c.Params = gin.Params{{Key: "uid", Value: "1"}}
	authContextSet(c, 1)

	svc.GetNearbyStranger(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) == 0 {
		t.Fatal("expected error, got success")
	}
	if resp["msg"] != "internal error" {
		t.Errorf("expected 'internal error', got '%v'", resp["msg"])
	}
}
