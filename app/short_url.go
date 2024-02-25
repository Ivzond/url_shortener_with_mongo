package app

import "time"

type ShortURL struct {
	ID       string
	URL      string
	ExpireAt time.Time
}
