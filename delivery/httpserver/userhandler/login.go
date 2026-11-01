package userhandler

import (
	"gameApp/param"
	"gameApp/pkg/httpmsg"
	"net/http"

	"github.com/labstack/echo/v4"
)


func (h Handler) userLogin(c echo.Context) error {

	var Req param.LoginRequest
	if err := c.Bind(&Req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if fieldErrors,err:= h.userValidator.ValidateLoginRequest(Req); err!=nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message":msg,
			"errors": fieldErrors,
		})
	}

	resp, err := h.userSvc.Login(Req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, resp)
}
