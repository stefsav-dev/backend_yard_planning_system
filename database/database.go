package database

import (
	"log"
	"os"

	"backend_yard_planning_system/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=yard_planning_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	log.Println("Connected to database")

	// Auto migrate
	err = db.AutoMigrate(
		&models.Yard{},
		&models.Block{},
		&models.YardPlan{},
		&models.Container{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	log.Println("Database migrated successfully")

	seedInitialData()
}

func seedInitialData() {
	var yardCount int64
	DB.Model(&models.Yard{}).Count(&yardCount)

	if yardCount == 0 {

		yard := models.Yard{Name: "YRD1"}
		DB.Create(&yard)

		block := models.Block{
			YardID:  yard.ID,
			Name:    "LC01",
			MaxSlot: 10,
			MaxRow:  5,
			MaxTier: 4,
		}
		DB.Create(&block)

		plans := []models.YardPlan{
			{
				BlockID:           block.ID,
				ContainerSize:     20,
				ContainerHeight:   8.6,
				ContainerType:     "DRY",
				StartSlot:         1,
				EndSlot:           3,
				StartRow:          1,
				EndRow:            5,
				PriorityDirection: "LEFT_TO_RIGHT",
			},
			{
				BlockID:           block.ID,
				ContainerSize:     40,
				ContainerHeight:   8.6,
				ContainerType:     "DRY",
				StartSlot:         4,
				EndSlot:           7,
				StartRow:          1,
				EndRow:            5,
				PriorityDirection: "LEFT_TO_RIGHT",
			},

			{
				BlockID:           block.ID,
				ContainerSize:     20,
				ContainerHeight:   9.6,
				ContainerType:     "DRY",
				StartSlot:         8,
				EndSlot:           10,
				StartRow:          1,
				EndRow:            3,
				PriorityDirection: "LEFT_TO_RIGHT",
			},
		}

		for _, plan := range plans {
			DB.Create(&plan)
		}

		log.Println("Initial data seeded successfully")
	}
}
