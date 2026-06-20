package test

import (
	"net/http/httptest"
	"testing"
	"user/constant"
	"user/model"
	"user/pkg/util"
	"user/service"
)

func newTestService() (*service.Service, *MockUserDao, *MockFriendsDao) {
	log := &MockLogger{}
	userDao := NewMockUserDao()
	friendsDao := NewMockFriendsDao()
	svc := service.NewService(log, userDao, friendsDao)
	return svc, userDao, friendsDao
}

func TestGetUser_Success(t *testing.T) {
	svc, userDao, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "bobby", Address: "shenzhen"}

	req := httptest.NewRequest("GET", "/user/1", nil)
	req = chiSetURLParam(req, "uid", "1")
	w := httptest.NewRecorder()

	data, err := svc.GetUser(w, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	user, ok := data.(*model.User)
	if !ok {
		t.Fatalf("expected *model.User, got %T", data)
	}
	if user.Id != 1 || user.Name != "bobby" {
		t.Errorf("got user %+v, want id=1 name=bobby", user)
	}
}

func TestGetUser_MissingUID(t *testing.T) {
	svc, _, _ := newTestService()
	req := httptest.NewRequest("GET", "/user/", nil)
	w := httptest.NewRecorder()

	_, err := svc.GetUser(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	codeErr, ok := err.(*util.CodeError)
	if !ok || codeErr.Code != constant.ERROR_PARAM_ERR {
		t.Errorf("expected CodeError with code %d, got %T code=%d", constant.ERROR_PARAM_ERR, err, codeErr.Code)
	}
}

func TestGetUser_InvalidUID(t *testing.T) {
	svc, _, _ := newTestService()
	req := httptest.NewRequest("GET", "/user/abc", nil)
	req = chiSetURLParam(req, "uid", "abc")
	w := httptest.NewRecorder()

	_, err := svc.GetUser(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetUser_NotFound(t *testing.T) {
	svc, _, _ := newTestService()
	req := httptest.NewRequest("GET", "/user/999", nil)
	req = chiSetURLParam(req, "uid", "999")
	w := httptest.NewRecorder()

	_, err := svc.GetUser(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "record not found" {
		t.Errorf("expected 'record not found', got '%v'", err)
	}
}
