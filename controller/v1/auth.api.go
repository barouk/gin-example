package v1

import (
	"gin-example/pkg/jwt"
	"gin-example/repository"
	"gin-example/serializers"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthAPI interface {
	Login(context *gin.Context)
	Register(context *gin.Context)
}

type authAPI struct {
	authrepository repository.AuthRepository
}

func NewAuthAPI(authrepository repository.AuthRepository) AuthAPI {
	return &authAPI{
		authrepository: authrepository,
	}
}

func (a authAPI) deferFunc(context *gin.Context) {
	if r := recover(); r != nil {
		PanicMessage, ok := r.(serializers.PanicMessage)
		if !ok {
			context.JSON(http.StatusBadRequest, gin.H{"fa_message": "خطایی پیش آمد ", "en_message": "error "})
			return
		}
		context.JSON(PanicMessage.Status, PanicMessage.Message)
	}
}

func (a authAPI) CheckSerializer(context *gin.Context, serializer interface{}) {
	err := context.ShouldBind(&serializer)
	if err != nil {
		panic(serializers.PanicMessage{400, gin.H{"fa_message": "خطا در داده های ورودی", "en_message": "Error in input data."}})

	}
}

func (a authAPI) Login(context *gin.Context) {
	defer a.deferFunc(context)
	var loginRequest serializers.LoginRequest
	a.CheckSerializer(context, &loginRequest)
	user := a.authrepository.LoginVerify(loginRequest)
	jwt := jwt.Jwt{}
	token := jwt.CreateToken(user)
	context.JSON(http.StatusOK, gin.H{"data": token})
	return
}

func (a authAPI) Register(context *gin.Context) {
	defer a.deferFunc(context)
	var serializer serializers.CreateUser
	a.CheckSerializer(context, &serializer)
	a.authrepository.CreateUser(serializer)
	context.JSON(http.StatusOK, gin.H{"fa_message": "کاربر با موفقیت ثبت شد", "en_message": "The user register"})
}
