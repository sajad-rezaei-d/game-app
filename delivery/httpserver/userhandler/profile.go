package userhandler

import (
	"gameApp/param"
	"gameApp/pkg/claim"
	"gameApp/pkg/httpmsg"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) userProfile(c echo.Context) error {

	// The first way (Parse sessionID) :
	// validate sessionID by database and get userID
	// sessionID := req.Header.Get("SessionID")

	// The second way (Parse jwt token) :
	//authToken := c.Request().Header.Get("Authorization")
	//claims, err := h.authSvc.ParseToken(authToken)
	//if err != nil {
	//	return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	//}

	// The third way (use echojwt middlware)
	//claims := c.Get(constant.AuthMiddlewareContextKey)
	//fmt.Println("claims", claims)
	//cl, ok := claims.(*authservice.Claims)
	//if ok {
	//	fmt.Println("cl", cl)
	//}else {
	//	fmt.Println("nok")
	//}
	// create func for third way :
	claims:= claim.GetClaimsFromEchoContext(c)
	resp, err := h.userSvc.Profile(param.ProfileRequest{UserID:claims.UserID})
	if err != nil {
		msg, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
	}
	return c.JSON(http.StatusOK, resp)
}
