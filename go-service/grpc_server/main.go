package main

import (
    "context"
    "log"
    "net"

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

    log.Println("gRPC server listening on :50051")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
