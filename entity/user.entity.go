package entity

type User struct {
	Id        uint   `gorm:"primaryKey; autoIncrement:true; not null; unique;"`
	Username  string `gorm:"not null;uniqueIndex;type:varchar(255) CHARACTER SET utf8 COLLATE utf8_persian_ci"`
	Firstname string `gorm:"not null;type:varchar(100) CHARACTER SET utf8 COLLATE utf8_persian_ci"`
	Lastname  string `gorm:"not null;type:varchar(100) CHARACTER SET utf8 COLLATE utf8_persian_ci"`
	Password  string `gorm:"not null;type:text CHARACTER SET utf8 COLLATE utf8_persian_ci"`
}

func (User) TableName() string {
	return "user"
}
