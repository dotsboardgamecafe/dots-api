package model

import (
	"dots-api/bootstrap"
	"errors"

	"github.com/sirupsen/logrus"
)

type (
	// Contract ...
	Contract struct {
		*bootstrap.App
	}
)

const (
	Cmd = "command"
)

func (c *Contract) errHandler(funcName string, err error, returnMsg string) error {
	c.Log.FromDefault().WithFields(logrus.Fields{
		"functionName": funcName,
		"error":        err,
	}).Errorf("Error message : %s", err.Error())
	return errors.New(returnMsg)
}
