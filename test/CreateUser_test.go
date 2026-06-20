package test

import (
	"errors"
	"strings"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateUser_Success(t *testing.T) {
	svc, userDao, _ := newTestService()
	body := `{"name":"bobby","dob":"1990-01-01","address":"shenzhen","description":"coder"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/user", strings.NewReader(body))

	svc.CreateUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
	data := resp["data"].(map[string]interface{})
	if data["name"] != "bobby" || data["address"] != "shenzhen" {
		t.Errorf("got %+v, want name=bobby address=shenzhen", data)
	}
	id := int(data["id"].(float64))
	if id == 0 {
		t.Error("expected non-zero id")
	}
	if _, exists := userDao.Users[id]; !exists {
		t.Error("user not stored in DAO")
	}
}

func TestCreateUser_MissingName(t *testing.T) {
	svc, _, _ := newTestService()
	body := `{"dob":"1990-01-01"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/user", strings.NewReader(body))

	svc.CreateUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -1 || resp["msg"] != "param name not set" {
		t.Errorf("got code=%v msg=%v, want code=-1 msg='param name not set'", resp["code"], resp["msg"])
	}
}

func TestCreateUser_WithLocation(t *testing.T) {
	svc, _, _ := newTestService()
	body := `{"name":"bobby","latitude":39.910934,"longitude":116.413385}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/user", strings.NewReader(body))

	svc.CreateUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
	data := resp["data"].(map[string]interface{})
	if data["latitude"] != 39.910934 || data["longitude"] != 116.413385 {
		t.Errorf("expected lat=39.910934 lng=116.413385, got %+v", data)
	}
	if data["loc_geohash"] == "" {
		t.Error("expected geohash to be computed")
	}
}

func TestCreateUser_NegativeLocation(t *testing.T) {
	svc, _, _ := newTestService()
	body := `{"name":"bobby","latitude":-10,"longitude":116.413385}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/user", strings.NewReader(body))

	svc.CreateUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
	data := resp["data"].(map[string]interface{})
	if data["loc_geohash"] != "" {
		t.Error("expected empty geohash for negative latitude")
	}
}

func TestCreateUser_DAOCreateError(t *testing.T) {
	svc, userDao, _ := newTestService()
	userDao.CreateUserErr = errors.New("db error")
	body := `{"name":"bobby"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/user", strings.NewReader(body))

	svc.CreateUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) == 0 {
		t.Fatal("expected error, got success")
	}
	if resp["msg"] != "db error" {
		t.Errorf("expected 'db error', got '%v'", resp["msg"])
	}
}
