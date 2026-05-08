package repo

type Channel struct {
	ID      int          `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	Name    string       `json:"name,omitempty"`
	BaseURL string       `json:"base_url,omitempty"`
	Keys    []ChannelKey `json:"keys,omitempty" gorm:"foreignKey:ChannelID;references:ID"`
}

type ChannelKey struct {
	ID        int    `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	ChannelID int    `json:"channel_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Metadata  string `json:"meatadata,omitempty"`
	ApiKey    string `json:"api_key,omitempty"`
	Status    string `json:"status,omitempty"`
}
