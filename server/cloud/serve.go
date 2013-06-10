package cloud

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"strings"
)

// Create default configuration first, then provide facilities for saving/loading of such

// TODO remove
var port = "8080"

// For autocomplete tests, does not usually work.
func Ping() bool {
	return true
}

// Niceties
func must_not(err error) {
	if err == nil {
		return
	}
	debug.PrintStack()
	log.Fatal(err.Error())
}

func say(w io.Writer, msg string) {
	io.WriteString(w, msg)
}

// Handlers
func info(w http.ResponseWriter, r *http.Request) {

	list_params := func(in map[string][]string) {
		for key, values := range in {
			say(w, key+":")
			for i, value := range values {
				if i != 0 {
					say(w, "|")
				}
				say(w, value)
			}
			say(w, "\n")
		}
	}

	// Main response:
	incoming, err := ioutil.ReadAll(r.Body) // Must read body first
	must_not(err)
	in_url := r.URL.String()
	say(w, "Request to: "+in_url+" \n")
	list_params(r.URL.Query())
	say(w, "Headers:\n")
	list_params(r.Header)
	say(w, "Input:--["+string(incoming)+"]--\n")
}

func user(param map[string][]string) *User {
	login := param["login"]
	password := param["password"]
	if len(login) == 0 || len(password) == 0 {
		return nil
	}
	return GetUser(login[0], password[0])
}

type CloudError string

func (err *CloudError) Error() string {
	return string(*err)
}

var file_not_specified = CloudError("File should be specified")
var file_list_is_empty = CloudError("No files are really specified")

func NewCloudError(reason string) error {
	failure := CloudError(reason)
	return &failure

}

func file(param map[string][]string, user string) (paths []CloudPath, err error) {
	none := []CloudPath{}
	files, ok := param["file"]
	if !ok {
		return none, &file_not_specified
	}
	for _, a_file := range files {
		if strings.Contains(a_file, "..") || a_file == "" {
			return none, NewCloudError( "Illegal name: " + a_file )
		}
		full_name := path.Join(user, a_file)
		paths = append(paths, CloudPath(full_name))
	}
	if len(paths) == 0 {
		return none, &file_list_is_empty
	}
	return
}

// user_and_file is a very common request
func user_and_file(param map[string][]string) (the_user *User, paths []CloudPath, err error) {
	if the_user = user(param); the_user == nil {
		err = NewCloudError("Failed to obtain a valid user")
		return // err
	}

	paths, err = file(param, the_user.Login)
	if err != nil {
		return // err
	}

	return 	// ok
}

func authorize(w http.ResponseWriter, r *http.Request) {
	u := user(r.URL.Query())
	if u == nil {
		say(w, "FAIL")
		return
	}
	say(w, "OK")
}

// TODO take out all the file dancing outside
func upload(w http.ResponseWriter, r *http.Request) {
	// Main response:
	incoming, err := ioutil.ReadAll(r.Body) // Must read body first
	if err != nil {
		say(w, "FAIL: Reading data: "+err.Error())
		return
	}

	_, file_path, err := user_and_file(r.URL.Query())
	if err != nil {
		say(w, "FAIL:"+err.Error())
		return
	}

	// Only use first path
	theCloud.Add(file_path[0], incoming)

	say(w, "OK")
}

func download(w http.ResponseWriter, r *http.Request) {
	// Main response:
	_, err := ioutil.ReadAll(r.Body) // Must read body first
	if err != nil {
		say(w, "FAIL:"+err.Error())
		return
	}

	_, file_path, err := user_and_file(r.URL.Query())
	if err != nil {
		say(w, "FAIL:"+err.Error())
		return
	}

	data, err := theCloud.GetContent(file_path[0])

	if err != nil {
		say(w, "FAIL:"+err.Error())
		return
	}
	// All seem okay.
	w.Write(data)
}

func list(w http.ResponseWriter, r *http.Request) {
	// Main response:
	_, err := ioutil.ReadAll(r.Body) // Must read body first
	if err != nil {
		say(w, "FAIL:"+err.Error())
		return
	}

	user, file_path, err := user_and_file(r.URL.Query())
	if err != nil {
		say(w, "FAIL:"+err.Error())
		return
	}

	paths, ids := theCloud.GotPrefix(file_path[0])

	var result string

	for i, path := range paths {
		user_path := user.ConvertToUserPath(path)
		result += fmt.Sprintf("%s\n%s\n", user_path, ids[i])
	}

	w.Write([]byte(result))
}

func remove(w http.ResponseWriter, r *http.Request) {
	// Main response:
	_, err := ioutil.ReadAll(r.Body) // Must read body first
	if err != nil {
		say(w, "FAIL:"+err.Error())
		return
	}

	q := r.URL.Query()
	u := user(q)

	if u == nil {
		say(w, "FAIL")
		return
	}

	file, ok := q["file"]
	if !ok {
		say(w, "FAIL: File name is not specified")
		return
	}
	full_name := path.Join(u.Login, file[0])
	dir := path.Dir(full_name)
	if strings.Contains(dir, "..") {
		say(w, "FAIL: .. is not allowed in the path")
		return
	}
	os.Remove(full_name)
	say(w, "OK")
}

// Server
func Serve() {
	http.HandleFunc("/info", info)
	http.HandleFunc("/authorize", authorize)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/download", download)
	http.HandleFunc("/list", list)
	http.HandleFunc("/delete", remove)
	http.ListenAndServe(":"+port, nil)

}

// Client
func body(resp *http.Response) []byte {
	data, err := ioutil.ReadAll(resp.Body)
	must_not(err)
	resp.Body.Close()
	return data
}

func Get(point string) []byte {
	resp, err := http.Get("http://localhost:8080/" + point)
	must_not(err)
	return body(resp)
}

func Post(point string, to_post []byte) []byte {
	resp, err := http.Post("http://localhost:8080/"+point, "application/octet-stream", bytes.NewReader(to_post))
	must_not(err)
	return body(resp)
}

// Convenient API
type Identity struct {
	Login    string
	Password string
}

func (i Identity) Authorize() string {
	return string(Get("authorize?login=" + i.Login + "&password=" + i.Password))
}

func (i Identity) Upload(remote string, data []byte) string {
	return string(Post("upload?login="+i.Login+"&password="+i.Password+"&file="+remote, data))
}

type FileID struct {
	File string
	FileID string
}

func ParseIdList( raw_list [] byte) []FileID {
	name_id_list := strings.Split(string(raw_list), "\n")
	var result []FileID
	for  n := 0; (n+1) < len(name_id_list); n+=2 {
		result = append(result, FileID{name_id_list[n], name_id_list[n+1]})
	}
	return result
}

func (i Identity) List(remote string) []FileID {
	return ParseIdList(Get("list?login="+i.Login+"&password="+i.Password+"&file="+remote))
}

func (i Identity) Download(remote string) []byte {
	return Get("download?login=" + i.Login + "&password=" + i.Password + "&file=" + remote)
}

func (i Identity) Delete(remote string) string {
	return string(Get("delete?login=" + i.Login + "&password=" + i.Password + "&file=" + remote))
}











