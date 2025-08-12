package model

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/forbearing/golib/types/consts"
)

type GormTime time.Time

func (t *GormTime) Scan(value any) error {
	localTime, err := time.Parse(consts.DATE_TIME_LAYOUT, string(value.([]byte)))
	if err != nil {
		return err
	}
	*t = GormTime(localTime)
	return nil
}

func (t GormTime) Value() (driver.Value, error) {
	return time.Time(t).Format(consts.DATE_TIME_LAYOUT), nil
}

func (t *GormTime) UnmarshalJSON(b []byte) error {
	// Trim quotes from the stringified JSON value
	s := strings.Trim(string(b), "\"")
	// Parse the time using the custom format
	parsedTime, err := time.Parse(consts.DATE_TIME_LAYOUT, s)
	if err != nil {
		return err
	}

	*t = GormTime(parsedTime)
	return nil
}

func (ct GormTime) MarshalJSON() ([]byte, error) {
	// Convert the time to the custom format and stringify it
	return []byte("\"" + time.Time(ct).Format(consts.DATE_TIME_LAYOUT) + "\""), nil
}

func (gs *GormStrings) UnmarshalJSON(data []byte) error {
	// format: ["a", "b", "c"]
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		*gs = arr
		return nil
	}
	// format: "a,b,c"
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		s = strings.TrimSpace(s)
		s = strings.Trim(s, ",")
		if s == "" {
			*gs = []string{}
		} else {
			*gs = strings.Split(s, ",")
		}
		return nil
	}
	return errors.New("GormStrings: invalid JSON type")
}

type GormStrings []string

func (gs *GormStrings) Scan(value any) error {
	if value == nil {
		*gs = nil
		return nil
	}
	switch v := value.(type) {
	case []byte:
		_v := bytes.TrimSpace(v)
		_v = bytes.Trim(_v, ",")
		*gs = strings.Split(string(_v), ",")
	case string:
		_v := strings.TrimSpace(v)
		_v = strings.Trim(_v, ",")
		*gs = strings.Split(_v, ",")
	default:
		return errors.New("GormStrings: unsupported type, expected []byte or string")
	}
	return nil
}
