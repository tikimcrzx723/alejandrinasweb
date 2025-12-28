package routes

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/tikimcrzx723/alejandrinasweb/controllers"
	"github.com/tikimcrzx723/alejandrinasweb/internal/env"
	"github.com/tikimcrzx723/alejandrinasweb/routes/middleware"
	"github.com/tikimcrzx723/alejandrinasweb/static"
)

type Routes struct {
	e *echo.Echo
}

func sessionKeyFromEnv(raw string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err == nil && len(decoded) >= 32 {
		return decoded
	}

	// Normalize arbitrary strings into a 32-byte key using SHA256.
	hash := sha256.Sum256([]byte(raw))
	return hash[:]
}

func NewRoutes() Routes {
	e := echo.New()

	authKey := sessionKeyFromEnv(env.GetString("SESSION_AUTH_KEY", "zRJdixjhVNDh..."))
	encKey := sessionKeyFromEnv(env.GetString("SESSION_ENC_KEY", "zRJdixjhVNDh..."))

	e.Use(
		session.Middleware(sessions.NewCookieStore(authKey, encKey)),
		controllers.RegisterAppContext,
		controllers.RegisterFlashMessageContext,
	)

	csrfSecure := env.GetBool("CSRF_COOKIE_SECURE", false)
	sameSiteMode := csrf.SameSiteDefaultMode
	if csrfSecure {
		sameSiteMode = csrf.SameSiteNoneMode
	}

	if !csrfSecure {
		// When running over HTTP (local/dev) we must disable the strict HTTPS-only referer checks.
		e.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				if !c.IsTLS() {
					c.SetRequest(csrf.PlaintextHTTPRequest(c.Request()))
				}
				return next(c)
			}
		})
	}

	csrfKey := env.GetString("CSRF_TOKEN_KEY", "32-byte-secret-key-minimo-32-chars!!")
	trustedOriginsValue := env.GetString("CSRF_TRUSTED_ORIGINS", "")
	var trustedOrigins []string
	for origin := range strings.SplitSeq(trustedOriginsValue, ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			trustedOrigins = append(trustedOrigins, origin)
		}
	}

	csrfOptions := []csrf.Option{
		csrf.Secure(csrfSecure),
		csrf.Path("/"),
		csrf.SameSite(sameSiteMode),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			fmt.Println("\n================ CSRF ERROR ================")
			fmt.Println("Method:", r.Method)
			fmt.Println("Path:", r.URL.Path)

			// 1) Revisar cookie CSRF
			c, err := r.Cookie("_gorilla_csrf") // <- nombre real de la cookie
			if err != nil {
				fmt.Println("Cookie _gorilla_csrf NO presente:", err)
			} else {
				fmt.Println("Cookie _gorilla_csrf =", c.Value)
			}

			// 2) Revisar header
			fmt.Println("Header X-CSRF-Token =", r.Header.Get("X-CSRF-Token"))

			// 3) Revisar formulario
			_ = r.ParseForm()
			fmt.Println("Form gorilla.csrf.Token =", r.Form.Get("gorilla.csrf.Token"))
			fmt.Println("Failure Reason:", csrf.FailureReason(r))

			fmt.Println("============================================")

			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte("csrf failed"))
		})),
	}

	if len(trustedOrigins) > 0 {
		csrfOptions = append(csrfOptions, csrf.TrustedOrigins(trustedOrigins))
	}

	csrfMiddleware := csrf.Protect([]byte(csrfKey), csrfOptions...)

	e.Use(echo.WrapMiddleware(csrfMiddleware))

	echo.MustSubFS(static.Files, "static")
	e.StaticFS("/static", static.Files)
	return Routes{e}
}

func (r Routes) Load() *echo.Echo {
	adminRoutes := r.e.Group("/admin", middleware.RequireAdminRole)
	adminRoutes.GET("/dashboard/product/register", func(c echo.Context) error {
		return controllers.RegisterProductPage(c)
	})
	adminRoutes.GET("/dashboard/category/register", func(c echo.Context) error {
		return controllers.CategoryPage(c)
	})

	adminRoutes.POST("/category/register", func(c echo.Context) error {
		return controllers.CreateCategory(c)
	})
	adminRoutes.POST("/product/register", func(c echo.Context) error {
		return controllers.CreateProduct(c)
	})
	adminRoutes.POST("/product/update", func(c echo.Context) error {
		return controllers.UpdateProduct(c)
	})
	// setup routes for diferents pages
	r.e.GET("", func(c echo.Context) error {
		return controllers.Home(c)
	})
	r.e.GET("/product/:sku", func(c echo.Context) error {
		return controllers.Product(c)
	})
	r.e.POST("/register", func(c echo.Context) error {
		return controllers.CreateUser(c)
	})
	r.e.GET("/login", func(c echo.Context) error {
		return controllers.LoginPage(c)
	}, middleware.RequireNoAuth)
	r.e.POST("/login", func(c echo.Context) error {
		return controllers.LoginUser(c)
	})
	r.e.GET("/logout", func(c echo.Context) error {
		return controllers.LogoutUser(c)
	}, middleware.RequireAuth)
	r.e.GET("/register", func(c echo.Context) error {
		return controllers.Register(c)
	}, middleware.RequireNoAuth)
	return r.e
}
