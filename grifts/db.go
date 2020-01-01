package grifts

import (
	"encoding/json"

	dao "github.com/mattwhip/icenine-database/daily_bonus"
	"github.com/mattwhip/icenine-service-daily_bonus/models"
	dailyBonus "github.com/mattwhip/icenine-services/daily_bonus"
	"github.com/markbates/grift/grift"
	"github.com/pkg/errors"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		// Check for existing configuration
		existingConfs := []dao.DbConfig{}
		if err := models.DB.All(&existingConfs); err != nil {
			if err != nil {
				return errors.Wrap(err, "failed to check for existing daily bonus configs")
			}
		}
		// Create config if one does not exist
		if len(existingConfs) < 1 {
			// Create wheels json
			wheels := dailyBonus.Wheels{
				Streak0: &dailyBonus.Wheel{
					Slices: []dailyBonus.Slice{
						dailyBonus.Slice{
							Value:  10000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  20000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  30000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  20000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  50000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  15000,
							Weight: 1,
						},
					},
				},
				Streak1: &dailyBonus.Wheel{
					Slices: []dailyBonus.Slice{
						dailyBonus.Slice{
							Value:  20000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  40000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  60000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  40000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  100000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  30000,
							Weight: 1,
						},
					},
				},
				Streak2: &dailyBonus.Wheel{
					Slices: []dailyBonus.Slice{
						dailyBonus.Slice{
							Value:  30000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  60000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  90000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  60000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  150000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  45000,
							Weight: 1,
						},
					},
				},
				Streak3: &dailyBonus.Wheel{
					Slices: []dailyBonus.Slice{
						dailyBonus.Slice{
							Value:  40000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  80000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  120000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  80000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  200000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  60000,
							Weight: 1,
						},
					},
				},
				Streak4: &dailyBonus.Wheel{
					Slices: []dailyBonus.Slice{
						dailyBonus.Slice{
							Value:  50000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  100000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  150000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  100000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  250000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  75000,
							Weight: 1,
						},
					},
				},
				Streak5: &dailyBonus.Wheel{
					Slices: []dailyBonus.Slice{
						dailyBonus.Slice{
							Value:  500000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  1000000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  1500000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  1000000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  2500000,
							Weight: 1,
						},
						dailyBonus.Slice{
							Value:  750000,
							Weight: 1,
						},
					},
				},
			}
			wheelsJSON, err := json.Marshal(wheels)
			if err != nil {
				return errors.Wrap(err, "failed to marshal wheels josn")
			}
			// Create DAO
			conf := &dao.DbConfig{
				// Default to 12 hours
				ResetSeconds: 60 * 60 * 12,
				// Default to 24 hours
				StreakBreakSeconds: 60 * 60 * 24,
				WheelsJSON:         string(wheelsJSON),
			}
			if err := models.DB.Create(conf); err != nil {
				return errors.Wrap(err, "failed to create new daily bonus config")
			}
		}
		return nil
	})

})
