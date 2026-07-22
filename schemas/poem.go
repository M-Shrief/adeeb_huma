package schemas

type Poem_IntroField struct {
	Intro string `json:"intro" doc:"Poem's intro" minLength:"4" maxLength:"256"`
}

type Poem_IntroField_Optional struct {
	Intro *string `json:"intro" doc:"Poem's intro" minLength:"4" maxLength:"256" required:"true"`
}
