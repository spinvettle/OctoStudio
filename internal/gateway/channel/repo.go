package channel

import "gorm.io/gorm"

type ChannelKeyStatus int

const (
	Enable ChannelKeyStatus = iota + 1
	Disable
	Colding
	Refreshing
)

type Channel struct {
	ID      int          `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	Name    string       `json:"name,omitempty"`
	BaseURL string       `json:"base_url,omitempty"`
	Keys    []ChannelKey `json:"keys,omitempty" gorm:"foreignKey:ChannelID;references:ID"`
	Models  []Model      `json:"models,omitempty" gorm:"many2many:channel_model;"`
}
type Model struct {
	ID       int    `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name,omitempty"`
	Provider string `json:"provider,omitempty"`
	Metadata string `json:"meatadata,omitempty"`
}

type ChannelModel struct {
	ChannelID int `json:"channel_id,omitempty" gorm:"primaryKey"`
	ModelID   int `json:"model_id,omitempty" gorm:"primaryKey"`
}
type ChannelKey struct {
	ID        int                `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	ChannelID int                `json:"channel_id,omitempty"`
	Name      string             `json:"name,omitempty"`
	Metadata  ChannelKeyMetaData `json:"meatadata,omitempty" gorm:"serializer:json"`
	ApiKey    string             `json:"api_key,omitempty"`
	Status    ChannelKeyStatus   `json:"status,omitempty"`
}

type ChannelKeyMetaData struct {
	CanRefresh            bool
	RefreshRequestPayload struct {
		IdClient     string
		GrantType    string
		RefreshToken string
	}
	Exp            int64
	RefreshBaseURL string
}

type ChannelRepo struct {
	DB *gorm.DB
}

func NewChannelRepo(db *gorm.DB) *ChannelRepo {
	return &ChannelRepo{
		DB: db,
	}
}

func (c *ChannelRepo) GetAllChannel() (*[]Channel, error) {
	var allChannels []Channel
	err := c.DB.Preload("keys").Preload("models").Find(&allChannels).Error
	return &allChannels, err
}

func (c *ChannelRepo) GetChannelByName(name string) (*Channel, error) {
	var ch Channel
	err := c.DB.Where("name = ?", name).Preload("Keys").Take(&ch).Error
	if err != nil {
		return nil, err
	}
	return &ch, nil
}

func (c *ChannelRepo) GetChannelByID(id int) (*Channel, error) {
	var ch Channel
	err := c.DB.Preload("Keys").First(&ch, id).Error
	if err != nil {
		return nil, err
	}
	return &ch, nil

}

func (c *ChannelRepo) AddChannel(ch *Channel) (int, error) {
	res := c.DB.Create(ch)
	if res.Error != nil {
		return 0, nil
	}
	return int(res.RowsAffected), nil

}

func (c *ChannelRepo) AddChannelKey(channelKey *ChannelKey) (int, error) {
	res := c.DB.Create(channelKey)
	if res.Error != nil {
		return 0, nil
	}
	return int(res.RowsAffected), nil
}

func (c *ChannelRepo) DeleteChannel(channelID int) error {
	err := c.DB.Delete(&Channel{}, channelID).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *ChannelRepo) UpdateChannelKey(channelKey *ChannelKey) error {
	return c.DB.Save(channelKey).Error

}
func (c *ChannelRepo) UpdateChannelKeyStatus(id int, status int) error {
	return c.DB.Model(&ChannelKey{}).Where("id = ?", id).
		Update("status", status).Error
}

func (c *ChannelRepo) UpdateChannel(ch *Channel) error {
	return c.DB.Save(ch).Error

}
