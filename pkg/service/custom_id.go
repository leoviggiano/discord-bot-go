package service

import (
	"errors"
	"log/slog"

	"megumin/pkg/lib/discord"
	"megumin/pkg/listener/interaction/context"
	"megumin/pkg/service/internal/repository"
)

var (
	ErrCommandExpired = errors.New("command expired")
)

type CustomID interface {
	Set(payload *context.Payload) discord.CustomID
	Get(customID discord.CustomID) (*context.Payload, error)
}

type customID struct {
	customIDRepository repository.CustomID
	log                *slog.Logger
}

func NewCustomID(logger *slog.Logger) CustomID {
	return customID{
		customIDRepository: repository.NewCustomID(),
		log:                logger,
	}
}

func (s customID) Set(payload *context.Payload) discord.CustomID {
	key := s.customIDRepository.Set(payload)
	return discord.CustomID(key)
}

func (s customID) Get(customID discord.CustomID) (*context.Payload, error) {
	payload, err := s.customIDRepository.Get(customID.Key())
	if err != nil {
		return nil, err
	}

	if event := customID.Event(); event != "" {
		payload.Event = event
	}

	if payload.ForceEvent != "" {
		payload.Event = payload.ForceEvent
	}

	return payload, nil
}
