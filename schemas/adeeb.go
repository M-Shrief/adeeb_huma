package schemas

import "adeeb_huma/database"

type Adeeb_Full struct {
	IDField
	Adeeb_NameField
	Adeeb_BioField
	Adeeb_TimePeriodField
	ReviewedField
	CreatedAtField
	UpdatedAtField
}

type Adeeb_Descriptive struct {
	IDField
	Adeeb_NameField
	Adeeb_BioField
	Adeeb_TimePeriodField
	ReviewedField
}

type Adeeb_Minimal struct {
	IDField
	Adeeb_NameField
}
type Adeeb_NameField struct {
	Name string `json:"name" doc:"Adeeb's name" minLength:"4" maxLength:"256"`
}

type Adeeb_NameField_Optional struct {
	Name *string `json:"name" doc:"Adeeb's name" minLength:"4" maxLength:"256" required:"false"`
}

type Adeeb_BioField struct {
	Bio string `json:"bio" doc:"Adeeb's bio" minLength:"4" maxLength:"1024"`
}

type Adeeb_BioField_Optional struct {
	Bio *string `json:"bio" doc:"Adeeb's bio" minLength:"4" maxLength:"1024" required:"false"`
}

type Adeeb_TimePeriodField struct {
	TimePeriod database.TimePeriodEnum `json:"time_period" enum:"غير محدد,جاهلي,أموي,عباسي,أندلسي,عثماني ومملوكي,حديث" doc:"Adeeb's time period"`
}

type Adeeb_TimePeriodField_Optional struct {
	TimePeriod *database.TimePeriodEnum `json:"time_period" enum:"غير محدد,جاهلي,أموي,عباسي,أندلسي,عثماني ومملوكي,حديث" doc:"Adeeb's time period" required:"false"`
}
