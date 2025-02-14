package helpers

import (
	"time"

	"github.com/labstack/gommon/log"
)

type Time interface {
	Now() time.Time
}

type UTC struct{}

func (t UTC) Now() time.Time {
	return time.Now().UTC()
}

func TimeFrom(s string) time.Time {
	v, err := time.Parse(time.RFC3339, s)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func TimePtrFrom(s string) *time.Time {
	v := TimeFrom(s)
	return &v
}
