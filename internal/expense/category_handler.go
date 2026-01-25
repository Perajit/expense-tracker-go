package expense

import (
	"errors"

	"github.com/Perajit/expense-tracker-go/internal/apperror"
	"github.com/Perajit/expense-tracker-go/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type CategoryHandler struct {
	categoryService CategoryService
	validate        *validator.Validate
}

func NewCategoryHandler(categoryService CategoryService, validate *validator.Validate) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		validate:        validate,
	}
}

func (h *CategoryHandler) RegisterRoutes(app *fiber.App) {
	group := app.Group("/expenses")
	group.Get("/", h.GetCategories)
	group.Get("/:id", h.GetCategoryByID)
	group.Post("/", h.CreateCategory)
	group.Patch("/:id", h.UpdateCategory)
	group.Delete("/:id", h.DeleteCategory)
}

func (h *CategoryHandler) GetCategories(c *fiber.Ctx) error {
	authUserID, errUserID := util.GetAuthUserID(c)
	if errUserID != nil {
		log.Error(errUserID)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": apperror.ErrUnauthorized.Error()})
	}

	categories, err := h.categoryService.GetCategories(authUserID)
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": apperror.ErrDefault.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(categories)
}

func (h *CategoryHandler) GetCategoryByID(c *fiber.Ctx) error {
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

	category, err := h.categoryService.GetCategoryByID(id, &authUserID)
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": apperror.ErrNotFound.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(category)
}

func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	userID, errUserID := util.GetAuthUserID(c)
	if errUserID != nil {
		log.Error(errUserID)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": apperror.ErrUnauthorized.Error()})
	}

	dto, errDTO := util.ExtractDto[CreateCategoryRequest](c, h.validate)
	if errDTO != nil {
		log.Error(errDTO)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": apperror.ErrInvalidRequest.Error()})
	}

	category, err := h.categoryService.CreateCategory(userID, dto)
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": apperror.ErrDefault.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(category)
}

func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
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

	dto, errDTO := util.ExtractDto[UpdateCategoryRequest](c, h.validate)
	if errDTO != nil {
		log.Error(errDTO)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": apperror.ErrInvalidRequest.Error()})
	}

	if err := h.categoryService.UpdateCategory(id, authUserID, dto); err != nil {
		log.Error(err)
		if errors.Is(err, apperror.ErrUnauthorized) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": apperror.ErrDefault.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
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

	if err := h.categoryService.DeleteCategory(id, authUserID); err != nil {
		log.Error(err)
		if errors.Is(err, apperror.ErrUnauthorized) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": apperror.ErrDefault.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}
