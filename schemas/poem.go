package schemas

type Poem_Full struct {
	IDField
	Poem_IntroField
	VersesField
	IsCoupletField
	ReviewedField

	CreatedAtField
	UpdatedAtField
	// Relation
	AdeebIDField
}

type Poem_Descriptive struct {
	IDField
	Poem_IntroField
	VersesField
	IsCoupletField
	ReviewedField

	// Relation
	AdeebIDField
}

type Poem_Minimal struct {
	IDField
	Poem_IntroField
}

type Poem_IntroField struct {
	Intro string `json:"intro" doc:"Poem's intro" minLength:"4" maxLength:"256"`
}

type Poem_IntroField_Optional struct {
	Intro *string `json:"intro" doc:"Poem's intro" minLength:"4" maxLength:"256" required:"true"`
}
