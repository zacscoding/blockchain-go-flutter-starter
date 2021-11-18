package wallet

import "time"

type Key struct {
	ID     uint   `gorm:"column:id;"`
	UserID uint   `gorm:"column:user_id;"`
	Tag    string `gorm:"column:tag;"`

	CreatedAt time.Time `gorm:"column:created_at;"`
	UpdatedAt time.Time `gorm:"column:updated_at;"`
}
