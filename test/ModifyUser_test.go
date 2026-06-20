package test

import (
	"strings"
	"net/http/httptest"
	"testing"
	"user/constant"
	"user/model"
	"user/pkg/util"
)

func TestModifyUser_Success(t *testing.T) {
	svc, userDao, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "old", Address: "beijing"}

	body := `{"id":1,"name":"new_name","address":"shanghai"}`
	req := httptest.NewRequest("PUT", "/user", strings.NewReader(body))
	w := httptest.NewRecorder()

	data, err := svc.ModifyUser(w, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	user, ok := data.(*model.User)
	if !ok {
		t.Fatalf("expected *model.User, got %T", data)
	}
	if user.Name != "new_name" || user.Address != "shanghai" {
		t.Errorf("got %+v, want name=new_name address=shanghai", user)
	}
}

func TestModifyUser_MissingID(t *testing.T) {
	svc, _, _ := newTestService()
	body := `{"name":"bobby"}`
	req := httptest.NewRequest("PUT", "/user", strings.NewReader(body))
	w := httptest.NewRecorder()

	_, err := svc.ModifyUser(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	codeErr, ok := err.(*util.CodeError)
	if !ok || codeErr.Code != constant.ERROR_PARAM_ERR {
		t.Errorf("expected CodeError code %d, got %T code=%d", constant.ERROR_PARAM_ERR, err, codeErr.Code)
	}
	if codeErr.Error() != "param id not set" {
		t.Errorf("expected 'param id not set', got '%s'", codeErr.Error())
	}
}

func TestModifyUser_WithLocation(t *testing.T) {
	svc, userDao, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "bobby"}

	body := `{"id":1,"latitude":39.91,"longitude":116.41}`
	req := httptest.NewRequest("PUT", "/user", strings.NewReader(body))
	w := httptest.NewRecorder()

	data, err := svc.ModifyUser(w, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	user := data.(*model.User)
	if user.Latitude != 39.91 || user.Longitude != 116.41 {
		t.Errorf("expected lat=39.91 lng=116.41, got %+v", user)
	}
	if user.LocGeohash == "" {
		t.Error("expected geohash to be computed")
	}
}

func TestModifyUser_NotFound(t *testing.T) {
	svc, _, _ := newTestService()
	body := `{"id":999,"name":"test"}`
	req := httptest.NewRequest("PUT", "/user", strings.NewReader(body))
	w := httptest.NewRecorder()

	_, err := svc.ModifyUser(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "record not found" {
		t.Errorf("expected 'record not found', got '%v'", err)
	}
}
