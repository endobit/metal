package main

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type logger interface {
	Logger() *slog.Logger
}

// unaryLoggingInterceptor logs the method and latency of a unary gRPC call.
func unaryLoggingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	var client string

	start := time.Now()

	svr, ok := info.Server.(logger)
	if !ok {
		return nil, errors.New("service does not implement logger interface")
	}
	logger := svr.Logger()

	if p, ok := peer.FromContext(ctx); ok {
		client = p.Addr.String()
	}

	resp, err := handler(ctx, req)
	if err != nil {
		logger.Error("Unary gRPC call", "err", err)
	}

	logger.Info("Unary gRPC call",
		"client", client,
		"method", info.FullMethod,
		"latency", time.Since(start))

	return resp, err
}

// streamLoggingInterceptor logs the method and latency of a streaming gRPC call.
func streamLoggingInterceptor(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	var client string

	start := time.Now()

	svr, ok := srv.(logger)
	if !ok {
		return errors.New("service does not implement logger interface")
	}
	logger := svr.Logger()

	err := handler(srv, stream)
	if err != nil {
		logger.Error("Stream gRPC call", "err", err)
	}

	logger.Info("Stream gRPC call",
		"client", client,
		"method", info.FullMethod,
		"latency", time.Since(start))

	return err
}
