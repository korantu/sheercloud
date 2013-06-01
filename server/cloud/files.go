package cloud

import (
  "path/filepath"
  "os"
)

type FileStore struct {
	Files []string
}

func NewFileStore( where string) ( result * FileStore, err error) {
	result = &FileStore{ []string{} }

	fn := func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			result.Files = append(result.Files, path)
		}
		return err
	}

	err = filepath.Walk( where, fn)
	return
}

func ( store * FileStore ) Size() int {
	return len( store.Files )
}




















