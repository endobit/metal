// Package main implements the stack service.
package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "modernc.org/sqlite" // sqlite driver

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	metald "endobit.io/metal"
	authpb "endobit.io/metal/gen/go/proto/auth/v1"
	metalpb "endobit.io/metal/gen/go/proto/metal/v1"
	"endobit.io/metal/internal/cert"
	"endobit.io/metal/internal/data"
	"endobit.io/metal/internal/svc/auth"
	"endobit.io/metal/internal/svc/metal"
	"endobit.io/metal/logging"
)

var version string

func main() {
	cmd := newRootCmd()
	cmd.Version = version

	if err := cmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func newRootCmd() *cobra.Command {
	var (
		dbDir         string
		useTailscale  bool
		port          int
		tokenTTL      time.Duration
		certDir       string
		adminPassword string
		logOpts       *logging.Options
	)
	cmd := cobra.Command{
		Use:   "metald",
		Short: "Metal Server",
		RunE: func(cmd *cobra.Command, _ []string) error {
			logger, err := logOpts.NewLogger()
			if err != nil {
				return err
			}

			logger.Info("Starting Metal Server", "version", cmd.Version, "port", port)

			dbPath := filepath.Join(dbDir, "metal.db")
			if err := setupDatabase(logger, dbPath); err != nil {
				return fmt.Errorf("failed to setup database: %w", err)
			}
			db, err := sql.Open("sqlite", dbPath)
			if err != nil {
				return fmt.Errorf("failed to open database: %w", err)
			}

			listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
			if err != nil {
				return fmt.Errorf("failed to listen: %w", err)
			}

			// Create a self-signed certificate if one does not exist.

			certPath := filepath.Join(certDir, "cert.pem")
			keyPath := filepath.Join(certDir, "key.pem")
			if err := setupCertificates(certPath, keyPath); err != nil {
				return err
			}

			cert, err := tls.LoadX509KeyPair(certPath, keyPath)
			if err != nil {
				return fmt.Errorf("failed to load certificate: %w", err)
			}

			ca := x509.NewCertPool()

			caBytes, err := os.ReadFile(certPath)
			if err != nil {
				return err
			}

			if !ca.AppendCertsFromPEM(caBytes) {
				return fmt.Errorf("failed to append CA certificate")
			}

			tlsConfig := tls.Config{
				Certificates: []tls.Certificate{cert},
				ClientAuth:   tls.NoClientCert,
				ClientCAs:    ca,
				MinVersion:   tls.VersionTLS12,
			}

			if adminPassword != "" {
				if err := setupAdminPassword(db, "admin", adminPassword); err != nil {
					return err
				}
			}

			store := data.NewStore(logger, db)

			authService := auth.NewService(
				auth.WithLogger(logger),
				auth.WithTTL(tokenTTL),
				auth.WithUser("admin", "admin"), // TODO: keep track of users in the database
				auth.WithStore(store),
			)

			metalService := metal.NewService(
				metal.WithLogger(logger),
				metal.WithStore(store),
			)

			server := grpc.NewServer(
				grpc.Creds(credentials.NewTLS(&tlsConfig)),
				grpc.ChainUnaryInterceptor(
					unaryLoggingInterceptor, // log the call
					authService.UnaryInterceptor(authpb.AuthService_Login_FullMethodName), // check for a valid token and authorization
				),
				grpc.ChainStreamInterceptor(
					streamLoggingInterceptor,        // log the call
					authService.StreamInterceptor(), // check for a valid token and authorization
				),
			)

			authpb.RegisterAuthServiceServer(server, authService)
			metalpb.RegisterMetalServiceServer(server, metalService)

			if err := server.Serve(listener); err != nil {
				return fmt.Errorf("failed to serve: %w", err)
			}

			return nil
		},
	}

	logOpts = logging.NewOptions(cmd.Flags())

	cmd.Flags().StringVar(&dbDir, "dbpath", ".", "Database directory")
	cmd.Flags().BoolVar(&useTailscale, "tailscale", false, "Get certificate from tailscale")
	cmd.Flags().IntVar(&port, "port", metald.DefaultPort, "port to listen on")
	cmd.Flags().DurationVar(&tokenTTL, "token-ttl", 5*time.Minute, "token time to live")
	cmd.Flags().StringVar(&certDir, "cert-path", ".", "path to cert.pem and key.pem files")
	cmd.Flags().StringVar(&adminPassword, "admin-password", "admin", "admin password")

	return &cmd
}

func setupCertificates(certPath, keyPath string) error {
	var missing int

	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		missing++
	}
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		missing++
	}
	if missing == 1 {
		return fmt.Errorf("missing cert or key file")
	}
	if missing == 0 {
		return nil
	}

	certFile, err := os.Create(certPath)
	if err != nil {
		return fmt.Errorf("failed to create cert file: %w", err)
	}
	defer certFile.Close()

	keyFile, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("failed to create key file: %w", err)
	}
	defer keyFile.Close()

	options := cert.NewOptions()
	return options.Create(certFile, keyFile)
}

func setupDatabase(logger *slog.Logger, dbPath string) error {
	store, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	defer store.Close()

	goose.SetLogger(logging.Legacy{Logger: logger})
	goose.SetBaseFS(data.Migrations)

	if err := goose.SetDialect("sqlite"); err != nil {
		return err
	}

	if err := goose.Up(store, "migrations"); err != nil {
		return err
	}

	return nil
}
