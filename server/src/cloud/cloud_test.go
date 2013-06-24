package cloud

import (
	"testing"
)

func TestAddUser(t *testing.T) {
	guys := Users{
		User{Login: "abc", Password: "123", Name: "Me"},
		User{Login: "asd", Password: "456", Name: "Him"},
		User{Login: "important", Password: "7890", Name: "Big CEO"},
	}

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
