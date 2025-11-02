package services

import (
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"backend_yard_planning_system/database"
	"backend_yard_planning_system/dto"
	"backend_yard_planning_system/models"

	"gorm.io/gorm"
)

type YardService struct {
	db *gorm.DB
}

func NewYardService() *YardService {
	return &YardService{db: database.DB}
}

func (s *YardService) GetSuggestion(req dto.SuggestionRequest) (*dto.SuggestionResponse, error) {
	log.Printf("🔍 Searching yard plan for: Yard=%s, Size=%d, Height=%.1f, Type=%s",
		req.Yard, req.ContainerSize, req.ContainerHeight, req.ContainerType)

	var yard models.Yard
	if err := s.db.Where("name = ?", req.Yard).First(&yard).Error; err != nil {
		log.Printf("❌ Yard not found: %s", req.Yard)
		return nil, errors.New("yard not found")
	}
	log.Printf("✅ Found yard: %s (ID: %d)", yard.Name, yard.ID)

	var yardPlan models.YardPlan
	query := s.db.Joins("JOIN blocks ON blocks.id = yard_plans.block_id").
		Where("blocks.yard_id = ?", yard.ID).
		Where("container_size = ?", req.ContainerSize).
		Where("container_type = ?", req.ContainerType).
		Where("ABS(container_height - ?) < 0.01", req.ContainerHeight).
		Preload("Block")

	err := query.First(&yardPlan).Error

	if err != nil {
		log.Printf("❌ No exact match found. Error: %v", err)

		var allPlans []models.YardPlan
		s.db.Joins("JOIN blocks ON blocks.id = yard_plans.block_id").
			Where("blocks.yard_id = ?", yard.ID).
			Preload("Block").
			Find(&allPlans)

		log.Printf("📋 All plans in database for yard %s:", req.Yard)
		for i, plan := range allPlans {
			log.Printf("  Plan %d: Size=%d, Height=%.10f, Type=%s",
				i+1, plan.ContainerSize, plan.ContainerHeight, plan.ContainerType)

			matchSize := plan.ContainerSize == req.ContainerSize
			matchType := plan.ContainerType == req.ContainerType
			matchHeight := math.Abs(plan.ContainerHeight-req.ContainerHeight) < 0.01

			log.Printf("    Match check - Size: %t, Height: %t, Type: %t",
				matchSize, matchHeight, matchType)
		}

		return nil, fmt.Errorf("no suitable yard plan found. Check server logs for details.")
	}

	log.Printf("✅ Found matching yard plan: Block=%s, Slots=%d-%d, Rows=%d-%d",
		yardPlan.Block.Name, yardPlan.StartSlot, yardPlan.EndSlot, yardPlan.StartRow, yardPlan.EndRow)

	var existingContainer models.Container
	if err := s.db.Where("container_number = ? AND is_placed = ?", req.ContainerNumber, true).First(&existingContainer).Error; err == nil {
		log.Printf("❌ Container already placed: %s", req.ContainerNumber)
		return nil, errors.New("container is already placed in the yard")
	}

	position, err := s.findAvailablePosition(yardPlan, req.ContainerSize)
	if err != nil {
		log.Printf("❌ No available position: %v", err)
		return nil, fmt.Errorf("no available position found: %v", err)
	}

	log.Printf("🎯 Suggested position: Block=%s, Slot=%d, Row=%d, Tier=%d",
		yardPlan.Block.Name, position.Slot, position.Row, position.Tier)

	return &dto.SuggestionResponse{
		SuggestedPosition: dto.Position{
			Block: yardPlan.Block.Name,
			Slot:  position.Slot,
			Row:   position.Row,
			Tier:  position.Tier,
		},
	}, nil
}

func (s *YardService) PlaceContainer(req dto.PlacementRequest) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var block models.Block
		err := tx.Joins("JOIN yards ON yards.id = blocks.yard_id").
			Where("yards.name = ? AND blocks.name = ?", req.Yard, req.Block).
			First(&block).Error

		if err != nil {
			log.Printf("❌ Block not found: Yard=%s, Block=%s, Error: %v", req.Yard, req.Block, err)
			return errors.New("block not found in specified yard")
		}

		log.Printf("✅ Found block: %s (ID: %d) in yard: %s", block.Name, block.ID, req.Yard)

		if req.Slot < 1 || req.Slot > block.MaxSlot ||
			req.Row < 1 || req.Row > block.MaxRow ||
			req.Tier < 1 || req.Tier > block.MaxTier {
			log.Printf("❌ Invalid position: Slot=%d/%d, Row=%d/%d, Tier=%d/%d",
				req.Slot, block.MaxSlot, req.Row, block.MaxRow, req.Tier, block.MaxTier)
			return errors.New("position exceeds block capacity")
		}

		var occupiedContainer models.Container
		if err := tx.Where("block_id = ? AND slot = ? AND row = ? AND tier = ? AND is_placed = ?",
			block.ID, req.Slot, req.Row, req.Tier, true).First(&occupiedContainer).Error; err == nil {
			log.Printf("❌ Position occupied: Block=%s, Slot=%d, Row=%d, Tier=%d by Container=%s",
				block.Name, req.Slot, req.Row, req.Tier, occupiedContainer.ContainerNumber)
			return errors.New("position is already occupied")
		}

		log.Printf("✅ Position available: Block=%s, Slot=%d, Row=%d, Tier=%d",
			block.Name, req.Slot, req.Row, req.Tier)

		var existingContainer models.Container
		if err := tx.Where("container_number = ?", req.ContainerNumber).First(&existingContainer).Error; err == nil {
			log.Printf("ℹ️ Container exists, updating: %s", req.ContainerNumber)

			existingContainer.BlockID = block.ID
			existingContainer.Slot = req.Slot
			existingContainer.Row = req.Row
			existingContainer.Tier = req.Tier
			existingContainer.IsPlaced = true
			existingContainer.PlacedAt = time.Now()
			existingContainer.PickedUpAt = nil

			if err := tx.Save(&existingContainer).Error; err != nil {
				log.Printf("❌ Failed to update container: %v", err)
				return err
			}

			log.Printf("✅ Container updated and placed: %s", req.ContainerNumber)
		} else {
			log.Printf("ℹ️ Creating new container: %s", req.ContainerNumber)

			container := models.Container{
				ContainerNumber: req.ContainerNumber,
				BlockID:         block.ID,
				ContainerSize:   20,    // Default value, bisa disesuaikan
				ContainerHeight: 8.6,   // Default value, bisa disesuaikan
				ContainerType:   "DRY", // Default value, bisa disesuaikan
				Slot:            req.Slot,
				Row:             req.Row,
				Tier:            req.Tier,
				IsPlaced:        true,
				PlacedAt:        time.Now(),
			}

			if err := tx.Create(&container).Error; err != nil {
				log.Printf("❌ Failed to create container: %v", err)
				return err
			}

			log.Printf("✅ New container created and placed: %s (Size: %d, Height: %.1f, Type: %s)",
				req.ContainerNumber, container.ContainerSize, container.ContainerHeight, container.ContainerType)
		}

		return nil
	})
}

func (s *YardService) PickupContainer(req dto.PickupRequest) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var container models.Container

		err := tx.Joins("JOIN blocks ON blocks.id = containers.block_id").
			Joins("JOIN yards ON yards.id = blocks.yard_id").
			Where("containers.container_number = ? AND yards.name = ?",
				req.ContainerNumber, req.Yard).
			First(&container).Error

		if err != nil {
			log.Printf(" Container not found: Number=%s, Yard=%s, Error: %v",
				req.ContainerNumber, req.Yard, err)
			return errors.New("container not found in specified yard")
		}

		log.Printf("✅ Found container: %s in yard: %s (Placed: %t)",
			container.ContainerNumber, req.Yard, container.IsPlaced)

		if !container.IsPlaced {
			log.Printf(" Container not placed: %s", req.ContainerNumber)
			return errors.New("container is not currently placed")
		}

		now := time.Now()
		container.IsPlaced = false
		container.PickedUpAt = &now

		if err := tx.Save(&container).Error; err != nil {
			log.Printf("Failed to pickup container: %v", err)
			return err
		}

		log.Printf("Container picked up successfully: %s", req.ContainerNumber)
		return nil
	})
}

type position struct {
	Slot int
	Row  int
	Tier int
}

func (s *YardService) findAvailablePosition(plan models.YardPlan, containerSize int) (*position, error) {
	var occupiedPositions []struct {
		Slot int
		Row  int
		Tier int
	}

	s.db.Model(&models.Container{}).Where(
		"block_id = ? AND is_placed = ? AND slot BETWEEN ? AND ? AND row BETWEEN ? AND ?",
		plan.BlockID, true, plan.StartSlot, plan.EndSlot, plan.StartRow, plan.EndRow,
	).Select("slot, row, tier").Find(&occupiedPositions)

	occupiedMap := make(map[string]bool)
	for _, pos := range occupiedPositions {
		key := getPositionKey(pos.Slot, pos.Row, pos.Tier)
		occupiedMap[key] = true
	}

	switch plan.PriorityDirection {
	case "LEFT_TO_RIGHT":
		for tier := 1; tier <= plan.Block.MaxTier; tier++ {
			for row := plan.StartRow; row <= plan.EndRow; row++ {
				for slot := plan.StartSlot; slot <= plan.EndSlot; slot++ {
					if containerSize == 40 {
						if slot+1 > plan.EndSlot {
							continue
						}
						if !isPositionOccupied(occupiedMap, slot, row, tier) &&
							!isPositionOccupied(occupiedMap, slot+1, row, tier) {
							return &position{Slot: slot, Row: row, Tier: tier}, nil
						}
					} else {
						if !isPositionOccupied(occupiedMap, slot, row, tier) {
							return &position{Slot: slot, Row: row, Tier: tier}, nil
						}
					}
				}
			}
		}
	case "BOTTOM_TO_TOP":
		for slot := plan.StartSlot; slot <= plan.EndSlot; slot++ {
			for row := plan.StartRow; row <= plan.EndRow; row++ {
				for tier := 1; tier <= plan.Block.MaxTier; tier++ {
					if containerSize == 40 {
						if slot+1 > plan.EndSlot {
							continue
						}
						if !isPositionOccupied(occupiedMap, slot, row, tier) &&
							!isPositionOccupied(occupiedMap, slot+1, row, tier) {
							return &position{Slot: slot, Row: row, Tier: tier}, nil
						}
					} else {
						if !isPositionOccupied(occupiedMap, slot, row, tier) {
							return &position{Slot: slot, Row: row, Tier: tier}, nil
						}
					}
				}
			}
		}
	}

	return nil, errors.New("no available position found in the planned area")
}

func isPositionOccupied(occupiedMap map[string]bool, slot, row, tier int) bool {
	key := getPositionKey(slot, row, tier)
	return occupiedMap[key]
}

func getPositionKey(slot, row, tier int) string {
	return string(rune(slot)) + string(rune(row)) + string(rune(tier))
}
