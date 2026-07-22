package poems

import (
	"adeeb_huma/database"
	"adeeb_huma/schemas"
)

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
	Body   OnePoem_Res
	Status int
}

type OnePoem_Res struct {
	schemas.IDField
	schemas.Poem_IntroField
	schemas.VersesField
	schemas.IsCoupletField
	schemas.ReviewedField
	// Relation
	schemas.AdeebIDField
}

func DBModel_To_ResModel(poem_model database.Poem) OnePoem_Res {
	// Because we embed structs we can't assign values to fields directly without a lot of boilerplate.
	// So we define the variable and it's type, then assign each field alone.
	var poem_res OnePoem_Res
	poem_res.ID = poem_model.ID
	poem_res.Intro = poem_model.Intro
	poem_res.Verses = poem_model.Verses
	poem_res.IsCouplet = poem_model.IsCouplet
	poem_res.Reviewed = poem_model.Reviewed

	poem_res.AdeebID = poem_model.AdeebID

	return poem_res
}
