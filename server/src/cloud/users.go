package cloud

import (
	"log"
	"strings"
)

// User state
type User struct {
	Name     string
	Login    string
	Password string
}

type Users []User

// Storage
var by_login = make(map[string]User)

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

	log.Printf("Worked")
	return &u
}

var test_guys = Users{
	User{Login: "abc", Password: "123", Name: "Me"},
	User{Login: "asd", Password: "456", Name: "Him"},
	User{Login: "important", Password: "7890", Name: "Big CEO"},
}

func DumpUsers() {
	for k, v := range by_login {
		log.Printf("User [%s] record [%v]", k, v)
	}
}

func Populate(ppl Users) {
	for _, user := range ppl {
		log.Printf("Adding [%s], [%s]", user.Login, user.Password)
		AddUser(user)
		DumpUsers()
	}
}
