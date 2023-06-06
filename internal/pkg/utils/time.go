package utils

import (
	"database/sql/driver"
	"time"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

type Time struct {
	time.Time
}

const (
	timeFormat = "2006-01-02T15:04:05.000Z"
	timeParse  = "2006-01-02T15:04:05.000-07:00"
)

func (t Time) MarshalEasyJSON(w *jwriter.Writer) {
	w.String(t.UTC().Format(timeFormat))
}

func (t *Time) UnmarshalEasyJSON(l *jlexer.Lexer) {
	var err error

	t.Time, err = time.Parse(timeParse, l.String())

	if err != nil {
		t.Time = time.Time{}
	}
}

func (t *Time) Scan(value interface{}) error {
	t.Time = value.(time.Time)
	return nil
}

func (t Time) Value() (driver.Value, error) {
	return t.Time, nil
}