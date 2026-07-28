package users

import (
	"adeeb_huma/database"
	"adeeb_huma/schemas"
)

type Signup_Req struct {
	Body struct {
		schemas.User_UsernameField
		schemas.User_PasswordField
		schemas.User_RolesField
	}
}

type Login_Req struct {
	Body struct {
		schemas.User_UsernameField
		schemas.User_PasswordField
	}
}

type UserAuthorized_Res struct {
	Body   UserAuthorized_Res_Body
	Status int
}

type UserAuthorized_Res_Body struct {
	User  schemas.User_Descriptive `json:"user" doc:"User's data"`
	Token string                   `json:"token"`
}

type GetAllUsers_Req struct {
	schemas.GetAll_Req
	schemas.AuthHeader
}

type GetCurrentUser_Req struct {
	schemas.AuthHeader
}

type GetUserByID_Req struct {
	schemas.AuthHeader
	schemas.IDPath
}

type GetOneUser_Res struct {
	Body   schemas.User_Descriptive
	Status int
}

type UpdateCurrentUser_Req struct {
	Body struct {
		schemas.User_UsernameField_Optional
		schemas.User_PasswordField_Optional
		schemas.User_RolesField_Optional
	}
	schemas.AuthHeader
}

type UpdateUserByID_Req struct {
	Body struct {
		schemas.User_UsernameField_Optional
		schemas.User_PasswordField_Optional
		schemas.User_RolesField_Optional
	}
	schemas.AuthHeader
	schemas.IDPath
}

func DBModel_To_DescriptiveSchema(user_model database.User) schemas.User_Descriptive {
	// Because we embed structs we can't assign values to fields directly without a lot of boilerplate.
	// So we define the variable and it's type, then assign each field alone.
	var user_res schemas.User_Descriptive
	user_res.ID = user_model.ID
	user_res.Username = user_model.Username
	user_res.Roles = user_model.Roles

	return user_res
}

func DBModels_To_DescriptiveSchemas(user_models []database.User) []schemas.User_Descriptive {
	var users []schemas.User_Descriptive

	for _, adeeb_model := range user_models {
		user_res := DBModel_To_DescriptiveSchema(adeeb_model)
		users = append(users, user_res)
	}

	return users
}
