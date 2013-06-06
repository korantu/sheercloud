package cloud_test

import (
	"cloud"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

func init() {
	print("Starting test server with test users...\n")
	cloud.Populate()
	go cloud.Serve()
	time.Sleep(100 * time.Millisecond)
	print("Done.\n")
}

func TestSimplePostGet(t *testing.T) {
	post_resp := string(cloud.Post("info", []byte("ABC123")))
	get_resp := string(cloud.Get("info?a=1&a=2&b=3"))
	t.Log(post_resp)
	t.Log(get_resp)
	switch true {
	case !strings.Contains(post_resp, "[ABC123]"):
		t.Error("Echo back posted file")
	case !strings.Contains(get_resp, "b:3"):
		t.Error("Parse single parameters")
	case !strings.Contains(get_resp, "a:1|2"):
		t.Error("Parse multiple parameters")
	}
}

var good_guy, bad_guy = cloud.Identity{"abc", "123"}, cloud.Identity{"bbq", "123"}

func TestLogin(t *testing.T) {
	// Raw
	good_result := string(cloud.Get("authorize?login=important&password=7890"))
	bad_result := string(cloud.Get("authorize?login=important&password=789"))
	switch false {
	case good_result == "OK":
		t.Error("Correct user")
	case bad_result == "FAIL":
		t.Error("Wrong password")
	case good_guy.Authorize() == "OK":
		t.Error("Correct user")
	case strings.Contains(bad_guy.Authorize(), "FAIL"):
		t.Error("Wrong user")
	}
}

func TestFileTransfer(t *testing.T) {
	// Raw
	uploaded := string(cloud.Post("upload?login=important&password=7890&file=numbers.txt", []byte("12345")))
	downloaded := string(cloud.Get("download?login=important&password=7890&file=numbers.txt"))
	t.Log(uploaded)
	t.Log(downloaded)
	switch false {
	case uploaded == "OK":
		t.Error("File upload")
	case downloaded == "12345":
		t.Error("File download")
	case strings.Contains(bad_guy.Upload("scene.txt", []byte("Act I")), "FAIL"):
		t.Error("Upload by a bad guy")
	case good_guy.Upload("scene.txt", []byte("Act I")) == "OK":
		t.Error("Upload")
	case string(good_guy.Download("scene.txt")) == "Act I":
		t.Error("Download")
	}
}

func TestFileDelete(t *testing.T) {
	switch false {
	case good_guy.Upload("to_remove/scene.txt", []byte("Act I")) == "OK":
		t.Error("Upload temporary")
	case good_guy.Delete("to_remove/scene.txt") == "OK":
		t.Error("Deletion")
	case good_guy.Delete("to_remove/not_scene.txt") == "OK":
		t.Error("Deletion of non-existing file")
	case strings.Contains(string(good_guy.Download("to_remove/not_scene.txt")), "FAIL"):
		t.Error("Download non-existing")
	case strings.Contains(string(good_guy.Download("to_remove/scene.txt")), "FAIL"):
		t.Error("Download deleted")
	}
}

func TestUsers(t *testing.T) {
	if ceo := cloud.GetUser("important", "7890"); ceo.Name != "Big CEO" {
		t.Error("Fail to get the user")
	}
}

type ConcreteStuff struct {
	PieceA, PieceB int
}

type AbstractConfig struct {
	ConcreteString    string
	RealConcreteStuff ConcreteStuff
}

func TestConfig(t *testing.T) {
	place := path.Join(os.TempDir(), "abstract.config")
	t.Log(place)
	a := AbstractConfig{"Entity", ConcreteStuff{42, time.Now().Nanosecond()}}
	var b AbstractConfig
	err := cloud.ConfigWrite(place, a)
	t.Log(err)
	if err != nil {
		t.Error("Saving")
	}
	err = cloud.ConfigRead(place, &b)
	t.Log(err)
	if err != nil {
		t.Error("Loading")
	}
	if b.RealConcreteStuff.PieceB != a.RealConcreteStuff.PieceB {
		t.Error("Check loading")
	}
}

func TestFileStoreHash(t *testing.T) {
	a, b := cloud.GetID([]byte("ABC")), cloud.GetID([]byte("CBA"))
	if a == b {
		t.Error("Hash smoketest")
	}
}

func make_file(store *cloud.FileStore, name string, contents string) {
	file_location := store.OsPath(cloud.CloudPath(name))
	os.MkdirAll(path.Dir(file_location), os.FileMode(0777))
	if err := ioutil.WriteFile(file_location, []byte(contents), os.FileMode(0666)); err != nil {
		log.Panic(err.Error())
	}
}

type TestGround string

var tg = TestGround(path.Join(os.TempDir(), "cloud"))

func (test_ground TestGround) create_files(t *testing.T) {
	location := string(test_ground)
	os.RemoveAll(location)
	os.MkdirAll(location, os.FileMode(0777))
	store, err := cloud.NewFileStore(location)
	if err != nil {
		t.Error("Failed to create store:" + err.Error())
	}

	initial_files := []string{"A.txt", "alot.txt", "of.txt", "files.txt", "a/usera.txt"}
	for _, name := range initial_files {
		make_file(store, name, name+" contains nothing useful")
	}
}

func (test_ground TestGround) get_store(t *testing.T) (store *cloud.FileStore) {
	location := string(test_ground)
	var err error
	if store, err = cloud.NewFileStore(location); err != nil {
		t.Fatalf("Unable to create file store: %s", err.Error())
	}
	return store
}

func TestTypes(t *testing.T) {
	store := tg.get_store(t)
	if store.Test("hi") != 2 {
		t.Error("Parameter check")
	}
}

func TestStoreCreation(t *testing.T) {
	tg.create_files(t)
	store := tg.get_store(t)
	store.Sync()
	if store.Size() != 5 {
		t.Errorf("Incorrect store size: %i", store.Size())
	}
}

func TestStoreRemove(t *testing.T) {
	tg.create_files(t)
	store := tg.get_store(t)

	if store.Sync(); store.Size() != 5 {
		t.Errorf("Incorrect store size: %i", store.Size())
	}

	new_file := cloud.CloudPath("cool/stuff/me.txt")

	store.Add( new_file, []byte("123"))
	if store.Sync(); store.Size() != 6 {
		t.Errorf("Incorrect store size: %i", store.Size())
	}

	if file, err := os.Stat(store.OsPath( new_file)); err != nil {
		t.Errorf("File %v got to exist, but not really: %v", file, err)
	}

	store.Remove(new_file)
	if store.Sync(); store.Size() != 5 {
		t.Errorf("Incorrect store size: %i", store.Size())
	}

	if file, err := os.Stat(store.OsPath(new_file)); err == nil {
		t.Errorf("File %v is not supposed to really exist", file)
	}

}

func TestFileStore(t *testing.T) {
	tg.create_files(t)
	store := tg.get_store(t)
	store.Sync()
	size := store.Size()
	useless := "even less useful"
	useless_id := cloud.GetID([]byte(useless)) // Check later
	make_file(store, "extra.txt", useless)

	store = tg.get_store(t)
	store.Sync()
	if store.Size() != size+1 {
		t.Error("Check that new file is there")
	}
	store.Add("real.txt", []byte("Now we are talking"))
	store.Sync()
	if store.Size() != size+2 {
		t.Error("Check that another new file is there")
	}
	if store.GotID(useless_id) == nil {
		t.Error("File should be locatable")
	}
	files, ids := store.GotPrefix("a/")
	if expected, got := 1, len(files); expected != got {
		t.Errorf("Matched %d instead of %d; %v|%v", got, expected, files, ids)
	}
}
