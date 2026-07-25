package schemas

type ChosenVerse_Full struct {
	IDField
	VersesField
	IsCoupletField
	TagsField
	ReviewedField
	CreatedAtField
	UpdatedAtField
	// Relations
	AdeebIDField
	PoemIDField
}

type ChosenVerse_Descriptive struct {
	IDField
	VersesField
	IsCoupletField
	TagsField
	ReviewedField
	// Relations
	AdeebIDField
	PoemIDField
}

type ChosenVerse_Minimal struct {
	IDField
}
