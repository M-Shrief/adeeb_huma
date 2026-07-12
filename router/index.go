package router

import (
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

	return API
}
