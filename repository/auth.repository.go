package repository

import (
	"errors"
	"gin-example/entity"
	"gin-example/serializers"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/mashingan/smapping"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthRepository interface {
	LoginVerify(loginRequest serializers.LoginRequest) entity.User
	CreateUser(registerRequest serializers.CreateUser)
}

type authRepository struct {
	conn *gorm.DB
}

func NewAuthRepository(connection *gorm.DB) AuthRepository {
	return &authRepository{
		conn: connection,
	}
}

func comparePasswords(hashedPwd string, plainPwd []byte) {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		panic(serializers.PanicMessage{400, gin.H{"fa_message": "نام کاربری یا رمز اشتباه است", "en_message": "The username or password is incorrect"}})
	}
}

func (a authRepository) LoginVerify(loginRequest serializers.LoginRequest) entity.User {
	var user entity.User
	res := a.conn.Where("username = ? ", loginRequest.Username).Take(&user)
	if res.Error != nil {
		panic(serializers.PanicMessage{400, gin.H{"fa_message": "نام کاربری یا رمز اشتباه است", "en_message": "The username or password is incorrect"}})
	}
	comparePasswords(user.Password, []byte(loginRequest.Password))
	return user

}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		panic(serializers.PanicMessage{400, gin.H{"fa_message": "خطا در ایجاد یوزر", "en_message": "error in create user"}})
	}
	return string(hash)
}

func (a authRepository) CreateUser(registerRequest serializers.CreateUser) {
	var mysqlErr *mysql.MySQLError
	var user entity.User
	mapped := smapping.MapFields(registerRequest)
	err := smapping.FillStruct(&user, mapped)
	if err != nil {
		panic(serializers.PanicMessage{400, gin.H{"fa_message": "خطا در ایجاد کاربر", "en_message": "error in create user"}})
	}
	user.Password = hashAndSalt([]byte(registerRequest.Password))
	if err := a.conn.Create(&user).Error; err != nil {
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			panic(serializers.PanicMessage{400, gin.H{"fa_message": "کاربر از قبل وجود دارد", "en_message": "user is exsist"}})
		}
		panic(serializers.PanicMessage{400, gin.H{"fa_message": "خطا در ایجاد یوزر", "en_message": "error in create user"}})
	}
}
