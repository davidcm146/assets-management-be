package service

import (
	"context"

	"github.com/davidcm146/assets-management-be.git/internal/error_middleware"
	"github.com/davidcm146/assets-management-be.git/internal/model"
	"github.com/davidcm146/assets-management-be.git/internal/repository"
)

type NotificationService interface {
	Send(ctx context.Context, n *model.Notification) error
	List(ctx context.Context, recipientID int, page int, limit int, isRead *bool) ([]*model.Notification, int, error)
	MarkAsRead(ctx context.Context, id int) error
	CountUnread(ctx context.Context, recipientID int) (int, error)
	BulkSend(ctx context.Context, notifications []*model.Notification) (*[]*model.Notification, error)
	SendEmails(ctx context.Context, notifications []*model.Notification)
}

type notificationService struct {
	repo         repository.NotificationRepository
	mailProvider MailProvider
}

func NewNotificationService(repo repository.NotificationRepository, mailProvider MailProvider) NotificationService {
	return &notificationService{
		repo:         repo,
		mailProvider: mailProvider,
	}
}

func (s *notificationService) Send(ctx context.Context, n *model.Notification) error {
	return s.repo.Create(ctx, n)
}

func (s *notificationService) List(ctx context.Context, recipientID int, page int, limit int, isRead *bool) ([]*model.Notification, int, error) {
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	items, total, err := s.repo.ListByRecipient(ctx, recipientID, offset, limit, isRead)

	if err != nil {
		return nil, 0, error_middleware.NewInternal("Lấy danh sách thông báo thất bại")
	}
	if items == nil {
		items = []*model.Notification{}
	}

	return items, total, nil
}

func (s *notificationService) MarkAsRead(ctx context.Context, id int) error {
	return s.repo.MarkAsRead(ctx, id)
}

func (s *notificationService) CountUnread(ctx context.Context, recipientID int) (int, error) {
	return s.repo.CountUnread(ctx, recipientID)
}

func (s *notificationService) SendEmails(ctx context.Context, notifications []*model.Notification) {
	for _, n := range notifications {
		err := s.mailProvider.Send(
			ctx,
			"Chiefsecurity@wyndhamgrand-phuquoc.com",
			n.Title,
			n.Content,
		)

		if err != nil {
			error_middleware.NewInternal("Gửi mail xảy ra lỗi")
			return
		}
	}
}

func (s *notificationService) BulkSend(ctx context.Context, notifications []*model.Notification) (*[]*model.Notification, error) {

	if len(notifications) == 0 {
		empty := []*model.Notification{}
		return &empty, nil
	}

	if err := s.repo.BulkCreate(ctx, notifications); err != nil {
		return nil, error_middleware.NewInternal("Tạo thông báo hàng loạt thất bại")
	}

	go s.SendEmails(context.Background(), notifications)

	return &notifications, nil
}
