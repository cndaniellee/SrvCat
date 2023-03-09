package router

import (
	"SrvCat/controller"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/_benchmarks/iris-mvc-templates/controllers"
	"github.com/kataras/iris/v12/mvc"
)

func InitRouter(app *iris.Application) {
	app.Use(CrossAccess)

	tmpl := iris.HTML("./web", ".html").Layout("index.html")
	app.RegisterView(tmpl)
	app.HandleDir("/", "./web")
	mvc.New(app.Party("/")).Handle(&controllers.HomeController{})

	// 公用Controller
	mvc.New(app.Party("/api")).Handle(controller.NewApiController())
}

func CrossAccess(ctx iris.Context) {
	ctx.ResponseWriter().Header().Add("Access-Control-Allow-Origin", "*")
	ctx.Next()
}
