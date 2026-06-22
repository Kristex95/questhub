package main

import (
	"sync"

	notificationv1 "github.com/questhub/notification/gen/notification/v1"
)

type store struct {
	mu     sync.Mutex
	nextID int64
	items  []*notificationv1.Notification
}

func newStore() *store {
	return &store{nextID: 1}
}

func (s *store) create(userID int64, message string) *notificationv1.Notification {
	s.mu.Lock()
	defer s.mu.Unlock()

	n := &notificationv1.Notification{
		Id:      s.nextID,
		UserId:  userID,
		Message: message,
		IsRead:  false,
	}
	s.items = append(s.items, n)
	s.nextID++
	return n
}

func (s *store) listByUser(userID int64) []*notificationv1.Notification {
	s.mu.Lock()
	defer s.mu.Unlock()

	res := make([]*notificationv1.Notification, 0)
	for _, n := range s.items {
		if n.UserId == userID {
			res = append(res, n)
		}
	}
	return res
}
