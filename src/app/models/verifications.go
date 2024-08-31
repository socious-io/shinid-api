package models

import (
	"context"
	"net/url"
	"shin/src/database"
	"shin/src/wallet"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
)

type Verification struct {
	ID            uuid.UUID      `db:"id" json:"id"`
	Name          string         `db:"name" json:"name"`
	Description   *string        `db:"description" json:"description"`
	SchemaID      uuid.UUID      `db:"schema_id" json:"schema_id"`
	Schema        *Schema        `db:"-" json:"schema"`
	UserID        uuid.UUID      `db:"user_id" json:"user_id"`
	User          *User          `db:"-" json:"user"`
	ConnectionID  *string        `db:"connection_id" json:"connection_id"`
	ConnectionURL *string        `db:"connection_url" json:"connection_url"`
	Body          types.JSONText `db:"body" json:"body"`

	Status CredentialStatusType `db:"status" json:"status"`

	ConnectionAt *time.Time `db:"connection_at" json:"connection_at"`
	VerifiedAt   *time.Time `db:"verified_at" json:"verified_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`

	UserJson   types.JSONText `db:"user" json:"-"`
	SchemaJson types.JSONText `db:"schema" json:"-"`
}

func (Verification) TableName() string {
	return "credential_verifications"
}

func (Verification) FetchQuery() string {
	return "credentials/fetch_verification"
}

func (v *Verification) Create(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"credentials/create_verification",
		v.Name, v.Description, v.UserID, v.SchemaID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(v); err != nil {
			return err
		}
	}
	return database.Fetch(v, v.ID)
}

func (v *Verification) Update(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"credentials/update_verification",
		v.ID, v.Name, v.Description, v.UserID, v.SchemaID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(v); err != nil {
			return err
		}
	}
	return database.Fetch(v, v.ID)
}

func (v *Verification) NewConnection(ctx context.Context, domain, callback string) error {
	conn, err := wallet.CreateConnection(callback)
	if err != nil {
		return err
	}
	connectURL, _ := url.JoinPath(domain, conn.ShortID)
	rows, err := database.Query(
		ctx,
		"credentials/update_connection_verification",
		v.ID, conn.ID, connectURL,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(v); err != nil {
			return err
		}
	}
	return database.Fetch(v, v.ID)
}

func (v *Verification) Delete(ctx context.Context) error {
	rows, err := database.Query(ctx, "credentials/delete_verification", v.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func GetVerification(id uuid.UUID) (*Verification, error) {
	v := new(Verification)

	if err := database.Fetch(v, id); err != nil {
		return nil, err
	}
	return v, nil
}

func GetVerificationByConnection(connectionID uuid.UUID) (*Verification, error) {
	v := new(Verification)

	if err := database.Get(v, "credentials/verification_by_connection", connectionID); err != nil {
		return nil, err
	}
	return v, nil
}

func GetVerifications(userId uuid.UUID, p database.Paginate) ([]Verification, int, error) {
	var (
		verifications = []Verification{}
		fetchList     []database.FetchList
		ids           []interface{}
	)

	if err := database.QuerySelect("credentials/get_verifications", &fetchList, userId, p.Limit, p.Offet); err != nil {
		return nil, 0, err
	}

	if len(fetchList) < 1 {
		return verifications, 0, nil
	}

	for _, f := range fetchList {
		ids = append(ids, f.ID)
	}

	if err := database.Fetch(&verifications, ids...); err != nil {
		return nil, 0, err
	}
	return verifications, fetchList[0].TotalCount, nil
}
