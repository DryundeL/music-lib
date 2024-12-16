package models

import "time"

type Song struct {
	ID         uint       `gorm:"primaryKey"`
	Name       string     `gorm:"not null" json:"name"`
	ArtistID   uint       `gorm:"not null;index" json:"artist_id"`
	Artist     Artist     `gorm:"foreignKey:ArtistID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"artist,omitempty"`
	SongDetail SongDetail `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"song_detail,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
