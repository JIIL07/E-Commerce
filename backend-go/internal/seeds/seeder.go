package seeds

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"ecommerce-backend/internal/config"
	"ecommerce-backend/internal/database"
	"ecommerce-backend/internal/utils"
)

type Seeder interface {
	Name() string
	Seed(db *sql.DB) error
	Priority() int
}

type SeedManager struct {
	db      *sql.DB
	seeders []Seeder
	logger  *utils.Logger
	config  *config.AppConfig
}

func NewSeedManager() (*SeedManager, error) {
	cfg, err := config.LoadConfig("")
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	err = database.InitDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	db := database.GetDB()
	logger := utils.NewLogger(utils.INFO, os.Stdout)

	utils.InitJWT(cfg.JWT.Secret, cfg.JWT.ExpiresIn, cfg.JWT.RefreshIn, cfg.JWT.Issuer, cfg.JWT.Audience)

	sm := &SeedManager{
		db:      db,
		logger:  logger,
		config:  cfg,
		seeders: make([]Seeder, 0),
	}

	sm.registerSeeders()

	return sm, nil
}

func (sm *SeedManager) registerSeeders() {
	sm.seeders = append(sm.seeders, &CategorySeeder{})
	sm.seeders = append(sm.seeders, &ProductSeeder{})
	sm.seeders = append(sm.seeders, &UserSeeder{})
	sm.seeders = append(sm.seeders, &OrderSeeder{})
	sm.seeders = append(sm.seeders, &ReviewSeeder{})
}
func (sm *SeedManager) Run() error {
	sm.logger.Info("Starting database seeding...")

	sm.sortSeedersByPriority()

	startTime := time.Now()
	successCount := 0
	errorCount := 0

	for _, seeder := range sm.seeders {
		sm.logger.Info("Seeding", "type", seeder.Name())

		seederStart := time.Now()
		if err := seeder.Seed(sm.db); err != nil {
			sm.logger.Error("Failed to seed", "type", seeder.Name(), "error", err)
			errorCount++
			continue
		}

		duration := time.Since(seederStart)
		sm.logger.Info("Successfully seeded", "type", seeder.Name(), "duration", duration)
		successCount++
	}

	totalDuration := time.Since(startTime)
	sm.logger.Info("Seeding completed",
		"success", successCount,
		"errors", errorCount,
		"total_duration", totalDuration)

	if errorCount > 0 {
		return fmt.Errorf("seeding completed with %d errors", errorCount)
	}

	return nil
}

func (sm *SeedManager) RunSpecific(seederNames []string) error {
	sm.logger.Info("Starting specific seeding", "seeders", seederNames)

	nameMap := make(map[string]bool)
	for _, name := range seederNames {
		nameMap[name] = true
	}

	var filteredSeeders []Seeder
	for _, seeder := range sm.seeders {
		if nameMap[seeder.Name()] {
			filteredSeeders = append(filteredSeeders, seeder)
		}
	}

	if len(filteredSeeders) == 0 {
		return fmt.Errorf("no seeders found for names: %v", seederNames)
	}

	sm.seeders = filteredSeeders
	sm.sortSeedersByPriority()

	return sm.Run()
}

func (sm *SeedManager) ListAvailableSeeders() []string {
	var names []string
	for _, seeder := range sm.seeders {
		names = append(names, seeder.Name())
	}
	return names
}

func (sm *SeedManager) sortSeedersByPriority() {
	for i := 0; i < len(sm.seeders)-1; i++ {
		for j := i + 1; j < len(sm.seeders); j++ {
			if sm.seeders[i].Priority() > sm.seeders[j].Priority() {
				sm.seeders[i], sm.seeders[j] = sm.seeders[j], sm.seeders[i]
			}
		}
	}
}

func (sm *SeedManager) Close() error {
	return database.CloseDatabase()
}
