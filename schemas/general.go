package schemas

import "github.com/google/uuid"

type IDPath struct {
	ID uuid.UUID `path:"id"`
}

type IDField struct {
	ID uuid.UUID `json:"id"`
}

type ReviewedField struct {
	Reviewed bool `json:"reviewed" doc:"is it reviewed?"`
}

type ReviewedField_Optional struct {
	Reviewed *bool `json:"reviewed" doc:"is it reviewed?" required:"false"`
}
