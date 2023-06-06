package delivery

import (
	"database/sql"
	"net/http"

	"forum/internal/pkg/domain"
	"forum/internal/pkg/repository"

	"github.com/labstack/echo/v4"
	"github.com/mailru/easyjson"
)

type UserHandler struct {
	users *repository.UserRepository
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{
		users: repository.NewUserRepository(db),
	}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	user := &domain.User{Nickname: c.Param("nickname")}

	err := easyjson.UnmarshalFromReader(c.Request().Body, user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, MsgBadJSON)
	}

	alreadyExisting, err := h.users.CreateUser(user)
	switch err {
	case nil:
		return c.JSON(http.StatusCreated, user)
	case domain.ErrAlreadyExists:
		return c.JSON(http.StatusConflict, alreadyExisting)
	default:
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}

func (h *UserHandler) GetUser(c echo.Context) error {
	nickname := c.Param("nickname")

	user, err := h.users.GetUser(nickname)

	switch err {
	case nil:
		return c.JSON(http.StatusOK, user)
	case domain.ErrNotFound:
		return echo.NewHTTPError(http.StatusNotFound, MsgUserNotFound)
	default:
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	user := &domain.User{Nickname: c.Param("nickname")}

	err := easyjson.UnmarshalFromReader(c.Request().Body, user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, MsgBadJSON)
	}

	err = h.users.UpdateUser(user)

	switch err {
	case nil:
		return c.JSON(http.StatusOK, user)
	case domain.ErrNotFound:
		return echo.NewHTTPError(http.StatusNotFound, MsgUserNotFound)
	case domain.ErrAlreadyExists:
		return echo.NewHTTPError(http.StatusConflict, MsgUserExists)
	default:
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}
