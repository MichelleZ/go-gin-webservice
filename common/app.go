package common

import (
	"log"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		log.Println(err.Key, err.Message)
	}
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func OutputRes(c *gin.Context, httpCode, errCode int, data interface{}) {
	c.JSON(httpCode, Response{
		Code: errCode,
		Msg:  GetMsg(errCode),
		Data: data,
	})

	return
}

func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	err := c.Bind(form)
	if err != nil {
		return http.StatusBadRequest, INVALID_PARAMS
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return http.StatusInternalServerError, ERROR
	}
	if !check {
		MarkErrors(valid.Errors)
		return http.StatusBadRequest, INVALID_PARAMS
	}

	return http.StatusOK, SUCCESS
}
