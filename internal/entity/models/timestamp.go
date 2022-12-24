package models

import "time"

type Timestamp struct {
	ID        int64     `db:"id"`
	Timestamp time.Time `json:"timestamp" db:"ts"`
	Value     int       `json:"value" db:"value"`
}
