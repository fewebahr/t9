package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"google.golang.org/grpc"

	"github.com/fewebahr/t9/assets"
	"github.com/fewebahr/t9/proto"
	"github.com/fewebahr/t9/service"
	"github.com/fewebahr/t9/t9"
)

func (s *server) instantiateAndRegisterGrpcHandler() error {
	t9, err := s.getLoadedT9()
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	t9service := service.New(t9, s.logger)
	proto.RegisterT9Server(grpcServer, t9service)
	s.grpcHandler = grpcServer

	return nil
}

func (s *server) getLoadedT9() (t9.T9, error) {
	dictionaryReader, err := getDictionaryReader(s.configuration.DictionaryFile)
	if err != nil {
		return nil, err
	}
	defer dictionaryReader.Close()

	t9, err := t9.NewCachingT9(s.configuration.CacheSize)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(dictionaryReader)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if len(word) == 0 {
			continue
		}

		s.logger.Debugf(`found dictionary word: '%s'`, word)

		if err := t9.InsertWord(word); err != nil {
			s.logger.Warn(`could not insert dictionary word: '%s'`, word)
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return t9, nil
}

func getDictionaryReader(dictionaryFile string) (io.ReadCloser, error) {
	if dictionaryFile == "" {
		reader := bytes.NewReader(assets.Dictionary)
		return io.NopCloser(reader), nil
	}

	file, err := os.Open(dictionaryFile)
	if err != nil {
		return nil, fmt.Errorf(`could not open dictionary file for reading: %w`, err)
	}

	return file, nil
}
