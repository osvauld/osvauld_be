package models

import (
	"database/sql/driver"
	"strings"
)

type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	return "{" + strings.Join(a, ",") + "}", nil
}

func (a *StringArray) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return nil
	}

	str := string(asBytes)
	*a = strings.Split(str[1:len(str)-1], ",")

	return nil
}
