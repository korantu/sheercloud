package cloud

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"runtime/debug"
	"strconv"
	"strings"
)

// --- Niceties

// must_not paicks if error happens. TODO to investigate
func must_not(err error) {
	if err == nil {
		return
	}
	debug.PrintStack()
	log.Fatal(err.Error())
}

// say writes a string to io.Writer; not too nice
func say(w io.Writer, msg string) {
	io.WriteString(w, msg)
}

// --- Request handling

// user picks a user from request if it is a valid one
func user(param map[string][]string) *User {
	login := param["login"]
	password := param["password"]
	if len(login) == 0 || len(password) == 0 {
		return nil
	}
	return GetUser(login[0], password[0])
}

// file picks out files from request
func file(param map[string][]string, user string) (paths []CloudPath, err error) {
	none := []CloudPath{}
	files, ok := param["file"]
	if !ok {
		return none, &CloudError{"File should be specified"}
	}
	for _, a_file := range files {
		if strings.Contains(a_file, "..") || strings.Contains(a_file, ":") || a_file == "" {
			return none, &CloudError{"Illegal name: " + a_file}
		}
		full_name := path.Join(user, a_file)
		paths = append(paths, CloudPath(full_name))
	}
	if len(paths) == 0 {
		return none, &CloudError{"No files are really specified"}
	}
	return
}

// user_and_file is a very common request
func user_and_file(param map[string][]string) (the_user *User, paths []CloudPath, err error) {
	if the_user = user(param); the_user == nil {
		err = &CloudError{"Failed to obtain a valid user"}
		return // err
	}

	paths, err = file(param, the_user.Login)
	if err != nil {
		return // err
	}

	return // ok
}

// --- Error handling

// CloudError handles errors in the couds
type CloudError struct {
	msg string
}

// Error is error interface
func (err *CloudError) Error() string {
	return err.msg
}

// catcher takes function which represents normal path through request.
// If main path fails, function returned by catcher handles the resulting error.
func catcher(a func(w http.ResponseWriter, r *http.Request) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := a(w, r); err != nil {
			w.Write([]byte("FAIL:" + err.Error()))
		}
	}
}

// --- API Handlers

// info provides test printout of the params of incloming request ***
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

// authorize checks that there is such user; not required for operation ***
func authorize(w http.ResponseWriter, r *http.Request) {
	u := user(r.URL.Query())
	log.Printf("checking user from %v", r.URL.Query())
	if u == nil {
		say(w, "FAIL")
		log.Printf("Authorization failed")
		return
	}
	say(w, "OK")
}

// job is API call to start a job ***
func job(w http.ResponseWriter, r *http.Request) {
	// Main response:
	_, err := ioutil.ReadAll(r.Body) // Must read body first
	if err != nil {
		say(w, "FAIL: Reading data: "+err.Error())
		return
	}

	_, file_path, err := user_and_file(r.URL.Query())
	if err != nil {
		say(w, "FAIL:"+err.Error())
		return
	}

	// Render
	result := DoJob(string(file_path[0]))

	say(w, fmt.Sprintf("OK:%s", result))
}

// progress returns the progress of a rendering job being done ***
func progress(w http.ResponseWriter, r *http.Request) {
	// Main response:
	_, err := ioutil.ReadAll(r.Body) // Must read body first
	if err != nil {
		say(w, "FAIL: Reading data: "+err.Error())
		return
	}

	the_user := user(r.URL.Query())
	if the_user == nil {
		say(w, "FAIL:"+err.Error())
		return
	}

	id, ok := r.URL.Query()["id"]

	if !ok {
		say(w, "FAIL:id parameter is required")
		return
	}

	// Render
	result, err := JobDone(JobID(id[0]))
	if err != nil {
		say(w, "FAIL:"+err.Error())
		return
	}

	if *result {
		say(w, fmt.Sprintf("OK:DONE"))
	} else {
		say(w, fmt.Sprintf("OK:PROGRESS"))
	}

}

// upload call puts a file in the cloud ***
func upload(w http.ResponseWriter, r *http.Request) {
	log.Print("Attempting upload")
	// Main response:
	incoming, err := ioutil.ReadAll(r.Body) // Must read body first
	if err != nil || r.ContentLength != int64(len(incoming)) {
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

// remove deletes a file from the cloud ***
func remove(w http.ResponseWriter, r *http.Request) {
	// Main response:
	_, err := ioutil.ReadAll(r.Body) // Must read body first
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
	theCloud.Remove(file_path[0])

	say(w, "OK")
}

// download get a file from the cloud ***
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
	w.Header().Set("Content-Length", strconv.FormatInt(int64(len(data)), 10))
	// All seem okay.
	w.Write(data)
}

// list returns a list of checksum and mtimes files from cloud ***
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
		result += fmt.Sprintf("%s\n%s\n%d\n", user_path, ids[i].MD5, ids[i].TimeStamp.Unix())
	}
	log.Print("Listing ", len(paths), " files at ", file_path[0], " prefix")
	w.Write([]byte(result))
}

/// version prints out a version ***
func version(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(Version))
}

// fail is an always-failing call, for testing relevant functions ***
func fail(w http.ResponseWriter, r *http.Request) error {
	return &CloudError{"OK"}
}

// --- Service entry points

// Serve starts all the API entry points
func Serve(port, static string) {
	http.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir(static))))
	http.HandleFunc("/error", catcher(fail))

	http.HandleFunc("/api/login", catcher(api_login))
	http.HandleFunc("/api/users", catcher(api_users))
	http.HandleFunc("/api/adduser", catcher(api_adduser))

	http.HandleFunc("/info", info)
	http.HandleFunc("/version", version)
	http.HandleFunc("/authorize", authorize)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/download", download)
	http.HandleFunc("/list", list)
	http.HandleFunc("/delete", remove)
	http.HandleFunc("/job", job)
	http.HandleFunc("/progress", progress)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Panic(err.Error())
	}
}

// --- Client for testing

// body reads all the request body for testing
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
	File     string
	FileID   string
	FileTime string
}

func ParseIdList(raw_list []byte) []FileID {
	name_id_list := strings.Split(string(raw_list), "\n")
	var result []FileID
	for n := 0; (n + 2) < len(name_id_list); n += 3 {
		result = append(result, FileID{name_id_list[n], name_id_list[n+1], name_id_list[n+2]})
	}
	return result
}

func (i Identity) List(remote string) []FileID {
	return ParseIdList(Get("list?login=" + i.Login + "&password=" + i.Password + "&file=" + remote))
}

func (i Identity) Download(remote string) []byte {
	return Get("download?login=" + i.Login + "&password=" + i.Password + "&file=" + remote)
}

func (i Identity) Delete(remote string) string {
	log.Print("Attempting to delete " + remote)
	return string(Get("delete?login=" + i.Login + "&password=" + i.Password + "&file=" + remote))
}

func (i Identity) Job(remote string) string {
	log.Print("Starting processing " + remote)
	return string(Post("job?login="+i.Login+"&password="+i.Password+"&file="+remote, []byte{}))
}

func (i Identity) Progress(id JobID) string {
	log.Print("Getting reslut of job " + id)
	return string(Get("progress?login=" + i.Login + "&password=" + i.Password + "&id=" + string(id)))
}
