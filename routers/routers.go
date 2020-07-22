package routers

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "github.com/miaozhang/webservice/docs"
	"github.com/miaozhang/webservice/middleware/jwt"
	"github.com/miaozhang/webservice/routers/api"
	v1 "github.com/miaozhang/webservice/routers/api/v1"
	"github.com/miaozhang/webservice/settings"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode(settings.ServerSetting.RunMode)

	r.GET("/auth", api.GetAuth)
	r.GET("/swagger/*ang", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiv1 := r.Group("api/v1")
	apiv1.Use(jwt.JWT())
	{
		apiv1.GET("/tags", v1.GetTags)
		apiv1.POST("/tags", v1.AddTag)
		apiv1.PUT("/tags/:id", v1.EditTag)
		apiv1.DELETE("/tags/:id", v1.DeleteTag)

		apiv1.GET("/articles", v1.GetArticles)
		apiv1.GET("/articles/:id", v1.GetArticle)
		apiv1.POST("/articles", v1.AddArticle)
		apiv1.PUT("/articles/:id", v1.EditArticle)
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	return r
}
