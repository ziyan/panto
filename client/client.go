package client

import (
	"io"

	"github.com/hashicorp/yamux"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("client")

type Client struct {
	rwc io.ReadWriteCloser
}

func NewClient(rwc io.ReadWriteCloser) *Client {
	return &Client{
		rwc,
	}
}

func (c *Client) Run() error {
	session, err := yamux.Client(c.rwc, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := session.Close(); err != nil {
			log.Errorf("failed to close session: %s", err)
		}
	}()

	if _, err := session.OpenStream(); err != nil {
		return err
	}

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

			if err := c.handleStream(stream); err != nil {
				log.Errorf("failed to handle stream: %s", err)
			}
		}()
	}
}

func (c *Client) handleStream(stream *yamux.Stream) error {
	log.Debugf("stream initiated from server")
	return nil
}
