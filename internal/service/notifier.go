package service

import (
	"context"
	"log"
)

type Notifier interface {
	Notify(ctx context.Context, userID int64, message string) error
}

type LocalNotifier struct{}

func NewLocalNotifier() *LocalNotifier {
	return &LocalNotifier{}
}

func (n *LocalNotifier) Notify(ctx context.Context, userID int64, message string) error {
	log.Printf("notification -> user %d: %s", userID, message)
	return nil
}
