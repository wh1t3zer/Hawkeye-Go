package router

import (
	"log"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/wh1t3zer/Hawkeye-Go/controller"
	"github.com/wh1t3zer/Hawkeye-Go/docs"

	"github.com/wh1t3zer/Hawkeye-Go/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter ...
func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	docs.SwaggerInfo.Title = lib.GetStringConf("base.swagger.title")
	docs.SwaggerInfo.Description = lib.GetStringConf("base.swagger.desc")
	docs.SwaggerInfo.Version = lib.GetStringConf("base.swagger.version")
	docs.SwaggerInfo.Host = lib.GetStringConf("base.swagger.host")
	docs.SwaggerInfo.BasePath = lib.GetStringConf("base.swagger.base_path")
	// docs.SwaggerInfo.Schemes = strings.Split(lib.GetStringConf("base.swagger.schemes"), ",")
	docs.SwaggerInfo.Schemes = lib.GetStringSliceConf("base.swagger.schemes")

	router := gin.Default()
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	adminLoginRouter := router.Group("/admin_login")
	store, err := sessions.NewRedisStore(10, "tcp", lib.GetStringConf("base.session.redis_server"), lib.GetStringConf("base.session.redis_password"), []byte("secret"))
	if err != nil {
		log.Fatalf("sessions.NewRedisStore err:%v", err)
	}
	adminLoginRouter.Use(
		sessions.Sessions("mysession", store),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.TranslationMiddleware())
	{
		controller.AdminLoginRegister(adminLoginRouter)
	}

	adminRouter := router.Group("/admin")
	adminRouter.Use(
		sessions.Sessions("mysession", store),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.SessionAuthMiddleware(),
		middleware.TranslationMiddleware())
	{
		controller.AdminRegister(adminRouter)
	}

	// 大盘
	DashboardRouter := router.Group("/dashboard")
	DashboardRouter.Use(
		sessions.Sessions("mysession", store),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		// middleware.SessionAuthMiddleware(),
		middleware.TranslationMiddleware(),
	)
	{
		controller.DashboardRegister(DashboardRouter)
	}

	// 任务
	TaskRouter := router.Group("/task")
	TaskRouter.Use(
		sessions.Sessions("mysession", store),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		// middleware.SessionAuthMiddleware(),
		middleware.TranslationMiddleware(),
	)
	{
		controller.TaskRegister(TaskRouter)
	}

	// 资产管理
	AssetRouter := router.Group("/asset")
	AssetRouter.Use(
		sessions.Sessions("mysession", store),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		// middleware.SessionAuthMiddleware(),
		middleware.TranslationMiddleware(),
	)
	{
		controller.AssetRegister(AssetRouter)
	}

	// 漏洞管理
	VulRouter := router.Group("/vul")
	VulRouter.Use(
		sessions.Sessions("mysession", store),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.SessionAuthMiddleware(),
		middleware.TranslationMiddleware(),
	)
	{
		controller.VulRegister(VulRouter)
	}

	// 蜜罐识别
	TrapRouter := router.Group("/trap")
	TrapRouter.Use(
		sessions.Sessions("mysession", store),
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		middleware.SessionAuthMiddleware(),
		middleware.TranslationMiddleware(),
	)
	{
		controller.TrapRegister(TrapRouter)
	}

	// trojan
	trojanRouter := router.Group("/trojan")
	trojanRouter.Use(
		middleware.RecoveryMiddleware(),
		middleware.RequestLog(),
		// middleware.IPAuthMiddleware(),
		middleware.TranslationMiddleware(),
	)
	controller.TrojanRegister(trojanRouter)
	return router
}
