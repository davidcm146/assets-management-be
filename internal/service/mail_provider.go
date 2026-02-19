package service

import (
	"context"
	"log"
)

type MailProvider interface {
	Send(ctx context.Context, to string, subject string, body string) error
}

type MultiMailProvider struct {
	primary  MailProvider
	fallback MailProvider
}

func NewMultiMailProvider(primary MailProvider, fallback MailProvider) MailProvider {
	return &MultiMailProvider{
		primary:  primary,
		fallback: fallback,
	}
}

func (m *MultiMailProvider) Send(ctx context.Context, to, subject, body string) error {

	err := m.primary.Send(ctx, to, subject, body)
	if err != nil {
		log.Println("Primary mail failed, using fallback:", err)
		return m.fallback.Send(ctx, to, subject, body)
	}

	return nil
}
