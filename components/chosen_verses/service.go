package chosen_verses

import (
	"adeeb_huma/internal/errs"
	"adeeb_huma/internal/logger"
	"adeeb_huma/schemas"
	"context"
	"net/http"
)

func GetAll_Service(input *schemas.GetAll_Req) (*schemas.GetAll_Res_Body[schemas.ChosenVerse_Descriptive], error) {
	result, err := GetAll_Repository(input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return nil, api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in GET /chosen_verses")
			return nil, api_error
		}
	}
	return result, nil
}

func GetOne_Service(ctx context.Context, input *GetOneChosenVerse_Req) (*GetOneChosenVerse_Res_Body, error) {
	result, err := GetOne_Repository(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return nil, api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in GET /chosen_verses/{id}")
			return nil, api_error
		}
	}
	return result, nil
}

func CreateOne_Service(ctx context.Context, input *CreateOneChosenVerse_Req) (*schemas.ChosenVerse_Descriptive, error) {
	result, err := CreateOne_Repository(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return nil, api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in POST /chosen_verses")
			return nil, api_error
		}
	}
	return result, nil
}

func CreateMany_Service(ctx context.Context, input *CreateManyChosenVerses_Req) (*schemas.CreateMany_Res_Body[schemas.ChosenVerse_Descriptive], error) {
	result, err := CreateMany_Repository(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return nil, api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in POST /chosen_verses/many")
			return nil, api_error
		}
	}
	return result, nil

}

func UpdateOne_Service(ctx context.Context, input *UpdateChosenVerse_Req) error {
	err := UpdateOne_Repository(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in PUT /chosen_verses/{id}")
			return api_error
		}
	}
	return nil
}

func DeleteOne_Service(ctx context.Context, input *DeleteChosenVerse_Req) error {
	err := DeleteOne_Repository(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in DELETE /chosen_verses/{id}")
			return api_error
		}
	}
	return nil
}
