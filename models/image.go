package models

type Image struct {
	Base
	Hash        string `gorm:"uniqueIndex;not null"` // SHA-256
	URL         string `gorm:"not null"`
	MimeType    string
	Size        int64
	Description string
}
