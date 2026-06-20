package test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"user/model"

	"github.com/gin-gonic/gin"
)

func TestDeleteUser_Success(t *testing.T) {
	svc, userDao, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "bobby"}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/user/1", nil)
	c.Params = gin.Params{{Key: "uid", Value: "1"}}

	svc.DeleteUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
	if resp["data"] != "delete succ" {
		t.Errorf("expected 'delete succ', got '%v'", resp["data"])
	}
	if _, exists := userDao.Users[1]; exists {
		t.Error("user still exists after delete")
	}
}

func TestDeleteUser_InvalidUID(t *testing.T) {
	svc, _, _ := newTestService()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/user/abc", nil)
	c.Params = gin.Params{{Key: "uid", Value: "abc"}}

	svc.DeleteUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) == 0 {
		t.Fatal("expected error, got success")
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	svc, _, _ := newTestService()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/user/999", nil)
	c.Params = gin.Params{{Key: "uid", Value: "999"}}

	svc.DeleteUser(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) == 0 {
		t.Fatal("expected error, got success")
	}
	if resp["msg"] != "record not found" {
		t.Errorf("expected 'record not found', got '%v'", resp["msg"])
	}
}
