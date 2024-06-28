package user

import (
	"fmt"
	"realTime/server/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service
	utils.XValidator
}

func NewHandler(s Service, x utils.XValidator) *Handler {
	return &Handler{
		s,
		x,
	}
}

func (h *Handler) RegisterRoute(route fiber.Router) {
	route.Post("/signup", h.createUsers)
	route.Post("/login", h.loginUser)
	route.Get("/logout", h.logOutUser)
}

func (h *Handler) createUsers(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")
	var Regis RegisUserRequest

	if err := c.BodyParser(&Regis); err != nil {
		return utils.WriteJson(c, 500, err.Error(), nil)
	}

	if errs := h.XValidator.Validate(&Regis); len(errs) > 0 && errs[0].Error {
		errMsg := make([]string, 0)
		for _, err := range errs {
			errMsg = append(errMsg, fmt.Sprintf("[%s:%v | need to implement '%s']", err.FailedField, err.Value, err.Tag))
		}
		return utils.WriteJson(c, 400, "validation failed ON :", errMsg)
	}

	response, err := h.Service.CreateUsers(c.Context(), &Regis)
	if err != nil {
		return utils.WriteJson(c, 400, err.Error(), nil)
	}

	return utils.WriteJson(c, 200, "success crete users", response)
}

func (h *Handler) loginUser(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")

	var payload LoginUserRequest
	if err := c.BodyParser(&payload); err != nil {
		return utils.WriteJson(c, 500, err.Error(), nil)
	}

	response, err := h.Service.LoginUser(c.Context(), &payload)
	if err != nil {
		return utils.WriteJson(c, 400, err.Error(), nil)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    response.accessToken,
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
		MaxAge:   3600,
		Path:     "/",
	})

	return utils.WriteJson(c, 200, "login success", &LoginUserResponse{
		ID:       response.ID,
		Username: response.Username,
	})

}

func (h *Handler) logOutUser(c *fiber.Ctx) error {

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Expires:  time.Now().Add(-1 * time.Hour),
		MaxAge:   -1,
		Secure:   false,
		HTTPOnly: true,
	})

	return utils.WriteJson(c, 200, "cookie deleted", nil)
}
