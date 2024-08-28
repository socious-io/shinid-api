package views

import (
	"shin/src/app/models"

	"github.com/google/uuid"
)

type OrganizationForm struct {
	Name        string     `json:"name" validate:"required,min=3,max=32"`
	Description string     `json:"description" validate:"required,min=3"`
	LogoID      *uuid.UUID `json:"logo_id"`
}

type SchemaForm struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Public      bool    `json:"public"`
	Attributes  []struct {
		Name        string               `json:"name"`
		Description *string              `json:"description"`
		Type        models.AttributeType `json:"type"`
	} `json:"attributes"`
}

type VerificationForm struct {
	Name        string    `json:"name" validate:"required,min=3,max=32"`
	Description *string   `json:"description" validate:"required,min=3"`
	SchemaID    uuid.UUID `json:"schema_id" validate:"required"`
}

type ProfileUpdateForm struct {
	Username  string  `json:"username"`
	JobTitle  *string `json:"job_title"`
	Bio       *string `json:"bio"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Phone     *string `json:"phone"`
}
