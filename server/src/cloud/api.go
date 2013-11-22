package cloud

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

/*

  This file defines apis needed.

  1. Login.
  2. List users.
  3. Update user (added if needed).

  All goes through post.


  Messages:
   post

*/

type SessionID string

type SessionInfo struct {
	UserName string
}

func (this *SessionInfo) Same(other *SessionInfo) bool {
	if this == nil || other == nil {
		return false
	}

	return this.UserName == other.UserName
}

// Common
type ApiResource struct {
	BytesAllowed, BytesUsed int
}

type ApiUser struct {
	Name  string
	Usage ApiResource
}

type ApiStatus struct {
	Success     bool
	Description string
}

//---> KdlApiDefinitions

// /api/login
type ApiLoginRequest struct {
	Username, Password string
}

type ApiLoginReply struct {
	ApiStatus
	Session SessionID
}

type SessionStorage map[SessionID]*SessionInfo

var sessions SessionStorage = make(SessionStorage)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

func GenerateSessionID() SessionID {
	return SessionID(fmt.Sprintf("%s", rand.Int()))
}

func (a SessionID) GetInfo() *SessionInfo {
	return sessions[a]
}

func (a SessionID) PutInfo(sessionInfo *SessionInfo) {
	sessions[a] = sessionInfo
}

func (a *ApiLoginRequest) Process() *ApiLoginReply {
	var user *User
	if user = GetUser(a.Username, a.Password); user == nil {
		return &ApiLoginReply{ApiStatus{false, "Unable to login"}, ""}
	}

	sess := GenerateSessionID()
	sess.PutInfo(&SessionInfo{user.Name})

	return &ApiLoginReply{
		ApiStatus{
			true,
			"Login Successful"},
		sess}
}

// /api/users
type ApiUsersRequest struct {
	Session string
}

type ApiUsersReply struct {
	ApiStatus
	Users []ApiUser
}

// /api/adduser
type ApiAddUserRequest struct {
	Session        string
	User, Password string
}

type ApiAddUsersReply struct {
	ApiStatus
}

// apis definition ***
func api_login(w http.ResponseWriter, r *http.Request) error {
	in := json.NewDecoder(r.Body)
	defer r.Body.Close()
	login := ApiLoginRequest{}
	if err := in.Decode(&login); err != nil {
		w.Write([]byte("failed"))
		return &CloudError{"Parsing failed"}
	}
	out := json.NewEncoder(w)
	out.Encode(&ApiLoginReply{ApiStatus{true, "OK"}, "12345"})
	return nil
}

func api_users(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("users"))
	return nil
}

func api_adduser(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("adduser"))
	return nil
}
