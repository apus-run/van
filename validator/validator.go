package validator

import (
	"context"
	"errors"
	"maps"
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhcn "github.com/go-playground/validator/v10/translations/zh"
)

var (
	once sync.Once
	vd   *Validator
)

type Validator struct {
	validate   *validator.Validate
	translator ut.Translator
}

func (v *Validator) Validate(s any) error {
	return v.valid(s)
}

func (v *Validator) ValidateContext(ctx context.Context, s any) error {
	return v.validCtx(ctx, s)
}

func (v *Validator) ValidatePartial(s any, partials map[string]struct{}) error {
	return v.validate.StructPartial(s, slices.Collect(maps.Keys(partials))...)
}

func (v *Validator) valid(obj any) error {
	if reflect.Indirect(reflect.ValueOf(obj)).Kind() != reflect.Struct {
		return nil
	}

	e := v.validate.Struct(obj)
	if e != nil {
		err, ok := e.(validator.ValidationErrors)
		if !ok {
			return e
		}
		return removeStructName(err.Translate(v.translator))
	}
	return nil
}

func (v *Validator) validCtx(ctx context.Context, obj any) error {
	if reflect.Indirect(reflect.ValueOf(obj)).Kind() != reflect.Struct {
		return nil
	}

	e := v.validate.StructCtx(ctx, obj)
	if e != nil {
		err, ok := e.(validator.ValidationErrors)
		if !ok {
			return e
		}
		return removeStructName(err.Translate(v.translator))
	}
	return nil
}

func removeStructName(fields map[string]string) error {
	errs := make([]string, 0, len(fields))
	for _, err := range fields {
		errs = append(errs, err)
	}
	return errors.New(strings.Join(errs, ";"))
}

func NewValidator(opts ...Option) *Validator {
	once.Do(func() {
		validate := validator.New(validator.WithRequiredStructEnabled())

		zhTrans := zh.New()
		trans, _ := ut.New(zhTrans, zhTrans).GetTranslator("zh")
		zhcn.RegisterDefaultTranslations(validate, trans)

		for _, f := range opts {
			f(validate, trans)
		}

		vd = &Validator{
			validate:   validate,
			translator: trans,
		}
	})

	return vd
}

// V 是 NewValidator 简写
func V(opts ...Option) *Validator {
	return NewValidator(opts...)
}

// ValidateStruct 验证结构体
func ValidateStruct(obj any) error {
	return vd.Validate(obj)
}

// ValidateStructCtx 验证结构体，带Context
func ValidateStructContext(ctx context.Context, obj any) error {
	return vd.ValidateContext(ctx, obj)
}

func ValidatePartial(obj any, partials map[string]struct{}) error {
	return vd.ValidatePartial(obj, partials)
}
