package cloud

// User state
type User struct {
	Name string
	Login string
	Password string
}

// Storage
var by_login = make(map[string] * User)

// Add a user
func AddUser( a * User) * User {
	if a.Login == "" || a.Password == "" {
		return nil
	}

	if _, ok := by_login[a.Login]; ok {
		return nil 
	}

	// Copy to avoid update to the map being used
	new_storage := make(map[string] * User)
	for k,v := range by_login {
		new_storage[k] = v
	}

	new_storage[a.Login] = a;
	by_login = new_storage;

	return a;
}

// Get a user
func GetUser( login string, password string) * User {
	u, ok := by_login[ login]
	if ! ok || u.Password != password {
		return nil
	}

	return u;
}

// Test users
func Populate() {
	AddUser( &User{Login: "abc", Password: "123", Name: "Me"} )
	AddUser( &User{Login: "asd", Password: "456", Name: "Him"} )
	AddUser( &User{Login: "important", Password: "7890", Name: "Big CEO"} )
}



