package bootstrap

import (
	"context"
	"dots-api/lib/utils"
	payment "dots-api/lib/xendit"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type CustomClaims struct {
	Code    string `json:"code"`
	Email   string `json:"email"`
	Channel string `json:"channel"`
	jwt.StandardClaims
}

const (
	ChannelApp = "app"
	ChannelCMS = "cms"
)

var (
	mustHeader = []string{"X-Channel", "Content-Type"}
	headerVal  = []string{ChannelApp, ChannelCMS, "application/json"}
)

func userContext(ctx context.Context, subject, id interface{}) context.Context {
	return context.WithValue(ctx, subject, id)
}

const pingReqURI string = "/v1/ping"

func isPingRequest(r *http.Request) bool {
	return r.RequestURI == pingReqURI
}

// NotfoundMiddleware A custom not found response.
func (app *App) NotfoundMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tctx := chi.NewRouteContext()
		rctx := chi.RouteContext(r.Context())

		if !rctx.Routes.Match(tctx, r.Method, r.URL.Path) {
			app.SendNotfound(w, utils.ErrNotFoundPage)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *App) Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				logEntry := middleware.GetLogEntry(r)
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				} else {
					debug.PrintStack()
				}

				app.Log.FromDefault().WithFields(logrus.Fields{
					"Panic": rvr,
				}).Errorf("Panic: %v \n %v", rvr, string(debug.Stack()))

				app.SendBadRequest(w, utils.ErrSystemError)
				return
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (app *App) VerifyJwtToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := &CustomClaims{}
		tokenAuth := r.Header.Get("Authorization")
		_, err := jwt.ParseWithClaims(tokenAuth, claims, func(token *jwt.Token) (interface{}, error) {
			if jwt.SigningMethodHS256 != token.Method {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			secret := app.Config.GetString("app.key")
			return []byte(secret), nil
		})

		if err != nil {
			msg := utils.ErrInvalidToken
			if mErr, ok := err.(*jwt.ValidationError); ok {
				if mErr.Errors == jwt.ValidationErrorExpired {
					msg = utils.ErrTokenExpired
				}
			}

			app.SendAuthError(w, msg)
			return
		}

		// check if token expired or not
		if claims.ExpiresAt <= time.Now().Unix() {
			app.SendAuthError(w, utils.ErrTokenExpired)
			return
		}

		ctx := userContext(r.Context(), "identifier", map[string]string{
			"code":    claims.Code,
			"email":   claims.Email,
			"channel": claims.Channel,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// HeaderCheckerMiddleware check the necesarry headers
func (app *App) HeaderCheckerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, v := range mustHeader {
			if len(r.Header.Get(v)) == 0 || !utils.Contains(headerVal, r.Header.Get(v)) {
				app.SendBadRequest(w, fmt.Sprintf("undefined %s header or wrong value of header", v))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (app *App) VerifyAccessRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx = context.TODO()
		)

		claims := &CustomClaims{}
		tokenAuth := r.Header.Get("Authorization")
		_, _ = jwt.ParseWithClaims(tokenAuth, claims, func(token *jwt.Token) (interface{}, error) {
			return nil, nil
		})

		// routePath := r.URL.Path
		routePattern := chi.RouteContext(r.Context()).RoutePattern()

		isPermitted, err := app.CheckPermission(app.DB, ctx, claims.Issuer, claims.Code, convertRoutePattern(routePattern), r.Method)
		if err != nil {
			app.SendForbidden(w, "")
			return
		}

		if !isPermitted {
			app.SendForbidden(w, "")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (c *App) CheckPermission(db *pgxpool.Pool, ctx context.Context, actorType, code, routePattern, routeMethod string) (bool, error) {
	var (
		rolePermissionID int
		query            string
	)

	queryUser := `select rp.id 
	from users u 
	inner join role_permissions rp on u.role_id = rp.role_id 
	inner join permissions p on rp.permission_id  = p.id
	where u.user_code = $1 and p.route_pattern = $2 and p.route_method = $3`

	queryAdmin := `select rp.id 
	from admins a 
	inner join role_permissions rp on a.role_id = rp.role_id 
	inner join permissions p on rp.permission_id  = p.id
	where a.admin_code = $1 and p.route_pattern = $2 and p.route_method = $3`

	switch actorType {
	case utils.User:
		query = queryUser
	case utils.Admin:
		query = queryAdmin
	default:
	}

	err := db.QueryRow(ctx, query, code, routePattern, routeMethod).Scan(
		&rolePermissionID,
	)
	if err != nil && err != pgx.ErrNoRows {
		return false, fmt.Errorf("%s, %s", "middleware.CheckPermission", utils.ErrCheckPermission)
	}

	if rolePermissionID < 1 {
		return false, nil
	}

	return true, nil

}

func convertRoutePattern(routePattern string) string {
	// Split the route pattern by '/'
	parts := strings.Split(routePattern, "/")

	// Iterate over the parts
	for i, part := range parts {
		// Replace '{parameter}' with '*'
		if strings.Contains(part, "{") && strings.Contains(part, "}") {
			parts[i] = "*"
		}
	}

	// Join the parts back into a string
	convertedPattern := strings.Join(parts, "/")

	return convertedPattern
}

func (app *App) VerifyXenditCallbackToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xenditCallbackToken := r.Header.Get("X-Callback-Token")
		if !payment.IsCallbackTokenVerified(xenditCallbackToken) {
			app.SendForbidden(w, "")
			return
		}

		next.ServeHTTP(w, r)
	})
}
