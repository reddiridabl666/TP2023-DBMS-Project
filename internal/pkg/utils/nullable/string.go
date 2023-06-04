package nullable

import (
	"database/sql"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

type String struct {
	sql.NullString
}

func NewString(i interface{}) *String {
	res := &String{}
	val, ok := i.(string)

	if !ok || val == "" {
		res.Valid = false
		return res
	}

	res.String = val
	res.Valid = true
	return res
}

func (s String) IsDefined() bool {
	return s.Valid
}

func (s String) MarshalEasyJSON(w *jwriter.Writer) {
	if s.Valid {
		w.String(s.String)
	} else {
		w.RawString("null")
	}
}

func (v *String) UnmarshalEasyJSON(l *jlexer.Lexer) {
	if l.IsNull() {
		v.Valid = false
		v.String = ""
		return
	}

	v.String = l.String()
	v.Valid = true
}
