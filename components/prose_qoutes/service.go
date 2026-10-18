package prose_qoutes

import (
	"adeeb_huma/internal/errs"
	"adeeb_huma/internal/logger"
	"adeeb_huma/schemas"
	"context"
	"net/http"
)

func GetAll_Service(input *schemas.GetAll_Req) (*schemas.GetAll_Res_Body[schemas.ProseQoute_Descriptive], error) {
	result, err := GetAll_Repository(input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return nil, api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in GET /prose_qoutes")
			return nil, api_error
		}
	}
	return result, nil
}

func GetOne_Service(ctx context.Context, input *GetOneProseQoute_Req) (*GetOneProseQoute_Res_Body, error) {
	result, err := GetOne_Repository(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return nil, api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in GET /prose_qoutes/{id}")
			return nil, api_error
		}
	}
	return result, nil
}

func CreateOne_Service(ctx context.Context, input *CreateOneProseQoute_Req) (*schemas.ProseQoute_Descriptive, error) {
	result, err := CreateOne_Repository(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return nil, api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in POST /prose_qoutes")
			return nil, api_error
		}
	}
	return result, nil
}

func CreateMany_Service(ctx context.Context, input *CreateManyProseQoutes_Req) (*schemas.CreateMany_Res_Body[schemas.ProseQoute_Descriptive], error) {
	result, err := CreateMany_Repository(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return nil, api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in POST /prose_qoutes/many")
			return nil, api_error
		}
	}
	return result, nil

}

func UpdateOne_Service(ctx context.Context, input *UpdateProseQoute_Req) error {
	err := UpdateOne_Repository(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in PUT /prose_qoutes/{id}")
			return api_error
		}
	}
	return nil
}

func DeleteOne_Service(ctx context.Context, input *DeleteProseQoute_Req) error {
	err := DeleteOne_Repository(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			return api_error
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "service", Err: err}
			logger.Error().Err(api_error).Msg("Error in DELETE /prose_qoutes/{id}")
			return api_error
		}
	}
	return nil
}
