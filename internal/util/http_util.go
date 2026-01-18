package util

import (
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ValidationError struct {
	Status int
	Err    error
}

func (e *ValidationError) Error() string {
	return e.Err.Error()
}

func ExtractDto[T any](c *fiber.Ctx, validate *validator.Validate) (T, error) {
	var dto T
	if err := c.BodyParser(&dto); err != nil {
		return dto, err
	}

	if err := validate.Struct(dto); err != nil {
		return dto, err
	}

	return dto, nil
}

func ExtractIDParam(c *fiber.Ctx) (uint, error) {
	param := c.Params("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil || id <= 0 {
		return 0, err
	}

	return uint(id), nil
}

func ExtractIDsParam(c *fiber.Ctx) ([]uint, error) {
	param := c.Params("ids")
	idStrs := strings.Split(param, ",")
	ids := []uint{}

	for range idStrs {
		id, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return nil, err
		}

		ids = append(ids, uint(id))
	}

	return ids, nil
}
