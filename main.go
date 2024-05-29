package main

import (
	"gin-example/middleware"
	"gin-example/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.NoRoute(middleware.NoRouteHandler())
	r.HandleMethodNotAllowed = true
	r.NoMethod(middleware.NoMethodHandler())

	routes.Urls(r)
	//other urls

	gin.SetMode(gin.ReleaseMode)
	err := r.Run("0.0.0.0:8000")
	if err != nil {
		return
	}
}
