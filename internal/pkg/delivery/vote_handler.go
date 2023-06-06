package delivery

import (
	"net/http"
	"strconv"

	"forum/internal/pkg/domain"
	"forum/internal/pkg/repository"
	"forum/internal/pkg/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	easyjson "github.com/mailru/easyjson"
)

type VoteHandler struct {
	votes   *repository.VoteRepository
	threads *usecase.ThreadUsecase
	users   *repository.UserRepository
}

func NewVoteHandler(db *pgxpool.Pool, threads *usecase.ThreadUsecase) *VoteHandler {
	return &VoteHandler{
		votes:   repository.NewVoteRepository(db),
		users:   repository.NewUserRepository(db),
		threads: threads,
	}
}

//easyjson:json
type VoteRequest struct {
	Nickname string
	Voice    int
}

func (h *VoteHandler) AddVote(c echo.Context) error {
	slugOrId := c.Param("slug_or_id")

	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		thread, err := h.threads.Get(slugOrId)
		switch err {
		case nil:
			break
		case domain.ErrNotFound:
			return echo.NewHTTPError(http.StatusNotFound, MsgThreadNotFound)
		default:
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		id = int(thread.Id)
	}

	req := &VoteRequest{}
	err = easyjson.UnmarshalFromReader(c.Request().Body, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, MsgBadJSON)
	}

	user, err := h.users.GetByNickname(req.Nickname)
	switch err {
	case nil:
		break
	case domain.ErrNotFound:
		return echo.NewHTTPError(http.StatusNotFound, MsgUserNotFound)
	default:
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = h.votes.Vote(&domain.Vote{
		UserId:   user.Id,
		ThreadId: id,
		Value:    req.Voice,
	})

	switch err {
	case nil:
		break
	case domain.ErrNotFound:
		return echo.NewHTTPError(http.StatusNotFound, MsgThreadNotFound)
	default:
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	thread, err := h.threads.GetById(id)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, thread)
}
