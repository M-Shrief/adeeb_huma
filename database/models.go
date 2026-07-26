package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func Migrate() {
	Conn.AutoMigrate(
		&User{},
		&Adeeb{},
		&Poem{},
		&ChosenVerse{},
		&ProseQoute{},
		&Order{},
		&Print{},
	)
}

type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp without time zone" json:"-"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp without time zone" json:"-"`
}

type User struct {
	BaseModel
	Username  string      `gorm:"size:256;not null;unique"`
	Passsword string      `gorm:"size:256;not null"`
	Roles     *[]RoleEnum `gorm:"type:role_enum"`

	// Relations
	Orders []Order `gorm:"foreignKey:UserID;references:ID"`
}

type Adeeb struct {
	BaseModel
	Name       string         `gorm:"size:256;not null;unique" json:"name"`
	TimePeriod TimePeriodEnum `gorm:"type:time_period_enum; not null" json:"time_period"`
	Bio        *string        `gorm:"size:1024" json:"bio"`
	Reviewed   bool           `gorm:"default:false" json:"reviewed"`

	// Relations
	// We explicitly define foreignKey & references,
	// and we don't rely on convention only

	Poems        []Poem        `gorm:"foreignKey:AdeebID;references:ID" json:"poems,omitempty"`
	ChosenVerses []ChosenVerse `gorm:"foreignKey:AdeebID;references:ID" json:"chosen_verses,omitempty"`
	ProseQoutes  []ProseQoute  `gorm:"foreignKey:AdeebID;references:ID" json:"prose_qoutes,omitempty"`
}

type Poem struct {
	BaseModel
	Intro     string         `gorm:"size:256;not null;unique" json:"intro"`
	Verses    pq.StringArray `gorm:"type:varchar(256)[]" json:"verses"`
	IsCouplet bool           `gorm:"default:true" json:"is_couplet"`
	Reviewed  bool           `gorm:"default:false" json:"reviewed"`

	// Relations
	AdeebID uuid.UUID `json:"adeeb_id"`
	Adeeb   Adeeb     `gorm:"foreignKey:AdeebID;references:ID" json:"adeeb"`

	ChosenVerses []ChosenVerse `gorm:"foreignKey:PoemID;references:ID" json:"chosen_verses,omitempty"`
}

type ChosenVerse struct {
	BaseModel
	Tags      pq.StringArray `gorm:"type:varchar(64)[];default:'{}'" json:"tags"`
	Verses    pq.StringArray `gorm:"type:varchar(256)[]" json:"verses"`
	IsCouplet bool           `gorm:"default:true" json:"is_couplet"`
	Reviewed  bool           `gorm:"default:false" json:"reviewed"`

	// Relations
	AdeebID uuid.UUID `json:"adeeb_id"`
	Adeeb   Adeeb     `gorm:"foreignKey:AdeebID;references:ID" json:"adeeb,omitempty"`

	PoemID uuid.UUID `json:"poem_id"`
	Poem   Poem      `gorm:"foreignKey:PoemID;references:ID" json:"poem"`
}

type ProseQoute struct {
	BaseModel
	Tags     pq.StringArray `gorm:"type:varchar(64)[];default:'{}'" json:"tags"`
	Qoute    string         `gorm:"size:512;not null" json:"qoute"`
	Source   *string        `gorm:"size:128" json:"source"`
	Reviewed bool           `gorm:"default:false" json:"reviewed"`

	// Relations
	AdeebID uuid.UUID `json:"adeeb_id"`
	Adeeb   Adeeb     `gorm:"foreignKey:AdeebID;references:ID" json:"adeeb,omitempty"`
}

type Order struct {
	BaseModel
	Name             string     `gorm:"size:256;not null" json:"name"`
	Phone            string     `gorm:"size:128;not null" json:"phone"`
	Address          string     `gorm:"size:256;not null" json:"address"`
	Reviewed         bool       `gorm:"default:false;not null" json:"reviewed"`
	IsUpdateable     bool       `gorm:"default:true;not null" json:"is_updateable"`
	Status           StatusEnum `gorm:"type:status_enum; not null" json:"status"`
	DeliverySchedule *time.Time `gorm:"type:timestamp without time zone;not null" json:"delivery_schedule"`

	// Relations
	UserID uuid.UUID `json:"user_id"`
	User   User      `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`

	Prints []Print `gorm:"foreignKey:OrderID;references:ID" json:"prints"`
}

type Print struct {
	BaseModel
	FontType    string         `gorm:"size:64;not null" json:"font_type"`
	FontColor   string         `gorm:"size:64;not null" json:"font_color"`
	OutfitColor string         `gorm:"size:64;not null" json:"outfit_color"`
	OutfitType  OutfitTypeEnum `gorm:"type:outfit_type_enum; not null" json:"outfit_type"`

	Qoute     *string         `gorm:"size:512" json:"qoute"`
	Verses    *pq.StringArray `gorm:"type:varchar(256)[]" json:"verses"`
	IsCouplet *bool           `gorm:"default:true" json:"is_couplet"`

	// Relations
	OrderID uuid.UUID `json:"order_id"`
	Order   Order     `gorm:"foreignKey:OrderID;references:ID" json:"order,omitempty"`

	UserID uuid.UUID `json:"user_id"`
	User   User      `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`

	// // Only for Analytics
	PoemID        *uuid.UUID `json:"poem_id,omitempty"`
	ChosenVerseID *uuid.UUID `json:"chosen_verse_id,omitempty"`
	ProseQouteID  *uuid.UUID `json:"prose_qoute_id,omitempty"`
}
