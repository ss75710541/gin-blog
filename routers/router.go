package routers

import (
	"gin-blog/middleware/jwt"
	"gin-blog/pkg/setting"
	"gin-blog/pkg/upload"
	"gin-blog/routers/api"
	v1 "gin-blog/routers/api/v1"
	v2 "gin-blog/routers/api/v2"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.ServerSetting.RunMode)

	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	r.StaticFS("/api/doc", http.Dir("docs/"))

	r.GET("/auth", api.GetAuth)

	// programatically set swagger info
	url := ginSwagger.URL("/api/doc/v1/swagger.json")
	r.GET("/swagger/v1/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,url))

	r.POST("/upload", api.UploadImage)

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{
		tags := apiv1.Group("/tags")
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
	url2 := ginSwagger.URL("/api/doc/v2/swagger.json")
	r.GET("/swagger/v2/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,url2))

	apiv2 := r.Group("/api/v2")
	apiv2.Use(jwt.JWT())
	{
		tags := apiv2.Group("/tags")
		{
			//获取标签列表
			tags.GET("", v2.GetTags)
			//新建标签
			tags.POST("", v2.AddTag)
			//更新指定标签
			tags.PUT(":id", v2.EditTag)
			//删除指定标签
			tags.DELETE(":id", v2.DeleteTag)
		}

		articles := apiv2.Group("/articles")
		{
			//获取文章列表
			articles.GET("", v2.GetArticles)
			//获取指定文章
			articles.GET(":id", v2.GetArticle)
			//新建文章
			articles.POST("", v2.AddArticle).Use(jwt.JWT())
			//更新指定文章
			articles.PUT(":id", v2.EditArticle).Use(jwt.JWT())
			//删除指定文章
			articles.DELETE(":id", v2.DeleteArticle).Use(jwt.JWT())
		}
	}

	return r
}
