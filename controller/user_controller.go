package controller

import (
	"echo-todo-app/model"
	"echo-todo-app/usecase"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

type UserController interface {
	SignUp(ctx echo.Context) error
	Login(ctx echo.Context) error
	LogOut(ctx echo.Context) error
	CsrfToken(ctx echo.Context) error
}

type userController struct {
	uu usecase.UserUsecaser
}

func NewUserController(uu usecase.UserUsecaser) UserController {
	return &userController{uu}
}

func (uc *userController) SignUp(ctx echo.Context) error {
	user := model.User{}
	if err := ctx.Bind(&user); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	res, err := uc.uu.SignUp(user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated, res)
}

func (uc *userController) Login(ctx echo.Context) error {
	user := model.User{}
	if err := ctx.Bind(&user); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	token, err := uc.uu.Login(user)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	ctx.SetCookie(cookie)
	return ctx.NoContent(http.StatusOK)
}

func (uc *userController) LogOut(ctx echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	ctx.SetCookie(cookie)
	return ctx.NoContent(http.StatusOK)
}

func (uc *userController) CsrfToken(ctx echo.Context) error {
	token := ctx.Get("csrf").(string)
	return ctx.JSON(http.StatusOK, echo.Map{
		"csrf_token": token,
	})
}
