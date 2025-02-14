package helpers

import (
	"time"
)

type Time interface {
	Now() time.Time
}

type UTC struct{}

func (t UTC) Now() time.Time {
	return time.Now().UTC()
}
