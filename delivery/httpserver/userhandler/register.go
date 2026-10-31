package userhandler

import (
	"gameApp/dto"
	"gameApp/pkg/httpmsg"
	"net/http"

	"github.com/labstack/echo/v4"
)




func (h Handler) userRegister(c echo.Context) error {

	var Req dto.RegisterRequest
	if err := c.Bind(&Req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())

	}

	if fieldErrors,err:= h.userValidator.ValidateRegisterRequest(Req); err!=nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message":msg,
			"errors": fieldErrors,
		})
	}
	resp, err := h.userSvc.Register(Req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}