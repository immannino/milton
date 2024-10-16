// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package orm

import (
	"database/sql"
)

type Tag struct {
	TagID     int64          `json:"tag_id"`
	WebsiteID int64          `json:"website_id"`
	TagType   sql.NullString `json:"tag_type"`
	Name      sql.NullString `json:"name"`
	Property  sql.NullString `json:"property"`
	Value     sql.NullString `json:"value"`
	Created   sql.NullTime   `json:"created"`
	Updated   sql.NullTime   `json:"updated"`
}

type Website struct {
	WebsiteID int64        `json:"website_id"`
	Url       string       `json:"url"`
	Created   sql.NullTime `json:"created"`
	Updated   sql.NullTime `json:"updated"`
}
