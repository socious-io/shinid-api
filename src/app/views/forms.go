package views

import "github.com/google/uuid"

type OrganizationForm struct {
	Name        string     `json:"name" validate:"required,min=3,max=32"`
	Description string     `json:"description" validate:"required,min=3,max=32"`
	LogoID      *uuid.UUID `json:"logo_id"`
}
