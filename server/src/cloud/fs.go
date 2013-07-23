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
	_ "os"
	_ "path"
	_ "strings"
)
