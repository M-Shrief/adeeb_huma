# adeeb_huma

An Iteration for Adeeb's RESTful API using Go, Huma and Postgres.

All iterations should apply the same characteristics, like:
- JWT authentication & authorization
- API documentation (mostly with Scalar)
- Input/Output Validation
- ...etc


## Tech stack
- Language: Go
- Framework: Huma & Chi
- Database: Postgres & Gorm
- Cache: ValKey
- Logger: Zerolog

## Characteristics
- Adherent to Open API industry standards, with auto generation for OpenAPI & JSON Schema.
- Interactive and Automated Documentation using Scalar.
 thanks to Huma.
- Validation for Requests and Responses with Huma automatic validation.
- Using structured logging with Zerolog.
- JWT Authentication and Authorization.
- Using stdlib net/http middlewares - mainly Chi's middlewares - for logging requests, rate limiting, CORS, ...etc.


## Folders structure
I'll default Go files, and general folders like `.gitignore`...etc.

- `internal` folder for packages that are specific for the project:
    - `auth` holds JWT authentication and authorization config and functionalities. 
    - `logger` folder holds Zerolog config and functionalities.
- `components` folder for the components that hold specific groups of routes and functionality, usually CRUD operations with JWT authentication & authorization. Each of them has:
    - `index.go`: the base file for the component, and have the main function to register it in the server: `RegisterAPI()` to be used in the router package.
    - `schema.go`: holds Request & Response schemas for each route
    - `router.go`: holds routes' handlers for each endpoint with different methods.
- `schemas`: used to contain shared schemas in the app, either general ones like APIs' specific schemas for certain requests or schemas for certain components. It's used to have one single source for the truth, without relying heavily on the Gorm models in the database package, and to reduce the boilerplate.
- `config`: holds the functionality to read ENVs, certificates, and PEM files.
- `database`: holds Gorm config and functionalities:
    - `enums`: have every enum declaration with `Value()` & `Scan()` function to deal with the database, and the SQL query used to create it if it doesn't exist in the app startup.
    - `models.go`: We define all database models here.
    - `index.go`: handle Database connection
- `cache`:
    - `utils.go`: shared utility used to deal with JSON data to reduce boilerplate
    - `index.go`: handle cache connection
- `router` folder holds Chi and Huma config and functionalities.
    - `config.go` holds Chi config, functionalities, and Middlewares.
    - `index.go` holds Huma's config, registering routes, and the initialization function for the router.
- `main.go` the entry point for the backend server.

## Notes
- if you've problems with Go's PATH: 
    ```sh
    $ export GO_PATH=~/go
    ```

    ```sh
    $ export PATH=$PATH:/$GO_PATH/bin
    ```