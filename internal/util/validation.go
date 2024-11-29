package util

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func Validate[T any](data T) map[string]string {
	err := validator.New().Struct(data)
	res := map[string]string{}
	if err != nil {
		for _, v := range err.(validator.ValidationErrors) {
			res[v.StructField()] = TranslateTag(v)
		}
	}
	return res
}

func TranslateTag(fd validator.FieldError) string {
	switch fd.ActualTag() {
	case "required":
		return fmt.Sprintf("field %s harus diisi", fd.StructField())
	case "unique":
		return fmt.Sprintf("field %s harus unik", fd.StructField())
	case "min":
		return fmt.Sprintf("field %s minimal %d karakter", fd.StructField(), fd.Param())
	}
	return "validasi gagal"
}
