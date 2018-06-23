# This makefile requires `mockgen`. Install with command:
# go get github.com/golang/mock/{gomock,mockgen}

all: mock_client.go mock_logger.go mock_server.go mock_t9.go

mock_client.go:
	mockgen -package mocks github.com/RobertGrantEllis/t9/client Client > mock_client.go
mock_logger.go:
	mockgen -package mocks github.com/RobertGrantEllis/t9/logger Logger > mock_logger.go
mock_server.go:
	mockgen -package mocks github.com/RobertGrantEllis/t9/server Server > mock_server.go
mock_t9.go:
	mockgen -package mocks github.com/RobertGrantEllis/t9/t9 T9 > mock_t9.go

clean:
	rm -f *.go