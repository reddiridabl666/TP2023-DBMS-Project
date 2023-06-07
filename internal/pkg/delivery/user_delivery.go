package delivery

import (
	"net/http"
	"strconv"

	"forum/internal/pkg/domain"
	"forum/internal/pkg/repository"
	"forum/internal/pkg/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/mailru/easyjson"
)

type UserHandler struct {
	users  *repository.UserRepository
	forums *usecase.ForumUsecase
}

func NewUserHandler(db *pgxpool.Pool, forums *usecase.ForumUsecase) *UserHandler {
	return &UserHandler{
		users:  repository.NewUserRepository(db),
		forums: forums,
	}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	user := &domain.User{Nickname: c.Param("nickname")}

	err := easyjson.UnmarshalFromReader(c.Request().Body, user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, MsgBadJSON)
	}

	alreadyExisting, err := h.users.Create(user)
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

	user, err := h.users.GetByNickname(nickname)

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

	err = h.users.Update(user)

	switch err {
	case nil:
		return c.JSON(http.StatusOK, user)
	case domain.ErrNotFound:
		return echo.NewHTTPError(http.StatusNotFound, MsgUserNotFound)
	case domain.ErrAlreadyExists:
		return echo.NewHTTPError(http.StatusConflict, MsgUserExists)
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}

func (h *UserHandler) GetByForum(c echo.Context) error {
	params := getUserListParams(c)

	forum, err := h.forums.Get(c.Param("slug"))
	switch err {
	case nil:
		break
	case domain.ErrNotFound:
		return echo.NewHTTPError(http.StatusNotFound, MsgForumNotFound)
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	params.ForumId = forum.Id

	users, err := h.users.GetByForum(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

func getUserListParams(c echo.Context) *domain.UserListParams {
	res := &domain.UserListParams{}

	limit := c.QueryParam("limit")
	res.Limit, _ = strconv.Atoi(limit)

	if res.Limit < 1 {
		res.Limit = 100
	}

	if c.QueryParam("desc") == "true" {
		res.Desc = true
	}

	res.Since = c.QueryParam("since")
	return res
}
