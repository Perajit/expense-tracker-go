package expense

import (
	"errors"

	"github.com/Perajit/expense-tracker-go/internal/apperror"
	"github.com/Perajit/expense-tracker-go/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type ExpenseHandler struct {
	expenseService  ExpenseService
	categoryService CategoryService
	tagService      TagService
	validate        *validator.Validate
}

func NewExpenseHandler(
	expenseService ExpenseService,
	categoryService CategoryService,
	tagService TagService,
	validate *validator.Validate,
) *ExpenseHandler {
	return &ExpenseHandler{
		expenseService:  expenseService,
		categoryService: categoryService,
		tagService:      tagService,
		validate:        validate,
	}
}

func (h *ExpenseHandler) RegisterRoutes(app *fiber.App) {
	group := app.Group("/expenses")
	group.Get("/", h.GetExpenses)
	group.Get("/:id", h.GetExpenseByID)
	group.Post("/", h.CreateExpense)
	group.Patch("/:id", h.UpdateExpense)
	group.Delete("/:id", h.DeleteExpense)
}

func (h *ExpenseHandler) GetExpenses(c *fiber.Ctx) error {
	authUserID, errUserID := util.GetAuthUserID(c)
	if errUserID != nil {
		log.Error(errUserID)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": apperror.ErrUnauthorized.Error()})
	}

	expenses, err := h.expenseService.GetExpenses(authUserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": apperror.ErrDefault.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(expenses)
}

func (h *ExpenseHandler) GetExpenseByID(c *fiber.Ctx) error {
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

	expense, err := h.expenseService.GetExpenseByID(id, authUserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": apperror.ErrNotFound.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(expense)
}

func (h *ExpenseHandler) CreateExpense(c *fiber.Ctx) error {
	authUserID, errUserID := util.GetAuthUserID(c)
	if errUserID != nil {
		log.Error(errUserID)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": apperror.ErrUnauthorized.Error()})
	}

	dto, errDTO := util.ExtractDto[CreateExpenseRequest](c, h.validate)
	if errDTO != nil {
		log.Error(errDTO)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": apperror.ErrInvalidRequest.Error()})
	}

	expense, err := h.expenseService.CreateExpense(authUserID, dto)
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": apperror.ErrDefault.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(expense)
}

func (h *ExpenseHandler) UpdateExpense(c *fiber.Ctx) error {
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

	dto, errDTO := util.ExtractDto[UpdateExpenseRequest](c, h.validate)
	if errDTO != nil {
		log.Error(errDTO)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": apperror.ErrInvalidRequest.Error()})
	}

	if err := h.expenseService.UpdateExpense(id, authUserID, dto); err != nil {
		log.Error(err)
		if errors.Is(err, apperror.ErrUnauthorized) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": apperror.ErrDefault.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func (h *ExpenseHandler) DeleteExpense(c *fiber.Ctx) error {
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

	if err := h.expenseService.DeleteExpense(id, authUserID); err != nil {
		log.Error(err)
		if errors.Is(err, apperror.ErrUnauthorized) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": apperror.ErrDefault.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}
