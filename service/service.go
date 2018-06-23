package service

import (
	"time"

	"golang.org/x/net/context"

	"github.com/RobertGrantEllis/t9/logger"
	"github.com/RobertGrantEllis/t9/proto"
	"github.com/RobertGrantEllis/t9/t9"
)

type Service interface {
	proto.T9Server
}

func New(t9 t9.T9, logger logger.Logger) Service {

	return &service{
		t9:     t9,
		logger: logger,
	}
}

type service struct {
	t9     t9.T9
	logger logger.Logger
}

func (service *service) Lookup(ctx context.Context, request *proto.LookupRequest) (*proto.LookupResponse, error) {

	// initialize response
	response := &proto.LookupResponse{
		Digits: request.Digits,
		Exact:  request.Exact,
	}

	start := time.Now()
	words, err := service.t9.GetWords(request.Digits, request.Exact)
	turnAroundTime := time.Now().Sub(start)

	if err != nil {
		response.Status = false
		response.Message = err.Error()
	} else {
		response.Status = true
		response.Words = words
	}

	service.logResponse(response, turnAroundTime)

	return response, nil
}

func (service *service) logResponse(response *proto.LookupResponse, turnAroundTime time.Duration) {

	if response.Status {
		service.logger.Infof(
			`turnaround=%s | digits=%s | exact=%t | words=%d`,
			turnAroundTime, response.Digits, response.Exact, len(response.Words))
	} else {
		service.logger.Warnf(
			`turnaround=%s | digits=%s | exact=%t | message=%s`,
			turnAroundTime, response.Digits, response.Exact, response.Message)
	}
}
