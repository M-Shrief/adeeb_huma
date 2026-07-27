package users

import "adeeb_huma/schemas"

type Signup_Req struct {
	Body struct {
		schemas.User_UsernameField
		schemas.User_PasswordField
		schemas.User_RolesField
	}
}

type UserAuthorized_Res struct {
	Body   UserAuthorized_Res_Body
	Status int
}

type UserAuthorized_Res_Body struct {
	User  UserData `json:"user" doc:"User's data"`
	Token string   `json:"token"`
}
type UserData struct {
	schemas.IDField
	schemas.User_UsernameField
	schemas.User_RolesField
}
