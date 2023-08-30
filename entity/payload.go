package entity

import "bot_test/lib/discord"

type Payload interface {
	GetKey() string
	SetKey(string)
	Data() BasePayload
}

type BasePayload struct {
	CustomID *discord.CustomID
	Reacters []string
}

func (bp *BasePayload) Data() BasePayload {
	return BasePayload{
		CustomID: bp.CustomID,
		Reacters: bp.Reacters,
	}
}

func (bp *BasePayload) GetKey() string {
	return bp.CustomID.PayloadKey
}

func (bp *BasePayload) SetKey(key string) {
	bp.CustomID.PayloadKey = key
}

func (gp *BasePayload) React(userID string) bool {
	if len(gp.Reacters) == 0 {
		return true
	}

	for _, id := range gp.Reacters {
		if id == userID {
			return true
		}
	}

	return false
}
