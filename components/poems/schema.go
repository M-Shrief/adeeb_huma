package poems

import (
	"adeeb_huma/database"
	"adeeb_huma/schemas"
)

type GetOnePoem_Req struct {
	schemas.IDPath
}

type GetOnePoem_Res struct {
	Body   GetOnePoem_Res_Body
	Status int
}
type GetOnePoem_Res_Body struct {
	schemas.Poem_Descriptive
	Adeeb        schemas.Adeeb_Minimal         `json:"adeeb"`
	ChosenVerses []schemas.ChosenVerse_Minimal `json:"chosen_verses"`
}

type CreateOnePoem_Req struct {
	Body CreateOnePoem_Req_Body
}

type CreateOnePoem_Req_Body struct {
	schemas.Poem_IntroField
	schemas.VersesField
	schemas.IsCoupletField
	schemas.ReviewedField

	// Relations
	schemas.AdeebIDField
}

type CreateOnePoem_Res struct {
	Body   schemas.Poem_Descriptive
	Status int
}

type CreateManyPoems_Req struct {
	Body []CreateOnePoem_Req_Body
}

type CreateManyPoems_Res struct {
	Body   schemas.CreateMany_Res_Body[schemas.Poem_Descriptive]
	Status int
}

type UpdatePoem_Req struct {
	schemas.IDPath
	Body UodatePoem_Req_Body
}

type UodatePoem_Req_Body struct {
	schemas.Poem_IntroField_Optional
	schemas.VersesField_Optional
	schemas.IsCoupletField_Optional
	schemas.ReviewedField_Optional

	// Relations
	schemas.AdeebIDField_Optional
}

type DeletePoem_Req struct {
	schemas.IDPath
}

func DBModel_To_ResModel(poem_model database.Poem) schemas.Poem_Descriptive {
	// Because we embed structs we can't assign values to fields directly without a lot of boilerplate.
	// So we define the variable and it's type, then assign each field alone.
	var poem_res schemas.Poem_Descriptive
	poem_res.ID = poem_model.ID
	poem_res.Intro = poem_model.Intro
	poem_res.Verses = poem_model.Verses
	poem_res.IsCouplet = poem_model.IsCouplet
	poem_res.Reviewed = poem_model.Reviewed

	poem_res.AdeebID = poem_model.AdeebID

	return poem_res
}

func DBModels_To_ResModels(poem_models []database.Poem) []schemas.Poem_Descriptive {
	var poems []schemas.Poem_Descriptive

	for _, adeeb_model := range poem_models {
		poem_res := DBModel_To_ResModel(adeeb_model)
		poems = append(poems, poem_res)
	}

	return poems
}

func ReqModel_To_DBModel(poem_req CreateOnePoem_Req_Body) database.Poem {
	poem_model := database.Poem{
		Intro:     poem_req.Intro,
		Verses:    poem_req.Verses,
		IsCouplet: poem_req.Reviewed,
		Reviewed:  poem_req.Reviewed,

		AdeebID: poem_req.AdeebID,
	}

	return poem_model
}

func ReqModels_To_DBModels(poems_req []CreateOnePoem_Req_Body) []database.Poem {
	var poem_models []database.Poem

	for _, poem_req := range poems_req {
		poem_model := ReqModel_To_DBModel(poem_req)
		poem_models = append(poem_models, poem_model)
	}

	return poem_models
}
