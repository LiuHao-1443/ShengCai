package model

import "time"

type CellData struct {
	ID           int       `gorm:"primaryKey;not null" json:"-"`
	SheetID      string    `gorm:"column:sheet_id;type:varchar(100)" json:"sheet_id"`
	Title        string    `gorm:"column:title;type:varchar(255)" json:"title"`
	Link         string    `gorm:"column:link;type:varchar(255)" json:"link"`
	ReleaseDate  string    `gorm:"column:release_date;type:varchar(50)" json:"release_date"`
	Content      string    `gorm:"column:content;type:text" json:"content"`
	Abstract     string    `gorm:"column:abstract;type:varchar(1000)" json:"abstract"`
	Keyword      string    `gorm:"column:keyword;type:varchar(1000)" json:"keyword"`
	SortNumber   int       `gorm:"column:sort_number;type:int" json:"sort_number"`
	CreatedAt    time.Time `gorm:"column:create_at;type:timestamp;not null;autoCreateTime" json:"-"`
	UpdatedAt    time.Time `gorm:"column:update_at;type:timestamp;not null;autoUpdateTime" json:"-"`
	Deleted      bool      `gorm:"column:deleted;type:tinyint(4);default:0" json:"-"`
	RuntimeState string    `gorm:"column:runtime_state;type:varchar(255)" json:"runtime_state"`
}

func (CellData) TableName() string {
	return "cell_data"
}
