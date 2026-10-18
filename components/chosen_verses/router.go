package chosen_verses

import (
	"adeeb_huma/internal/errs"
	"adeeb_huma/internal/logger"
	"adeeb_huma/schemas"
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func GetAllChosenVerses_Handler(ctx context.Context, input *schemas.GetAll_Req) (*schemas.GetAll_Res[schemas.ChosenVerse_Descriptive], error) {
	result, err := GetAll_Service(input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in GET /chosen_verses")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}
	res := &schemas.GetAll_Res[schemas.ChosenVerse_Descriptive]{
		Body:   *result,
		Status: http.StatusOK,
	}

	return res, nil
}

func GetOneChosenVerse_Handler(ctx context.Context, input *GetOneChosenVerse_Req) (*GetOneChosenVerse_Res, error) {
	result, err := GetOne_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusNotFound:
				return nil, huma.Error404NotFound("ChosenVerse is not found")
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in GET /chosen_verses/{id}")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	res := &GetOneChosenVerse_Res{
		Body:   *result,
		Status: http.StatusOK,
	}

	return res, nil
}

func CreateOneChosenVerse_Handler(ctx context.Context, input *CreateOneChosenVerse_Req) (*CreateOneChosenVerse_Res, error) {
	result, err := CreateOne_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusConflict:
				return nil, huma.Error409Conflict(api_error.Message)
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in POST /chosen_verses")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}
	return &CreateOneChosenVerse_Res{Body: *result, Status: http.StatusCreated}, nil
}

func CreateManyChosenVerses_Handler(ctx context.Context, input *CreateManyChosenVerses_Req) (*CreateManyChosenVerses_Res, error) {
	result, err := CreateMany_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in POST /chosen_verses/many")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	return &CreateManyChosenVerses_Res{
		Body:   *result,
		Status: http.StatusCreated,
	}, nil

}

func UpdateChosenVerse_Handler(ctx context.Context, input *UpdateChosenVerse_Req) (*schemas.Update_Res, error) {
	err := UpdateOne_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusNotFound:
				return nil, huma.Error404NotFound("ChosenVerse is not found")
			case http.StatusConflict:
				return nil, huma.Error409Conflict(api_error.Message)
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in PUT /chosen_verses/{id}")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	res := &schemas.Update_Res{
		Status: http.StatusNoContent,
	}

	return res, nil
}

func DeleteChosenVerseHandler(ctx context.Context, input *DeleteChosenVerse_Req) (*schemas.Delete_Res, error) {
	err := DeleteOne_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusConflict:
				return nil, huma.Error409Conflict("ChosenVerse is referenced in other tables")
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in DELETE /chosen_verses/{id}")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	res := &schemas.Delete_Res{
		Status: http.StatusNoContent,
	}

	return res, nil
}
