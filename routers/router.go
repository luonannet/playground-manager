package routers

import (
	"github.com/omnibuildplatform/omni-manager/controllers"
	"github.com/omnibuildplatform/omni-manager/docs"
	"github.com/omnibuildplatform/omni-manager/models"
	"github.com/omnibuildplatform/omni-manager/util"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

//InitRouter init router
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(util.LoggerToFile())
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Title = util.GetConfig().AppName
	docs.SwaggerInfo.Description = "set token name: 'Authorization' at header "
	auth := r.Group(docs.SwaggerInfo.BasePath)
	{
		auth.GET("/v1/auth/loginok", controllers.AuthingLoginOk)
		auth.GET("/v1/auth/getDetail/:authingUserId", controllers.AuthingGetUserDetail)
		auth.Use(models.Authorize()) //
		auth.POST("/v1/auth/createUser", controllers.AuthingCreateUser)
	}
	//version 1 . call k8s api
	v1 := r.Group(docs.SwaggerInfo.BasePath + "/v1")
	{
		v1.Use(models.Authorize()) //

	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}
