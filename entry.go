package cache

import "time"

type entry struct {
	Data    []byte
	Expires time.Time
}
