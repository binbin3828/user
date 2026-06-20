package test

import (
	"net/http/httptest"
	"testing"
	"user/model"

	"github.com/gorilla/mux"
)

func TestGetFriendsList_Success(t *testing.T) {
	svc, userDao, friendsDao := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob"}
	// create bidirectional friendship
	friendsDao.AddFriend(1, 2)

	req := httptest.NewRequest("GET", "/friends/1", nil)
	req = mux.SetURLVars(req, map[string]string{"uid": "1"})
	w := httptest.NewRecorder()

	data, err := svc.GetFriendsList(w, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result, ok := data.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", data)
	}
	if result["uid"].(int) != 1 {
		t.Errorf("expected uid=1, got %v", result["uid"])
	}
	list, ok := result["list"].([]*model.RetListFriends)
	if !ok {
		t.Fatalf("expected []*RetListFriends, got %T", result["list"])
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 friend, got %d", len(list))
	}
	if list[0].FriUid != 2 {
		t.Errorf("expected fri_uid=2, got %d", list[0].FriUid)
	}
}

func TestGetFriendsList_MissingUID(t *testing.T) {
	svc, _, _ := newTestService()
	req := httptest.NewRequest("GET", "/friends/", nil)
	w := httptest.NewRecorder()

	_, err := svc.GetFriendsList(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetFriendsList_UserNotFound(t *testing.T) {
	svc, _, _ := newTestService()
	req := httptest.NewRequest("GET", "/friends/999", nil)
	req = mux.SetURLVars(req, map[string]string{"uid": "999"})
	w := httptest.NewRecorder()

	_, err := svc.GetFriendsList(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetFriendsList_EmptyList(t *testing.T) {
	svc, userDao, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}

	req := httptest.NewRequest("GET", "/friends/1", nil)
	req = mux.SetURLVars(req, map[string]string{"uid": "1"})
	w := httptest.NewRecorder()

	data, err := svc.GetFriendsList(w, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := data.(map[string]interface{})
	list := result["list"].([]*model.RetListFriends)
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d items", len(list))
	}
}
