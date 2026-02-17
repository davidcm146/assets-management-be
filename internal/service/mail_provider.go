package service

import "context"

type MailProvider interface {
	Send(ctx context.Context, to string, subject string, body string) error
}
