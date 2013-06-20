package cloud

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type waiter chan func() error

type ID string

type CloudPath string

var hasher = md5.New()

//GetID computes an ID of the given byte sequence; MD5 in this case.
func GetID(data []byte) ID {
	hasher.Reset()
	hasher.Write(data)
	sum := fmt.Sprintf("%x", hasher.Sum(nil))
	return ID(sum)
}

type FileStore struct {
	location          string
	files             map[CloudPath]ID
	queue, meta_queue waiter
}

var theCloud *FileStore

// populateFromDisk() Reads all the files in the disk in the folder and makes sure they are in the store
func (store *FileStore) populateFromDisk(location string) (err error) {
	fn := func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			var bytes []byte
			bytes, err = ioutil.ReadFile(path)
			store.NoteContent(CloudPath(path[len(location)+1:]), bytes)
		}
		return err
	}

	err = filepath.Walk(location, fn)
	return
}

// NewFileStore starts processing of queue and populates files
func NewFileStore(where string) (result *FileStore, err error) {
	result = &FileStore{
		where,
		make(map[CloudPath]ID),
		make(waiter, 1),
		make(waiter, 1)}

	do_update := func(a waiter) {
		for deed := range a {
			err := deed()
			if err != nil {
				log.Fatal(err.Error())
			}
		}
	}

	// Queue
	go do_update(result.queue)
	go do_update(result.meta_queue)

	// Get existing store
	err = result.populateFromDisk(where)
	return
}

func (store *FileStore) Size() int {
	return len(store.files)
}

func (store *FileStore) Place() string {
	return store.location
}

func (store *FileStore) OsPath(name CloudPath) string {
	return path.Join(store.location, string(name))
}

func (store *FileStore) NoteID(where CloudPath, id ID) {
	store.meta_queue <- func() (err error) {
		store.files[where] = id
		return
	}
}

func (store *FileStore) NoteContent(where CloudPath, content []byte) {
	store.NoteID(where, GetID(content))
}

func (store *FileStore) UnNote(where CloudPath){
	store.meta_queue <- func() (err error) {
		delete(store.files, where)
		return
	}
}

func (store *FileStore) KeepContent( where CloudPath, content []byte) {
	store.queue <- func() (err error) {
		full_name := store.OsPath(where)
		file_dir := path.Dir(full_name)
		os.MkdirAll(file_dir, os.FileMode(0777))
		err = ioutil.WriteFile(full_name, content, os.FileMode(0666))
		return
	}
}

func (store *FileStore) Add(where CloudPath, content []byte) (err error) {
	store.KeepContent(where, content)
	store.NoteContent(where, content)

	return
}

func (store *FileStore) GetContent(where CloudPath) (content []byte, err error) {
	full_name := store.OsPath(where)
	info, err := os.Stat(full_name)
	if err != nil {
		return
	}
	if info.IsDir() {
		return nil, NewCloudError("FAIL: Unable to download directory")
	}

	content, err = ioutil.ReadFile(full_name)
	return
}

func (store *FileStore) GotID(id ID) *CloudPath {
	for name, cloud_id := range store.files {
		if cloud_id == id {
			return &name
		}
	}
	return nil
}

func (store *FileStore) GotName(name CloudPath) * ID {
	if id, ok := store.files[name]; ok {
		return &id
	} 
	return nil
}

func (store *FileStore) GotPrefix(prefix CloudPath) (names []CloudPath, ids []ID ) {
	done := make( chan bool, 1)
	store.meta_queue <- func() (err error) {
		names, ids = []CloudPath{}, []ID{}
		for name, id := range store.files {
			if strings.HasPrefix(string(name), string(prefix)) {
				names = append(names, name)
				ids = append(ids, id)
			}
		}
		done <- true
		return
	}
	<-done
	return
}

func (store *FileStore) Link(new_name CloudPath, id ID) (err error) {
	done := make( chan bool, 1)
	store.queue <-  func () (err error) {
		old_name := store.GotID(id)
		if old_name != nil {
			err = os.Link( store.OsPath( *old_name), store.OsPath( new_name))
			if err == nil {
				store.NoteID( new_name, id)
				done <-true
			}
		} 	
		done <- false
		return
	}
	if was_linked := <-done; ! was_linked {
		err = os.ErrNotExist
	}
	return
}

func (store *FileStore) Remove(the_name CloudPath) {
	store.queue <- func() (err error) {
		file := store.GotName(the_name)
		if file != nil {
			err = os.Remove(store.OsPath(the_name))
			if err == nil {
				store.UnNote(the_name)
			}
		}
		return
	}
}

// Make sure all queues are done
func (store *FileStore) Sync() {
	done := make(chan bool, 2)

	fn := func() (err error) {
		done <- true
		return
	}

	store.queue <- fn
	store.meta_queue <- fn

	<-done
	<-done
	return
}

func (store *FileStore) Test(name CloudPath) int {
	fmt.Println("%v", name)
	return len(name)
}