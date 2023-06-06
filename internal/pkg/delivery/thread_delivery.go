package delivery

import (
	"database/sql"
	"net/http"

	"forum/internal/pkg/domain"
	"forum/internal/pkg/usecase"

	"github.com/labstack/echo/v4"
	easyjson "github.com/mailru/easyjson"
)

type ThreadHandler struct {
	threads *usecase.ThreadUsecase
}

func NewThreadHandler(db *sql.DB) *ThreadHandler {
	return &ThreadHandler{
		threads: usecase.NewThreadUsecase(db),
	}
}

func (h *ThreadHandler) Create(c echo.Context) error {
	thread := &domain.Thread{}

	err := easyjson.UnmarshalFromReader(c.Request().Body, thread)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, MsgBadJSON)
	}

	err = h.threads.Create(thread)
	switch err {
	case nil:
		return c.JSON(http.StatusCreated, thread)
	case domain.ErrNotFound:
		return echo.NewHTTPError(http.StatusNotFound, MsgUserNotFound)
	case domain.ErrAlreadyExists:
		return echo.NewHTTPError(http.StatusConflict, thread)
	default:
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}

func (h *ThreadHandler) Get(c echo.Context) error {
	thread, err := h.threads.Get(c.Param("slug_or_id"))
	switch err {
	case nil:
		return c.JSON(http.StatusOK, thread)
	case domain.ErrNotFound:
		return echo.NewHTTPError(http.StatusNotFound, MsgThreadNotFound)
	default:
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}

func (h *ThreadHandler) Update(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented)
}
