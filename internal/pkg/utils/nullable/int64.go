package nullable

import (
	"database/sql"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

type Int64 struct {
	sql.NullInt64
}

func NewInt64(i interface{}) *Int64 {
	res := &Int64{}
	val, ok := i.(int64)

	if !ok {
		res.Valid = false
		return res
	}

	res.Int64 = val
	res.Valid = true
	return res
}

func (i Int64) IsDefined() bool {
	return i.Valid
}

func (i Int64) MarshalEasyJSON(w *jwriter.Writer) {
	if i.Valid {
		w.Int64(i.Int64)
	} else {
		w.RawString("null")
	}
}

func (v *Int64) UnmarshalEasyJSON(l *jlexer.Lexer) {
	if l.IsNull() {
		v.Valid = false
		v.Int64 = 0
		return
	}

	v.Int64 = l.Int64()

	if v.Int64 != 0 {
		v.Valid = true
	}
}
