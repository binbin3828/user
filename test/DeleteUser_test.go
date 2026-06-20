package test

import (
	"net/http/httptest"
	"testing"
	"user/model"
)

func TestDeleteUser_Success(t *testing.T) {
	svc, userDao, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "bobby"}

	req := httptest.NewRequest("DELETE", "/user/1", nil)
	req = chiSetURLParam(req, "uid", "1")
	w := httptest.NewRecorder()

	data, err := svc.DeleteUser(w, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data != "delete succ" {
		t.Errorf("expected 'delete succ', got '%v'", data)
	}
	if _, exists := userDao.Users[1]; exists {
		t.Error("user still exists after delete")
	}
}

func TestDeleteUser_InvalidUID(t *testing.T) {
	svc, _, _ := newTestService()
	req := httptest.NewRequest("DELETE", "/user/abc", nil)
	req = chiSetURLParam(req, "uid", "abc")
	w := httptest.NewRecorder()

	_, err := svc.DeleteUser(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	svc, _, _ := newTestService()
	req := httptest.NewRequest("DELETE", "/user/999", nil)
	req = chiSetURLParam(req, "uid", "999")
	w := httptest.NewRecorder()

	_, err := svc.DeleteUser(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "record not found" {
		t.Errorf("expected 'record not found', got '%v'", err)
	}
}
