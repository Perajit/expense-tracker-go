package expense

import (
	"errors"

	"github.com/Perajit/expense-tracker-go/internal/apperror"
	"github.com/Perajit/expense-tracker-go/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type TagHandler struct {
	tagService TagService
	validate   *validator.Validate
}

func NewTagHandler(tagService TagService, validate *validator.Validate) *TagHandler {
	return &TagHandler{
		tagService: tagService,
		validate:   validate,
	}
}

func (h *TagHandler) RegisterRoutes(app *fiber.App) {
	group := app.Group("/expenses")
	group.Get("/", h.GetTags)
	group.Get("/:ids", h.GetTagByIDs)
	group.Post("/", h.CreateTag)
	group.Patch("/:id", h.UpateTag)
	group.Delete("/:id", h.DeleteTag)
}

func (h *TagHandler) GetTags(c *fiber.Ctx) error {
	authUserID, errUserID := util.GetAuthUserID(c)
	if errUserID != nil {
		log.Error(errUserID)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": apperror.ErrUnauthorized.Error()})
	}

	tag, err := h.tagService.GetTags(authUserID)
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": apperror.ErrDefault.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(tag)
}

func (h *TagHandler) GetTagByIDs(c *fiber.Ctx) error {
	ids, errID := util.ExtractIDsParam(c)
	if errID != nil {
		log.Error(errID)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": apperror.ErrInvalidRequest.Error()})
	}

	authUserID, errUserID := util.GetAuthUserID(c)
	if errUserID != nil {
		log.Error(errUserID)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": apperror.ErrUnauthorized.Error()})
	}

	tag, err := h.tagService.GetTagsByIDs(ids, authUserID)
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": apperror.ErrDefault.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(tag)
}

func (h *TagHandler) CreateTag(c *fiber.Ctx) error {
	authUserID, errUserID := util.GetAuthUserID(c)
	if errUserID != nil {
		log.Error(errUserID)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": apperror.ErrUnauthorized.Error()})
	}

	dto, errDTO := util.ExtractDto[CreateTagRequest](c, h.validate)
	if errDTO != nil {
		log.Error(errDTO)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": apperror.ErrInvalidRequest.Error()})
	}

	tag, err := h.tagService.CreateTag(authUserID, dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": apperror.ErrDefault.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(tag)
}

func (h *TagHandler) UpateTag(c *fiber.Ctx) error {
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

	dto, errDTO := util.ExtractDto[UpdateTagRequest](c, h.validate)
	if errDTO != nil {
		log.Error(errDTO)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": apperror.ErrInvalidRequest.Error()})
	}

	if err := h.tagService.UpdateTag(id, authUserID, dto); err != nil {
		log.Error(err)
		if errors.Is(err, apperror.ErrUnauthorized) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": apperror.ErrDefault.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func (h *TagHandler) DeleteTag(c *fiber.Ctx) error {
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

	if err := h.tagService.DeleteTag(id, authUserID); err != nil {
		log.Error(err)
		if errors.Is(err, apperror.ErrUnauthorized) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": apperror.ErrDefault.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}
