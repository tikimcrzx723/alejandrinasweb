package controllers

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/gosimple/slug"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/tikimcrzx723/alejandrinasweb/internal/api"
	"github.com/tikimcrzx723/alejandrinasweb/internal/dtos"
	"github.com/tikimcrzx723/alejandrinasweb/internal/env"
	"github.com/tikimcrzx723/alejandrinasweb/routes/contexts"
	"github.com/tikimcrzx723/alejandrinasweb/views"
)

const (
	AuthSessionName       = "authSessionCookie"
	authUserIDKey         = "USER_ID"
	authUserEmailKey      = "USER_EMAIL"
	AuthUserAuthenticated = "USER_AUTHENTICATED"
	flashSessionName      = "flashSession"
)

func createAuthSession(
	ctx echo.Context,
	user dtos.LoginResponse,
	extendSession bool,
) error {
	s, err := session.Get(AuthSessionName, ctx)
	fmt.Println(s)
	if err != nil {
		return err
	}

	maxAge := 604800
	if extendSession {
		maxAge = maxAge * 2
	}

	s.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   true,
	}

	s.Values[authUserIDKey] = user.Data.User.ID
	s.Values[authUserEmailKey] = user.Data.User.Email
	s.Values["TOKEN_KEY"] = user.Data.AccessToken
	s.Values[AuthUserAuthenticated] = true
	s.Values["ROLE"] = user.Data.User.Role

	return s.Save(ctx.Request(), ctx.Response())
}

func setAppCtx(ctx echo.Context) context.Context {
	appCtxkey := contexts.AppKey{}
	appCtx := ctx.Get(appCtxkey.String())
	withAppCtx := context.WithValue(
		ctx.Request().Context(),
		appCtxkey,
		appCtx,
	)

	flashCtxKey := contexts.FlashKey{}
	flashCtx := ctx.Get(flashCtxKey.String())

	return context.WithValue(
		withAppCtx,
		flashCtxKey,
		flashCtx,
	)
}

func LogoutUser(ctx echo.Context) error {
	s, err := session.Get(AuthSessionName, ctx)
	if err != nil {
		return err
	}

	s.Options.MaxAge = -1

	if err := s.Save(ctx.Request(), ctx.Response()); err != nil {
		return err
	}

	return ctx.Redirect(http.StatusSeeOther, "/")
}

func RegisterProductPage(c echo.Context) error {
	token := csrf.Token(c.Request())

	products, err := api.GetProducts(c.Request().Context(), env.GetString("API_URL", "http://localhost:8080/api/v1/"))
	if err != nil {
		return err
	}

	c.Response().Header().Set("Cache-Control", "no-store")
	return views.RegisterProduct("Alejandrinas - Registro de Producto", token, products).
		Render(renderArgs(c))
}

// func addFlash(ctx echo.Context, flashType contexts.FlashType, msg string) error {
// 	s, err := session.Get(flashSessionName, ctx)
// 	if err != nil {
// 		return err
// 	}

// 	flash := contexts.FlashMessage{
// 		ID:        uuid.New(),
// 		Type:      flashType,
// 		CreatedAt: time.Now(),
// 		Message:   msg,
// 	}

// 	s.AddFlash(flash, flashSessionName)

// 	if err := s.Save(ctx.Request(), ctx.Response()); err != nil {
// 		return err
// 	}

// 	// También lo agregamos al contexto actual para que esté disponible en esta misma respuesta.
// 	if existing, ok := ctx.Get(contexts.FlashKey{}.String()).([]contexts.FlashMessage); ok {
// 		ctx.Set(contexts.FlashKey{}.String(), append(existing, flash))
// 	} else {
// 		ctx.Set(contexts.FlashKey{}.String(), []contexts.FlashMessage{flash})
// 	}

// 	return nil
// }

func RegisterFlashMessageContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if strings.HasPrefix(c.Request().URL.Path, "/static") {
			return next(c)
		}

		s, err := session.Get(flashSessionName, c)
		if err != nil {
			return next(c)
		}

		flashMessages := []contexts.FlashMessage{}
		if flashes := s.Flashes(flashSessionName); len(flashes) > 0 {
			for _, flash := range flashes {
				if msg, ok := flash.(contexts.FlashMessage); ok {
					flashMessages = append(flashMessages, contexts.FlashMessage{
						ID:        msg.ID,
						Type:      msg.Type,
						CreatedAt: msg.CreatedAt,
						Message:   msg.Message,
					})
				}
			}
		}

		if err := s.Save(c.Request(), c.Response()); err != nil {
			slog.ErrorContext(
				c.Request().Context(),
				"could not save flash session",
				"err",
				err,
			)
			return next(c)
		}

		c.Set(contexts.FlashKey{}.String(), flashMessages)

		return next(c)
	}
}

func RegisterAppContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get(AuthSessionName, c)
		if err != nil {
			return err
		}

		isAuth, _ := sess.Values[AuthUserAuthenticated].(bool)
		userID := 0
		switch v := sess.Values[authUserIDKey].(type) {
		case int:
			userID = v
		case int64:
			userID = int(v)
		case float64:
			userID = int(v)
		}
		token := ""
		switch v := sess.Values["TOKEN_KEY"].(type) {
		case string:
			token = v
		case []byte:
			token = string(v)
		}

		role := ""
		switch v := sess.Values["ROLE"].(type) {
		case string:
			role = v
		case []byte:
			role = string(v)
		}

		appContext := contexts.App{
			Context:         c,
			UserID:          userID,
			IsAuthenticated: isAuth,
			Token:           token,
			Role:            role,
		}

		c.Set(contexts.AppKey{}.String(), appContext)
		reqCtx := context.WithValue(c.Request().Context(), contexts.AppKey{}, appContext)
		c.SetRequest(c.Request().WithContext(reqCtx))
		return next(c)
	}
}

func SessionNew(ctx echo.Context) error {
	token := csrf.Token(ctx.Request())

	ctx.Response().Header().Set("Cache-Control", "no-store")

	return views.LoginPage("Login", token).Render(renderArgs(ctx))
}

func renderArgs(ctx echo.Context) (context.Context, io.Writer) {
	return setAppCtx(ctx), ctx.Response().Writer
}

func Home(c echo.Context) error {
	products, err := api.GetProducts(c.Request().Context(), env.GetString("API_URL", "http://localhost:8080/api/v1/"))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return views.HomePage("Alejandrinas - Inicio", products.Product).
		Render(renderArgs(c))
}

func Product(c echo.Context) error {
	sku := c.Param("sku")
	product, err := api.GetProductBySKU(c.Request().Context(), env.GetString("API_URL", "http://localhost:8080/api/v1/"), sku)
	if err != nil {
		return views.ErrorPage(
			views.WithErrPageTitle("El producto no existe o fue eliminado"),
			views.WithErrPageMsg("El producto que buscas no fue encontrado"),
		).Render(renderArgs(c))
	}
	return views.ProductPage("Alejandrinas - Detalle Producto", product.Product).
		Render(renderArgs(c))
}

func Register(c echo.Context) error {
	token := csrf.Token(c.Request())

	c.Response().Header().Set("Cache-Control", "no-store")

	return views.RegisterPage("Alejandrinas - Registro", token).
		Render(renderArgs(c))
}

func LoginPage(c echo.Context) error {
	token := csrf.Token(c.Request())

	c.Response().Header().Set("Cache-Control", "no-store")

	return views.LoginPage("Alejandrinas - Iniciar Sesión", token).
		Render(renderArgs(c))
}

func LoginUser(c echo.Context) error {
	var payload dtos.LoginUserForm
	if err := c.Bind(&payload); err != nil {
		return err
	}

	user, err := api.Login(c.Request().Context(), env.GetString("API_URL", "http://localhost:8080/api/v1/"), dtos.LoginRequest{
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		return err
	}

	if err := createAuthSession(c, user, false); err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func CreateUser(c echo.Context) error {
	var payload dtos.RegisterUserForm
	if err := c.Bind(&payload); err != nil {
		return err
	}

	_, err := api.Register(c.Request().Context(), env.GetString("API_URL", "http://localhost:8080/api/v1/"), dtos.RegisterRequest{
		Email:     payload.Email,
		Password:  payload.Password,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Phone:     payload.Phone,
	})
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, "/")
}

func CreateCategory(c echo.Context) error {
	var payload dtos.CreateCategoryForm
	if err := c.Bind(&payload); err != nil {
		return err
	}

	token := contexts.ExtractToken(c.Request().Context())

	_, err := api.CreateCategory(
		c.Request().Context(),
		env.GetString("API_URL", "http://localhost:8080/api/v1/"),
		dtos.CreateCategoryRequest{
			Name:        payload.Name,
			Description: payload.Description},
		token,
	)

	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/admin/dashboard/category/register")
}

func CreateProduct(c echo.Context) error {
	var payload dtos.CreateProductForm
	if err := c.Bind(&payload); err != nil {
		return err
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	var images []*multipart.FileHeader
	if form != nil && len(form.File["images"]) > 0 {
		images = append(images, form.File["images"]...)
	}

	token := contexts.ExtractToken(c.Request().Context())

	product, err := api.CreateProduct(
		c.Request().Context(),
		env.GetString("API_URL", "http://localhost:8080/api/v1/"),
		dtos.CreateProductRequest{
			Name:        payload.Name,
			Description: payload.Description,
			Price:       payload.Price,
			CategoryID:  payload.CategoryID,
			Stock:       payload.Stock,
			SKU:         slug.Make(payload.Name),
		},
		token,
	)

	if err != nil {
		return err
	}

	if len(images) > 0 {
		_, err := api.AddProductImages(
			c.Request().Context(),
			env.GetString("API_URL", "http://localhost:8080/api/v1/"),
			product.Product.ID,
			images,
			token,
		)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return c.Redirect(http.StatusSeeOther, "/admin/dashboard/product/register")
}

func UpdateProduct(c echo.Context) error {
	var payload dtos.CreateProductForm
	if err := c.Bind(&payload); err != nil {
		return err
	}
	token := contexts.ExtractToken(c.Request().Context())
	_, err := api.UpdateProduct(c.Request().Context(), env.GetString("API_URL", "http://localhost:8080/api/v1/"), token, payload.ID, dtos.UpdateProductRequest{
		Name:        payload.Name,
		Description: payload.Description,
		Price:       payload.Price,
		CategoryID:  payload.CategoryID,
		Stock:       payload.Stock,
	})
	if err != nil {
		fmt.Print(err)
		return c.Redirect(http.StatusSeeOther, "/admin/dashboard/product/register")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/dashboard/product/register")
}

func CategoryPage(c echo.Context) error {
	token := csrf.Token(c.Request())

	c.Response().Header().Set("Cache-Control", "no-store")
	return views.RegisterCategory("Alejandrinas - Registro de Categorias", token).
		Render(renderArgs(c))
}
