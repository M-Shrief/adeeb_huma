package prose_qoutes

import (
	"adeeb_huma/database"
	"adeeb_huma/schemas"
)

type GetOneProseQoute_Req struct {
	schemas.IDPath
}

type GetOneProseQoute_Res struct {
	Body   GetOneProseQoute_Res_Body
	Status int
}
type GetOneProseQoute_Res_Body struct {
	schemas.ProseQoute_Descriptive
	Adeeb schemas.Adeeb_Minimal `json:"adeeb"`
}

type CreateOneProseQoute_Req struct {
	Body CreateOneProseQoute_Req_Body
}

type CreateOneProseQoute_Req_Body struct {
	schemas.ProseQoute_QouteField
	schemas.ProseQoute_SourceField_Optional
	schemas.TagsField
	schemas.ReviewedField

	// Relations
	schemas.AdeebIDField
}

type CreateOneProseQoute_Res struct {
	Body   schemas.ProseQoute_Descriptive
	Status int
}

type CreateManyProseQoutes_Req struct {
	Body []CreateOneProseQoute_Req_Body
}

type CreateManyProseQoutes_Res struct {
	Body   schemas.CreateMany_Res_Body[schemas.ProseQoute_Descriptive]
	Status int
}

type UpdateProseQoute_Req struct {
	schemas.IDPath
	Body UodateProseQoute_Req_Body
}

type UodateProseQoute_Req_Body struct {
	schemas.ProseQoute_QouteField_Optional
	schemas.ProseQoute_SourceField_Optional
	schemas.TagsField_Optional
	schemas.ReviewedField_Optional

	// Relations
	schemas.AdeebIDField_Optional
}

type DeleteProseQoute_Req struct {
	schemas.IDPath
}

func DBModel_To_ResSchema(prose_qoute_model database.ProseQoute) schemas.ProseQoute_Descriptive {
	// Because we embed structs we can't assign values to fields directly without a lot of boilerplate.
	// So we define the variable and it's type, then assign each field alone.
	var prose_qoute_res schemas.ProseQoute_Descriptive
	prose_qoute_res.ID = prose_qoute_model.ID
	prose_qoute_res.Qoute = prose_qoute_model.Qoute
	prose_qoute_res.Source = prose_qoute_model.Source
	prose_qoute_res.Tags = prose_qoute_model.Tags
	prose_qoute_res.Reviewed = prose_qoute_model.Reviewed

	prose_qoute_res.AdeebID = prose_qoute_model.AdeebID

	return prose_qoute_res
}

func DBModels_To_ResSchemas(prose_qoute_models []database.ProseQoute) []schemas.ProseQoute_Descriptive {
	var prose_qoutes []schemas.ProseQoute_Descriptive

	for _, adeeb_model := range prose_qoute_models {
		prose_qoute_res := DBModel_To_ResSchema(adeeb_model)
		prose_qoutes = append(prose_qoutes, prose_qoute_res)
	}

	return prose_qoutes
}

type ProseQouteWithTotalCount struct {
	database.ProseQoute
	TotalCount int64 `json:"total_count"`
}

func DistillDBModelsWithCount(models_with_count []ProseQouteWithTotalCount) ([]schemas.ProseQoute_Descriptive, int64) {
	// Because we embed structs we can't assign values to fields directly without a lot of boilerplate.
	// So we define the variable and it's type, then assign each field alone.
	var prose_qoutes_res []schemas.ProseQoute_Descriptive
	var total_count int64
	for i, model := range models_with_count {
		var prose_qoute_res schemas.ProseQoute_Descriptive
		prose_qoute_res.ID = model.ID
		prose_qoute_res.Qoute = model.Qoute
		prose_qoute_res.Source = model.Source
		prose_qoute_res.Tags = model.Tags
		prose_qoute_res.Reviewed = model.Reviewed

		prose_qoute_res.AdeebID = model.AdeebID

		prose_qoutes_res = append(prose_qoutes_res, prose_qoute_res)
		if i == 0 {
			total_count = model.TotalCount
		}
	}
	return prose_qoutes_res, total_count
}

func ReqSchema_To_DBModel(prose_qoute_req CreateOneProseQoute_Req_Body) database.ProseQoute {
	prose_qoute_model := database.ProseQoute{
		Qoute:    prose_qoute_req.Qoute,
		Source:   prose_qoute_req.Source,
		Tags:     prose_qoute_req.Tags,
		Reviewed: prose_qoute_req.Reviewed,

		AdeebID: prose_qoute_req.AdeebID,
	}

	return prose_qoute_model
}

func ReqSchemas_To_DBModels(prose_qoutes_req []CreateOneProseQoute_Req_Body) []database.ProseQoute {
	var prose_qoute_models []database.ProseQoute

	for _, prose_qoute_req := range prose_qoutes_req {
		prose_qoute_model := ReqSchema_To_DBModel(prose_qoute_req)
		prose_qoute_models = append(prose_qoute_models, prose_qoute_model)
	}

	return prose_qoute_models
}
