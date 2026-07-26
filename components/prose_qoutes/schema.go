package prose_qoutes

import (
	"adeeb_huma/database"
	"adeeb_huma/schemas"
)

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

func DBModel_To_DescriptiveSchema(prose_qoute_model database.ProseQoute) schemas.ProseQoute_Descriptive {
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

func DBModels_To_DescriptiveSchemas(prose_qoute_models []database.ProseQoute) []schemas.ProseQoute_Descriptive {
	var prose_qoutes []schemas.ProseQoute_Descriptive

	for _, adeeb_model := range prose_qoute_models {
		prose_qoute_res := DBModel_To_DescriptiveSchema(adeeb_model)
		prose_qoutes = append(prose_qoutes, prose_qoute_res)
	}

	return prose_qoutes
}

func ReqModel_To_DBModel(prose_qoute_req CreateOneProseQoute_Req_Body) database.ProseQoute {
	prose_qoute_model := database.ProseQoute{
		Qoute:    prose_qoute_req.Qoute,
		Source:   prose_qoute_req.Source,
		Tags:     prose_qoute_req.Tags,
		Reviewed: prose_qoute_req.Reviewed,

		AdeebID: prose_qoute_req.AdeebID,
	}

	return prose_qoute_model
}

func ReqModels_To_DBModels(prose_qoutes_req []CreateOneProseQoute_Req_Body) []database.ProseQoute {
	var prose_qoute_models []database.ProseQoute

	for _, prose_qoute_req := range prose_qoutes_req {
		prose_qoute_model := ReqModel_To_DBModel(prose_qoute_req)
		prose_qoute_models = append(prose_qoute_models, prose_qoute_model)
	}

	return prose_qoute_models
}
