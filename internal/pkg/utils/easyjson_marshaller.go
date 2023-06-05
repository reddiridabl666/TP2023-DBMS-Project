package utils

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/mailru/easyjson"
)

type EasyJSONSerializer struct{}

func (e EasyJSONSerializer) Serialize(c echo.Context, i interface{}, _ string) error {
	obj, ok := i.(easyjson.Marshaler)
	if !ok {
		return errors.New("does not implement easyjson.Marshaller interface")
	}
	_, _, err := easyjson.MarshalToHTTPResponseWriter(obj, c.Response().Writer)
	return err
}

func (e EasyJSONSerializer) Deserialize(c echo.Context, i interface{}) error {
	obj, ok := i.(easyjson.Unmarshaler)
	if !ok {
		return errors.New("does not implement easyjson.Unmarshaller interface")
	}
	return easyjson.UnmarshalFromReader(c.Request().Body, obj)
}
