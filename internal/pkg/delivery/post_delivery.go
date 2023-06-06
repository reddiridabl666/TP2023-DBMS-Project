package delivery

import (
	"database/sql"
	"net/http"

	"forum/internal/pkg/domain"
	"forum/internal/pkg/repository"
	"forum/internal/pkg/usecase"

	"github.com/labstack/echo/v4"
	easyjson "github.com/mailru/easyjson"
)

type PostHandler struct {
	posts   *repository.PostRepository
	threads *usecase.ThreadUsecase
}

func NewPostHandler(db *sql.DB, threads *usecase.ThreadUsecase) *PostHandler {
	return &PostHandler{
		posts:   repository.NewPostRepository(db),
		threads: threads,
	}
}

func (h *PostHandler) AddPosts(c echo.Context) error {
	posts := domain.PostBatch([]*domain.Post{})

	err := easyjson.UnmarshalFromReader(c.Request().Body, &posts)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, MsgBadJSON)
	}

	thread, err := h.threads.Get(c.Param("slug_or_id"))
	switch err {
	case nil:
		break
	case domain.ErrNotFound:
		return echo.NewHTTPError(http.StatusNotFound, MsgThreadNotFound)
	default:
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if len(posts) == 0 {
		return c.JSON(http.StatusCreated, posts)
	}

	err = h.posts.AddPosts(thread, posts)
	switch err {
	case nil:
		return c.JSON(http.StatusCreated, posts)
	case domain.ErrNoParent, domain.ErrInvalidParent:
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	case domain.ErrNotFound:
		return echo.NewHTTPError(http.StatusNotFound, MsgUserNotFound)
	default:
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}
