package validation

import (
	"os"
	"reflect"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func registrationFunc(tag, translation string) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) error {
		err := ut.Add(tag, translation, true)

		return err
	}
}

func translateFunc(ut ut.Translator, filedErr validator.FieldError) string {
	t, err := ut.T(filedErr.Tag(), reflect.ValueOf(filedErr.Value()).String())
	if err != nil {
		return filedErr.(error).Error()
	}

	return t
}

func validateDir(fl validator.FieldLevel) bool {
	path := fl.Field().String()
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return true
	}

	return false
}

func validateFile(fl validator.FieldLevel) bool {
	path := fl.Field().String()
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return true
	}

	return false
}

func validateDescription(fl validator.FieldLevel) bool {
	description := fl.Field().String()

	return len(description) <= maxDescriptionLength
}

func validateName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	if errs := IsQualifiedName(name); len(errs) > 0 {
		return false
	}

	return true
}

func prefixEach(msgs []string, prefix string) []string {
	for i := range msgs {
		msgs[i] = prefix + msgs[i]
	}
	return msgs
}
