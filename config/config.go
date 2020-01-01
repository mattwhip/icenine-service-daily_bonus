package config

import (
	"encoding/json"
	"sync"

	dao "github.com/mattwhip/icenine-database/daily_bonus"
	"github.com/mattwhip/icenine-service-daily_bonus/models"
	dailyBonus "github.com/mattwhip/icenine-services/daily_bonus"
	"github.com/pkg/errors"
)

const (
	// MaxStreak is the maximum daily bonus streak
	MaxStreak = 5
)

// Get retrieves the default config
func Get() (dailyBonus.Config, error) {
	if err := lazyInit(); err != nil {
		return dailyBonus.Config{}, errors.Wrap(err, "failed to lazy init")
	}
	return cachedConfig, nil
}

var cachedConfig dailyBonus.Config
var initialized bool
var configMutex *sync.Mutex

func init() {
	initialized = false
	configMutex = &sync.Mutex{}
}

func loadConfig() error {
	// Load config from database
	c := &dao.DbConfig{}
	if err := models.DB.First(c); err != nil {
		return errors.Wrap(err, "failed to load DbConfig from database")
	}
	// Parse wheels JSON
	wheels := &dailyBonus.Wheels{}
	if err := json.Unmarshal([]byte(c.WheelsJSON), wheels); err != nil {
		return errors.Wrap(err, "failed to unmarshal wheels json")
	}
	// Initialize population for sampling of each wheel
	allWheels := []*dailyBonus.Wheel{
		wheels.Streak0, wheels.Streak1, wheels.Streak2, wheels.Streak3, wheels.Streak4, wheels.Streak5}
	for wheelIndex, wheel := range allWheels {
		if err := wheel.InitializePopulation(); err != nil {
			return errors.Wrapf(err, "failed to create population for wheel at streak %v", wheelIndex)
		}
	}
	// Create config
	cachedConfig = dailyBonus.Config{
		ResetSeconds:       c.ResetSeconds,
		StreakBreakSeconds: c.StreakBreakSeconds,
		Wheels:             wheels,
	}
	return nil
}

func lazyInit() error {
	if !initialized {
		configMutex.Lock()
		defer configMutex.Unlock()
		if !initialized {
			// Load config from database
			if err := loadConfig(); err != nil {
				return errors.Wrap(err, "failed to load initial config")
			}
			// Flag initialized
			initialized = true
		}
	}
	return nil
}
