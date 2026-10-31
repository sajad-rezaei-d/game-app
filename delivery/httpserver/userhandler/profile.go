package userhandler

import (
	"gameApp/dto"
	"gameApp/pkg/httpmsg"
	"net/http"

	"github.com/labstack/echo/v4"
)


func (h Handler) userProfile(c echo.Context) error {

	// The first way (Parse sessionID) :
	// validate sessionID by database and get userID
	// sessionID := req.Header.Get("SessionID")

	// The second way (Parse jwt token) :
	authToken := c.Request().Header.Get("Authorization")
	claims, err := h.authSvc.ParseToken(authToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}


	resp, err := h.userSvc.Profile(dto.ProfileRequest{UserID: claims.UserID})
	if err != nil {
		msg, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, resp)
}
