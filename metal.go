package metal

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	authpb "endobit.io/metal/gen/go/proto/auth/v1"
	metalpb "endobit.io/metal/gen/go/proto/metal/v1"
)

const (
	// DefaultPort is the default port to listen on. It can be overridden with the
	// --port flag.
	DefaultPort = 8080

	// AuthorizationMetaData is the key to be used by the client when storing
	// the authorization token in context metadata.
	AuthorizationMetaData = "authorization"
)

// Client is a gRPC client for the metal and auth services.
type Client struct {
	metalpb.MetalServiceClient
	Logger   *slog.Logger
	Auth     authpb.AuthServiceClient
	username string
	password string
	token    string
	conn     *grpc.ClientConn
}

func NewClient(conn *grpc.ClientConn, logger *slog.Logger) *Client {
	if logger == nil {
		logger = slog.Default()
	}

	return &Client{
		MetalServiceClient: metalpb.NewMetalServiceClient(conn),
		Logger:             logger,
		Auth:               authpb.NewAuthServiceClient(conn),
	}
}

// Close closes the gRPC connection if it is open. This should be called when
// the client is no longer needed to avoid resource leaks.
func (c *Client) Close() error {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}
		c.conn = nil
	}

	return nil
}

// Context returns a context with the authorization token if it exists. This
// context should be used for all metal gRPC requests.
func (c *Client) Context() context.Context {
	ctx := context.Background()
	if c.token != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, AuthorizationMetaData, c.token)
	}

	return ctx
}

// Authorize logs in to the metal service and stores the token for future requests.
func (c *Client) Authorize(username, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := authpb.LoginRequest_builder{
		Username: &username,
		Password: &password,
	}.Build()

	resp, err := c.Auth.Login(ctx, req)
	if err != nil {
		return err
	}

	c.username = username
	c.password = password
	c.token = resp.GetToken()

	return nil
}

// ReAuthorize logs in to the metal service using the stored credentials.
func (c *Client) ReAuthorize() error {
	return c.Authorize(c.username, c.password)
}

// RefreshTokenMiddleware pings the grpc service and if the credentials fail it
// re-authorizes the connection.
func (c *Client) RefreshTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := c.Auth.Ping(c.Context(), &emptypb.Empty{}); err != nil {
			if status.Code(err) == codes.Unauthenticated {
				if err := c.ReAuthorize(); err != nil {
					http.Error(w, "failed to re-authorize with metald", http.StatusUnauthorized)
					return
				}
			} else {
				http.Error(w, "failed to ping auth service", http.StatusInternalServerError)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
