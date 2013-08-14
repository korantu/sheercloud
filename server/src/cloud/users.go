package cloud

import (
	"strings"
)

/*

*/

// User state
type User struct {
	Name     string
	Login    string
	Password string
}

type Users []User

// Storage
var by_login = make(map[string]User)

func NumberOfUsers() int {
	return len(by_login)
}

func ListUsers() []string {
	result := []string{}
	for login, _ := range by_login {
		result = append(result, login)
	}
	return result
}

// ResetUsers removes all existing users
func ResetUsers() {
	by_login = make(map[string]User)
}

func (a *User) CloudPathPrefix() CloudPath {
	return CloudPath(a.Login + "/")
}

func (a *User) ConvertToUserPath(full_path CloudPath) string {
	return strings.Replace(string(full_path), string(a.CloudPathPrefix()), "", 1)
}

// AddUser adds a user to the userlist
func AddUser(a User) {
	// Copy to avoid update to the map being used
	new_storage := make(map[string]User)
	for k, v := range by_login {
		new_storage[k] = v

	}
	new_storage[a.Login] = a
	by_login = new_storage
}

// Get a user
func GetUser(login, password string) *User {
	DumpUsers()
	u, ok := by_login[login]
	if !ok || u.Password != password {
		return nil
	}
	return &u
}

var test_guys = Users{
	User{Login: "sheer/abc", Password: "123", Name: "Me"},
	User{Login: "sheer/asd", Password: "456", Name: "Him"},
	User{Login: "sheer/important", Password: "7890", Name: "Big CEO"},
}

func DumpUsers() (result Users) {
	result = Users{}
	for _, user := range by_login {
		result = append(result, user)
	}
	return
}

func Populate(ppl Users) {
	for _, user := range ppl {
		AddUser(user)
		DumpUsers()
	}
}
