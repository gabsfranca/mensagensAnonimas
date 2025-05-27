package models

type Admin struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"uniqueIndex;nol null"`
	Password string `gorm:"not null"`
}
