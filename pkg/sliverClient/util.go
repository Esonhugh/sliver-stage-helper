package sliverClient

import (
	"context"

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

func (c *Client) GetImplantProfileByName(name string) *clientpb.ImplantProfile {
	pbProfiles, err := c.ImplantProfiles(context.Background(), &commonpb.Empty{})
	if err != nil {
		c.log.Errorf("Error getting implant profiles: %v", err)
		return nil
	}
	for _, profile := range pbProfiles.Profiles {
		if profile.Name == name {
			return profile
		}
	}
	return nil
}
