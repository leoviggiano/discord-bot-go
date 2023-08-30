package service

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"bot_test/entity"
)

var (
	ErrCommandExpired = errors.New("command expired")
)

type CustomID interface {
	Set(payload entity.Payload) (string, error)
	Update(payload entity.Payload) error
	Get(key string, target entity.Payload) error
}

type customID struct {
	customIDs map[string][]byte
	log       logrus.FieldLogger
}

func NewCustomID(logger logrus.FieldLogger) CustomID {
	return customID{
		customIDs: make(map[string][]byte),
		log:       logger,
	}
}

func (s customID) Set(payload entity.Payload) (string, error) {
	key := uuid.New().String()
	payload.SetKey(key)

	jsonValue, err := json.Marshal(payload)
	if err != nil {
		s.log.Errorf("[Set] Error in set CustomID [%s]:[%s]\n%s", key, payload, err)
		return "", err
	}

	s.customIDs[key] = jsonValue
	return key, nil
}

func (s customID) Update(payload entity.Payload) error {
	if payload == nil {
		return nil
	}

	key := payload.GetKey()

	jsonValue, err := json.Marshal(payload)
	if err != nil {
		s.log.Errorf("[Set] Error in Update CustomID [%s]:[%s]\n%s", key, payload, err)
		return err
	}

	s.customIDs[key] = jsonValue

	return nil
}

func (s customID) Get(key string, target entity.Payload) error {
	customID, ok := s.customIDs[key]
	if !ok {
		return ErrCommandExpired
	}

	return json.Unmarshal(customID, &target)
}
