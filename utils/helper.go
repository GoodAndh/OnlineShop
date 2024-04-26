package utils

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

func ValidationField(validate *validator.Validate, payload any) (map[any]any, error) {

	if validate == nil {
		return map[any]any{
			"error": "validate equal to nil",
		}, errors.New("validate equal to nil")
	}

	if err := validate.Struct(payload); err != nil {
		errorList := make(map[any]any)
		errors := err.(validator.ValidationErrors)
		for _, er := range errors {
			va := reflect.TypeOf(payload).Elem()
			field, _ := va.FieldByName(er.StructField())
			fieldName := field.Tag.Get("json")
			var errMsg string
			switch er.Tag() {
			case "number":
				errMsg = fmt.Sprintf("form %v hanya bisa di isi dengan nomor", er.Field())
			case "required":
				errMsg = fmt.Sprintf("form %v harus diisi", er.Field())
			case "min":
				errMsg = fmt.Sprintf("form %v minimal memiliki %v karakter", er.Field(), er.Param())
			case "max":
				errMsg = fmt.Sprintf("form %v maximal memiliki %v karakter", er.Field(), er.Param())
			case "email":
				errMsg = fmt.Sprintf("form %v tidak sesuai dengan format email", er.Field())
			case "eqfield":
				errMsg = fmt.Sprintln("pastikan konfirmasi password sama dengan password")
			case "oneof":
				errMsg = fmt.Sprintf("silahkan pilih salah satu dari %v ", er.Param())
			case "gt":
				errMsg=fmt.Sprintf("form %v harus memiliki minal %v angka",er.Field(),er.Param())
			}
			errorList["error"+fieldName] = errMsg
		}

		return errorList, err
	}
	return map[any]any{}, nil
}
