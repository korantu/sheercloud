package cloud

import (
	"os"
	"path"
	"testing"
)

var guys = Users{
	User{Login: "sheer/abc", Password: "123", Name: "Me"},
	User{Login: "sheer/asd", Password: "456", Name: "Him"},
	User{Login: "sheer/important", Password: "7890", Name: "Big CEO"},
}

func TestAddUser(t *testing.T) {
	// repro
	by_login_test := make(map[string]*User)

	for _, v := range guys {
		tmp := User{Login: v.Login, Password: v.Password, Name: v.Name}
		tmp_a := v
		AddUser(tmp)
		by_login_test[v.Login] = &tmp_a
	}

	for _, v := range guys {
		if user := GetUser(v.Login, v.Password); user == nil {
			t.Logf("Adding user %s failed", v.Login)
		}
	}
}

func TestUserLoading(t *testing.T) {
	if Populate(guys); NumberOfUsers() != 3 {
		t.Logf("Users: %v", ListUsers())
		t.Error("Unexpected number of users")
	}

	user_place := path.Join(os.TempDir(), "users.json")
	os.RemoveAll(user_place)

	if err := ConfigWrite(user_place, guys); err != nil {
		t.Error(err.Error())
	}

	old_guys := Users{}
	if ResetUsers(); NumberOfUsers() != 0 {
		t.Error("Number of users should be reset")
	}

	if err := ConfigRead(user_place, &old_guys); err != nil || len(old_guys) != len(guys) {
		t.Errorf("Failed to load config [%v]", err)
	}

	if Populate(guys); NumberOfUsers() != 3 {
		t.Error("Unexpected number of users")
	}
}

func TestInitialConfig(t *testing.T) {
	the_place := path.Join(os.TempDir(), "cloud")
	Configure(the_place)
	AddUser(User{Login: "newer", Password: "secret", Name: "007"})
	SaveUsers()
	ResetUsers()
	if Configure(the_place); NumberOfUsers() != 4 {
		t.Errorf("Expected 4 users, got %d", NumberOfUsers())
	}

	if user := GetUser("newer", "secret"); user == nil {
		t.Error("user has to be present")
	}
}
