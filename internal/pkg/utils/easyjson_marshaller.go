package utils

import (
	"github.com/labstack/echo/v4"
	"github.com/mailru/easyjson"
)

type EasyJSONSerializer struct{}

var fallback = echo.DefaultJSONSerializer{}

func (e EasyJSONSerializer) Serialize(c echo.Context, i interface{}, indent string) error {
	obj, ok := i.(easyjson.Marshaler)
	if !ok {
		return fallback.Serialize(c, i, indent)
	}
	_, err := easyjson.MarshalToWriter(obj, c.Response())
	return err
}

func (e EasyJSONSerializer) Deserialize(c echo.Context, i interface{}) error {
	obj, ok := i.(easyjson.Unmarshaler)
	if !ok {
		return fallback.Deserialize(c, i)
	}
	return easyjson.UnmarshalFromReader(c.Request().Body, obj)
}
