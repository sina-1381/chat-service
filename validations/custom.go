package validations

import (
	"errors"
	"fmt"
	"ginGorm/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"net/http"
	"strings"
)

var Uniq validator.Func = func(fl validator.FieldLevel) bool {
	date := fl.Field().Interface()
	keyParam:= strings.Split(fl.Param(),".")
	var count int64
	models.DB.Table(keyParam[0]).Select("id").Where(keyParam[1]+" = ?", date).Count(&count)
	if count != 0{
		return false
	}else {
		return true
	}
}

var Exists validator.Func = func(fl validator.FieldLevel) bool {
	date := fl.Field().Interface()
	keyParam:= strings.Split(fl.Param(),".")
	var count int64
	models.DB.Table(keyParam[0]).Select("id").Where(keyParam[1]+" = ?", date).Count(&count)
	if count != 0{
		return true
	}else {
		return false
	}
}

func Simple(verr validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)

	for _, f := range verr {
		err := f.ActualTag()
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}
		errs[f.Field()] = err
	}
	return errs
}

func CheckValidate(input interface{}, c *gin.Context) {
	if err := c.ShouldBindJSON(input); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": Simple(verr)})
			panic("validation error")
		}
	}
}
