package database

import "github.com/google/uuid"

type Model interface {
	TableName() string
	// Scan(any) error
	FetchQuery() string
}

type Paginate struct {
	Limit int
	Offet int
}

type FetchList struct {
	ID         uuid.UUID `db:"id"`
	TotalCount int       `db:"total_count"`
}
