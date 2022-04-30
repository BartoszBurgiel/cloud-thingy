package client

import (
	"fmt"
	"time"
)

// stats represents the stats of a successful upload
type stats struct {
	size int
	dur  time.Duration
	time time.Time
}

func (s stats) toCSV() string {
	return fmt.Sprintf("%d,%d,%s\n", s.size, s.dur, s.time.Format("2006-01-02,15:04:05"))
}
