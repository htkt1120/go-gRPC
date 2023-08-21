// Package main implements a server for Greeter service.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"grpc/cmd/server/greeter"
	"grpc/cmd/server/interceptor"
	pb "grpc/pkg/grpc"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			interceptor.UnaryServerInterceptor1),
		grpc.StreamInterceptor(
			interceptor.StreamServerInterceptor1),
	)

	pb.RegisterGreeterServer(s, greeter.NewGreaterServer())
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
