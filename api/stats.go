package api

import (
	"net/http"

	"github.com/Sudo-Ivan/fusionx/server"
	"github.com/labstack/echo/v4"
)

type statsAPI struct {
	srv *server.Stats
}

func newStatsAPI(srv *server.Stats) *statsAPI {
	return &statsAPI{
		srv: srv,
	}
}

func (s statsAPI) Get(c echo.Context) error {
	resp, err := s.srv.Get(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
