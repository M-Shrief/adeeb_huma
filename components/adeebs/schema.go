package adeebs

import (
	"adeeb_huma/database"
	"adeeb_huma/schemas"
)

type GetOneAdeeb_Req struct {
	schemas.IDPath
}

type GetOneAdeeb_Res struct {
	Body   schemas.Adeeb_Descriptive
	Status int
}

type CreateOneAdeeb_Req struct {
	Body CreateOneAdeeb_Req_Body
}

type CreateOneAdeeb_Req_Body struct {
	schemas.Adeeb_NameField
	schemas.Adeeb_TimePeriodField
	schemas.Adeeb_BioField
	schemas.ReviewedField
}

type CreateOneAdeeb_Res struct {
	Body   schemas.Adeeb_Descriptive
	Status int
}

type CreateManyAdeebs_Req struct {
	Body []CreateOneAdeeb_Req_Body
}

type CreateManyAdeebs_Res struct {
	Body   schemas.CreateMany_Res_Body[schemas.Adeeb_Descriptive]
	Status int
}

type UpdateAdeeb_Req struct {
	schemas.IDPath
	Body UpdateAdeeb_Req_Body
}

type UpdateAdeeb_Req_Body struct {
	schemas.Adeeb_NameField_Optional
	schemas.Adeeb_BioField_Optional
	schemas.Adeeb_TimePeriodField_Optional
	schemas.ReviewedField_Optional
}

type DeleteAdeeb_Req struct {
	schemas.IDPath
}

func DBModel_To_ResSchema(adeeb_model database.Adeeb) schemas.Adeeb_Descriptive {
	// Because we embed structs we can't assign values to fields directly without a lot of boilerplate.
	// So we define the variable and it's type, then assign each field alone.
	var adeeb_res schemas.Adeeb_Descriptive
	adeeb_res.ID = adeeb_model.ID
	adeeb_res.Name = adeeb_model.Name
	adeeb_res.Bio = *adeeb_model.Bio
	adeeb_res.TimePeriod = adeeb_model.TimePeriod
	adeeb_res.Reviewed = adeeb_model.Reviewed

	return adeeb_res
}

func DBModels_To_ResSchemas(adeeb_models []database.Adeeb) []schemas.Adeeb_Descriptive {
	var adeebs []schemas.Adeeb_Descriptive

	for _, adeeb_model := range adeeb_models {
		adeeb_res := DBModel_To_ResSchema(adeeb_model)
		adeebs = append(adeebs, adeeb_res)
	}

	return adeebs
}

func ReqSchema_To_DBModel(adeeb_req CreateOneAdeeb_Req_Body) database.Adeeb {
	adeeb_model := database.Adeeb{
		Name:       adeeb_req.Name,
		Bio:        &adeeb_req.Bio,
		TimePeriod: adeeb_req.TimePeriod,
		Reviewed:   adeeb_req.Reviewed,
	}

	return adeeb_model
}

func ReqSchemas_To_DBModels(adeebs_req []CreateOneAdeeb_Req_Body) []database.Adeeb {
	var adeeb_models []database.Adeeb

	for _, adeeb_req := range adeebs_req {
		adeeb_model := ReqSchema_To_DBModel(adeeb_req)
		adeeb_models = append(adeeb_models, adeeb_model)
	}

	return adeeb_models
}
