package cloud

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

type waiter chan func() error

type FileStore struct {
	location               string
	files                  []string
	update_info, update_fs waiter
}

func NewFileStore(where string) (result *FileStore, err error) {
	result = &FileStore{
		where,
		[]string{},
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

	// Queues
	go do_update(result.update_fs)
	go do_update(result.update_info)

	fn := func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			result.files = append(result.files, path[len(where)+1:])
		}
		return err
	}

	err = filepath.Walk(where, fn)
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

func (store *FileStore) Add(where string, content []byte, wait_till_done bool) (err error) {
	completed := make(chan bool, 1)

	store.update_info <- func() (err error) {
		store.files = append(store.files, where)
		return
	}

	store.update_fs <- func() (err error) {
		full_name := store.OsPath(where)
		file_dir := path.Dir(full_name)
		os.MkdirAll(file_dir, os.FileMode(0777))
		err = ioutil.WriteFile(full_name, content, os.FileMode(0666))
		completed <- true
		return
	}
	if wait_till_done {
		<-completed
	}
	return
}
