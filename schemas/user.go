package schemas

import "adeeb_huma/database"

type User_Full struct {
	IDField
	User_UsernameField
	// We don't return passwords in responses.
	// User_PasswordField
	User_RolesField

	CreatedAtField
	UpdatedAtField
}

type User_Descriptive struct {
	IDField
	User_UsernameField
	User_RolesField
}

type User_Minimal struct {
	IDField
	User_UsernameField
}

type User_UsernameField struct {
	Username string `json:"username" doc:"User's name" minLength:"4" maxLength:"256"`
}

type User_UsernameField_Optional struct {
	Username *string `json:"username" doc:"User's name" minLength:"4" maxLength:"256" required:"false"`
}

type User_PasswordField struct {
	Password string `json:"password" doc:"Password's name" minLength:"4" maxLength:"256"`
}

type User_PasswordField_Optional struct {
	Password *string `json:"password" doc:"Password's name" minLength:"4" maxLength:"256" required:"false"`
}

type User_RolesField struct {
	Roles database.RolesType `json:"roles" enum:"Normal,Banned,Management,DBA,Analytics" doc:"User's roles"`
}

type User_RolesField_Optional struct {
	Roles *database.RolesType `json:"roles" enum:"Normal,Banned,Management,DBA,Analytics" doc:"User's roles" required:"false"`
}
