package status

import (
	"time"

	dao "github.com/mattwhip/icenine-database/daily_bonus"
	dailyBonus "github.com/mattwhip/icenine-services/daily_bonus"
	pb "github.com/mattwhip/icenine-services/generated/daily_bonus"
)

const (
	// leeway is used to ensure multiple independently measured times do not cause
	// race conditions/comparison errors when checking for daily bonus availability
	availabilityComparisonLeeway time.Duration = 1 * time.Second
)

// GetStatus returns a Status for the user
func GetStatus(config dailyBonus.Config, user *dao.DbUser) *pb.Status {
	now := time.Now()
	lenientNow := now.Add(availabilityComparisonLeeway)
	lastPlayed := user.LastPlayed
	availableAt := lastPlayed.Add(time.Duration(int64(config.ResetSeconds) * time.Second.Nanoseconds()))
	isAvailable := availableAt.Before(lenientNow)
	streakBreakThreshold := lastPlayed.Add(time.Duration(int64(config.StreakBreakSeconds) * time.Second.Nanoseconds()))
	isStreakBroken := streakBreakThreshold.Before(now)
	streak := user.Streak
	secondsUntilAvailable := int32(0)
	if !isAvailable {
		diff := availableAt.Sub(now)
		secondsUntilAvailable = int32(diff.Seconds())
	}
	if isStreakBroken {
		streak = 0
	}
	return &pb.Status{
		IsAvailable:           isAvailable,
		SecondsUntilAvailable: secondsUntilAvailable,
		Streak:                streak,
	}
}
