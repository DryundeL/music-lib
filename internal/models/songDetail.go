package models

import "time"

type SongDetail struct {
	ID          uint      `gorm:"primaryKey"`
	SongID      uint      `gorm:"unique;not null;index" json:"song_id"`
	Text        string    `gorm:"type:text;not null" json:"text"`
	ReleaseDate time.Time `gorm:"type:date;not null" json:"release_date"`
	Link        string    `gorm:"type:varchar(255)" json:"link,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
