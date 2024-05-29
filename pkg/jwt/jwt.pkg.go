package jwt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"gin-example/entity"
	"gin-example/serializers"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"io"
	"os"
	"time"
)

type JwtService interface {
	CreateToken(user entity.User) serializers.Token
	ValidateToken(accessToken string) (entity.User, error)
	ValidateRefreshToken(model string) (entity.User, error)
	CreateRefreshToken(token serializers.Token) serializers.Token
}

type Jwt struct {
}

func (j Jwt) CreateToken(user entity.User) serializers.Token {
	var err error
	claims := jwt.MapClaims{}
	claims["id"] = user.Id
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwt := serializers.Token{}
	jwt.AccessToken, err = token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return jwt
	}
	return j.createRefreshToken(jwt)
}

func (Jwt) ValidateToken(accessToken string) (entity.User, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	user := entity.User{}

	if err != nil {
		return user, err
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		user.Username = payload["username"].(string)
		return user, nil
	}
	return user, errors.New("invalid token")
}

func (Jwt) ValidateRefreshToken(model serializers.RefreshToken) (entity.User, error) {
	sha1 := sha1.New()
	io.WriteString(sha1, os.Getenv("SECRET_KEY"))
	user := entity.User{}
	salt := string(sha1.Sum(nil))[0:16]
	block, err := aes.NewCipher([]byte(salt))
	if err != nil {
		return user, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return user, err
	}
	data, err := base64.URLEncoding.DecodeString(model.RefreshToken)
	if err != nil {
		return user, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return user, err
	}
	claims := jwt.MapClaims{}
	parser := jwt.Parser{}
	token, _, err := parser.ParseUnverified(string(plain), claims)
	if err != nil {
		return user, err
	}
	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return user, errors.New("invalid token")
	}
	user.Username = payload["username"].(string)
	return user, nil
}

func (Jwt) createRefreshToken(token serializers.Token) serializers.Token {
	sha1 := sha1.New()
	io.WriteString(sha1, os.Getenv("SECRET_KEY"))
	salt := string(sha1.Sum(nil))[0:16]
	block, err := aes.NewCipher([]byte(salt))
	if err != nil {
		panic(serializers.PanicMessage{400, gin.H{"fa_message": "نام کاربری یا رمز اشتباه است", "en_message": "The username or password is incorrect"}})
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(serializers.PanicMessage{400, gin.H{"fa_message": "نام کاربری یا رمز اشتباه است", "en_message": "The username or password is incorrect"}})
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		panic(serializers.PanicMessage{400, gin.H{"fa_message": "نام کاربری یا رمز اشتباه است", "en_message": "The username or password is incorrect"}})
	}
	token.RefreshToken = base64.URLEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(token.AccessToken), nil))
	return token
}
