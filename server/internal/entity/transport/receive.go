package transport

import (
	"encoding/json"
	"io"
)

type SendRequest struct {
	Timestamp int64 `json:"timestamp"`
	Value     int   `json:"value"`
}

func (x *SendRequest) Decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(x)
}
