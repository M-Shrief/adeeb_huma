package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Migrate() {
	Conn.AutoMigrate(
		&User{},
		&Adeeb{},
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

type User struct {
	BaseModel
	Username  string      `gorm:"size:256;not null;unique"`
	Passsword string      `gorm:"size:256;not null"`
	Roles     *[]RoleEnum `gorm:"type:role_enum"`
}

type Adeeb struct {
	BaseModel
	Name       string         `gorm:"size:256;not null;unique"`
	TimePeriod TimePeriodEnum `gorm:"type:time_period_enum; not null"`
	Bio        *string        `gorm:"size:1024"`
	Reviewed   bool           `gorm:"default:false"`
}
