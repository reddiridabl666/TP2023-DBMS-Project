package delivery

import (
	"database/sql"
	"net/http"

	"forum/internal/pkg/domain"
	"forum/internal/pkg/repository"

	"github.com/labstack/echo/v4"
	"github.com/mailru/easyjson"
)

type ForumHandler struct {
	forums *repository.ForumRepository
}

func NewForumHandler(db *sql.DB) *ForumHandler {
	return &ForumHandler{
		forums: repository.NewForumRepository(db),
	}
}

func (h *ForumHandler) Create(c echo.Context) error {
	forum := &domain.Forum{}

	err := easyjson.UnmarshalFromReader(c.Request().Body, forum)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, MsgBadJSON)
	}

	err = h.forums.Create(forum)
	switch err {
	case nil:
		return c.JSON(http.StatusCreated, forum)
	case domain.ErrNotFound:
		return echo.NewHTTPError(http.StatusNotFound, MsgUserNotFound)
	case domain.ErrAlreadyExists:
		return echo.NewHTTPError(http.StatusConflict, forum)
	default:
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}

func (h *ForumHandler) Get(c echo.Context) error {
	forum, err := h.forums.Get(c.Param("slug"))
	switch err {
	case nil:
		return c.JSON(http.StatusOK, forum)
	case domain.ErrNotFound:
		return echo.NewHTTPError(http.StatusNotFound, MsgForumNotFound)
	default:
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}
