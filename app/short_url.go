package app

import "time"

type ShortURL struct {
	ID       string     `bson:"_id"`
	URL      string     `bson:"url"`
	ExpireAt *time.Time `bson:"expireAt,omitempty"`
}
