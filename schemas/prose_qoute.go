package schemas

type ProseQoute_QouteField struct {
	Qoute string `json:"qoute" doc:"Qoute" minLength:"4" maxLength:"512"`
}

type ProseQoute_QouteField_Optional struct {
	Qoute *string `json:"qoute" doc:"Qoute" minLength:"4" maxLength:"512" required:"false"`
}

type ProseQoute_SourceField_Optional struct {
	Source *string `json:"source" doc:"Qoute's Source" minLength:"4" maxLength:"128" required:"false"`
}

type ProseQoute_Full struct {
	IDField
	ProseQoute_QouteField
	ProseQoute_SourceField_Optional
	TagsField
	ReviewedField

	CreatedAtField
	UpdatedAtField
	// Relation
	AdeebIDField
}

type ProseQoute_Descriptive struct {
	IDField
	ProseQoute_QouteField
	ProseQoute_SourceField_Optional
	TagsField
	ReviewedField

	// Relation
	AdeebIDField
}

type ProseQoute_Minimal struct {
	IDField
	ProseQoute_QouteField
}
