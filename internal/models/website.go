package models

import "time"

type WebsiteDnsAccount struct {
	Id            uint      `gorm:"primarykey;AUTO_INCREMENT" json:"id"`
	Name          string    `gorm:"type:varchar(64);not null" json:"name"`
	Type          string    `gorm:"type:varchar(64);not null" json:"type"`
	Authorization string    `gorm:"type:varchar(256);not null" json:"-"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type WebsiteAcmeAccount struct {
	Id         uint      `gorm:"primarykey;AUTO_INCREMENT" json:"id"`
	Email      string    `gorm:"not null" json:"email"`
	URL        string    `gorm:"not null" json:"url"`
	PrivateKey string    `gorm:"not null" json:"-"`
	Type       string    `gorm:"not null;default:letsencrypt" json:"type"`
	EabKid     string    `gorm:"default:null;" json:"eabKid"`
	EabHmacKey string    `gorm:"default:null" json:"eabHmacKey"`
	KeyType    string    `gorm:"not null;default:2048" json:"keyType"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// type WebsiteCA struct {
// 	Id         uint   `gorm:"primarykey;AUTO_INCREMENT" json:"id"`
// 	CSR        string `gorm:"not null;" json:"csr"`
// 	Name       string `gorm:"not null;" json:"name"`
// 	PrivateKey string `gorm:"not null" json:"privateKey"`
// 	KeyType    string `gorm:"not null;default:2048" json:"keyType"`
// }

func (w WebsiteAcmeAccount) TableName() string {
	return "website_acme_accounts"
}

func (w WebsiteDnsAccount) TableName() string {
	return "website_dns_accounts"
}
