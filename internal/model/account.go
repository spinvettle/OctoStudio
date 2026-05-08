package model

type Account struct {
	ID           string `gorm:"primaryKey;column:id"`
	Name         string `gorm:"name"`
	AccessToken  string `gorm:"access_token"`
	RefreshToken string `gotm:"refresh_token"`
	TokenExp     int64
	Status       int
	UsagePercent float64
}
