package serializers

import "github.com/gin-gonic/gin"

type PanicMessage struct {
	Status  int
	Message gin.H
}
