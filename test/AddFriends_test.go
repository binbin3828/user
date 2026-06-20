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

func TestAddFriend_Success(t *testing.T) {
	svc, userDao, friendsDao := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob"}

	body := `{"uid":1,"fri":2}`
	req := httptest.NewRequest("POST", "/friends", strings.NewReader(body))
	w := httptest.NewRecorder()

	data, err := svc.AddFriend(w, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, ok := data.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", data)
	}
	if int(result["uid"].(int)) != 1 || int(result["friend_id"].(int)) != 2 {
		t.Errorf("got uid=%v friend_id=%v, want uid=1 friend_id=2", result["uid"], result["friend_id"])
	}
	// verify bidirectional friendship
	found := 0
	for _, f := range friendsDao.Friends {
		if f.Uid == 1 && f.FriendID == 2 {
			found++
		}
		if f.Uid == 2 && f.FriendID == 1 {
			found++
		}
	}
	if found != 2 {
		t.Errorf("expected 2 friend records, found %d", found)
	}
}

func TestAddFriend_MissingUID(t *testing.T) {
	svc, _, _ := newTestService()
	body := `{"fri":2}`
	req := httptest.NewRequest("POST", "/friends", strings.NewReader(body))
	w := httptest.NewRecorder()

	_, err := svc.AddFriend(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	codeErr, ok := err.(*util.CodeError)
	if !ok || codeErr.Code != constant.ERROR_PARAM_ERR {
		t.Errorf("expected CodeError code %d, got %T code=%d", constant.ERROR_PARAM_ERR, err, codeErr.Code)
	}
	if codeErr.Error() != "param uid not set" {
		t.Errorf("expected 'param uid not set', got '%s'", codeErr.Error())
	}
}

func TestAddFriend_MissingFri(t *testing.T) {
	svc, userDao, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}

	body := `{"uid":1}`
	req := httptest.NewRequest("POST", "/friends", strings.NewReader(body))
	w := httptest.NewRecorder()

	_, err := svc.AddFriend(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	codeErr, ok := err.(*util.CodeError)
	if !ok || codeErr.Code != constant.ERROR_PARAM_ERR {
		t.Errorf("expected CodeError code %d, got %T code=%d", constant.ERROR_PARAM_ERR, err, codeErr.Code)
	}
	if codeErr.Error() != "param fri not set" {
		t.Errorf("expected 'param fri not set', got '%s'", codeErr.Error())
	}
}

func TestAddFriend_UserNotFound(t *testing.T) {
	svc, _, _ := newTestService()
	body := `{"uid":999,"fri":2}`
	req := httptest.NewRequest("POST", "/friends", strings.NewReader(body))
	w := httptest.NewRecorder()

	_, err := svc.AddFriend(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAddFriend_FriendNotFound(t *testing.T) {
	svc, userDao, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}

	body := `{"uid":1,"fri":999}`
	req := httptest.NewRequest("POST", "/friends", strings.NewReader(body))
	w := httptest.NewRecorder()

	_, err := svc.AddFriend(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAddFriend_DAOAddError(t *testing.T) {
	svc, userDao, friendsDao := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob"}
	friendsDao.AddFriendErr = errors.New("db error")

	body := `{"uid":1,"fri":2}`
	req := httptest.NewRequest("POST", "/friends", strings.NewReader(body))
	w := httptest.NewRecorder()

	_, err := svc.AddFriend(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "db error" {
		t.Errorf("expected 'db error', got '%v'", err)
	}
}
