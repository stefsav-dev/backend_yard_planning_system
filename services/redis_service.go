// services/redis_service.go
package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"backend_yard_planning_system/config"
	"backend_yard_planning_system/models"
)

type RedisService struct{}

func NewRedisService() *RedisService {
	return &RedisService{}
}

// Cache keys
const (
	CacheKeyYardPlans      = "yard_plans:%s"             // yard_plans:YRD1
	CacheKeyBlockOccupancy = "block_occupancy:%s"        // block_occupancy:LC01
	CacheKeyContainer      = "container:%s"              // container:ALFI000001
	CacheKeySuggestions    = "suggestions:%s:%d:%.1f:%s" // suggestions:YRD1:20:8.6:DRY
)

// Cache durations
var (
	CacheDurationYardPlans      = 5 * time.Minute
	CacheDurationBlockOccupancy = 2 * time.Minute
	CacheDurationContainer      = 10 * time.Minute
	CacheDurationSuggestions    = 1 * time.Minute
)

func (r *RedisService) GetYardPlans(yardName string) ([]models.YardPlan, error) {
	cacheKey := fmt.Sprintf(CacheKeyYardPlans, yardName)

	// Try to get from cache
	cachedData, err := config.RedisClient.Get(config.Ctx, cacheKey).Result()
	if err == nil {
		var plans []models.YardPlan
		if err := json.Unmarshal([]byte(cachedData), &plans); err == nil {
			log.Printf("✅ Cache HIT for yard plans: %s", yardName)
			return plans, nil
		}
	}

	log.Printf("❌ Cache MISS for yard plans: %s", yardName)
	return nil, err
}

func (r *RedisService) SetYardPlans(yardName string, plans []models.YardPlan) error {
	cacheKey := fmt.Sprintf(CacheKeyYardPlans, yardName)

	jsonData, err := json.Marshal(plans)
	if err != nil {
		return err
	}

	err = config.RedisClient.Set(config.Ctx, cacheKey, jsonData, CacheDurationYardPlans).Err()
	if err != nil {
		log.Printf("❌ Failed to cache yard plans: %v", err)
		return err
	}

	log.Printf("✅ Cached yard plans: %s", yardName)
	return nil
}

func (r *RedisService) GetBlockOccupancy(blockID uint) (map[string]bool, error) {
	cacheKey := fmt.Sprintf(CacheKeyBlockOccupancy, fmt.Sprintf("%d", blockID))

	cachedData, err := config.RedisClient.Get(config.Ctx, cacheKey).Result()
	if err == nil {
		var occupancy map[string]bool
		if err := json.Unmarshal([]byte(cachedData), &occupancy); err == nil {
			log.Printf("✅ Cache HIT for block occupancy: %d", blockID)
			return occupancy, nil
		}
	}

	log.Printf("❌ Cache MISS for block occupancy: %d", blockID)
	return nil, err
}

func (r *RedisService) SetBlockOccupancy(blockID uint, occupancy map[string]bool) error {
	cacheKey := fmt.Sprintf(CacheKeyBlockOccupancy, fmt.Sprintf("%d", blockID))

	jsonData, err := json.Marshal(occupancy)
	if err != nil {
		return err
	}

	err = config.RedisClient.Set(config.Ctx, cacheKey, jsonData, CacheDurationBlockOccupancy).Err()
	if err != nil {
		log.Printf("❌ Failed to cache block occupancy: %v", err)
		return err
	}

	log.Printf("✅ Cached block occupancy: %d", blockID)
	return nil
}

func (r *RedisService) InvalidateBlockOccupancy(blockID uint) {
	cacheKey := fmt.Sprintf(CacheKeyBlockOccupancy, fmt.Sprintf("%d", blockID))
	err := config.RedisClient.Del(config.Ctx, cacheKey).Err()
	if err != nil {
		log.Printf("❌ Failed to invalidate block occupancy cache: %v", err)
	} else {
		log.Printf("✅ Invalidated block occupancy cache: %d", blockID)
	}
}

func (r *RedisService) InvalidateYardPlans(yardName string) {
	cacheKey := fmt.Sprintf(CacheKeyYardPlans, yardName)
	err := config.RedisClient.Del(config.Ctx, cacheKey).Err()
	if err != nil {
		log.Printf("❌ Failed to invalidate yard plans cache: %v", err)
	} else {
		log.Printf("✅ Invalidated yard plans cache: %s", yardName)
	}
}

func (r *RedisService) CacheSuggestion(yardName string, containerSize int, containerHeight float64, containerType string, position interface{}) error {
	cacheKey := fmt.Sprintf(CacheKeySuggestions, yardName, containerSize, containerHeight, containerType)

	jsonData, err := json.Marshal(position)
	if err != nil {
		return err
	}

	err = config.RedisClient.Set(config.Ctx, cacheKey, jsonData, CacheDurationSuggestions).Err()
	if err != nil {
		log.Printf("❌ Failed to cache suggestion: %v", err)
		return err
	}

	log.Printf("✅ Cached suggestion: %s", cacheKey)
	return nil
}

func (r *RedisService) GetCachedSuggestion(yardName string, containerSize int, containerHeight float64, containerType string) (interface{}, error) {
	cacheKey := fmt.Sprintf(CacheKeySuggestions, yardName, containerSize, containerHeight, containerType)

	cachedData, err := config.RedisClient.Get(config.Ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}

	var position interface{}
	if err := json.Unmarshal([]byte(cachedData), &position); err != nil {
		return nil, err
	}

	log.Printf("✅ Cache HIT for suggestion: %s", cacheKey)
	return position, nil
}

// Health check
func (r *RedisService) HealthCheck() error {
	_, err := config.RedisClient.Ping(config.Ctx).Result()
	return err
}
