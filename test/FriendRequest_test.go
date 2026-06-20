package test

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"user/model"

	"github.com/gin-gonic/gin"
)

func TestSendFriendRequest_Success(t *testing.T) {
	svc, userDao, _, friendReqDao, _, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob"}

	body := `{"to_uid":2}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/friend-requests", strings.NewReader(body))
	authContextSet(c, 1)

	svc.SendFriendRequest(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
	data := resp["data"].(map[string]interface{})
	if data["status"] != "pending" {
		t.Errorf("expected status=pending, got %v", data["status"])
	}
	if len(friendReqDao.Requests) != 1 {
		t.Errorf("expected 1 request created, got %d", len(friendReqDao.Requests))
	}
}

func TestSendFriendRequest_Self(t *testing.T) {
	svc, userDao, _, _, _, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}

	body := `{"to_uid":1}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/friend-requests", strings.NewReader(body))
	authContextSet(c, 1)

	svc.SendFriendRequest(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -1 || resp["msg"] != "cannot friend yourself" {
		t.Errorf("got code=%v msg=%v, want code=-1 msg='cannot friend yourself'", resp["code"], resp["msg"])
	}
}

func TestSendFriendRequest_MissingToUID(t *testing.T) {
	svc, _, _, _, _, _ := newTestService()

	body := `{}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/friend-requests", strings.NewReader(body))

	svc.SendFriendRequest(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -1 || resp["msg"] != "param touid not set" {
		t.Errorf("got code=%v msg=%v, want code=-1 msg='param touid not set'", resp["code"], resp["msg"])
	}
}

func TestSendFriendRequest_TargetNotFound(t *testing.T) {
	svc, userDao, _, _, _, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}

	body := `{"to_uid":999}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/friend-requests", strings.NewReader(body))
	authContextSet(c, 1)

	svc.SendFriendRequest(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) == 0 {
		t.Fatal("expected error, got success")
	}
	if resp["msg"] != "target user not found" {
		t.Errorf("expected 'target user not found', got '%v'", resp["msg"])
	}
}

func TestSendFriendRequest_AlreadyFriends(t *testing.T) {
	svc, userDao, _, friendReqDao, _, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob"}
	friendReqDao.AlreadyFriends = true

	body := `{"to_uid":2}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/friend-requests", strings.NewReader(body))
	authContextSet(c, 1)

	svc.SendFriendRequest(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != -5 || resp["msg"] != "already friends" {
		t.Errorf("got code=%d msg=%v, want code=-5 msg='already friends'", code, resp["msg"])
	}
}

func TestSendFriendRequest_ExistingPending(t *testing.T) {
	svc, userDao, _, friendReqDao, _, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob"}
	friendReqDao.HasPending = true

	body := `{"to_uid":2}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/friend-requests", strings.NewReader(body))
	authContextSet(c, 1)

	svc.SendFriendRequest(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != -6 || resp["msg"] != "friend request already sent" {
		t.Errorf("got code=%d msg=%v, want code=-6 msg='friend request already sent'", code, resp["msg"])
	}
}

func TestSendFriendRequest_ReversePendingAutoAccept(t *testing.T) {
	svc, userDao, friendsDao, friendReqDao, _, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob"}
	friendReqDao.HasPending = false
	friendReqDao.Requests = append(friendReqDao.Requests, &model.FriendRequest{
		Id:      1,
		FromUID: 2,
		ToUID:   1,
		Status:  "pending",
	})

	body := `{"to_uid":2}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/friend-requests", strings.NewReader(body))
	authContextSet(c, 1)

	svc.SendFriendRequest(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
	data := resp["data"].(map[string]interface{})
	if data["status"] != "accepted" {
		t.Errorf("expected auto-accepted, got status=%v", data["status"])
	}
	found := 0
	for _, f := range friendsDao.Friends {
		if (f.Uid == 1 && f.FriendID == 2) || (f.Uid == 2 && f.FriendID == 1) {
			found++
		}
	}
	if found != 2 {
		t.Errorf("expected 2 friend records, found %d", found)
	}
}

func TestSendFriendRequest_DAOCreateError(t *testing.T) {
	svc, userDao, _, friendReqDao, _, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob"}
	friendReqDao.CreateErr = errors.New("db error")

	body := `{"to_uid":2}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/friend-requests", strings.NewReader(body))
	authContextSet(c, 1)

	svc.SendFriendRequest(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) == 0 {
		t.Fatal("expected error, got success")
	}
	if resp["msg"] != "internal error" {
		t.Errorf("expected 'internal error', got '%v'", resp["msg"])
	}
}

func TestGetIncomingFriendRequests_Success(t *testing.T) {
	svc, userDao, _, friendReqDao, _, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob"}

	friendReqDao.CreateRequest(nil, 2, 1)
	friendReqDao.CreateRequest(nil, 3, 1)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/friend-requests/incoming", nil)
	authContextSet(c, 1)

	svc.GetIncomingFriendRequests(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d", code)
	}
	pagination := resp["pagination"].(map[string]interface{})
	if int(pagination["total"].(float64)) != 2 {
		t.Errorf("expected total=2, got %v", pagination["total"])
	}
}

func TestGetIncomingFriendRequests_Empty(t *testing.T) {
	svc, _, _, _, _, _ := newTestService()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/friend-requests/incoming", nil)
	authContextSet(c, 1)

	svc.GetIncomingFriendRequests(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
	data := resp["data"]
	if data == nil {
		t.Error("expected data to be present")
	}
}

func TestGetOutgoingFriendRequests_Success(t *testing.T) {
	svc, userDao, _, friendReqDao, _, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob"}

	friendReqDao.CreateRequest(nil, 1, 2)
	friendReqDao.CreateRequest(nil, 1, 3)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/friend-requests/outgoing", nil)
	authContextSet(c, 1)

	svc.GetOutgoingFriendRequests(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d", code)
	}
	pagination := resp["pagination"].(map[string]interface{})
	if int(pagination["total"].(float64)) != 2 {
		t.Errorf("expected total=2, got %v", pagination["total"])
	}
}

func TestAcceptFriendRequest_Success(t *testing.T) {
	svc, userDao, friendsDao, friendReqDao, _, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob"}

	req, _ := friendReqDao.CreateRequest(nil, 2, 1)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/friend-requests/1/accept", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	authContextSet(c, 1)

	svc.AcceptFriendRequest(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
	if req.Status != "accepted" {
		t.Errorf("expected status=accepted, got %v", req.Status)
	}
	friendCount := len(friendsDao.Friends)
	if friendCount != 2 {
		t.Errorf("expected 2 friend records, got %d", friendCount)
	}
}

func TestAcceptFriendRequest_NotFound(t *testing.T) {
	svc, _, _, _, _, _ := newTestService()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/friend-requests/999/accept", nil)
	c.Params = gin.Params{{Key: "id", Value: "999"}}
	authContextSet(c, 1)

	svc.AcceptFriendRequest(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != -7 || resp["msg"] != "friend request not found" {
		t.Errorf("got code=%d msg=%v, want code=-7 msg='friend request not found'", code, resp["msg"])
	}
}

func TestAcceptFriendRequest_NotRecipient(t *testing.T) {
	svc, userDao, _, friendReqDao, _, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob"}
	userDao.Users[3] = &model.User{Id: 3, Name: "charlie"}

	friendReqDao.CreateRequest(nil, 1, 2)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/friend-requests/1/accept", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	authContextSet(c, 3)

	svc.AcceptFriendRequest(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != -4 || resp["msg"] != "permission denied" {
		t.Errorf("got code=%d msg=%v, want code=-4 msg='permission denied'", code, resp["msg"])
	}
}

func TestAcceptFriendRequest_NotPending(t *testing.T) {
	svc, userDao, _, friendReqDao, _, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}

	req, _ := friendReqDao.CreateRequest(nil, 2, 1)
	req.Status = "accepted"

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/friend-requests/1/accept", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	authContextSet(c, 1)

	svc.AcceptFriendRequest(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != -8 || resp["msg"] != "request is no longer pending" {
		t.Errorf("got code=%d msg=%v, want code=-8 msg='request is no longer pending'", code, resp["msg"])
	}
}

func TestRejectFriendRequest_Success(t *testing.T) {
	svc, userDao, _, friendReqDao, _, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}

	req, _ := friendReqDao.CreateRequest(nil, 2, 1)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/friend-requests/1/reject", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	authContextSet(c, 1)

	svc.RejectFriendRequest(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	code := int(resp["code"].(float64))
	if code != 0 {
		t.Fatalf("expected code=0, got %d: %s", code, resp["msg"])
	}
	if req.Status != "rejected" {
		t.Errorf("expected status=rejected, got %v", req.Status)
	}
}

func TestRejectFriendRequest_NotFound(t *testing.T) {
	svc, _, _, _, _, _ := newTestService()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/friend-requests/999/reject", nil)
	c.Params = gin.Params{{Key: "id", Value: "999"}}
	authContextSet(c, 1)

	svc.RejectFriendRequest(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -7 || resp["msg"] != "friend request not found" {
		t.Errorf("got code=%v msg=%v, want code=-7 msg='friend request not found'", resp["code"], resp["msg"])
	}
}

func TestRejectFriendRequest_NotRecipient(t *testing.T) {
	svc, userDao, _, friendReqDao, _, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}
	userDao.Users[2] = &model.User{Id: 2, Name: "bob"}
	userDao.Users[3] = &model.User{Id: 3, Name: "charlie"}

	friendReqDao.CreateRequest(nil, 1, 2)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/friend-requests/1/reject", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	authContextSet(c, 3)

	svc.RejectFriendRequest(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -4 || resp["msg"] != "permission denied" {
		t.Errorf("got code=%v msg=%v, want code=-4 msg='permission denied'", resp["code"], resp["msg"])
	}
}

func TestRejectFriendRequest_NotPending(t *testing.T) {
	svc, userDao, _, friendReqDao, _, _ := newTestService()
	userDao.Users[1] = &model.User{Id: 1, Name: "alice"}

	req, _ := friendReqDao.CreateRequest(nil, 2, 1)
	req.Status = "rejected"

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/friend-requests/1/reject", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	authContextSet(c, 1)

	svc.RejectFriendRequest(c)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["code"].(float64)) != -8 || resp["msg"] != "request is no longer pending" {
		t.Errorf("got code=%v msg=%v, want code=-8 msg='request is no longer pending'", resp["code"], resp["msg"])
	}
}
