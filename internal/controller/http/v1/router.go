package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"os"
	"wb-l-zero/internal/service"
)

func NewRouter(handler *echo.Echo, service *service.Service) {
	handler.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}", "method":"${method}","uri":"${uri}", "status":${status},"error":"${error}"}` + "\n",
		Output: setLogsFile("./logs/requests.log"),
	}))
	handler.Use(middleware.Recover())
	//handler.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := handler.Group("/api/v1")
	{
		newOrderRoutes(v1, service.Order)
	}

}

func setLogsFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		logrus.Fatal(err)
	}
	return file
}
