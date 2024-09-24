package model

import (
	"dots-api/bootstrap"
	"dots-api/lib/utils"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

type Contract struct {
	*bootstrap.App
}

const (
	ChannelApp = "app"
	ChannelCMS = "cms"
)

func (c *Contract) errHandler(funcName string, err error, returnMsg string) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return errors.New(utils.EmptyData)
	}

	c.Log.FromDefault().WithFields(logrus.Fields{
		"functionName": funcName,
		"error":        err,
	}).Errorf("Error message : %s", err.Error())

	return errors.New(returnMsg)
}
