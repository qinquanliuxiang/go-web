package validator

import (
	"context"
	"errors"
	"fmt"
	"qqlx/base/logger"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CheckReqInterface interface {
	CheckReq(ctx context.Context, value any) (errMsg string, err error)
}

type Validator struct {
	validate *validator.Validate
}

type FormErrorField struct {
	ErrorField string `json:"errorField"`
	ErrorMsg   string `json:"errorMsg"`
}

func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func (receive *Validator) CheckReq(ctx context.Context, value any) (errMsg string, err error) {
	err = receive.validate.Struct(value)
	if err != nil {
		var valErrors validator.ValidationErrors
		if !errors.As(err, &valErrors) {
			logger.WithContext(ctx, true).Errorf("validate check %#v error: %v", value, err)
			return "", errors.New("validate check exception")
		}
		errFields := make([]FormErrorField, len(valErrors))
		for index, fieldError := range valErrors {
			errField := FormErrorField{
				ErrorField: fieldError.Field(),
				ErrorMsg:   fieldError.Tag(),
			}
			structNamespace := fieldError.StructNamespace()
			_, fieldName, found := strings.Cut(structNamespace, ".")
			if found {
				originalTag := getObjectTagByFieldName(ctx, value, fieldName)
				if len(originalTag) > 0 {
					errField.ErrorField = originalTag
				}
			}
			errFields[index] = errField
		}

		if len(errFields) > 0 {
			msgArr := make([]string, len(errFields))
			for i, v := range errFields {
				msgArr[i] = fmt.Sprintf("%s: %s", v.ErrorField, v.ErrorMsg)
			}
			return strings.Join(msgArr, ","), nil
		}
	}
	return "", nil
}

func getObjectTagByFieldName(ctx context.Context, obj any, fieldName string) (tag string) {
	defer func() {
		if err := recover(); err != nil {
			logger.WithContext(ctx, true).Error(err)
		}
	}()

	objT := reflect.TypeOf(obj)
	objT = objT.Elem()

	structField, exists := objT.FieldByName(fieldName)
	if !exists {
		return ""
	}
	tag = structField.Tag.Get("json")
	if len(tag) == 0 {
		tag = structField.Tag.Get("form")
	}
	if len(tag) == 0 {
		tag = structField.Tag.Get("uri")
	}
	return tag
}
