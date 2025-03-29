// Package stack implements the stack command line.
package metal

import (
	"log/slog"

	pb "endobit.io/metal/gen/go/proto/metal/v1"
	"endobit.io/metal/internal/data"
)

//go:generate go tool "github.com/dmarkham/enumer" -type AttrScope -linecomment -text

type AttrScope int

const (
	GlobalScope      AttrScope = iota // global
	ModelScope                        // model
	ZoneScope                         // zone
	ClusterScope                      // cluster
	RackScope                         // rack
	ApplianceScope                    // appliance
	EnvironmentScope                  // environment
	HostScope                         // host
	SwitchScope                       // switch
	BMCScope                          // bmc
)

// Service implements the stackd grpc service.
type Service struct {
	pb.UnimplementedMetalServiceServer
	logger *slog.Logger
	store  *data.Store
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
	}
}

// NewService returns a new stack service.
func NewService(opts ...func(*Service)) *Service {
	svc := Service{
		logger: slog.Default(),
	}

	for _, opt := range opts {
		opt(&svc)
	}

	return &svc
}

// Logger returns the logger.
func (s Service) Logger() *slog.Logger {
	return s.logger
}
