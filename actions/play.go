package actions

import (
	"fmt"
	"io"
	"time"

	"github.com/mattwhip/icenine-service-daily_bonus/config"
	"github.com/mattwhip/icenine-service-daily_bonus/status"

	dao "github.com/mattwhip/icenine-database/daily_bonus"
	pb "github.com/mattwhip/icenine-services/generated/daily_bonus"
	"github.com/mattwhip/icenine-services/middleware"
	userData "github.com/mattwhip/icenine-services/user_data"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/pop"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

// PlayHandler handles a play request
func PlayHandler(c buffalo.Context) error {

	// Retrieve daily bonus user/status from DB
	tx := c.Value("tx").(*pop.Connection)
	userID := c.Value(middleware.CtxKeyUserID).(string)
	users := []dao.DbUser{}

	// Select and lock row for user
	rq := fmt.Sprintf("SELECT * FROM db_users WHERE u_id = '%s' FOR UPDATE", userID)
	if err := tx.RawQuery(rq).All(&users); err != nil {
		return err
	}

	// Ensure a single user was returned from the query
	if len(users) != 1 {
		return fmt.Errorf("expected a single user for u_id '%v' but received %v", userID, len(users))
	}
	user := &users[0]

	// Get daily bonus config
	conf, err := config.Get()
	if err != nil {
		return errors.Wrap(err, "failed to get daily bonus config")
	}

	// Get current status
	status := status.GetStatus(conf, user)

	// Make sure the daily bonus is available for this user
	if !status.IsAvailable {
		return errors.New("daily bonus not available")
	}

	// Get the wheel to be played at this streak
	wheel := conf.GetWheel(int(status.Streak))

	// Sample wheel
	item, err := wheel.Population.Sample()
	if err != nil {
		return errors.Wrap(err, "failed to sample wheel population")
	}
	awardValue := item.Value.(int64)
	sliceIndex := item.Index

	// Increase the streak by one
	streak := min(status.Streak+1, config.MaxStreak)

	// Create the status to be returned after play
	newStatus := &pb.Status{
		IsAvailable:           false,
		SecondsUntilAvailable: int32(conf.ResetSeconds),
		Streak:                streak,
	}

	// Snapshot time for last played
	now := time.Now()

	// Store daily bonus user updates to database
	user.Streak = streak
	user.LastPlayed = now
	if err := tx.Save(user); err != nil {
		return errors.Wrap(err, "failed to save daily bonus user updates to database")
	}

	// Award player credits, get update player balance
	balances, err := userData.Get().DoCoinTransaction(map[string]int64{
		user.UID: awardValue,
	})
	if err != nil {
		return errors.Wrap(err, "failed to execute coin transaction for player daily bonus award")
	}

	// Render successful response with protobuf payload
	return c.Render(200, r.Func("application/proto", func(w io.Writer, d render.Data) error {
		pbresp := &pb.PlayResponse{
			Status: newStatus,
			Wheel: &pb.Wheel{
				Values: wheel.GetAwardValues(),
			},
			Index:   int32(sliceIndex),
			Balance: balances[user.UID],
		}
		serializedProto, err := proto.Marshal(pbresp)
		if err != nil {
			return errors.Wrap(err, "failed to serialize protobuf")
		}
		_, err = w.Write(serializedProto)
		return err
	}))
}

func min(a int32, b int32) int32 {
	if a > b {
		return b
	}
	return a
}
