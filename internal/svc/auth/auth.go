package auth

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	"endobit.io/metal"
	pb "endobit.io/metal/gen/go/proto/auth/v1"
	"endobit.io/metal/internal/data"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user exists")
)

// Service implements the auth grpc service.
type Service struct {
	pb.UnimplementedAuthServiceServer
	logger      *slog.Logger
	store       *data.Store
	users       *userStore
	tokens      *tokenManager
	interceptor *interceptor
}

// WithStore is an option setting function for NewService. It connects the data
// store to the service.
func WithStore(store *data.Store) func(*Service) {
	return func(s *Service) {
		s.store = store
	}
}

// WithLogger is an option setting function for NewService. It sets the logger
// to l.
func WithLogger(l *slog.Logger) func(*Service) {
	return func(s *Service) {
		s.logger = l
		s.interceptor.logger = l
	}
}

// WithUser is an option setting function for NewService. It adds user
// credentials to the service.
func WithUser(username, password string) func(*Service) {
	return func(s *Service) {
		user, err := newUser(username, password, true)
		if err != nil {
			s.logger.Error("cannot create user", "err", err)
			return
		}

		_ = s.users.Save(*user)
	}
}

// WithTTL is an option setting function for NewService. It sets the token TTL
// to t. The default token time to live is 5 minutes.
func WithTTL(t time.Duration) func(*Service) {
	return func(s *Service) {
		s.tokens.tokenTTL = t
	}
}

// NewService returns a new auth service.
func NewService(opts ...func(*Service)) *Service {
	tokens := newTokenManager(5 * time.Minute)

	svc := Service{
		logger: slog.Default(),
		users:  newUserStore(),
		tokens: tokens,
		interceptor: &interceptor{
			logger: slog.Default(),
			tokens: tokens,
		},
	}

	for _, opt := range opts {
		opt(&svc)
	}

	return &svc
}

// Logger returns the logger.
func (s *Service) Logger() *slog.Logger {
	return s.logger
}

// UnaryInterceptor returns a grpc.UnaryServerInterceptor that checks for a
// valid token and authorization.
func (s *Service) UnaryInterceptor(skip ...string) grpc.UnaryServerInterceptor {
	return s.interceptor.UnaryInterceptor(skip)
}

// StreamInterceptor returns a grpc.StreamServerInterceptor that checks for a
// valid token and authorization.
func (s *Service) StreamInterceptor(skip ...string) grpc.StreamServerInterceptor {
	return s.interceptor.StreamInterceptor(skip)
}

func (s *Service) Ping(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *Service) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.store.ReadUser(ctx, req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid username or password")
	}

	roles, err := s.store.ReadRolesForUser(ctx, req.GetUsername())
	if err != nil {
		return nil, err
	}

	var admin bool
	// Check if the user has the admin role.
	for _, role := range roles {
		if role.Name == "admin" {
			admin = true
		}
	}

	if user.PasswordHash == nil || *user.PasswordHash != req.GetPassword() { // TODO: hash the password
		return nil, status.Errorf(codes.Unauthenticated, "invalid password")
	}

	token, err := s.tokens.Generate(user.Name, admin)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	res := pb.LoginResponse_builder{
		Token: proto.String(token),
	}.Build()

	return res, nil
}

type interceptor struct {
	logger *slog.Logger
	tokens *tokenManager
}

func (i *interceptor) UnaryInterceptor(skip []string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if err := i.authorize(ctx, info.FullMethod, skip); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (i *interceptor) StreamInterceptor(skip []string) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := i.authorize(ss.Context(), info.FullMethod, skip); err != nil {
			return err
		}

		return handler(srv, ss)
	}
}

func (i *interceptor) authorize(ctx context.Context, method string, skip []string) error {
	for _, s := range skip { // things like login skip the token check
		if s == method {
			return nil
		}
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata not found")
	}

	vals := md.Get(metal.AuthorizationMetaData)
	if len(vals) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization metadata not found")
	}

	claims, err := i.tokens.Verify(vals[0])
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "invalid token")
	}

	// TODO: Check if method requires admin, need a dictionary of method names.
	// At this point we have a valid token so for now just treat everyone as an
	// admin.

	admin, ok := claims[tokenAdmin]
	if !ok {
		return nil
	}

	_, ok = admin.(bool)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "invalid token (invalid admin claim)")
	}

	return nil
}
