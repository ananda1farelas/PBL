package helper

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	_ = validate.RegisterValidation("nim", func(f1 validator.FieldLevel) bool {
		nim := strings.TrimSpace(f1.Field().String())
		matched, _ := regexp.MatchString(`^[0-9]{8,12}$`, nim)
		return matched
	})

	_ = validate.RegisterValidation("notblank", func(f1 validator.FieldLevel) bool {
		return strings.TrimSpace(f1.Field().String()) != ""
	})

}

func ValidateStruct(data any) map[string]string {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	errs := make(map[string]string)

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		errs["request"] = "data tidak valid"
		return errs
	}

	for _, fieldErr := range validationErrors {
		field := strings.ToLower(fieldErr.Field())
		errs[field] = messageFor(fieldErr)
	}
	return errs
}

func messageFor(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "field wajib diisi"
	case "notblank":
		return "field tidak boleh berisi spasi"
	case "min":
		return "panjang minimal " + fe.Param() + " karakter"
	case "max":
		return "panjang maksimal " + fe.Param() + " karakter"
	case "nim":
		return "NIM harus berupa 8 sampai 12 digit angka"
	default:
		return "field tidak memenuhi aturan validasi"
	}
}
