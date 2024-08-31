package models

import (
	"database/sql/driver"
	"fmt"
)

type AttributeType string

// ENUM values
const (
	Text     AttributeType = "TEXT"
	Number   AttributeType = "NUMBER"
	Boolean  AttributeType = "BOOLEAN"
	Url      AttributeType = "URL"
	Datetime AttributeType = "DATETIME"
	Email    AttributeType = "EMAIL"
)

func (a *AttributeType) Scan(value interface{}) error {
	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to scan attribute type: %v", value)
	}
	*a = AttributeType(strValue)
	return nil
}

func (a AttributeType) Value() (driver.Value, error) {
	return string(a), nil
}

type CredentialStatusType string

// ENUM values
const (
	StatusRequested CredentialStatusType = "REQUESTED"
	StatusVerfied   CredentialStatusType = "VERIFIED"
	StatusFailed    CredentialStatusType = "FAILED"
)

func (c *CredentialStatusType) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*c = CredentialStatusType(string(v))
	case string:
		*c = CredentialStatusType(v)
	default:
		return fmt.Errorf("failed to scan credential type: %v", value)
	}
	return nil
}

func (c CredentialStatusType) Value() (driver.Value, error) {
	return string(c), nil
}
