package routes

import (
	v1 "gin-example/controller/v1"
	"gin-example/db"
	"gin-example/middleware"
	"gin-example/pkg/jwt"
	"gin-example/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	Db *gorm.DB = db.ConnectDb()

	jwtAuth        jwt.Jwt
	authRepository repository.AuthRepository = repository.NewAuthRepository(Db)

	authAPI v1.AuthAPI = v1.NewAuthAPI(authRepository)
	userAPI v1.UserAPI = v1.NewUserAPI(authRepository)
)

func Urls(r *gin.Engine) *gin.Engine {
	apiV1 := r.Group("api/v1/")
	{
		auth := apiV1.Group("auth")
		{
			auth.POST("/login/", authAPI.Login)
			auth.POST("/create/User/", authAPI.Register)
		}

		user := apiV1.Group("user", middleware.AthorizationJWT(jwtAuth))
		{
			user.GET("/info/", userAPI.Info)
		}

	}

	return r
}
