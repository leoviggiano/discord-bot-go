package service

import (
	"github.com/sirupsen/logrus"

)

type All struct {
	CustomID CustomID
}

func GetAll(logger logrus.FieldLogger) All {
	return All{
		CustomID: NewCustomID(logger),
	}
}
