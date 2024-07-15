package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

const (
	missingJWTErrorMessage       = "Requires authentication"
	invalidJWTErrorMessage       = "Bad credentials"
	permissionDeniedErrorMessage = "Permission denied"
)

type customClaims struct {
	Permissions []string `json:"permissions"`
}

func (c customClaims) Validate(ctx context.Context) error {
	return nil
}

func (c customClaims) hasPermissions(expectedClaims []string) bool {
	if len(expectedClaims) == 0 {
		return false
	}
	for _, scope := range expectedClaims {
		if !contains(c.Permissions, scope) {
			return false
		}
	}
	return true
}

func ValidatePermissions(expectedClaims []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
		claims := token.CustomClaims.(*customClaims)
		if !claims.hasPermissions(expectedClaims) {
			errorMessage := errorMessage{Message: permissionDeniedErrorMessage}
			if err := writeJSON(w, http.StatusForbidden, errorMessage); err != nil {
				slog.Error("Failed to write error message: " + err.Error())
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ValidateJWT(audience, domain string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		issuerURL, err := url.Parse("https://" + domain + "/")
		if err != nil {
			slog.Error("Failed to parse the issuer url: " + err.Error())
			return
		}

		provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

		jwtValidator, err := validator.New(
			provider.KeyFunc,
			validator.RS256,
			issuerURL.String(),
			[]string{audience},
			validator.WithCustomClaims(func() validator.CustomClaims {
				return new(customClaims)
			}),
		)
		if err != nil {
			slog.Error("Failed to set up the jwt validator")
			return
		}

		if authHeaderParts := strings.Fields(r.Header.Get("Authorization")); len(authHeaderParts) > 0 && strings.ToLower(authHeaderParts[0]) != "bearer" {
			errorMessage := errorMessage{Message: invalidJWTErrorMessage}
			if err := writeJSON(w, http.StatusUnauthorized, errorMessage); err != nil {
				slog.Error("Failed to write error message: " + err.Error())
			}
			return
		}

		errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
			slog.Error("Encountered error while validating JWT: " + err.Error())
			if errors.Is(err, jwtmiddleware.ErrJWTMissing) {
				errorMessage := errorMessage{Message: missingJWTErrorMessage}
				if err := writeJSON(w, http.StatusUnauthorized, errorMessage); err != nil {
					slog.Error("Failed to write error message: " + err.Error())
				}
				return
			}
			if errors.Is(err, jwtmiddleware.ErrJWTInvalid) {
				errorMessage := errorMessage{Message: invalidJWTErrorMessage}
				if err := writeJSON(w, http.StatusUnauthorized, errorMessage); err != nil {
					slog.Error("Failed to write error message: " + err.Error())
				}
				return
			}
			serverError(w, err)
		}

		middleware := jwtmiddleware.New(
			jwtValidator.ValidateToken,
			jwtmiddleware.WithErrorHandler(errorHandler),
		)

		middleware.CheckJWT(next).ServeHTTP(w, r)
	})
}
