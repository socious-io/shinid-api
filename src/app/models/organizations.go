package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Organization struct {
	ID          uuid.UUID `db:"id" json:"id"`
	DID         string    `db:"did" json:"did"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	LogoID      uuid.UUID `db:"logo_id" json:"logo_id"`
	IsVerified  bool      `db:"is_verified" json:"is_verified"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

func (*Organization) TableName() string {
	return "organizations"
}

func (*Organization) FetchQuery() string {
	return "organizations/fetch"
}

func (o *Organization) Scan(rows *sqlx.Rows) error {
	return rows.StructScan(o)
}
