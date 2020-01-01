package rpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	dao "github.com/mattwhip/icenine-database/daily_bonus"
	"github.com/mattwhip/icenine-service-daily_bonus/config"
	"github.com/mattwhip/icenine-service-daily_bonus/models"
	"github.com/mattwhip/icenine-service-daily_bonus/status"
	pb "github.com/mattwhip/icenine-services/generated/daily_bonus"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"

	"google.golang.org/grpc"
)

// Serve starts the GRPC server
func Serve() error {
	dailyBonusServer := &Server{}
	grpcServer := grpc.NewServer()
	pb.RegisterDailyBonusServer(grpcServer, dailyBonusServer)
	rpcListenPort := os.Getenv("RPC_LISTEN_PORT")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", rpcListenPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	return grpcServer.Serve(listener)
}

// Server implements protobuf generated DailyBonusServer interface.
type Server struct{}

// GetStatus gets the daily bonus status for a given user ID
func (s *Server) GetStatus(ctx context.Context, req *pb.UserRequest) (*pb.StatusResponse, error) {
	// Get configuration
	conf, err := config.Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get daily bonus config")
	}

	userID := req.UserID
	var statusResponse *pb.StatusResponse
	if err := models.DB.Transaction(func(tx *pop.Connection) error {
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

		// Get the current status and wheel
		status := status.GetStatus(conf, user)
		wheel := conf.GetWheel(int(status.Streak))

		// Create status response
		statusResponse = &pb.StatusResponse{
			Status: status,
			Wheel: &pb.Wheel{
				Values: wheel.GetAwardValues(),
			},
		}
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "failed to execute transaction to retrieve user status")
	}
	return statusResponse, nil
}

// InitNewUser initializes the daily bonus for a new user
func (s *Server) InitNewUser(ctx context.Context, req *pb.UserRequest) (*pb.StatusResponse, error) {
	user := &dao.DbUser{
		UID:        req.UserID,
		LastPlayed: time.Time{},
		Streak:     0,
	}
	if err := models.DB.Create(user); err != nil {
		return nil, errors.Wrapf(err, "failed to initialize new Daily Bonus user with UserID %v", req.UserID)
	}
	// Get configuration
	conf, err := config.Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get daily bonus config")
	}
	wheel := conf.GetWheel(0)
	return &pb.StatusResponse{
		Status: &pb.Status{
			IsAvailable:           true,
			SecondsUntilAvailable: 0,
			Streak:                user.Streak,
		},
		Wheel: &pb.Wheel{
			Values: wheel.GetAwardValues(),
		},
	}, nil
}
