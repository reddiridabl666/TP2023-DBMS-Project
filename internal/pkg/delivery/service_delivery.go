package delivery

import (
	"net/http"

	"forum/internal/pkg/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type ServiceHandler struct {
	service *repository.ServiceRepository
}

func NewServiceHandler(db *pgxpool.Pool) *ServiceHandler {
	return &ServiceHandler{
		service: repository.NewServiceRepository(db),
	}
}

func (h *ServiceHandler) Clear(c echo.Context) error {
	err := h.service.Clear()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

func (h *ServiceHandler) Status(c echo.Context) error {
	info, err := h.service.Status()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, info)
}
