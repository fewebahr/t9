package server

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/RobertGrantEllis/t9/assets"
	"github.com/RobertGrantEllis/t9/proto"
	"github.com/RobertGrantEllis/t9/service"
	"github.com/RobertGrantEllis/t9/t9"
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
		return nil, errors.WithStack(err)
	}

	if closer, ok := dictionaryReader.(io.Closer); ok {
		defer closer.Close() // clean up if necessary
	}

	return t9, nil
}

func getDictionaryReader(dictionaryFile string) (io.Reader, error) {

	if dictionaryFile == "" {
		return bytes.NewReader(assets.Dictionary), nil
	}

	file, err := os.Open(dictionaryFile)
	if err != nil {
		return nil, errors.Wrap(err, `could not open dictionary file for reading`)
	}

	return file, nil
}
