package models

import (
	"encoding/json"
	"io"
	"time"
)

type Group struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

// FromJSON deserializes the object from JSON string, in an io.Reader, to the given interface
func (c *Group) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(c)
}
