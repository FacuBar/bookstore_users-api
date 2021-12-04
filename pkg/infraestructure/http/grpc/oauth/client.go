package oauth_grpc

import (
	"github.com/FacuBar/bookstore_users-api/pkg/infraestructure/http/grpc/oauth/oauthpb"
	"google.golang.org/grpc"
)

type Client struct {
	CC *grpc.ClientConn
	C  oauthpb.OauthServiceClient
}

func NewClient(address string) (*Client, error) {
	cc, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c := oauthpb.NewOauthServiceClient(cc)

	client := &Client{
		CC: cc,
		C:  c,
	}

	return client, nil
}
