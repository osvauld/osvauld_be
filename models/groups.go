package models

import (
	"database/sql/driver"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type Group struct {
	BaseModel
	Name      string    `gorm:"size:255;column:name"`
	Members   UUIDArray `gorm:"type:uuid[];column:members"`
	CreatedBy uuid.UUID `gorm:"type:uuid;column:created_by"`
}

func (u *Group) TableName() string {
	return "groups"
}

type UUIDArray []uuid.UUID

// Scan value into UUIDArray
func (a *UUIDArray) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return errors.New("Failed to unmarshal UUID value into UUIDArray")
	}
	strValue := string(asBytes)
	values := strings.Split(strings.Trim(strValue, "{}"), ",")
	for _, v := range values {
		v = strings.Trim(v, `"`) // trimming potential quotes
		*a = append(*a, uuid.MustParse(v))
	}
	return nil
}

// Value returns UUIDArray as a driver.Value
func (a UUIDArray) Value() (driver.Value, error) {
	out := "{"
	for i, u := range a {
		out += u.String()
		if i < len(a)-1 {
			out += ","
		}
	}
	out += "}"
	return out, nil
}
