package cloud

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

const ui_dir = "./ui"

func init() {
	log.Print("Starting test server ...\n")
	// Real server should probably configured away from the default location.
	// make sure the place is new.
	Configure("/tmp/cloud_testing/" + fmt.Sprint(time.Now().Unix()))
	go Serve("8080", ui_dir)
	time.Sleep(100 * time.Millisecond)
	log.Print("Test server is up and running.\n")
}

var content_counter = time.Now().Nanosecond()

// some_content returns a string and a byte representation of some random bytes.
func some_content() (as_string string, as_bytes []byte) {
	content_counter++
	as_string = fmt.Sprintf("<<<%d>>>", content_counter)
	as_bytes = []byte(as_string)
	return
}

func TestSomeContent(t *testing.T) {
	a_string, a_bytes := some_content()
	b_string, b_bytes := some_content()
	if bytes.Equal(a_bytes, b_bytes) {
		t.Errorf("Must be different: <%v> <%v>", a_bytes, b_bytes)
	}
	if !bytes.Equal(b_bytes, []byte(b_string)) || !bytes.Equal(a_bytes, []byte(a_string)) {
		t.Error("Bytes and string must be same")
	}

}

// TestFail verifies that failing API fails.
func TestFail(t *testing.T) {
	if response := string(Get("error")); !strings.Contains(response, "FAIL") {
		t.Log(response)
		t.Error("FAIL expected.")
	}
}

func TestSimplePostGet(t *testing.T) {
	a_string, a_bytes := some_content()
	post_resp := string(Post("info", a_bytes))
	get_resp := string(Get("info?a=1&a=2&b=3"))
	t.Log(post_resp)
	t.Log(get_resp)
	switch true {
	case !strings.Contains(post_resp, "["+a_string+"]"):
		t.Error("Echo back posted file")
	case !strings.Contains(get_resp, "b:3"):
		t.Error("Parse single parameters")
	case !strings.Contains(get_resp, "a:1|2"):
		t.Error("Parse multiple parameters")
	}
}

var good_guy, bad_guy = Identity{"sheer/abc", "123"}, Identity{"sheer/bbq", "123"}

func TestLogin(t *testing.T) {
	// Raw
	good_result := string(Get("authorize?login=sheer/important&password=7890"))
	bad_result := string(Get("authorize?login=sheer/important&password=789"))
	Log("Should be good: " + good_result)
	Log("Should be bad: " + bad_result)

	switch {
	case good_result != "OK":
		t.Error("Correct user was not authorized")
	case !strings.Contains(bad_result, "FAIL"):
		t.Error("Wrong password")
	case good_guy.Authorize() != "OK":
		t.Error("Correct user was not authorized")
	case !strings.Contains(bad_guy.Authorize(), "FAIL"):
		t.Error("Wrong user not failed")
	}
}

func TestUploadDownload(t *testing.T) {
	_, test_bytes := some_content()
	t.Log(string(Post("upload?login=sheer/important&password=7890&file=numbers.txt", test_bytes)))
	t.Log(string(Post("upload?login=sheer/abc&password=123&file=scene.txt", test_bytes)))
	t.Log(string(Get("download?login=sheer/important&password=7890&file=numbers.txt")))
	t.Log(string(Get("download?login=sheer/abc&password=123&file=scene.txt")))
	t.Log(string(Post("upload?login=sheer/abc&password=123&file=scene.txt", test_bytes)))
	//	t.Fail()
}

func TestFileTransfer(t *testing.T) {
	// Get some data for test
	test_string, test_bytes := some_content()
	// Raw
	uploaded := string(Post("upload?login=sheer/important&password=7890&file=numbers.txt", test_bytes))
	downloaded := string(Get("download?login=sheer/important&password=7890&file=numbers.txt"))
	bad_guy_uploaded := bad_guy.Upload("scene.txt", []byte("Act I"))
	t.Log("Normal upload " + uploaded)
	t.Log("Normal download " + downloaded)
	t.Log("Bad guy attempted upload " + bad_guy_uploaded)
	switch {
	case uploaded != "OK":
		t.Error("File upload")
	case downloaded != test_string:
		t.Error("File download")
	case !strings.Contains(bad_guy_uploaded, "FAIL"):
		t.Error("Upload by a bad guy")
	case good_guy.Upload("scene.txt", test_bytes) != "OK":
		t.Error("Upload")
	case string(good_guy.Download("scene.txt")) != test_string:
		t.Error("Download")
	}
}

func TestFileList(t *testing.T) {
	_, test_bytes := some_content()
	names := []string{"a", "b", "c"}
	checked := make(map[string]bool)
	for _, name := range names {
		a_full_name := "to/list/" + name + ".txt"
		checked[a_full_name] = false
		uploaded := string(Post("upload?login=sheer/important&password=7890&file="+a_full_name, test_bytes))
		t.Log(uploaded)
	}
	got := Get("list?login=sheer/important&password=7890&file=to/list")
	t.Log(string(got))
	list := ParseIdList(got)
	for _, to_check := range list {
		if was_checked, ok := checked[to_check.File]; !ok || was_checked == true {
			t.Errorf("List %v does not match reference %v", list, checked)
		}
		checked[to_check.File] = true
	}
}

func TestFileListInterface(t *testing.T) {
	content := func() (some_bytes []byte) {
		_, some_bytes = some_content()
		return
	}

	type checkpoint struct {
		data        []byte
		was_checked bool
	}

	var files_to_send = map[string]checkpoint{"to_list/scene.txt": {content(), false},
		"to_list/scene88.txt":               {content(), false},
		"to_list/scene_more/cool.stuff.txt": {content(), false}}
	for name, checking := range files_to_send {
		if good_guy.Upload(name, checking.data) != "OK" {
			t.Errorf("Failed to upload %s", name)
		}
	}
	files_and_ids := good_guy.List("to_list")
	if len(files_and_ids) != len(files_to_send) {
		t.Errorf("Supposed to be of the same length: %v and %v", files_and_ids, files_to_send)
	}
	for _, the_file := range files_and_ids {
		checkdata, ok := files_to_send[the_file.File]
		if !ok {
			t.Errorf("File %s was not really sent: %v", the_file.File, files_to_send)
		}
		if checkdata.was_checked {
			t.Errorf("Duplicate file %s detected", the_file.File)
		}
		if MD5(checkdata.data) != the_file.FileID {
			t.Errorf("Checksum mistmatch for %#v", the_file)
		}
		checkdata.was_checked = true
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
	if ceo := GetUser("sheer/important", "7890"); ceo.Name != "Big CEO" {
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
	err := ConfigWrite(place, a)
	t.Log(err)
	if err != nil {
		t.Error("Saving")
	}
	err = ConfigRead(place, &b)
	t.Log(err)
	if err != nil {
		t.Error("Loading")
	}
	if b.RealConcreteStuff.PieceB != a.RealConcreteStuff.PieceB {
		t.Error("Check loading")
	}
}

func TestFileStoreHash(t *testing.T) {
	a, b := MD5([]byte("ABC")), MD5([]byte("CBA"))
	if a == b {
		t.Error("Hash smoketest")
	}
}

func make_file(store *FileStore, name string, contents string) {
	file_location := store.OsPath(CloudPath(name))
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
	store, err := NewFileStore(location)
	if err != nil {
		t.Error("Failed to create store:" + err.Error())
	}

	initial_files := []string{"A.txt", "alot.txt", "of.txt", "files.txt", "a/usera.txt"}
	for _, name := range initial_files {
		make_file(store, name, name+" contains nothing useful")
	}
}

func (test_ground TestGround) get_store(t *testing.T) (store *FileStore) {
	location := string(test_ground)
	var err error
	if store, err = NewFileStore(location); err != nil {
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

	new_file := CloudPath("cool/stuff/me.txt")
	new_content := []byte("123")

	store.Add(new_file, new_content)
	if store.Sync(); store.Size() != 6 {
		t.Errorf("Incorrect store size: %i", store.Size())
	}

	if content, err := store.GetContent(new_file); !bytes.Equal(content, new_content) || err != nil {
		t.Errorf("File content does not match; Othewise, bad (%v) things happened.", err)
	}

	if file, err := os.Stat(store.OsPath(new_file)); err != nil {
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

func NotTestFileStore(t *testing.T) {
	tg.create_files(t)
	store := tg.get_store(t)
	store.Sync()
	size := store.Size()
	useless := "even less useful"
	useless_id := GetID([]byte(useless), time.Now()) // Check later
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

func TestJobs(t *testing.T) {
	started := string(Post("job?login=sheer/important&password=7890&file=scene.txt", []byte{}))
	if strings.Contains(started, "FAIL") {
		t.Error(started)
	}
	id := JobID(started[3:])

	done_so_far := string(Get("progress?login=sheer/important&password=7890&id=" + string(id)))
	if strings.Contains(started, "FAIL") {
		t.Error(started)
	}

	if done_so_far[3:] != "PROGRESS" {
		t.Error(done_so_far)
	}

	time.Sleep(time.Second + 100*time.Millisecond)
	done_so_far = string(Get("progress?login=sheer/important&password=7890&id=" + string(id)))
	if done_so_far[3:] != "DONE" {
		t.Error(done_so_far)
	}

}

// TODO: can be simpler
func TestJobsSimple(t *testing.T) {
	started := good_guy.Job("scene.txt")
	if strings.Contains(started, "FAIL") {
		t.Error(started)
	}
	id := JobID(started[3:])

	done_so_far := good_guy.Progress(id)
	if strings.Contains(started, "FAIL") {
		t.Error(started)
	}

	if done_so_far[3:] != "PROGRESS" {
		t.Error(done_so_far)
	}

	time.Sleep(time.Second + 100*time.Millisecond)

	if done_so_far = good_guy.Progress(id); done_so_far[3:] != "DONE" {
		t.Error(done_so_far)
	}

}
