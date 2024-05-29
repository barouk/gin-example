package serializers

type CreateUser struct {
	Username  string `json:"username" form:"username" binding:"required"`
	Firstname string `json:"firstname" form:"firstname" binding:"required"`
	Lastname  string `json:"lastname" form:"lastname" binding:"required"`
	Password  string `json:"password" form:"password" binding:"required,min=4"`
}
