package repository

import (
	"errors"
	"time"

	"megumin/pkg/lib/discord"
	"megumin/pkg/lib/snowflake"
	"megumin/pkg/listener/interaction/context"
)

const (
	expirationTime = 10 * time.Minute
)

var (
	ErrKeyNotFound = errors.New("customID: key not found")
)

type CustomID interface {
	Set(payload *context.Payload) string
	Get(key string) (*context.Payload, error)
}

type customID struct {
	customIDs map[string]*context.Payload
	snowflake snowflake.SnowflakeGenerator
}

func NewCustomID() CustomID {
	return customID{
		customIDs: make(map[string]*context.Payload),
		snowflake: snowflake.NewSnowflakeGenerator(),
	}
}

func (r customID) Set(payload *context.Payload) string {
	key := r.snowflake.Generate().String()

	r.customIDs[key] = payload
	payload.CustomID = discord.CustomID(key)

	time.AfterFunc(expirationTime, func() {
		delete(r.customIDs, key)
	})

	return key
}

func (r customID) Get(key string) (*context.Payload, error) {
	payload, hasKey := r.customIDs[key]
	if !hasKey {
		return nil, ErrKeyNotFound
	}

	return payload, nil
}
