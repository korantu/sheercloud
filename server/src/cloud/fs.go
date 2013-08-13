package cloud

/*

  Need - random file.
  File set, which can list, add a file, give stats, etc.
  Then - list of entities, so far of 1 type.

  user:
    login,password,name,parent,limits
*/

/*

  Timeline: if we don't know of a file, we create meta-info for it.
  We can also verify metainfo(?)

*/

import (
	"log"
	"os"
	"path"
	"strings"
	"time"
)

type FileListError string

func (a FileListError) Error() string {
	return string(a)
}

type Location struct {
	path string
}

type FileInfo struct {
	LocalPath Location
	FullPath  string
	MD5       string
	Created   time.Time
}

// FileList holds per-file info for a user.
type FileList struct {
	Base  string
	Files map[string]*FileInfo
}

func ensure_dir(name string) error {
	if info, err := os.Stat(name); err == nil {
		if !info.IsDir() {
			return FileListError("Directory expected")
		} else {
			return nil
		}
	}
	log.Printf("Creating storage folder %s", name)
	return os.MkdirAll(name, 0777)
}

// NewFileList creates a file list for the given location.
func NewFileList(base string) (*FileList, error) {
	if err := ensure_dir(base); err != nil {
		return nil, err
	}

	fl := &FileList{
		Base:  base,
		Files: map[string]*FileInfo{},
	}
	// TODO populate if needed.
	return fl, nil
}

func (a *FileList) full_path(some Location) string {
	return path.Join(a.Base, some.path)
}

// acquireFileAs physically moves the source file into the store
func (a *FileList) acquireFileAs(local Location, file string) error {
	new_path := a.full_path(local)
	base := path.Dir(new_path)
	log.Printf("Creating path %s for file %s", base, new_path)
	if err := os.MkdirAll(base, 0777); err != nil {
		return err
	}
	if err := os.Rename(file, new_path); err != nil {
		return err
	}
	return nil // All ok
}

// Returns proper file structure;
// Involves getting info and MD5 sum.
func (a *FileList) getFileInfo(local Location) (*FileInfo, error) {
	var details os.FileInfo
	var err error

	if details, err = os.Stat(a.full_path(local)); err != nil {
		return nil, err
	}

	return &FileInfo{
		LocalPath: local,
		FullPath:  a.full_path(local),
		MD5:       "", //TBD
		Created:   details.ModTime(),
	}, nil
}

// Add adds existing file from outside of the store into local store.
func (a *FileList) Add(local Location, file string) error {
	if err := a.acquireFileAs(local, file); err != nil {
		return err
	}

	var info *FileInfo
	var err error
	if info, err = a.getFileInfo(local); err != nil {
		os.Remove(a.full_path(local))
		return err
	}
	a.Files[local.path] = info
	return nil // All is fine
}

// Delete physically removes the file and corresponding structure from the store.
func (a *FileList) Delete(local Location) error {

	exists := func(some string) bool {
		_, err := os.Stat(some)
		return err == nil
	}
	if info, ok := a.Files[local.path]; !ok && !exists(info.FullPath) {
		return nil // Does not exist, anyway
	}

	if err := os.Remove(a.full_path(local)); err != nil {
		return err
	}

	delete(a.Files, local.path)
	return nil // All is ok
}

// List returns matching Files; to return everything, use "*".
func (a *FileList) List(prefix string) []*FileInfo {
	list := []*FileInfo{}
	for local, info := range a.Files {
		if strings.HasPrefix(local, prefix) || prefix == "*" {
			list = append(list, info)
		}
	}
	return list
}
