package schemas

import (
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
	Verses *pq.StringArray `json:"verses" minItems:"1" gorm:"type:varchar(256)[]" doc:"Verses" required:"false"` // nullable:"true"
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

// Relations' Fields

type AdeebIDField struct {
	AdeebID uuid.UUID `json:"adeeb_id"`
}

type AdeebIDField_Optional struct {
	AdeebID *uuid.UUID `json:"adeeb_id" required:"false"`
}
