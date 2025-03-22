package sliverClient

import (
	"context"
	"io"

	"github.com/Esonhugh/sliver-stage-helper/pkg/sliverClient/protobuf/clientpb"
	"github.com/Esonhugh/sliver-stage-helper/pkg/sliverClient/protobuf/commonpb"
	"github.com/Esonhugh/sliver-stage-helper/pkg/sliverClient/protobuf/rpcpb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type EventHandler func(c *Client, event *clientpb.Event) error

type Client struct {
	rpcpb.SliverRPCClient
	rpcpb.SliverRPC_EventsClient
	eventHandlers map[string]EventHandler
	grpcConn      *grpc.ClientConn

	log *log.Entry
}

func NewClient(config *ClientConfig) (*Client, error) {
	// connect to the server
	rpc, ln, err := MTLSConnect(config)
	if err != nil {
		return nil, err
	}
	log.Info("[*] Connected to sliver server")

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
		log:                    log.WithField("client", "sliver"),
		grpcConn:               ln,
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

func (c *Client) Close() {
	c.SliverRPC_EventsClient.CloseSend()
	c.grpcConn.Close()
}
