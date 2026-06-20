package test

import (
	"errors"
	"strings"
	"net/http/httptest"
	"testing"
	"user/constant"
	"user/model"
	"user/pkg/util"
)

func TestCreateUser_Success(t *testing.T) {
	svc, userDao, _ := newTestService()
	body := `{"name":"bobby","dob":"1990-01-01","address":"shenzhen","description":"coder"}`
	req := httptest.NewRequest("POST", "/user", strings.NewReader(body))
	w := httptest.NewRecorder()

	data, err := svc.CreateUser(w, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	user, ok := data.(*model.User)
	if !ok {
		t.Fatalf("expected *model.User, got %T", data)
	}
	if user.Name != "bobby" || user.Address != "shenzhen" {
		t.Errorf("got %+v, want name=bobby address=shenzhen", user)
	}
	if user.Id == 0 {
		t.Error("expected non-zero id")
	}
	if _, exists := userDao.Users[user.Id]; !exists {
		t.Error("user not stored in DAO")
	}
}

func TestCreateUser_MissingName(t *testing.T) {
	svc, _, _ := newTestService()
	body := `{"dob":"1990-01-01"}`
	req := httptest.NewRequest("POST", "/user", strings.NewReader(body))
	w := httptest.NewRecorder()

	_, err := svc.CreateUser(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	codeErr, ok := err.(*util.CodeError)
	if !ok || codeErr.Code != constant.ERROR_PARAM_ERR {
		t.Errorf("expected CodeError code %d, got %T code=%d", constant.ERROR_PARAM_ERR, err, codeErr.Code)
	}
	if codeErr.Error() != "param name not set" {
		t.Errorf("expected 'param name not set', got '%s'", codeErr.Error())
	}
}

func TestCreateUser_WithLocation(t *testing.T) {
	svc, _, _ := newTestService()
	body := `{"name":"bobby","latitude":39.910934,"longitude":116.413385}`
	req := httptest.NewRequest("POST", "/user", strings.NewReader(body))
	w := httptest.NewRecorder()

	data, err := svc.CreateUser(w, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	user, ok := data.(*model.User)
	if !ok {
		t.Fatalf("expected *model.User, got %T", data)
	}
	if user.Latitude != 39.910934 || user.Longitude != 116.413385 {
		t.Errorf("expected lat=39.910934 lng=116.413385, got %+v", user)
	}
	if user.LocGeohash == "" {
		t.Error("expected geohash to be computed")
	}
}

func TestCreateUser_NegativeLocation(t *testing.T) {
	svc, _, _ := newTestService()
	body := `{"name":"bobby","latitude":-10,"longitude":116.413385}`
	req := httptest.NewRequest("POST", "/user", strings.NewReader(body))
	w := httptest.NewRecorder()

	data, err := svc.CreateUser(w, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	user := data.(*model.User)
	if user.LocGeohash != "" {
		t.Error("expected empty geohash for negative latitude")
	}
}

func TestCreateUser_DAOCreateError(t *testing.T) {
	svc, userDao, _ := newTestService()
	userDao.CreateUserErr = errors.New("db error")
	body := `{"name":"bobby"}`
	req := httptest.NewRequest("POST", "/user", strings.NewReader(body))
	w := httptest.NewRecorder()

	_, err := svc.CreateUser(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "db error" {
		t.Errorf("expected 'db error', got '%v'", err)
	}
}
