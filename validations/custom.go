package validations

import (
	"chat-service/models"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

// this custom validation function checks the database if a value is unique if we dont want
// duplicates of the same value . for example email of the user should be unique
var Uniq validator.Func = func(fl validator.FieldLevel) bool {
	date := fl.Field().Interface()
	keyParam := strings.Split(fl.Param(), ".")
	var count int64
	models.DB.Table(keyParam[0]).Select("id").Where(keyParam[1]+" = ?", date).Count(&count)
	if count != 0 {
		return false
	} else {
		return true
	}
}

// checks if this value exists in the database . opposite of the uniq validation
var Exists validator.Func = func(fl validator.FieldLevel) bool {
	date := fl.Field().Interface()
	keyParam := strings.Split(fl.Param(), ".")
	var count int64
	models.DB.Table(keyParam[0]).Select("id").Where(keyParam[1]+" = ?", date).Count(&count)
	if count != 0 {
		return true
	} else {
		return false
	}
}

// generating simple validation errors . this is used in CheckValidate func below
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

// this is where every validation is handled . every api in this project calls this function
// values of the request json body are recived here and checked based on rules set in validations
func CheckValidate(input interface{}, c *gin.Context) {
	if err := c.ShouldBindJSON(input); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": Simple(verr)})
			panic("validation error")
		}
	}
}
