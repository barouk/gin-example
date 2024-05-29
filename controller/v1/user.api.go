package v1

import (
	"fmt"
	"gin-example/repository"
	"gin-example/serializers"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserAPI interface {
	Info(context *gin.Context)
}

type userAPI struct {
	authrepository repository.AuthRepository
}

func NewUserAPI(authrepository repository.AuthRepository) *userAPI {
	return &userAPI{
		authrepository: authrepository,
	}
}
func (a userAPI) GetUser(context *gin.Context) string {
	user, exists := context.Get("username")

	fmt.Println(user)
	if !exists {
		panic(serializers.PanicMessage{400, gin.H{"en_message": "User not authenticated", "fa_message": "کاربر احراز هویت نشد"}})
	}
	return user.(string)
}

func (a userAPI) deferFunc(context *gin.Context) {
	if r := recover(); r != nil {
		PanicMessage, ok := r.(serializers.PanicMessage)
		if !ok {
			context.JSON(http.StatusBadRequest, gin.H{"fa_message": "خطایی پیش آمد ", "en_message": "error "})
			return
		}
		context.JSON(PanicMessage.Status, PanicMessage.Message)
	}
}

func (a userAPI) Info(context *gin.Context) {
	defer a.deferFunc(context)
	user := a.GetUser(context)
	context.JSON(http.StatusOK, gin.H{"user": user})
	return
}
