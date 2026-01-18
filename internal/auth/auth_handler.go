package auth

import (
	"github.com/Perajit/expense-tracker-go/internal/apperror"
	"github.com/Perajit/expense-tracker-go/internal/user"
	"github.com/Perajit/expense-tracker-go/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type AuthHandler struct {
	authService AuthService
	userService user.UserService
	validate    *validator.Validate
}

func NewAuthHandler(authService AuthService, userService user.UserService, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
		validate:    validate,
	}
}

func (h *AuthHandler) RegisterRoutes(app *fiber.App) {
	group := app.Group("/auth")
	group.Post("/login", h.Login)
	group.Post("/refresh", h.Refresh)
	group.Post("/logout", h.Logout)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	dto, errDTO := util.ExtractDto[LoginRequest](c, h.validate)
	if errDTO != nil {
		log.Error(errDTO)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": apperror.ErrInvalidRequest.Error()})
	}

	tokens, err := h.authService.Login(dto, h.userService)
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": apperror.ErrInvalidCredentials.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(tokens)
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	dto, errDTO := util.ExtractDto[RefreshRequest](c, h.validate)
	if errDTO != nil {
		log.Error(errDTO)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": apperror.ErrInvalidRequest.Error()})
	}

	tokens, err := h.authService.Refresh(dto.RefreshToken)
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": apperror.ErrDefault.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(tokens)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	authUserID, errUserID := util.GetAuthUserID(c)
	if errUserID != nil {
		log.Error(errUserID)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": apperror.ErrUnauthorized.Error()})
	}

	err := h.authService.Logout(authUserID)
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": apperror.ErrDefault.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}
