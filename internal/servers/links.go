package servers

import (
	"context"
	"main/internal/interfaces"
	"main/internal/models"
	pb "main/internal/proto"
)

// NewLinksServer creates a new LinksServer instance with the given service and host
func NewLinksServer(s interfaces.LinksService, host string) *LinksServer {
	return &LinksServer{
		host:         host,
		linksService: s,
	}
}

// LinksServer implements the gRPC server for link operations
type LinksServer struct {
	host string
	pb.UnimplementedLinksServer
	linksService interfaces.LinksService
}

// GetLink retrieves the original URL for a given short URL
func (h *LinksServer) GetLink(ctx context.Context, in *pb.RedirectLinkRequest) (*pb.RedirectLinkResponse, error) {
	var response pb.RedirectLinkResponse

	originLink, err := h.linksService.Get(ctx, in.ShortUrl)
	if err != nil {
		return nil, err
	}

	response.OriginalUrl = originLink
	return &response, nil

}

// AddLinks processes a batch of link creation requests
func (h *LinksServer) AddLinks(ctx context.Context, in *pb.AddLinksRequest) (*pb.AddLinksResponse, error) {
	var response pb.AddLinksResponse

	var originLinks []models.OriginLink
	for _, req := range in.Links {
		originLink := models.OriginLink{
			CorrelationID: req.CorrelationId,
			URL:           req.OriginalUrl,
		}
		originLinks = append(originLinks, originLink)
	}

	results, err := h.linksService.AddBatch(ctx, originLinks, h.host)
	if err != nil {
		return nil, err
	}

	for _, result := range results {
		response.Links = append(response.Links, &pb.LinkResponse{
			CorrelationId: result.CorrelationID,
			ShortUrl:      result.Result,
		})
	}

	return &response, nil
}

// AddLink handles a single link creation request
func (h *LinksServer) AddLink(ctx context.Context, in *pb.AddLinkRequest) (*pb.AddLinkResponse, error) {
	var response pb.AddLinkResponse
	var err error

	originLink := models.OriginLink{
		URL: in.Url,
	}

	response.Result, err = h.linksService.Add(ctx, originLink, h.host)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
