package v1

import (
	"net/http"

	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"github.com/miaozhang/webservice/common"
	"github.com/miaozhang/webservice/service/article_service"
	"github.com/miaozhang/webservice/service/tag_service"
	"github.com/miaozhang/webservice/settings"
	"github.com/miaozhang/webservice/util"
)

// @Summary Get a single article
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} common.Response
// @Failure 500 {object} common.Response
// @Router /api/v1/articles/{id} [get]
func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID > 0")

	if valid.HasErrors() {
		common.MarkErrors(valid.Errors)
		common.OutputRes(c, http.StatusBadRequest, common.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		common.OutputRes(c, http.StatusInternalServerError, common.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		common.OutputRes(c, http.StatusOK, common.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	article, err := articleService.Get()
	if err != nil {
		common.OutputRes(c, http.StatusInternalServerError, common.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}

	common.OutputRes(c, http.StatusOK, common.SUCCESS, article)
}

// @Summary Get multiple articles
// @Produce  json
// @Param tag_id body int false "TagID"
// @Param state body int false "State"
// @Param created_by body int false "CreatedBy"
// @Success 200 {object} common.Response
// @Failure 500 {object} common.Response
// @Router /api/v1/articles [get]
func GetArticles(c *gin.Context) {
	valid := validation.Validation{}
	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state")
	}

	tagId := -1
	if arg := c.PostForm("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		valid.Min(tagId, 1, "tag_id")
	}

	if valid.HasErrors() {
		common.MarkErrors(valid.Errors)
		common.OutputRes(c, http.StatusBadRequest, common.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{
		TagID:    tagId,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: settings.AppSetting.PageSize,
	}

	total, err := articleService.Count()
	if err != nil {
		common.OutputRes(c, http.StatusInternalServerError, common.ERROR_COUNT_ARTICLE_FAIL, nil)
		return
	}

	articles, err := articleService.GetAll()
	if err != nil {
		common.OutputRes(c, http.StatusInternalServerError, common.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}

	data := make(map[string]interface{})
	data["lists"] = articles
	data["total"] = total

	common.OutputRes(c, http.StatusOK, common.SUCCESS, data)
}

type AddArticleForm struct {
	TagID     int    `form:"tag_id" valid:"Required;Min(1)"`
	Title     string `form:"title" valid:"Required;MaxSize(100)"`
	Desc      string `form:"desc" valid:"Required;MaxSize(255)"`
	Content   string `form:"content" valid:"Required;MaxSize(65535)"`
	CreatedBy string `form:"created_by" valid:"Required;MaxSize(100)"`
	State     int    `form:"state" valid:"Range(0,1)"`
}

// @Summary Add article
// @Produce  json
// @Param tag_id body int true "TagID"
// @Param title body string true "Title"
// @Param desc body string true "Desc"
// @Param content body string true "Content"
// @Param created_by body string true "CreatedBy"
// @Param state body int true "State"
// @Success 200 {object} common.Response
// @Failure 500 {object} common.Response
// @Router /api/v1/articles [post]
func AddArticle(c *gin.Context) {
	var form AddArticleForm

	httpCode, errCode := common.BindAndValid(c, &form)
	if errCode != common.SUCCESS {
		common.OutputRes(c, httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{ID: form.TagID}
	exists, err := tagService.ExistByID()
	if err != nil {
		common.OutputRes(c, http.StatusInternalServerError, common.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		common.OutputRes(c, http.StatusOK, common.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	articleService := article_service.Article{
		TagID:     form.TagID,
		Title:     form.Title,
		Desc:      form.Desc,
		Content:   form.Content,
		State:     form.State,
		CreatedBy: form.CreatedBy,
	}
	if err := articleService.Add(); err != nil {
		common.OutputRes(c, http.StatusInternalServerError, common.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}

	common.OutputRes(c, http.StatusOK, common.SUCCESS, nil)
}

type EditArticleForm struct {
	ID         int    `form:"id" valid:"Required;Min(1)"`
	TagID      int    `form:"tag_id" valid:"Required;Min(1)"`
	Title      string `form:"title" valid:"Required;MaxSize(100)"`
	Desc       string `form:"desc" valid:"Required;MaxSize(255)"`
	Content    string `form:"content" valid:"Required;MaxSize(65535)"`
	ModifiedBy string `form:"modified_by" valid:"Required;MaxSize(100)"`
	State      int    `form:"state" valid:"Range(0,1)"`
}

// @Summary Update article
// @Produce  json
// @Param id path int true "ID"
// @Param tag_id body string false "TagID"
// @Param title body string false "Title"
// @Param desc body string false "Desc"
// @Param content body string false "Content"
// @Param modified_by body string true "ModifiedBy"
// @Param state body int false "State"
// @Success 200 {object} common.Response
// @Failure 500 {object} common.Response
// @Router /api/v1/articles/{id} [put]
func EditArticle(c *gin.Context) {
	form := EditArticleForm{ID: com.StrTo(c.Param("id")).MustInt()}

	httpCode, errCode := common.BindAndValid(c, &form)
	if errCode != common.SUCCESS {
		common.OutputRes(c, httpCode, errCode, nil)
		return
	}

	articleService := article_service.Article{
		ID:         form.ID,
		TagID:      form.TagID,
		Title:      form.Title,
		Desc:       form.Desc,
		Content:    form.Content,
		ModifiedBy: form.ModifiedBy,
		State:      form.State,
	}
	exists, err := articleService.ExistByID()
	if err != nil {
		common.OutputRes(c, http.StatusInternalServerError, common.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		common.OutputRes(c, http.StatusOK, common.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	tagService := tag_service.Tag{ID: form.TagID}
	exists, err = tagService.ExistByID()
	if err != nil {
		common.OutputRes(c, http.StatusInternalServerError, common.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		common.OutputRes(c, http.StatusOK, common.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = articleService.Edit()
	if err != nil {
		common.OutputRes(c, http.StatusInternalServerError, common.ERROR_EDIT_ARTICLE_FAIL, nil)
		return
	}

	common.OutputRes(c, http.StatusOK, common.SUCCESS, nil)
}

// @Summary Delete article
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} common.Response
// @Failure 500 {object} common.Response
// @Router /api/v1/articles/{id} [delete]
func DeleteArticle(c *gin.Context) {
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID > 0")

	if valid.HasErrors() {
		common.MarkErrors(valid.Errors)
		common.OutputRes(c, http.StatusOK, common.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		common.OutputRes(c, http.StatusInternalServerError, common.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		common.OutputRes(c, http.StatusOK, common.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	err = articleService.Delete()
	if err != nil {
		common.OutputRes(c, http.StatusInternalServerError, common.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}

	common.OutputRes(c, http.StatusOK, common.SUCCESS, nil)
}
