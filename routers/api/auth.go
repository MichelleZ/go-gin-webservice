package api

import (
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"github.com/miaozhang/webservice/common"
	"github.com/miaozhang/webservice/service/auth_service"
	"github.com/miaozhang/webservice/util"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	a := auth{Username: username, Password: password}

	valid := validation.Validation{}
	ok, _ := valid.Valid(&a)
	if !ok {
		common.MarkErrors(valid.Errors)
		common.OutputRes(c, http.StatusBadRequest, common.INVALID_PARAMS, nil)
		return
	}

	authService := auth_service.Auth{Username: username, Password: password}
	isExist, err := authService.Check()
	if err != nil {
		common.OutputRes(c, http.StatusInternalServerError, common.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if !isExist {
		common.OutputRes(c, http.StatusUnauthorized, common.ERROR_AUTH, nil)
		return
	}

	token, err := util.GenerateToken(username, password)
	if err != nil {
		common.OutputRes(c, http.StatusInternalServerError, common.ERROR_AUTH_TOKEN, nil)
		return
	}

	common.OutputRes(c, http.StatusOK, common.SUCCESS, map[string]string{
		"token": token,
	})
}
