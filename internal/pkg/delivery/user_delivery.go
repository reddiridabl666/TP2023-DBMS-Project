package delivery

import (
	"net/http"

	"forum/internal/pkg/domain"
	"forum/internal/pkg/repository"

	"github.com/labstack/echo/v4"
	"github.com/mailru/easyjson"
)

type UserHandler struct {
	users *repository.UserRepository
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	user := &domain.User{}
	err := easyjson.UnmarshalFromReader(c.Request().Body, user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, MsgBadJSON)
	}

	alreadyExisting, err := h.users.CreateUser(user)
	switch err {
	case nil:
		return c.JSON(http.StatusCreated, user)
	case domain.ErrUniqueViolation:
		return c.JSON(http.StatusConflict, alreadyExisting)
	default:
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, MsgInternalError)
	}
}

func (h *UserHandler) GetUser(c echo.Context) error {
	nickname := c.Param("nickname")

	user, err := h.users.GetUser(nickname)

	switch err {
	case nil:
		return c.JSON(http.StatusOK, user)
	case domain.ErrNotFound:
		return c.JSON(http.StatusNotFound, MsgUserNotFound)
	default:
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, MsgInternalError)
	}
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	user := &domain.User{}
	err := easyjson.UnmarshalFromReader(c.Request().Body, user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, MsgBadJSON)
	}

	err = h.users.UpdateUser(user)

	switch err {
	case nil:
		return c.JSON(http.StatusOK, user)
	case domain.ErrNotFound:
		return c.JSON(http.StatusNotFound, MsgUserNotFound)
	case domain.ErrUniqueViolation:
		return c.JSON(http.StatusConflict, MsgUserExists)
	default:
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, MsgInternalError)
	}
}
