package delivery

import (
	"net/http"
	"strconv"
	"strings"

	"forum/internal/pkg/domain"
	"forum/internal/pkg/repository"
	"forum/internal/pkg/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	easyjson "github.com/mailru/easyjson"
)

type PostHandler struct {
	posts   *repository.PostRepository
	users   *repository.UserRepository
	threads *usecase.ThreadUsecase
	forums  *usecase.ForumUsecase
}

func NewPostHandler(db *pgxpool.Pool, threads *usecase.ThreadUsecase, forums *usecase.ForumUsecase) *PostHandler {
	return &PostHandler{
		posts:   repository.NewPostRepository(db),
		threads: threads,
		forums:  forums,
		users:   repository.NewUserRepository(db),
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

//easyjson:json
type PostResponse struct {
	Post   *domain.Post
	Author *domain.User   `json:"author,omitempty"`
	Thread *domain.Thread `json:"thread,omitempty"`
	Forum  *domain.Forum  `json:"forum,omitempty"`
}

func (h *PostHandler) GetPost(c echo.Context) error {
	strId := c.Param("id")

	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	post, err := h.posts.GetPost(id)
	switch err {
	case nil:
		break
	case domain.ErrNotFound:
		return echo.NewHTTPError(http.StatusNotFound, MsgPostNotFound)
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp := &PostResponse{Post: post}

	related := strings.Split(c.QueryParam("related"), ",")

	for _, obj := range related {
		switch obj {
		case "user":
			if resp.Author != nil {
				break
			}
			user, err := h.users.GetByNickname(post.Author)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			resp.Author = user

		case "forum":
			if resp.Forum != nil {
				break
			}
			forum, err := h.forums.Get(post.Forum)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			resp.Forum = forum
		case "thread":
			if resp.Thread != nil {
				break
			}
			thread, err := h.threads.GetById(int(post.Thread))
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			resp.Thread = thread
		}
	}

	return c.JSON(http.StatusOK, resp)
}

//easyjson:json
type PostMessage struct {
	Message string
}

func (h *PostHandler) UpdatePost(c echo.Context) error {
	strId := c.Param("id")

	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	msg := &PostMessage{}
	err = easyjson.UnmarshalFromReader(c.Request().Body, msg)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, MsgBadJSON)
	}

	post := &domain.Post{Id: id, Message: msg.Message}
	err = h.posts.Update(post)

	switch err {
	case nil:
		return c.JSON(http.StatusOK, post)
	case domain.ErrNotFound:
		return echo.NewHTTPError(http.StatusNotFound, MsgPostNotFound)
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}
