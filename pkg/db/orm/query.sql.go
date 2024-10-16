// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package orm

import (
	"context"
	"database/sql"
)

const getWebsites = `-- name: GetWebsites :many
SELECT website_id, url, created, updated FROM website
`

func (q *Queries) GetWebsites(ctx context.Context) ([]Website, error) {
	rows, err := q.db.QueryContext(ctx, getWebsites)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Website
	for rows.Next() {
		var i Website
		if err := rows.Scan(
			&i.WebsiteID,
			&i.Url,
			&i.Created,
			&i.Updated,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const upsertTag = `-- name: UpsertTag :one
INSERT INTO tag (website_id, tag_type, name, property, value)
VALUES (?, ?, ?, ?, ?)
RETURNING tag_id, website_id, tag_type, name, property, value, created, updated
`

type UpsertTagParams struct {
	WebsiteID int64          `json:"website_id"`
	TagType   sql.NullString `json:"tag_type"`
	Name      sql.NullString `json:"name"`
	Property  sql.NullString `json:"property"`
	Value     sql.NullString `json:"value"`
}

func (q *Queries) UpsertTag(ctx context.Context, arg UpsertTagParams) (Tag, error) {
	row := q.db.QueryRowContext(ctx, upsertTag,
		arg.WebsiteID,
		arg.TagType,
		arg.Name,
		arg.Property,
		arg.Value,
	)
	var i Tag
	err := row.Scan(
		&i.TagID,
		&i.WebsiteID,
		&i.TagType,
		&i.Name,
		&i.Property,
		&i.Value,
		&i.Created,
		&i.Updated,
	)
	return i, err
}

const upsertWebsite = `-- name: UpsertWebsite :one
INSERT INTO website (url) VALUES (?) 
ON CONFLICT (url) DO UPDATE SET
 updated = CURRENT_TIMESTAMP
RETURNING website_id, url, created, updated
`

func (q *Queries) UpsertWebsite(ctx context.Context, url string) (Website, error) {
	row := q.db.QueryRowContext(ctx, upsertWebsite, url)
	var i Website
	err := row.Scan(
		&i.WebsiteID,
		&i.Url,
		&i.Created,
		&i.Updated,
	)
	return i, err
}
