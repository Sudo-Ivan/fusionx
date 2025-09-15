package api

import (
	"net/http"

	"github.com/0x2e/fusion/server"
	"github.com/labstack/echo/v4"
)

type configAPI struct {
	srv *server.Config
}

func newConfigAPI(srv *server.Config) *configAPI {
	return &configAPI{
		srv: srv,
	}
}

func (c configAPI) Get(ctx echo.Context) error {
	resp, err := c.srv.Get(ctx.Request().Context())
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (c configAPI) Update(ctx echo.Context) error {
	var req server.ReqConfigUpdate
	if err := bindAndValidate(&req, ctx); err != nil {
		return err
	}

	if err := c.srv.Update(ctx.Request().Context(), &req); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}
