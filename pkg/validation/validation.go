package validation

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

const (
	maxDescriptionLength = 255
)

type Validator struct {
	val   *validator.Validate
	data  interface{}
	trans ut.Translator
}
