package schemas

import "adeeb_huma/database"

type Print_Full struct {
	IDField
	Print_FontTypeField
	Print_FontColorField
	Print_OutfitColorField
	Print_OutfitTypeField

	CreatedAtField
	UpdatedAtField
	// Relations
	ProseQoute_QouteField_Optional
	VersesField_Optional
	IsCoupletField_Optional

	PoemIDField_Optional
	ChosenVerseIDField_Optional
	ProseQouteIDField_Optional

	OrderIDField
	UserIDField_Optional
}

type Print_Descriptive struct {
	IDField
	Print_FontTypeField
	Print_FontColorField
	Print_OutfitColorField
	Print_OutfitTypeField
	// Relations
	ProseQoute_QouteField_Optional
	VersesField_Optional
	IsCoupletField_Optional

	PoemIDField_Optional
	ChosenVerseIDField_Optional
	ProseQouteIDField_Optional

	OrderIDField
	UserIDField_Optional
}

type Print_Minimal struct {
	IDField
	// Relations
	OrderIDField
	UserIDField_Optional
}

type Print_FontTypeField struct {
	FontType string `json:"font_type" doc:"Print's font type" minLength:"4" maxLength:"64"`
}
type Print_FontTypeField_Optional struct {
	FontType *string `json:"font_type" doc:"Print's font type" minLength:"4" maxLength:"64" required:"false"`
}

type Print_FontColorField struct {
	FontColor string `json:"font_color" doc:"Print's font color" minLength:"4" maxLength:"64"`
}
type Print_FontColorField_Optional struct {
	FontColor *string `json:"font_color" doc:"Print's font color" minLength:"4" maxLength:"64" required:"false"`
}

type Print_OutfitColorField struct {
	OutfitColor string `json:"outfit_color" doc:"Print's outfit color" minLength:"4" maxLength:"64"`
}
type Print_OutfitColorField_Optional struct {
	OutfitColor *string `json:"outfit_color" doc:"Print's outfit color" minLength:"4" maxLength:"64" required:"false"`
}

type Print_OutfitTypeField struct {
	OutfitType database.OutfitTypeEnum `json:"status" enum:"تيشيرت - لياقة 7,تيشيرت - نص لياقة ,تشيرت - لياقة بولو,جاكيت,سويت شيرت,بلوفر" doc:"Print's outfit type"`
}
type Print_OutfitTypeField_Optional struct {
	OutfitType *database.OutfitTypeEnum `json:"status" enum:"تيشيرت - لياقة 7,تيشيرت - نص لياقة ,تشيرت - لياقة بولو,جاكيت,سويت شيرت,بلوفر" doc:"Print's outfit type" required:"false"`
}
