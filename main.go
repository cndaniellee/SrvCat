package main

import (
	"SrvCat/config"
	_ "SrvCat/forward"
	_ "SrvCat/logger"
	"SrvCat/router"
	"SrvCat/util"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

func main() {
	app := iris.New()
	// 中间件
	app.Use(recover.New())
	app.Use(logger.New())
	// 初始化路由
	router.InitRouter(app)
	// 添加配置
	app.Configure(iris.WithOptimizations)
	app.Configure(iris.WithoutServerError(iris.ErrServerClosed))
	// 启动项目
	err := app.Run(iris.Addr(config.Config.Server.Host + ":" + config.Config.Server.Port))
	util.FailOnException("An error occurred while starting the service", err)
}
