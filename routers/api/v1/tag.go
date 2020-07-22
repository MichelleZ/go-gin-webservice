package v1

import (
	"log"
	"net/http"

	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"github.com/miaozhang/webservice/common"
	"github.com/miaozhang/webservice/service/tag_service"
	"github.com/miaozhang/webservice/settings"
	"github.com/miaozhang/webservice/util"
)

// @Summary Get multiple article tags
// @Produce  json
// @Param name query string false "Name"
// @Param state query int false "State"
// @Success 200 {object} common.Response
// @Failure 500 {object} common.Response
// @Router /api/v1/tags [get]
func GetTags(c *gin.Context) {
	name := c.Query("name")
	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:     name,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: settings.AppSetting.PageSize,
	}

	tags, err := tagService.GetAll()
	if err != nil {
		common.OutputRes(c, http.StatusInternalServerError, common.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	count, err := tagService.Count()
	if err != nil {
		common.OutputRes(c, http.StatusInternalServerError, common.ERROR_COUNT_TAG_FAIL, nil)
		return
	}

	common.OutputRes(c, http.StatusOK, common.SUCCESS, map[string]interface{}{
		"lists": tags,
		"total": count,
	})
}

type AddTagForm struct {
	Name      string `form:"name" valid:"Required;MaxSize(100)"`
	CreatedBy string `form:"created_by" valid:"Required;MaxSize(100)"`
	State     int    `form:"state" valid:"Range(0,1)"`
}

// @Summary Add new tag
// @Produce json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Success 200 {object} common.Response
// @Failure 500 {object} common.Response
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	var form AddTagForm
	httpCode, errCode := common.BindAndValid(c, &form)
	if errCode != common.SUCCESS {
		common.OutputRes(c, httpCode, errCode, nil)
		return
	}
	log.Printf("hell form: %v", form)

	tagService := tag_service.Tag{
		Name:      form.Name,
		CreatedBy: form.CreatedBy,
		State:     form.State,
	}

	exists, err := tagService.ExistByName()
	if err != nil {
		common.OutputRes(c, http.StatusInternalServerError, common.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if exists {
		common.OutputRes(c, http.StatusOK, common.ERROR_EXIST_TAG, nil)
		return
	}

	err = tagService.Add()
	if err != nil {
		common.OutputRes(c, http.StatusInternalServerError, common.ERROR_ADD_TAG_FAIL, nil)
		return
	}

	common.OutputRes(c, http.StatusOK, common.SUCCESS, nil)
}

type EditTagForm struct {
	ID         int    `form:"id" valid:"Required;Min(1)"`
	Name       string `form:"name" valid:"Required;MaxSize(100)"`
	ModifiedBy string `form:"modified_by" valid:"Required;MaxSize(100)"`
	State      int    `form:"state" valid:"Range(0,1)"`
}

// @Summary Edit tag
// @Produce json
// @Param id path int true "ID"
// @Param state query int false "State"
// @Param modified_by query string true "ModifiedBy"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags/{id} [put]
func EditTag(c *gin.Context) {
	form := EditTagForm{ID: com.StrTo(c.Param("id")).MustInt()}

	httpCode, errCode := common.BindAndValid(c, &form)
	if errCode != common.SUCCESS {
		common.OutputRes(c, httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{
		ID:         form.ID,
		Name:       form.Name,
		ModifiedBy: form.ModifiedBy,
		State:      form.State,
	}

	exists, err := tagService.ExistByID()
	if err != nil {
		common.OutputRes(c, http.StatusInternalServerError, common.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		common.OutputRes(c, http.StatusOK, common.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = tagService.Edit()
	if err != nil {
		common.OutputRes(c, http.StatusInternalServerError, common.ERROR_EDIT_TAG_FAIL, nil)
		return
	}

	common.OutputRes(c, http.StatusOK, common.SUCCESS, nil)
}

// @Summary Delete tag
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} common.Response
// @Failure 500 {object} common.Response
// @Router /api/v1/tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		common.MarkErrors(valid.Errors)
		common.OutputRes(c, http.StatusBadRequest, common.INVALID_PARAMS, nil)
	}

	tagService := tag_service.Tag{ID: id}
	exists, err := tagService.ExistByID()
	if err != nil {
		common.OutputRes(c, http.StatusInternalServerError, common.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		common.OutputRes(c, http.StatusOK, common.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	if err := tagService.Delete(); err != nil {
		common.OutputRes(c, http.StatusInternalServerError, common.ERROR_DELETE_TAG_FAIL, nil)
		return
	}

	common.OutputRes(c, http.StatusOK, common.SUCCESS, nil)
}
