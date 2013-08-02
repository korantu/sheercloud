package cloud

/*
  File system should be a primary citizen.
  Data store has a root.
  Then, the rest is:
  user + "/" + local_path

  Explicitly scan users... and verify they make sense, and keep file lists.
  And keep sizes, too. Okay.

  Also, rescanning the whole fs is fast.

*/

import (
	"io"
)

type UserName string
type LocalPath string
type OsPath string
type FullPath string
type File struct {
	Path, ID string
	Size     int
}

type FS interface {
	Add(UserName, LocalPath, OsPath) error
	Open(UserName, LocalPath) (io.Reader, error)
	Del(UserName, LocalPath) error
	List(UserName, LocalPath) ([]File, error)
	Usage(UserName) int
}
