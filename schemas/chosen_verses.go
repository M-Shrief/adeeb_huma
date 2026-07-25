package schemas

type ChosenVerses_Full struct {
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

type ChosenVerses_Descriptive struct {
	IDField
	VersesField
	IsCoupletField
	TagsField
	ReviewedField
	// Relations
	AdeebIDField
	PoemIDField
}

type ChosenVerses_Minimal struct {
	IDField
}
