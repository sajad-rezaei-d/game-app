package httpmsg

import (
	"gameApp/pkg/errmsg"
	"gameApp/pkg/richerror"
	"net/http"
)



func Error(err error) (message string, code int) {
	switch e:= err.(type) {
	case richerror.RichError:
		msg := e.Message()
		// we should not expose unexpected error message
		code := mapKindToHTTPStatusCode(e.Kind())
		if code >=500 {
			msg = errmsg.ErrorMsgSomethingWentWrong

		}
		return msg,code
	default:
		return	err.Error(),http.StatusBadRequest
	}

}

func mapKindToHTTPStatusCode (kind richerror.Kind) int {
	switch kind {
	case richerror.KindInvalid:
		return http.StatusUnprocessableEntity
	case richerror.KindNotFound:
		return http.StatusNotFound
	case richerror.KindForbidden:
		return http.StatusForbidden
	case richerror.KindUnexpected:
		return http.StatusInternalServerError		
	default:
		return http.StatusBadRequest		
	}
}	