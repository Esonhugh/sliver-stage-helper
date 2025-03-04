package sliverClient

import (
	"context"
	"io"

	"github.com/bishopfox/sliver/client/assets"
	"github.com/bishopfox/sliver/client/transport"
	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"github.com/bishopfox/sliver/protobuf/rpcpb"
	log "github.com/sirupsen/logrus"
)

type EventHandler func(c *Client, event *clientpb.Event) error

type Client struct {
	rpcpb.SliverRPCClient
	rpcpb.SliverRPC_EventsClient
	eventHandlers map[string]EventHandler
}

func ReadConfig(path string) (*assets.ClientConfig, error) {
	return assets.ReadConfig(path)
}

func NewClient(config *assets.ClientConfig) (*Client, error) {
	// connect to the server
	rpc, ln, err := transport.MTLSConnect(config)
	if err != nil {
		return nil, err
	}
	log.Info("[*] Connected to sliver server")
	defer ln.Close()

	// Open the event stream to be able to collect all events sent by  the server
	eventStream, err := rpc.Events(context.Background(), &commonpb.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	// infinite loop
	return &Client{
		SliverRPCClient:        rpc,
		SliverRPC_EventsClient: eventStream,
		eventHandlers:          make(map[string]EventHandler),
	}, nil
}

func (c *Client) startEventHandler() {
	for {
		event, err := c.Recv()
		if err == io.EOF || event == nil {
			return
		}
		// Trigger event based on type
		if handler, ok := c.eventHandlers[event.EventType]; ok {
			err := handler(c, event)
			if err != nil {
				log.Error(err)
			}
		} else {
			log.Tracef("No handler for event type: %s", event.EventType)
		}
	}
}

func (c *Client) RegisterEventHandler(eventType string, handler EventHandler) {
	c.eventHandlers[eventType] = handler
}
