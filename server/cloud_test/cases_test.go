package cloud_test

import (
	"cloud"
	"strings"
	"testing"
	"time"
	"os"
	"path"
	"io/ioutil"
)

func init() {
	print("Starting test server with test users...\n")
	cloud.Populate()
	go cloud.Serve()
	time.Sleep(100 * time.Millisecond)
	print("Done.\n")
}

func Must(t *testing.T, be_true bool, reason string) {
	if !be_true {
		t.Log(reason)
		t.Fail()
	}
}

func TestSimplePostGet(t *testing.T) {
	post_resp := string(cloud.Post("info", []byte("ABC123")))
	get_resp := string(cloud.Get("info?a=1&a=2&b=3"))
	t.Log(post_resp)
	t.Log(get_resp)
	Must(t, strings.Contains(post_resp, "[ABC123]"), "Echo back posted file")
	Must(t, strings.Contains(get_resp, "b:3"), "Parse single parameters")
	Must(t, strings.Contains(get_resp, "a:1|2"), "Parse multiple parameters")
}

var good_guy, bad_guy = cloud.Identity{"abc", "123"}, cloud.Identity{"bbq", "123"}

func TestLogin(t *testing.T) {
	// Raw
	good_result := string(cloud.Get("authorize?login=important&password=7890"))
	bad_result := string(cloud.Get("authorize?login=important&password=789"))
	Must(t, good_result == "OK", "Correct user")
	Must(t, bad_result == "FAIL", "Wrong password")
	// Nicer
	Must(t, good_guy.Authorize() == "OK", "Correct user")
	Must(t, strings.Contains(bad_guy.Authorize(), "FAIL"), "Wrong user")
}

func TestFileTransfer(t *testing.T) {
	// Raw
	uploaded := string(cloud.Post("upload?login=important&password=7890&file=numbers.txt", []byte("12345")))
	downloaded := string(cloud.Get("download?login=important&password=7890&file=numbers.txt"))
	t.Log(uploaded)
	t.Log(downloaded)
	Must(t, uploaded == "OK", "File upload")
	Must(t, downloaded == "12345", "File download")
	// Nicer
	Must(t, strings.Contains(bad_guy.Upload("scene.txt", []byte("Act I")), "FAIL"), "Upload by a bad guy")
	// Flow
	Must(t, good_guy.Upload("scene.txt", []byte("Act I")) == "OK", "Upload")
	Must(t, string(good_guy.Download("scene.txt")) == "Act I", "Download")
}

func TestFileDelete(t *testing.T) {
	Must(t, good_guy.Upload("to_remove/scene.txt", []byte("Act I")) == "OK", "Upload temporary")
	Must(t, good_guy.Delete("to_remove/scene.txt") == "OK", "Deletion")
	Must(t, good_guy.Delete("to_remove/not_scene.txt") == "OK", "Deletion of non-existing file")
	Must(t, strings.Contains(string(good_guy.Download("to_remove/not_scene.txt")), "FAIL"), "Download non-existing")
	Must(t, strings.Contains(string(good_guy.Download("to_remove/scene.txt")), "FAIL"), "Download deleted")
}

func TestUsers(t *testing.T) {
	ceo := cloud.GetUser("important", "7890")
	Must(t, ceo.Name == "Big CEO", "Get the user")
}

type ConcreteStuff struct {
	PieceA, PieceB int 
}

type AbstractConfig struct {
	ConcreteString string
	RealConcreteStuff ConcreteStuff
}

func TestConfig(t *testing.T) {
	place := path.Join( os.TempDir(), "abstract.config")
	t.Log( place)
	a := AbstractConfig{ "Entity", ConcreteStuff {42, time.Now().Nanosecond() } }
	var b AbstractConfig
	err := cloud.ConfigWrite( place, a)
	t.Log( err)
	Must( t, err == nil, "Saving" )
	err = cloud.ConfigRead( place, &b)
	t.Log( err)
	Must( t, err == nil, "Loading" )
	Must( t, b.RealConcreteStuff.PieceB == a.RealConcreteStuff.PieceB, "Check loading") 
}

func TestFileStore(t * testing.T){
	location := path.Join( os.TempDir(), "cloud")
	
	os.Mkdir( location, os.FileMode( 0777))

	make_file := func( name string, contents string) {
		file_location := path.Join( location, name)
		ioutil.WriteFile( file_location, []byte(contents), os.FileMode( 0666))
	}

	initial_files := []string {"A.txt", "lot.txt", "of.txt", "files.txt"}
	for _, name := range initial_files {
		make_file( name, name + " contains nothing useful")
	}

	store, err := cloud.NewFileStore( location)
	Must( t, err == nil, "Create store")
	t.Log( err)
	size := store.Size() 
	Must( t, size == len( initial_files), "Number of entries")
	make_file( "extra.txt", "even less useful")
	store, err = cloud.NewFileStore( location)
	Must(t, err == nil, "Re-check")
	Must(t, store.Size() == size + 1, "Check that new file is there")
}



















