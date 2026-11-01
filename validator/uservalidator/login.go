package uservalidator

import (
	"fmt"
	"gameApp/param"
	"gameApp/pkg/errmsg"
	"gameApp/pkg/richerror"
	"regexp"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateLoginRequest(req param.LoginRequest) (map[string]string, error) {
	const op = "uservalidator.ValidateLoginRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.PhoneNumber,
			 validation.Required,
			//if you want to use custom error message
			 validation.Match(regexp.MustCompile(phoneNumberRegex)).
			 Error(errmsg.ErrorMsgPhoneNumberIsNotValid), 
			 //with By method in validation packge you can define custom validation
			 validation.By(v.doesPhoneNumberExist)),
		validation.Field(&req.Password, validation.Required), 
	); err != nil {
		fieldErrors := make(map[string]string)
		
		errV, ok := err.(validation.Errors)
		if ok {
			for key,value := range errV {
				if value!= nil {
					fieldErrors[key] = value.Error()
				}
			}
		}
		return fieldErrors,richerror.New(op).WithMessage(errmsg.ErrormsgInvalidInput).
			WithKind(richerror.KindInvalid).
			WithMeta(map[string]interface{}{"req": req}).WithErr(err)
	}

	return nil,nil
}


func (v Validator) doesPhoneNumberExist(value interface{}) error {
	phoneNumber := value.(string)
	_, err := v.repo.GetUserByPhoneNumber(phoneNumber)
	if err!=nil {
		return fmt.Errorf(errmsg.ErrorMsgNotFound)
	}

	return nil
}
