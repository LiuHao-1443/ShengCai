package model

import "time"

type SheetInfo struct {
	ID           int       `gorm:"primaryKey;not null" json:"-"`
	SheetID      string    `gorm:"column:sheet_id;type:varchar(255)" json:"sheet_id"`
	SheetName    string    `gorm:"column:sheet_name;type:varchar(255)" json:"sheet_name"`
	UpdateLog    string    `gorm:"column:update_log;type:varchar(255)" json:"update_log"`
	CreatedAt    time.Time `gorm:"column:create_at;type:timestamp;not null;autoCreateTime" json:"-"`
	UpdatedAt    time.Time `gorm:"column:update_at;type:timestamp;not null;autoUpdateTime" json:"-"`
	Deleted      bool      `gorm:"column:deleted;type:tinyint(4);default:0" json:"-"`
	RuntimeState string    `gorm:"column:runtime_state;type:varchar(255)" json:"runtime_state"`
}

func (SheetInfo) TableName() string {
	return "sheet_info"
}
