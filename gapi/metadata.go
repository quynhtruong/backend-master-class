package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	xForwardedForHeader        = "x-forwarded-for"
	userAgentHeader            = "user-agent"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		userAgents := md.Get(grpcGatewayUserAgentHeader)
		if len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}
		userAgents = md.Get(userAgentHeader)
		if len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}
		clientIps := md.Get(xForwardedForHeader)
		if len(clientIps) > 0 {
			mtdt.ClientIP = clientIps[0]
		}
	}
	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}
	return mtdt
}
