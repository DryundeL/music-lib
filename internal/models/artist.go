package models

type Artist struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"unique;not null;index" json:"name"`
	IsGroup bool   `json:"is_group"`
	Songs   []Song `gorm:"foreignKey:ArtistID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"songs,omitempty"`
}
