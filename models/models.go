package models

import (
	"time"
)

type Yard struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null" json:"name"`
	Blocks    []Block   `json:"blocks,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Block struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	YardID     uint        `gorm:"not null" json:"yard_id"`
	Yard       Yard        `gorm:"foreignKey:YardID" json:"yard,omitempty"`
	Name       string      `gorm:"not null" json:"name"`
	MaxSlot    int         `gorm:"not null" json:"max_slot"`
	MaxRow     int         `gorm:"not null" json:"max_row"`
	MaxTier    int         `gorm:"not null" json:"max_tier"`
	Plans      []YardPlan  `json:"plans,omitempty"`
	Containers []Container `json:"containers,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type YardPlan struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	BlockID           uint      `gorm:"not null" json:"block_id"`
	Block             Block     `gorm:"foreignKey:BlockID" json:"block,omitempty"`
	ContainerSize     int       `gorm:"not null" json:"container_size"`   // 20 or 40
	ContainerHeight   float64   `gorm:"not null" json:"container_height"` // 8.6 or 9.6
	ContainerType     string    `gorm:"not null" json:"container_type"`   // DRY, REEFER, OPEN_TOP
	StartSlot         int       `gorm:"not null" json:"start_slot"`
	EndSlot           int       `gorm:"not null" json:"end_slot"`
	StartRow          int       `gorm:"not null" json:"start_row"`
	EndRow            int       `gorm:"not null" json:"end_row"`
	PriorityDirection string    `gorm:"not null" json:"priority_direction"` // LEFT_TO_RIGHT, BOTTOM_TO_TOP
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type Container struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	ContainerNumber string     `gorm:"uniqueIndex;not null" json:"container_number"`
	BlockID         uint       `gorm:"not null" json:"block_id"`
	Block           Block      `gorm:"foreignKey:BlockID" json:"block,omitempty"`
	ContainerSize   int        `gorm:"not null" json:"container_size"`
	ContainerHeight float64    `gorm:"not null" json:"container_height"`
	ContainerType   string     `gorm:"not null" json:"container_type"`
	Slot            int        `gorm:"not null" json:"slot"`
	Row             int        `gorm:"not null" json:"row"`
	Tier            int        `gorm:"not null" json:"tier"`
	IsPlaced        bool       `gorm:"not null;default:true" json:"is_placed"`
	PlacedAt        time.Time  `json:"placed_at"`
	PickedUpAt      *time.Time `json:"picked_up_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
