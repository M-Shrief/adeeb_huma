package chosen_verses

import (
	"adeeb_huma/database"
	"adeeb_huma/schemas"
)

type CreateOneChosenVerse_Req struct {
	Body CreateOneChosenVerse_Req_Body
}

type CreateOneChosenVerse_Req_Body struct {
	schemas.VersesField
	schemas.IsCoupletField
	schemas.TagsField
	schemas.ReviewedField

	// Relations
	schemas.AdeebIDField
	schemas.PoemIDField
}

type CreateOneChosenVerse_Res struct {
	Body   schemas.ChosenVerse_Descriptive
	Status int
}

type CreateManyChosenVerses_Req struct {
	Body []CreateOneChosenVerse_Req_Body
}

type CreateManyChosenVerses_Res struct {
	Body   schemas.CreateMany_Res_Body[schemas.ChosenVerse_Descriptive]
	Status int
}

func DBModel_To_ResModel(chosen_verse_model database.ChosenVerse) schemas.ChosenVerse_Descriptive {
	// Because we embed structs we can't assign values to fields directly without a lot of boilerplate.
	// So we define the variable and it's type, then assign each field alone.
	var chosen_verse_res schemas.ChosenVerse_Descriptive
	chosen_verse_res.ID = chosen_verse_model.ID
	chosen_verse_res.Verses = chosen_verse_model.Verses
	chosen_verse_res.IsCouplet = chosen_verse_model.IsCouplet
	chosen_verse_res.Tags = chosen_verse_model.Tags
	chosen_verse_res.Reviewed = chosen_verse_model.Reviewed

	chosen_verse_res.AdeebID = chosen_verse_model.AdeebID
	chosen_verse_res.PoemID = chosen_verse_model.PoemID

	return chosen_verse_res
}

func DBModels_To_ResModels(chosen_verse_models []database.ChosenVerse) []schemas.ChosenVerse_Descriptive {
	var poems []schemas.ChosenVerse_Descriptive

	for _, adeeb_model := range chosen_verse_models {
		chosen_verse_res := DBModel_To_ResModel(adeeb_model)
		poems = append(poems, chosen_verse_res)
	}

	return poems
}

func ReqModel_To_DBModel(chosen_verse_req CreateOneChosenVerse_Req_Body) database.ChosenVerse {
	chosen_verse_model := database.ChosenVerse{
		Verses:    chosen_verse_req.Verses,
		IsCouplet: chosen_verse_req.Reviewed,
		Tags:      chosen_verse_req.Tags,
		Reviewed:  chosen_verse_req.Reviewed,

		AdeebID: chosen_verse_req.AdeebID,
		PoemID:  chosen_verse_req.PoemID,
	}

	return chosen_verse_model
}

func ReqModels_To_DBModels(chosen_verses_req []CreateOneChosenVerse_Req_Body) []database.ChosenVerse {
	var chosen_verse_models []database.ChosenVerse

	for _, chosen_verse_req := range chosen_verses_req {
		chosen_verse_model := ReqModel_To_DBModel(chosen_verse_req)
		chosen_verse_models = append(chosen_verse_models, chosen_verse_model)
	}

	return chosen_verse_models
}
