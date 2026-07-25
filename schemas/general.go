package schemas

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type IDPath struct {
	ID uuid.UUID `path:"id"`
}

type IDField struct {
	ID uuid.UUID `json:"id"`
}

type VersesField struct {
	Verses pq.StringArray `json:"verses" minItems:"1" gorm:"type:varchar(256)[]" doc:"Verses"`
}

type VersesField_Optional struct {
	Verses *pq.StringArray `json:"verses" minItems:"1" gorm:"type:varchar(256)[]" doc:"Verses" required:"false"`
}

type TagsField struct {
	Tags pq.StringArray `json:"tags" gorm:"type:varchar(64)[];default:'{}'" doc:"Tags"`
}

type TagsField_Optional struct {
	Tags *pq.StringArray `json:"tags" gorm:"type:varchar(64)[];default:'{}'" doc:"Tags" required:"false"`
}

type IsCoupletField struct {
	IsCouplet bool `json:"is_couplet" doc:"is it couplet?"`
}

type IsCoupletField_Optional struct {
	IsCouplet *bool `json:"reviis_coupletewed" doc:"is it couplet?" required:"false"`
}

type ReviewedField struct {
	Reviewed bool `json:"reviewed" doc:"is it reviewed?"`
}

type ReviewedField_Optional struct {
	Reviewed *bool `json:"reviewed" doc:"is it reviewed?" required:"false"`
}

type CreatedAtField struct {
	CreatedAt time.Time `json:"created_at"`
}
type UpdatedAtField struct {
	UpdatedAt time.Time `json:"updated_at"`
}

// Relations' Fields

type AdeebIDField struct {
	AdeebID uuid.UUID `json:"adeeb_id"`
}

type AdeebIDField_Optional struct {
	AdeebID *uuid.UUID `json:"adeeb_id" nullable:"true" required:"false"`
}

type PoemIDField struct {
	PoemID uuid.UUID `json:"poem_id"`
}

type PoemIDField_Optional struct {
	PoemID *uuid.UUID `json:"poem_id" nullable:"true" required:"false"`
}
