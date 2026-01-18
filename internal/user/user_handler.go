package user

import (
	"errors"

	"github.com/Perajit/expense-tracker-go/internal/apperror"
	"github.com/Perajit/expense-tracker-go/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UserHandler struct {
	userService UserService
	validate    *validator.Validate
}

func NewUserHandler(userService UserService, validate *validator.Validate) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validate,
	}
}

func (h *UserHandler) RegisterRoutes(app *fiber.App, authMiddleware fiber.Handler) {
	group := app.Group("/users")
	group.Get("/:id", authMiddleware, h.GetUserByID)
	group.Post("/", h.CreateUser)
	group.Patch("/:id", authMiddleware, h.UpdateUser)
	group.Delete("/:id", authMiddleware, h.DeleteUser)
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, errID := util.ExtractIDParam(c)
	if errID != nil {
		log.Error(errID)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": apperror.ErrInvalidRequest.Error()})
	}

	authUserID, errUserID := util.GetAuthUserID(c)
	if errUserID != nil {
		log.Error(errUserID)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": apperror.ErrUnauthorized.Error()})
	}

	user, err := h.userService.GetUserByID(id, authUserID)
	if err != nil {
		if errors.Is(err, apperror.ErrUnauthorized) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": apperror.ErrNotFound.Error()})
	}

	userResponse := UserResponse{}.FromEntity(*user)

	return c.Status(fiber.StatusOK).JSON(userResponse)
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	dto, errDTO := util.ExtractDto[CreateUserRequest](c, h.validate)
	if errDTO != nil {
		log.Error(errDTO)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": apperror.ErrInvalidRequest.Error()})
	}

	user, err := h.userService.CreateUser(dto)
	if err != nil {
		log.Error(err)
		if errors.Is(err, apperror.ErrUserDuplication) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": apperror.ErrDefault.Error()})
	}

	userResponse := UserResponse{}.FromEntity(*user)

	return c.Status(fiber.StatusCreated).JSON(userResponse)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	dto, errDTO := util.ExtractDto[UpdateUserRequest](c, h.validate)
	if errDTO != nil {
		log.Error(errDTO)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errDTO.Error()})
	}

	id, errID := util.ExtractIDParam(c)
	if errID != nil {
		log.Error(errID)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": apperror.ErrInvalidRequest.Error()})
	}

	authUserID, errUserID := util.GetAuthUserID(c)
	if errUserID != nil {
		log.Error(errUserID)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": apperror.ErrUnauthorized.Error()})
	}

	err := h.userService.UpdateUser(id, authUserID, dto)
	if err != nil {
		log.Error(err)
		if errors.Is(err, apperror.ErrUnauthorized) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": apperror.ErrDefault.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, errID := util.ExtractIDParam(c)
	if errID != nil {
		log.Error(errID)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": apperror.ErrInvalidRequest.Error()})
	}

	authUserID, errUserID := util.GetAuthUserID(c)
	if errUserID != nil {
		log.Error(errUserID)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": apperror.ErrUnauthorized.Error()})
	}

	err := h.userService.DeleteUser(id, authUserID)
	if err != nil {
		log.Error(err)
		if errors.Is(err, apperror.ErrUnauthorized) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": apperror.ErrDefault.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}
