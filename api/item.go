package api

import (
	"net/http"

	"github.com/Sudo-Ivan/fusionx/server"

	"github.com/labstack/echo/v4"
)

type itemAPI struct {
	srv *server.Item
}

func newItemAPI(srv *server.Item) *itemAPI {
	return &itemAPI{
		srv: srv,
	}
}

func (i itemAPI) List(c echo.Context) error {
	var req server.ReqItemList
	if err := bindAndValidate(&req, c); err != nil {
		return err
	}

	resp, err := i.srv.List(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (i itemAPI) Get(c echo.Context) error {
	var req server.ReqItemGet
	if err := bindAndValidate(&req, c); err != nil {
		return err
	}

	resp, err := i.srv.Get(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (i itemAPI) Delete(c echo.Context) error {
	var req server.ReqItemDelete
	if err := bindAndValidate(&req, c); err != nil {
		return err
	}

	if err := i.srv.Delete(c.Request().Context(), &req); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (i itemAPI) UpdateUnread(c echo.Context) error {
	var req server.ReqItemUpdateUnread
	if err := bindAndValidate(&req, c); err != nil {
		return err
	}

	if err := i.srv.UpdateUnread(c.Request().Context(), &req); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (i itemAPI) UpdateBookmark(c echo.Context) error {
	var req server.ReqItemUpdateBookmark
	if err := bindAndValidate(&req, c); err != nil {
		return err
	}

	if err := i.srv.UpdateBookmark(c.Request().Context(), &req); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
