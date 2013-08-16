package cloud

/*

  This file defines apis needed.

  1. Login.
  2. List users.
  3. Update user (added if needed).

  All goes through post.


  Messages:
   post

*/

// Common
type ApiResource struct {
	BytesAllowed, BytesUsed int
}

type ApiUser struct {
	Name  string
	Usage ApiResource
}

type ApiStatus struct {
	Success bool
}

//---> KdlApiDefinitions

// /ui/login
type ApiLoginRequest struct {
	Username, Password string
}

type ApiLoginReply struct {
	ApiStatus
	Session string
}

// /ui/users
type ApiUsersRequest struct {
	Session string
}

type ApiUsersReply struct {
	ApiStatus
	Users []ApiUser
}

// /ui/adduser
type ApiAddUserRequest struct {
	Session        string
	User, Password string
}

type ApiAddUsersReply struct {
	ApiStatus
}
