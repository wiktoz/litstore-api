package models

type Image struct {
	Base
	Hash        string `gorm:"uniqueIndex;not null" json:"-"` // SHA-256
	URL         string `gorm:"not null" json:"url"`
	MimeType    string `json:"mime_type"`
	Size        int64  `json:"size"`
	Description string `json:"description,omitempty"`
}
