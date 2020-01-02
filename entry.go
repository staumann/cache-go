package cache

import "time"

type entry struct {
	Data    []byte
	Expires time.Time
}

type updateStruct struct {
	Entry    entry
	ID       string
	Delete   bool
	Shutdown bool
}
