package httpserver

import (
	"gameApp/dto"
	"gameApp/pkg/httpmsg"
	"gameApp/service/userservice"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s Server) userLogin(c echo.Context) error {

	var Req userservice.LoginRequest
	if err := c.Bind(&Req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resp, err := s.userSvc.Login(Req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, resp)
}

func (s Server) userRegister(c echo.Context) error {

	var Req dto.RegisterRequest
	if err := c.Bind(&Req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())

	}

	if fieldErrors,err:= s.userValidator.ValidateRegisterRequest(Req); err!=nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message":msg,
			"errors": fieldErrors,
		})
	}
	resp, err := s.userSvc.Register(Req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

func (s Server) userProfile(c echo.Context) error {

	// The first way (Parse sessionID) :
	// validate sessionID by database and get userID
	// sessionID := req.Header.Get("SessionID")

	// The second way (Parse jwt token) :
	authToken := c.Request().Header.Get("Authorization")
	claims, err := s.authSvc.ParseToken(authToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}


	resp, err := s.userSvc.Profile(userservice.ProfileRequest{UserID: claims.UserID})
	if err != nil {
		msg, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, resp)
}
