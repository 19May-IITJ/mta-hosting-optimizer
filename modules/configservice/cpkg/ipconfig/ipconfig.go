package ipconfig

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

type IPConfigData struct {
	Hostname    string `json:"hostname" validate:"required,string"`
	IPAddresses string `json:"ipAddresses" validate:"required,string"`
	Status      bool   `json:"status"`
}

func IsString(fl validator.FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()
	return kind == reflect.String
}
