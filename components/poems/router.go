package poems

import (
	"adeeb_huma/internal/errs"
	"adeeb_huma/internal/logger"
	"adeeb_huma/schemas"
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func GetAllPoems_Handler(ctx context.Context, input *schemas.GetAll_Req) (*schemas.GetAll_Res[schemas.Poem_Descriptive], error) {
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
			logger.Error().Err(api_error).Msg("Error in GET /poems")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}
	res := &schemas.GetAll_Res[schemas.Poem_Descriptive]{
		Body:   *result,
		Status: http.StatusOK,
	}

	return res, nil
}

func GetOnePoem_Handler(ctx context.Context, input *GetOnePoem_Req) (*GetOnePoem_Res, error) {
	result, err := GetOne_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusNotFound:
				return nil, huma.Error404NotFound("Poem is not found")
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in GET /poems/{id}")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	res := &GetOnePoem_Res{
		Body:   *result,
		Status: http.StatusOK,
	}

	return res, nil
}

func CreateOnePoem_Handler(ctx context.Context, input *CreateOnePoem_Req) (*CreateOnePoem_Res, error) {
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
			logger.Error().Err(api_error).Msg("Error in POST /poems")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}
	return &CreateOnePoem_Res{Body: *result, Status: http.StatusCreated}, nil
}

func CreateManyPoems_Handler(ctx context.Context, input *CreateManyPoems_Req) (*CreateManyPoems_Res, error) {
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
			logger.Error().Err(api_error).Msg("Error in POST /poems/many")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	return &CreateManyPoems_Res{
		Body:   *result,
		Status: http.StatusCreated,
	}, nil

}

func UpdatePoem_Handler(ctx context.Context, input *UpdatePoem_Req) (*schemas.Update_Res, error) {
	err := UpdateOne_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusNotFound:
				return nil, huma.Error404NotFound("Poem is not found")
			case http.StatusConflict:
				return nil, huma.Error409Conflict(api_error.Message)
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in PUT /poems/{id}")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	res := &schemas.Update_Res{
		Status: http.StatusNoContent,
	}

	return res, nil
}

func DeletePoemHandler(ctx context.Context, input *DeletePoem_Req) (*schemas.Delete_Res, error) {
	err := DeleteOne_Service(ctx, input)
	if err != nil {
		if api_error, ok := err.(*errs.APIError); ok {
			switch api_error.Status {
			case http.StatusConflict:
				return nil, huma.Error409Conflict("Poem is referenced in other tables")
			case http.StatusBadRequest:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			default:
				return nil, huma.Error400BadRequest("Bad request, try again later")
			}
		} else {
			api_error := &errs.APIError{Status: http.StatusBadRequest, CausedIn: "router", Err: err}
			logger.Error().Err(api_error).Msg("Error in DELETE /poems/{id}")
			return nil, huma.Error400BadRequest("Bad request, try again later")
		}
	}

	res := &schemas.Delete_Res{
		Status: http.StatusNoContent,
	}

	return res, nil
}
