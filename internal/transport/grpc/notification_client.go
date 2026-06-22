package grpc

import (
	"context"
	"fmt"

	notificationv1 "github.com/questhub/notification/gen/notification/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NotificationClient struct {
	conn   *grpc.ClientConn
	client notificationv1.NotificationServiceClient
}

func NewNotificationClient(addr string) (*NotificationClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("dial notification service: %w", err)
	}

	return &NotificationClient{
		conn:   conn,
		client: notificationv1.NewNotificationServiceClient(conn),
	}, nil
}

func (c *NotificationClient) Notify(ctx context.Context, userID int64, message string) error {
	_, err := c.client.CreateNotification(ctx, &notificationv1.CreateNotificationRequest{
		UserId:  userID,
		Message: message,
	})
	if err != nil {
		return fmt.Errorf("create notification: %w", err)
	}
	return nil
}

func (c *NotificationClient) Close() error {
	return c.conn.Close()
}
