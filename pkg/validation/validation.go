package validation

import (
	"fmt"

	english "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"

	"github.com/pachirode/iam_study/pkg/validation/field"
)

const (
	maxDescriptionLength = 255
)

type Validator struct {
	val   *validator.Validate
	data  interface{}
	trans ut.Translator
}

func NewValidator(data interface{}) *Validator {
	result := validator.New()

	result.RegisterValidation("dir", validateDir)
	result.RegisterValidation("file", validateFile)
	result.RegisterValidation("description", validateDescription)
	result.RegisterValidation("name", validateName)

	// translations
	eng := english.New()
	uni := ut.New(eng, eng)
	trans, _ := uni.GetTranslator("en")
	err := en.RegisterDefaultTranslations(result, trans)
	if err != nil {
		panic(err)
	}

	translations := []struct {
		tag         string
		translation string
	}{
		{
			tag:         "dir",
			translation: "{0} must point to an existing directory, but found '{1}'",
		},
		{
			tag:         "file",
			translation: "{0} must point to an existing file, but found '{1}'",
		},
		{
			tag:         "description",
			translation: fmt.Sprintf("must be less than %d", maxDescriptionLength),
		},
		{
			tag:         "name",
			translation: "is not a invalid name",
		},
	}

	for _, t := range translations {
		err = result.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation), translateFunc)
		if err != nil {
			panic(err)
		}
	}

	return &Validator{
		val:   result,
		data:  data,
		trans: trans,
	}
}

func (v *Validator) Validate() field.ErrorList {
	err := v.val.Struct(v.data)
	if err == nil {
		return nil
	}

	if _, ok := err.(*validator.InvalidValidationError); ok {
		return field.ErrorList{field.Invalid(field.NewPath(""), err.Error(), "")}
	}

	allErrs := field.ErrorList{}

	vErrors, _ := err.(validator.ValidationErrors)
	for _, vErr := range vErrors {
		allErrs = append(allErrs, field.Invalid(field.NewPath(vErr.Namespace()), vErr.Translate(v.trans), ""))
	}

	return allErrs
}
