package models

import (
	"context"
	"shin/src/database"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OTP struct {
	ID         uuid.UUID `db:"id" json:"id"`
	UserID     uuid.UUID `db:"user_id" json:"user_id"`
	Code       int       `db:"code" json:"code"`
	IsVerified bool      `db:"is_verified" json:"is_verified"`
	ExpiresAt  time.Time `db:"expired_at" json:"expired_at"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

func (*OTP) TableName() string {
	return "otps"
}

func (*OTP) FetchQuery() string {
	return "otps/fetch"
}

func (o *OTP) Scan(rows *sqlx.Rows) error {
	return rows.StructScan(o)
}

func (o *OTP) Create(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"otp/create",
		o.UserID, o.Code,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := o.Scan(rows); err != nil {
			return err
		}
	}
	return nil
}

func (o *OTP) Verify(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"otp/verify",
		o.UserID, o.Code,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := o.Scan(rows); err != nil {
			return err
		}
	}
	return nil
}

func GetOTPByUserID(user_id uuid.UUID) (*OTP, error) {
	o := new(OTP)
	if err := database.Get(o, "otp/fetch_by_userid", user_id); err != nil {
		return nil, err
	}
	return o, nil
}
