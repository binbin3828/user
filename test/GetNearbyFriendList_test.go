package test

import (
	"context"
	"net/http/httptest"
	"testing"
	"user/model"
)

func TestGetNearbyFriend_Success(t *testing.T) {
	svc, userDao, friendsDao := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice", LocGeohash: "wx4g0" + "xxxxx"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob", LocGeohash: "wx4g0" + "yyyyy"}
	friendsDao.AddFriend(context.Background(), 1, 2)

	req := httptest.NewRequest("GET", "/nearbyfriends/1", nil)
	req = chiSetURLParam(req, "uid", "1")
	w := httptest.NewRecorder()

	data, err := svc.GetNearbyFriend(w, req)
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
	list, ok := result["list"].([]*model.RetNearbyFriendsList)
	if !ok {
		t.Fatalf("expected []*RetNearbyFriendsList, got %T", result["list"])
	}
	if len(list) == 0 {
		t.Fatal("expected at least 1 nearby friend")
	}
	if list[0].FriUid != 2 {
		t.Errorf("expected fri_uid=2, got %d", list[0].FriUid)
	}
}

func TestGetNearbyFriend_MissingUID(t *testing.T) {
	svc, _, _ := newTestService()
	req := httptest.NewRequest("GET", "/nearbyfriends/", nil)
	w := httptest.NewRecorder()

	_, err := svc.GetNearbyFriend(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetNearbyFriend_UserNotFound(t *testing.T) {
	svc, _, _ := newTestService()
	req := httptest.NewRequest("GET", "/nearbyfriends/999", nil)
	req = chiSetURLParam(req, "uid", "999")
	w := httptest.NewRecorder()

	_, err := svc.GetNearbyFriend(w, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetNearbyFriend_EmptyList(t *testing.T) {
	svc, userDao, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice", LocGeohash: "wx4g0xxxxx"}

	req := httptest.NewRequest("GET", "/nearbyfriends/1", nil)
	req = chiSetURLParam(req, "uid", "1")
	w := httptest.NewRecorder()

	data, err := svc.GetNearbyFriend(w, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := data.(map[string]interface{})
	list := result["list"].([]*model.RetNearbyFriendsList)
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d items", len(list))
	}
}
