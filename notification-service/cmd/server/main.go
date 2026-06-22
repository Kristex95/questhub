package main

import (
	"context"
	"log"
	"net"

	notificationv1 "github.com/questhub/notification/gen/notification/v1"
	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotificationServer struct {
	notificationv1.UnimplementedNotificationServiceServer
	store *store
}

func NewNotificationServer(s *store) *NotificationServer {
	return &NotificationServer{store: s}
}

func (s *NotificationServer) CreateNotification(ctx context.Context, req *notificationv1.CreateNotificationRequest) (*notificationv1.Notification, error) {
	if req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.GetMessage() == "" {
		return nil, status.Error(codes.InvalidArgument, "message is required")
	}

	n := s.store.create(req.GetUserId(), req.GetMessage())
	log.Printf("created notification id=%d for user=%d", n.GetId(), n.GetUserId())
	return n, nil
}

func (s *NotificationServer) ListNotifications(ctx context.Context, req *notificationv1.ListNotificationsRequest) (*notificationv1.ListNotificationsResponse, error) {
	if req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	items := s.store.listByUser(req.GetUserId())
	return &notificationv1.ListNotificationsResponse{Notifications: items}, nil
}

func main() {
	const addr = ":50051"

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen %s: %v", addr, err)
	}

	grpcServer := grpc.NewServer()
	notificationv1.RegisterNotificationServiceServer(grpcServer, NewNotificationServer(newStore()))
	reflection.Register(grpcServer)

	log.Printf("notification gRPC server listening on %s", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
