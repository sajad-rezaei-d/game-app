package uservalidator

import (
	"fmt"
	"gameApp/param"
	"gameApp/pkg/errmsg"
	"gameApp/pkg/richerror"
	"regexp"

	"github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateRegisterRequest(req param.RegisterRequest) (map[string]string, error) {
	const op = "uservalidator.ValidateRegisterRequest"

	if err := validation.ValidateStruct(&req,
		//TODO: add 3 to config
		validation.Field(&req.Name,
			 validation.Required,
			  validation.Length(3, 50)),

		validation.Field(&req.Password,
			validation.Required,
			validation.Match(regexp.MustCompile(`^[A-Za-z0-9!@#^&*]{8,}$`))),

		validation.Field(&req.PhoneNumber,
			 validation.Required,
			//if you want to use custom error message
			 validation.Match(regexp.MustCompile(phoneNumberRegex)).Error(errmsg.ErrorMsgPhoneNumberIsNotValid),//.Error("your cusotm error message"), 
			 //with By method in validation packge you can define custom validation
			 validation.By(v.checkPhoneNumberUniqueness)),
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

func (v Validator) checkPhoneNumberUniqueness(value interface{}) error {
	phoneNumber := value.(string)
	fmt.Println(phoneNumber)
	if isUnique, err := v.repo.IsPhoneNumberUnique(phoneNumber); err != nil || !isUnique {
		if err != nil {
			return err
		}

		if !isUnique {
			return fmt.Errorf(errmsg.ErrorMsgPhoneNumberIsNotUnique)
		}
	}

	return nil
}
