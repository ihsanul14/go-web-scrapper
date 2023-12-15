package model

type Data struct {
	Id          string `gorm:"id,primaryKey"`
	Name        string `gorm:"name,omitempty"`
	Description string `gorm:"description,omitempty"`
	ImageLink   string `gorm:"image_link,omitempty"`
	Price       string `gorm:"price,omitempty"`
	Ratings     string `gorm:"ratings,omitempty"`
	MerchatName string `gorm:"merchant_name,omitempty"`
}
