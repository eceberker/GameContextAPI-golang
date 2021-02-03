package models

type User struct {
	ID      int64  `json:"id,omitempty"`
	Name    string `json:"display_name,omitempty"`
	Country string `json:"country,omitempty"`
	Points  int64  `json:"points,omitempty"`
}
