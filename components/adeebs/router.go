package adeebs

import (
	"adeeb_huma/internal/errs"
	"adeeb_huma/internal/logger"
	"adeeb_huma/schemas"
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func GetAllAdeebs_Handler(ctx context.Context, input *schemas.GetAll_Req) (*schemas.GetAll_Res[schemas.Adeeb_Descriptive], error) {
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
			logger.Error().Err(api_error).Msg("Error in GET /adeebs")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	resp := &schemas.GetAll_Res[schemas.Adeeb_Descriptive]{
		Body:   *result,
		Status: http.StatusOK,
	}

	return resp, nil
}

func GetOneAdeeb_Handler(ctx context.Context, input *GetOneAdeeb_Req) (*GetOneAdeeb_Res, error) {
	result, err := GetOne_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusNotFound:
				return nil, huma.Error404NotFound("Adeeb is not found")
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in GET /adeebs/{id}")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	return &GetOneAdeeb_Res{*result, http.StatusOK}, nil
}

func CreateOneAdeeb_Handler(ctx context.Context, input *CreateOneAdeeb_Req) (*CreateOneAdeeb_Res, error) {
	result, err := CreateOne_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusConflict:
				return nil, huma.Error409Conflict("Adeeb is not found")
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in POST /adeebs")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}
	return &CreateOneAdeeb_Res{Body: *result, Status: http.StatusCreated}, nil
}

func CreateManyAdeeb_Handler(ctx context.Context, input *CreateManyAdeebs_Req) (*CreateManyAdeebs_Res, error) {
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
			logger.Error().Err(api_error).Msg("Error in POST /adeebs/many")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	return &CreateManyAdeebs_Res{
		Body:   *result,
		Status: http.StatusCreated,
	}, nil

}

func UpdateAdeeb_Handler(ctx context.Context, input *UpdateAdeeb_Req) (*schemas.Update_Res, error) {
	err := UpdateOne_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusNotFound:
				return nil, huma.Error404NotFound("Adeeb is not found")
			case http.StatusConflict:
				return nil, huma.Error409Conflict("Adeeb already exists")
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in PUT /adeebs/{id}")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	res := &schemas.Update_Res{
		Status: http.StatusNoContent,
	}

	return res, nil
}

func DeleteAdeeb_Handler(ctx context.Context, input *DeleteAdeeb_Req) (*schemas.Delete_Res, error) {
	err := DeleteOne_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusConflict:
				return nil, huma.Error409Conflict("Adeeb is referenced in other tables")
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in DELETE /adeebs/{id}")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	res := &schemas.Delete_Res{
		Status: http.StatusNoContent,
	}

	return res, nil
}
