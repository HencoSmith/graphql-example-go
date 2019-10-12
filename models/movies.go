package models

import "time"

// Movie - Basic information about a movie
type Movie struct {
	ID          string     `json:"id"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	UsersID     string     `json:"users_id"`
	Name        string     `json:"name"`
	ReleaseYear int64      `json:"release_year"`
	Description string     `json:"description,omitempty"`
	Rating      float64    `json:"rating"`
	ReviewCount int64      `json:"review_count"`
}
