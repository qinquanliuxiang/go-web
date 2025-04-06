package validator

import (
	"context"
	"errors"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/zh"
	"qqlx/base/logger"
	"strings"
)

var trans ut.Translator

type CheckReqInterface interface {
	CheckReq(ctx context.Context, value any) (errMsg string, err error)
}

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()
	zhTrans := zh.New()
	uni := ut.New(zhTrans, zhTrans)
	trans, _ = uni.GetTranslator("zh")
	_ = translations.RegisterDefaultTranslations(v, trans)

	return &Validator{
		validate: v,
	}
}

func (receive *Validator) CheckReq(ctx context.Context, value any) (errMsg string, err error) {
	err = receive.validate.Struct(value)
	if err == nil {
		return "", nil
	}

	var valErrors validator.ValidationErrors
	if !errors.As(err, &valErrors) {
		logger.WithContext(ctx, true).Errorf("validate check %#v error: %v", value, err)
		return "", errors.New("validate check exception")
	}

	// 使用翻译器
	msgArr := make([]string, 0, len(valErrors))
	for _, e := range valErrors {
		msg := e.Translate(trans)
		msgArr = append(msgArr, msg)
	}

	return strings.Join(msgArr, "; "), nil
}
