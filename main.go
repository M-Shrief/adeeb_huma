package main

import (
	"adeeb_huma/config"
	"adeeb_huma/database"
	"adeeb_huma/logger"
	"adeeb_huma/router"
	"fmt"
	"net/http"
)

func main() {
	// Load Config Variables
	config.LoadENV()

	// Init Logger
	logger.Init()

	_, err := database.NewDatabase()
	if err != nil {
		logger.Error().Err(err).Msg("Couldn't initialize db connection")
	}
	// 2. Create DB Enums
	database.CreateEnumsIfNotExists()
	// 3. Auto Migrate Models (Optional)
	database.Migrate()

	// Router & API
	r := router.NewRouter()
	router.UseMiddlewares()
	router.Init()

	logger.Info().Msgf("Starting Server at %v:%v", config.HOST, config.PORT)
	http.ListenAndServe(
		fmt.Sprintf("%v:%v", config.HOST, config.PORT),
		r,
	)

}
