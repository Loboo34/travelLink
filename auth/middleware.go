package auth

import (
	"context"
	"net/http"
	"strings"

	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/utils"
	"go.uber.org/zap"
)

func Authenticate(jwtManager *JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := extractClaims(r, jwtManager)
			if err != nil {
				utils.RespondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}

			ctx := context.WithValue(r.Context(), utils.ContextKeyUserID, claims.UserID)
			ctx = context.WithValue(ctx, utils.ContextKeyRole, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAdmin(jwtManager *JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			utils.Logger.Info("RequireAdmin middleware called", zap.String("path", r.URL.Path))
			claims, err := extractClaims(r, jwtManager)
			if err != nil {
				utils.Logger.Error("Failed to extract claims", zap.Error(err))
				utils.RespondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}

			utils.Logger.Info("Claims extracted", zap.String("userID", claims.UserID), zap.String("role", string(claims.Role)))

			if claims.Role != model.UserRoleAdmin {
				utils.Logger.Error("User is not admin", zap.String("role", string(claims.Role)))
				utils.RespondWithError(w, http.StatusForbidden, "admin access required")
				return
			}

			ctx := context.WithValue(r.Context(), utils.ContextKeyUserID, claims.UserID)
			ctx = context.WithValue(ctx, utils.ContextKeyRole, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractClaims(r *http.Request, jwtManager *JWTManager) (*Claims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, ErrInvalidToken
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, ErrInvalidToken
	}

	return jwtManager.Verify(parts[1])
}
