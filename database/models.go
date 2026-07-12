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
		&Poem{},
		&ChosenVerse{},
		&ProseQoute{},
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

	// Relations
	// We explicitly define foreignKey & references,
	// and we don't rely on convention only

	Poems        []Poem        `gorm:"foreignKey:AdeebID;references:ID"`
	ChosenVerses []ChosenVerse `gorm:"foreignKey:AdeebID;references:ID"`
	ProseQoutes  []ProseQoute  `gorm:"foreignKey:AdeebID;references:ID"`
}

type Poem struct {
	BaseModel
	Intro     string   `gorm:"size:256;not null;unique"`
	Verses    []string `gorm:"type:varchar(256)[]"`
	IsCouplet bool     `gorm:"default:true"`
	Reviewed  bool     `gorm:"default:false"`

	// Relations
	AdeebID uuid.UUID
	Adeeb   Adeeb `gorm:"foreignKey:AdeebID;references:ID"`
}

type ChosenVerse struct {
	BaseModel
	Tags      []string `gorm:"type:varchar(256)[];default:'{}'"`
	Verses    []string `gorm:"type:varchar(256)[]"`
	IsCouplet bool     `gorm:"default:true"`
	Reviewed  bool     `gorm:"default:false"`

	// Relations
	AdeebID uuid.UUID
	Adeeb   Adeeb `gorm:"foreignKey:AdeebID;references:ID"`

	PoemID uuid.UUID
	Poem   Poem `gorm:"foreignKey:PoemID;references:ID"`
}

type ProseQoute struct {
	BaseModel
	Tags     []string `gorm:"type:varchar(256)[];default:'{}'"`
	Qoute    string   `gorm:"size:512;not null"`
	Source   *string  `gorm:"size:128"`
	Reviewed bool     `gorm:"default:false"`

	// Relations
	AdeebID uuid.UUID
	Adeeb   Adeeb `gorm:"foreignKey:AdeebID;references:ID"`
}
