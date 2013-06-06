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

type File struct {
	Name CloudPath
	Id   ID
}

type FileStore struct {
	location string
	files    []File
	queue    waiter
}

// populateFromDisk() Reads all the files in the disk in the folder and makes sure they are in the store
func (store *FileStore) populateFromDisk(location string) (err error) {
	fn := func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			var bytes []byte
			bytes, err = ioutil.ReadFile(path)
			store.Add(CloudPath(path[len(location)+1:]), bytes)
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
		[]File{},
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

func (store *FileStore) OsPath(name string) string {
	return path.Join(store.location, name)
}

func (store *FileStore) Add(where string, content []byte) (err error) {
	store.queue <- func() (err error) {
		store.files = append(store.files, file{where, GetID(content)})
		full_name := store.OsPath(where)
		file_dir := path.Dir(full_name)
		os.MkdirAll(file_dir, os.FileMode(0777))
		err = ioutil.WriteFile(full_name, content, os.FileMode(0666))
		return
	}
	return
}

func (store *FileStore) Got( id ID) bool {
	return false
}

func (store *FileStore) Sync() {
	done := make(chan bool, 1)
	store.queue <- func() (err error) {
		done <- true
		return
	}
	<-done
	return
}
