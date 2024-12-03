package v1

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"wb-l-zero/internal/service"
)

type orderRoutes struct {
	orderService service.Order
}

func newOrderRoutes(g *echo.Group, orderService service.Order) {
	r := &orderRoutes{orderService: orderService}

	g.GET("/orders/:orderUID", r.getByID)
}

func (r *orderRoutes) getByID(c echo.Context) error {
	orderUID := c.Param("orderUID")

	order, err := r.orderService.GetOrderDetails(c.Request().Context(), orderUID)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			return newErrorResponse(c, http.StatusNotFound, err)
		}

		return newErrorResponse(c, http.StatusInternalServerError, err)
	}

	return newSuccessResponse(c, "order retrieved", order)
}
