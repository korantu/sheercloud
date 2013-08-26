package cloud

/*

 Problem : no time.

*/

import (
	"io"
)

type Entity struct {
	Name, FullName, Secret string
}

type EntityToken interface {
	Delete()
	Update(Entity)
}

type EntitiesAccessor interface {
	Add(Entity)
	Find(func(Entity) bool) []EntityToken
}

type Stater interface {
	Save(io.Writer)
	Load(io.Writer)
}
