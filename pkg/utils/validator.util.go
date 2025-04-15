// pkg/utils/validator.go
package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	
	// 구조체 필드 이름 대신 json 태그 이름 사용
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return fld.Name
		}
		return name
	})
}

// Validate는 구조체의 유효성을 검증합니다.
func Validate(data interface{}) error {
	if data == nil {
		return fmt.Errorf("데이터가 nil입니다")
	}

	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	// 검증 오류 메시지 포맷팅
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	errorMessages := make([]string, 0, len(validationErrors))
	for _, e := range validationErrors {
		switch e.Tag() {
		case "required":
			errorMessages = append(errorMessages, fmt.Sprintf("%s 필드는 필수입니다", e.Field()))
		case "oneof":
			errorMessages = append(errorMessages, fmt.Sprintf("%s 필드는 %s 중 하나여야 합니다", e.Field(), e.Param()))
		case "required_if":
			params := strings.Split(e.Param(), " ")
			if len(params) >= 2 {
				errorMessages = append(errorMessages, fmt.Sprintf("%s 필드는 %s가 %s일 때 필수입니다", e.Field(), params[0], params[1]))
			} else {
				errorMessages = append(errorMessages, fmt.Sprintf("%s 필드는 조건에 따라 필수입니다", e.Field()))
			}
		default:
			errorMessages = append(errorMessages, fmt.Sprintf("%s 필드가 %s 규칙을 만족하지 않습니다", e.Field(), e.Tag()))
		}
	}

	return fmt.Errorf("%s", strings.Join(errorMessages, ", "))
}