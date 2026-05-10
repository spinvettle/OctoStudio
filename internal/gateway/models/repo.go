package model

type Model struct {
	ID       int    `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name,omitempty"`
	Provider string `json:"provider,omitempty"`
	Metadata string `json:"meatadata,omitempty"`
}
