package main

import (
    "context"
    "log"
    "net"
    "os"
    "os/signal"
    "syscall"
    "time"

    "lab10-web-go-python/go-service/proto"

    "google.golang.org/grpc"
)

type server struct {
    proto.UnimplementedCalculatorServer
}

func (s *server) Add(ctx context.Context, req *proto.AddRequest) (*proto.AddResponse, error) {
    return &proto.AddResponse{Result: req.A + req.B}, nil
}

func (s *server) Multiply(ctx context.Context, req *proto.MultiplyRequest) (*proto.MultiplyResponse, error) {
    return &proto.MultiplyResponse{Result: req.A * req.B}, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    proto.RegisterCalculatorServer(grpcServer, &server{})

    go func() {
        log.Println("gRPC server listening on :50051")
        if err := grpcServer.Serve(lis); err != nil {
            log.Fatalf("failed to serve: %v", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down gRPC server...")

    done := make(chan struct{})
    go func() {
        grpcServer.GracefulStop()
        close(done)
    }()

    select {
    case <-done:
        log.Println("gRPC server gracefully stopped")
    case <-time.After(5 * time.Second):
        grpcServer.Stop()
        log.Println("gRPC server forced to stop")
    }
}
