package routers

import (
	"gin-blog/middleware/jwt"
	"gin-blog/pkg/setting"
	"gin-blog/routers/api"
	v1 "gin-blog/routers/api/v1"
	"github.com/gin-gonic/gin"

	"gin-blog/docs" // docs is generated by Swag CLI, you have to import it.
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	r.GET("/auth", api.GetAuth)

	apiv1 := r.Group("/api/v1")
	{
		tags := apiv1.Group("/tags").Use(jwt.JWT())
		{
			//获取标签列表
			tags.GET("", v1.GetTags)
			//新建标签
			tags.POST("", v1.AddTag)
			//更新指定标签
			tags.PUT(":id", v1.EditTag)
			//删除指定标签
			tags.DELETE(":id", v1.DeleteTag)
		}

		articles := apiv1.Group("/articles")
		{
			//获取文章列表
			articles.GET("", v1.GetArticles)
			//获取指定文章
			articles.GET(":id", v1.GetArticle)
			//新建文章
			articles.POST("", v1.AddArticle).Use(jwt.JWT())
			//更新指定文章
			articles.PUT(":id", v1.EditArticle).Use(jwt.JWT())
			//删除指定文章
			articles.DELETE(":id", v1.DeleteArticle).Use(jwt.JWT())
		}
	}

	// programatically set swagger info
	docs.SwaggerInfo.Title = "gin-blog API"
	docs.SwaggerInfo.Description = "gin-blog api"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
