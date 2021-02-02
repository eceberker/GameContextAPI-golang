package models

import (
	"encoding/json"
)

type User struct {
	ID      int64  `json:"id"`
	Name    string `json:"display_name"`
	Country string `json:"country"`
	Points  int64  `json:"points"`
}

// MarshalBinary -
func (e *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

// UnmarshalBinary -
func (e *User) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &e); err != nil {
		return err
	}

	return nil
}
