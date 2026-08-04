package auth

import (
	"adeeb_huma/internal/utils"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// A middleware used to check if adminstrator authority in WRITE operations in domain components,
// like: POST/PUT/DELETE operations in adeebs, poems component...etc
func VerifyAdminstratorMiddleware(api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		auth_header := ctx.Header("Authorization")

		if len(auth_header) < 10 {
			// Condition met: Stop request and return response
			// Do NOT call next(ctx)
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Not Authorizaed")
			return
		}

		jwt_claims, err := VerifyJWT(auth_header)
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Not Authorizaed")
			return
		}
		user_permissions, err := utils.InterfaceToStringSlice(jwt_claims["permissions"])
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Not Authorizaed")
			return
		}

		is_adminstrator := CheckAdminstration(user_permissions, OPEnum_Read)
		if is_adminstrator == false {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Not Authorizaed")
			return
		}

		// Condition not met: Continue to next middleware/handler
		next(ctx)
	}
}

// // The same middleware but for Chi
// func VerifyAdminstratorMiddleware() func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			auth_header := r.Header.Get("Authorization")

// 			jwt_claims, err := VerifyJWT(auth_header)
// 			if err != nil {
// 				// Condition met: Stop request and return response
// 				http.Error(w, "Not Authorizaed", http.StatusUnauthorized)
// 				return // Crucial: Do not call next.ServeHTTP()
// 			}
// 			user_permissions, err := utils.InterfaceToStringSlice(jwt_claims["permissions"])
// 			if err != nil {
// 				http.Error(w, "Not Authorizaed", http.StatusUnauthorized)
// 				return
// 			}

// 			is_adminstrator := CheckAdminstration(user_permissions, OPEnum_Read)
// 			if is_adminstrator == false {
// 				http.Error(w, "Not Authorizaed", http.StatusUnauthorized)
// 				return
// 			}

// 			// Condition not met: Continue to the next handler
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }
