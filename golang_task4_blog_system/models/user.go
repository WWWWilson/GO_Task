package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string `gorm:"size:50;uniqueIndex;not null" json:"username" binding:"required"`
	Email    string `gorm:"size:100;uniqueIndex;not null" json:"email" binding:"required,email"`
	Password string `gorm:"size:255;not null" json:"-" binding:"required,min=6"` // json:"-" 表示不序列化到JSON

	// 关联关系
	Posts    []Post    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"`
	Comments []Comment `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"`
}

// HashPassword 加密密码
func (u *User) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
