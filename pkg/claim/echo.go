package claim

import (
	"gameApp/service/authservice"
	cfg "gameApp/config"
	"github.com/labstack/echo/v4"
)

func GetClaimsFromEchoContext(c echo.Context) *authservice.Claims {
	// let it crash
	return c.Get(cfg.AuthMiddlewareContextKey).(*authservice.Claims)
	/* defensive programing
	claims := c.Get(constant.AuthMiddlewareContextKey)
	fmt.Println("claims", claims)
	// convert claims of middlware to authservice.Claims
	cl, ok := claims.(*authservice.Claims)
	if !ok {
		panic("claim is not found")
	}
	return cl
	*/
}