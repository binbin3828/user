package test

import (
	"strings"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"user/model"

	"github.com/gin-gonic/gin"
)

func TestModifyUser_Success(t *testing.T) {
	svc, userDao, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "old", Address: "beijing"}

	body := `{"id":1,"name":"new_name","address":"shanghai"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/user", strings.NewReader(body))
	authContextSet(c, 1)

	svc.ModifyUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
	data := resp["data"].(map[string]interface{})
	if data["name"] != "new_name" || data["address"] != "shanghai" {
		t.Errorf("got %+v, want name=new_name address=shanghai", data)
	}
}

func TestModifyUser_MissingID(t *testing.T) {
	svc, _, _ := newTestService()
	body := `{"name":"bobby"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/user", strings.NewReader(body))

	svc.ModifyUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -1 || resp["msg"] != "param id not set" {
		t.Errorf("got code=%v msg=%v, want code=-1 msg='param id not set'", resp["code"], resp["msg"])
	}
}

func TestModifyUser_WithLocation(t *testing.T) {
	svc, userDao, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "bobby"}

	body := `{"id":1,"latitude":39.91,"longitude":116.41}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/user", strings.NewReader(body))
	authContextSet(c, 1)

	svc.ModifyUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
	data := resp["data"].(map[string]interface{})
	if data["latitude"] != 39.91 || data["longitude"] != 116.41 {
		t.Errorf("expected lat=39.91 lng=116.41, got %+v", data)
	}
	if data["loc_geohash"] == "" {
		t.Error("expected geohash to be computed")
	}
}

func TestModifyUser_NotFound(t *testing.T) {
	svc, _, _ := newTestService()
	body := `{"id":999,"name":"test"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/user", strings.NewReader(body))
	authContextSet(c, 999)

	svc.ModifyUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) == 0 {
		t.Fatal("expected error, got success")
	}
	if resp["msg"] != "internal error" {
		t.Errorf("expected 'internal error', got '%v'", resp["msg"])
	}
}
