package controllers

import (
	"backend_yard_planning_system/database"
	"backend_yard_planning_system/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"backend_yard_planning_system/dto"
	"backend_yard_planning_system/services"
)

type YardController struct {
	yardService *services.YardService
	validate    *validator.Validate
}

func NewYardController() *YardController {
	validate := validator.New()
	dto.RegisterCustomValidations(validate)

	return &YardController{
		yardService: services.NewYardService(),
		validate:    validate,
	}
}

func (c *YardController) GetSuggestion(ctx *fiber.Ctx) error {
	var req dto.SuggestionRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	if err := c.validate.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": dto.GetValidationError(err),
		})
	}

	response, err := c.yardService.GetSuggestion(req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(response)
}

func (c *YardController) PlaceContainer(ctx *fiber.Ctx) error {
	var req dto.PlacementRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	if err := c.validate.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": dto.GetValidationError(err),
		})
	}

	if err := c.yardService.PlaceContainer(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(dto.MessageResponse{
		Message: "Success",
	})
}

func (c *YardController) PickupContainer(ctx *fiber.Ctx) error {
	var req dto.PickupRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	if err := c.validate.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": dto.GetValidationError(err),
		})
	}

	if err := c.yardService.PickupContainer(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(dto.MessageResponse{
		Message: "Success",
	})
}

func (c *YardController) GetYardPlans(ctx *fiber.Ctx) error {
	yardName := ctx.Query("yard", "YRD1")

	var yardPlans []models.YardPlan
	err := database.DB.
		Joins("JOIN blocks ON blocks.id = yard_plans.block_id").
		Joins("JOIN yards ON yards.id = blocks.yard_id").
		Where("yards.name = ?", yardName).
		Preload("Block").
		Find(&yardPlans).Error

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"yard":  yardName,
		"plans": yardPlans,
	})
}
