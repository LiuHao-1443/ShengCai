package service

import (
	"github.com/spf13/viper"
	"shengcai/internal/repository"
	"shengcai/pkg/jwt"
	"shengcai/pkg/log"
	"shengcai/pkg/sid"
)

type Service struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     repository.Transaction
	conf   *viper.Viper
}

func NewService(
	tm repository.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
	jwt *jwt.JWT,
	conf *viper.Viper,
) *Service {
	return &Service{
		logger: logger,
		sid:    sid,
		jwt:    jwt,
		tm:     tm,
		conf:   conf,
	}
}
