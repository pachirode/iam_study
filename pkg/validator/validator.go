package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/pachirode/iam_study/pkg/validation"
)

func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if errs := validation.IsQualifiedName(username); len(errs) > 0 {
		return false
	}

	return true
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if err := validation.IsValidPassword(password); err != nil {
		return false
	}

	return true
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("username", validateUsername)
		_ = v.RegisterValidation("password", validatePassword)
	}
}
