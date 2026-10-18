package prose_qoutes

import (
	"adeeb_huma/internal/errs"
	"adeeb_huma/internal/logger"
	"adeeb_huma/schemas"
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func GetAllProseQoutes_Handler(ctx context.Context, input *schemas.GetAll_Req) (*schemas.GetAll_Res[schemas.ProseQoute_Descriptive], error) {
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
			logger.Error().Err(api_error).Msg("Error in GET /prose_qoutes")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	res := &schemas.GetAll_Res[schemas.ProseQoute_Descriptive]{
		Body:   *result,
		Status: http.StatusOK,
	}

	return res, nil
}

func GetOneProseQoute_Handler(ctx context.Context, input *GetOneProseQoute_Req) (*GetOneProseQoute_Res, error) {
	result, err := GetOne_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusNotFound:
				return nil, huma.Error404NotFound("ProseQoute is not found")
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in GET /prose_qoutes/{id}")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	res := &GetOneProseQoute_Res{
		Body:   *result,
		Status: http.StatusOK,
	}

	return res, nil
}

func CreateOneProseQoute_Handler(ctx context.Context, input *CreateOneProseQoute_Req) (*CreateOneProseQoute_Res, error) {
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
			logger.Error().Err(api_error).Msg("Error in POST /prose_qoutes")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}
	return &CreateOneProseQoute_Res{Body: *result, Status: http.StatusCreated}, nil
}

func CreateManyProseQoutes_Handler(ctx context.Context, input *CreateManyProseQoutes_Req) (*CreateManyProseQoutes_Res, error) {
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
			logger.Error().Err(api_error).Msg("Error in POST /prose_qoutes/many")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	return &CreateManyProseQoutes_Res{
		Body:   *result,
		Status: http.StatusCreated,
	}, nil

}

func UpdateProseQoute_Handler(ctx context.Context, input *UpdateProseQoute_Req) (*schemas.Update_Res, error) {
	err := UpdateOne_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusNotFound:
				return nil, huma.Error404NotFound("ProseQoute is not found")
			case http.StatusConflict:
				return nil, huma.Error409Conflict(api_error.Message)
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in PUT /prose_qoutes/{id}")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	res := &schemas.Update_Res{
		Status: http.StatusNoContent,
	}

	return res, nil
}

func DeleteProseQouteHandler(ctx context.Context, input *DeleteProseQoute_Req) (*schemas.Delete_Res, error) {
	err := DeleteOne_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusConflict:
				return nil, huma.Error409Conflict("ProseQoute is referenced in other tables")
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in DELETE /prose_qoutes/{id}")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	res := &schemas.Delete_Res{
		Status: http.StatusNoContent,
	}

	return res, nil
}
