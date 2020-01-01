package actions

import (
	"fmt"
	"io"

	dao "github.com/mattwhip/icenine-database/daily_bonus"
	"github.com/mattwhip/icenine-service-daily_bonus/config"
	"github.com/mattwhip/icenine-service-daily_bonus/status"
	pb "github.com/mattwhip/icenine-services/generated/daily_bonus"
	"github.com/mattwhip/icenine-services/middleware"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/pop"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

// StatusHandler handles a status request
func StatusHandler(c buffalo.Context) error {
	userID := c.Value(middleware.CtxKeyUserID).(string)
	tx := c.Value("tx").(*pop.Connection)

	// Get user
	users := []dao.DbUser{}
	if err := tx.Where("u_id = ?", userID).All(&users); err != nil {
		return errors.Wrapf(err, "failed to find user with userID %v", userID)
	}
	if len(users) == 0 {
		return fmt.Errorf("failed to find daily bonus user with ID %v", userID)
	}
	if len(users) > 1 {
		return fmt.Errorf("found multiple users with UserID %v", userID)
	}
	user := &users[0]

	// Get daily bonus config
	conf, err := config.Get()
	if err != nil {
		return errors.Wrap(err, "failed to get daily bonus config")
	}

	// Get the current status and wheel
	status := status.GetStatus(conf, user)
	wheel := conf.GetWheel(int(status.Streak))

	// Render successful response with protobuf payload
	return c.Render(200, r.Func("application/proto", func(w io.Writer, d render.Data) error {
		serializedProto, err := proto.Marshal(&pb.StatusResponse{
			Status: status,
			Wheel: &pb.Wheel{
				Values: wheel.GetAwardValues(),
			},
		})
		if err != nil {
			return errors.Wrap(err, "failed to serialize protobuf")
		}
		_, err = w.Write(serializedProto)
		return err
	}))
}
