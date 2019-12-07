package server

import (
	"io"

	"github.com/hashicorp/yamux"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("server")

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run(rwc io.ReadWriteCloser) error {
	session, err := yamux.Server(rwc, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := session.Close(); err != nil {
			log.Errorf("failed to close session: %s", err)
		}
	}()

	for {
		stream, err := session.AcceptStream()
		if err != nil {
			return err
		}

		go func() {
			defer func() {
				if err := stream.Close(); err != nil {
					log.Errorf("failed to close stream: %s", err)
				}
			}()

			if err := s.handleStream(stream); err != nil {
				log.Errorf("failed to handle stream: %s", err)
			}
		}()
	}
}

func (s *Server) handleStream(stream *yamux.Stream) error {
	return nil
}
