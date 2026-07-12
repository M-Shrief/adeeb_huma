package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Migrate() {
	Conn.AutoMigrate(
		&UserModel{},
	)
}

type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp without time zone"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp without time zone"`
}

// BeforeCreate hook generates UUID before insertion
func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}

type UserModel struct {
	BaseModel
	Username  string      `gorm:"size:256;not null;unique"`
	Passsword string      `gorm:"size:256;not null"`
	Roles     *[]RoleEnum `gorm:"type:public.role_enum"`
}
