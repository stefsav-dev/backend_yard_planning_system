package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidations(validate *validator.Validate) {
	validate.RegisterValidation("container_height", validateContainerHeight)
}

func validateContainerHeight(fl validator.FieldLevel) bool {
	if height, ok := fl.Field().Interface().(float64); ok {
		return height == 8.6 || height == 9.6
	}
	return false
}

type SuggestionRequest struct {
	Yard            string  `json:"yard" validate:"required"`
	ContainerNumber string  `json:"container_number" validate:"required"`
	ContainerSize   int     `json:"container_size" validate:"required,oneof=20 40"`
	ContainerHeight float64 `json:"container_height" validate:"required,container_height"`
	ContainerType   string  `json:"container_type" validate:"required,oneof=DRY REEFER OPEN_TOP"`
}

type PlacementRequest struct {
	Yard            string `json:"yard" validate:"required"`
	ContainerNumber string `json:"container_number" validate:"required"`
	Block           string `json:"block" validate:"required"`
	Slot            int    `json:"slot" validate:"required,min=1"`
	Row             int    `json:"row" validate:"required,min=1"`
	Tier            int    `json:"tier" validate:"required,min=1"`
}

type PickupRequest struct {
	Yard            string `json:"yard" validate:"required"`
	ContainerNumber string `json:"container_number" validate:"required"`
}

type Position struct {
	Block string `json:"block"`
	Slot  int    `json:"slot"`
	Row   int    `json:"row"`
	Tier  int    `json:"tier"`
}

type SuggestionResponse struct {
	SuggestedPosition Position `json:"suggested_position"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

func GetValidationError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			switch fieldError.Field() {
			case "Yard":
				return "yard is required"
			case "ContainerNumber":
				return "container_number is required"
			case "ContainerSize":
				return "container_size must be either 20 or 40"
			case "ContainerHeight":
				return "container_height must be either 8.6 or 9.6"
			case "ContainerType":
				return "container_type must be one of: DRY, REEFER, OPEN_TOP"
			case "Block":
				return "block is required"
			case "Slot":
				return "slot must be at least 1"
			case "Row":
				return "row must be at least 1"
			case "Tier":
				return "tier must be at least 1"
			default:
				return fmt.Sprintf("validation failed for field %s", fieldError.Field())
			}
		}
	}
	return err.Error()
}
