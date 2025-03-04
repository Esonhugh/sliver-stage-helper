package sliverClient

import (
	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
)

func MakeRequest(session *clientpb.Session) *commonpb.Request {
	if session == nil {
		return nil
	}
	timeout := int64(60)
	return &commonpb.Request{
		SessionID: session.ID,
		Timeout:   timeout,
	}
}
