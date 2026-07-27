package router

import (
	"adeeb_huma/components/adeebs"
	"adeeb_huma/components/chosen_verses"
	"adeeb_huma/components/poems"
	"adeeb_huma/components/prose_qoutes"
	"adeeb_huma/components/users"
	"adeeb_huma/schemas"
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

var API huma.API

func Init() huma.API {
	config := huma.DefaultConfig("Adeeb Huma", "0.1.0")
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearer": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}

	// Choosing Scalar to render /docs
	config.DocsRenderer = huma.DocsRendererScalar

	// // disable the route by setting
	// config.OpenAPIPath = ""

	API = humachi.New(R, config)

	RegisterAllRoutes()

	return API
}

func RegisterAllRoutes() {

	huma.Register(
		API,
		huma.Operation{
			OperationID: "index",
			Method:      http.MethodGet,
			Path:        "/",
			Summary:     "Server's metadata",
			Description: "Respond with server's metadata",
			// Tags:          []string{},
			DefaultStatus: http.StatusOK,
		},
		IndexRouteHandler,
	)

	huma.Register(
		API,
		huma.Operation{
			OperationID: "ping",
			Method:      http.MethodGet,
			Path:        "/ping",
			Summary:     "Ping the server",
			Description: "Ping the server to check it's status",
			// Tags:          []string{},
			DefaultStatus: http.StatusOK,
		},
		PingRouteHandler,
	)

	adeebs.RegisterAPI(API)
	poems.RegisterAPI(API)
	chosen_verses.RegisterAPI(API)
	prose_qoutes.RegisterAPI(API)

	users.RegisterAPI(API)
}

type IndexResponse_JSONBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Summary     string `json:"summary"`
	Version     string `json:"version"`
	Docs_url    string `json:"docs-url"`
}
type IndexResponse struct {
	Body   IndexResponse_JSONBody
	Status int
}

func IndexRouteHandler(ctx context.Context, input *struct{}) (*IndexResponse, error) {
	body := IndexResponse_JSONBody{
		Title:       "Adeeb Huma",
		Description: "An Iteration for Adeeb's RESTful API using Go, Huma and Postgres.",
		Summary:     "An Iteration for Adeeb's RESTful API using Go",
		Version:     "0.1.0",
		Docs_url:    "/docs",
	}
	return &IndexResponse{Body: body, Status: http.StatusOK}, nil
}

func PingRouteHandler(ctx context.Context, input *struct{}) (*schemas.BaseResponse, error) {
	return &schemas.BaseResponse{Body: schemas.BaseResponse_JSONBody{Message: "Pong"}, Status: http.StatusOK}, nil
}
