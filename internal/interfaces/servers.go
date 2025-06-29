package interfaces

import (
	"context"
	pb "main/internal/proto"
)

// LinksServer aggregates gRPC dealing with link manipulation (creation, retrieval).
type LinksServer interface {
	AddLink(ctx context.Context, in *pb.AddLinkRequest) (*pb.AddLinkResponse, error)
	AddLinks(ctx context.Context, in *pb.AddLinksRequest) (*pb.AddLinksResponse, error)
	GetLink(ctx context.Context, in *pb.RedirectLinkRequest) (*pb.RedirectLinkResponse, error)
}
