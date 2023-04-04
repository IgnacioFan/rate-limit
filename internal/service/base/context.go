package base

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Ctx struct {
	context.Context
	logrus.FieldLogger
}

func Background() Ctx {
	return Ctx{
		Context:     context.Background(),
		FieldLogger: logrus.StandardLogger(),
	}
}
