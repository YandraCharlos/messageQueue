package service

import (
	"context"
	"messageQueue/models"
	"messageQueue/repository"
)

type service struct {
	qR repository.QueueRepository
}
type QueueService interface {
	DeleteQueue(ctx context.Context, key string) error
	PushMessage(ctx context.Context, key string, message models.MessageIn) error
	GetMessage(ctx context.Context, key string) ([]string, error)
}

func NewService(queueRepo repository.QueueRepository) QueueService {
	return &service{queueRepo}
}

func (b *service) DeleteQueue(ctx context.Context, key string) error {
	return b.qR.Delete(key)
}

func (b *service) PushMessage(ctx context.Context, key string, message models.MessageIn) error {
	return b.qR.PushMessage(key, message)
}

func (b *service) GetMessage(ctx context.Context, key string) ([]string, error) {
	message, err := b.qR.GetMessagesFromKeys(key)
	if err != nil {
		return message, err
	}
	return message, nil
}
