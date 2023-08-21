package greeter

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	pb "grpc/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func NewGreaterServer() *greeter_server {
	return &greeter_server{}
}

// server is used to implement helloworld.GreeterServer.
type greeter_server struct {
	pb.UnimplementedGreeterServer
}

// UnaryHello implements helloworld.GreeterServer
func (s *greeter_server) UnaryHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Println(md)
	}

	// metadata
	headerMD := metadata.New(map[string]string{"type": "unary", "from": "server", "in": "header"})
	if err := grpc.SetHeader(ctx, headerMD); err != nil {
		return nil, err
	}

	trailerMD := metadata.New(map[string]string{"type": "unary", "from": "server", "in": "trailer"})
	if err := grpc.SetTrailer(ctx, trailerMD); err != nil {
		return nil, err
	}

	return &pb.HelloReply{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

func (s *greeter_server) HelloServerStream(req *pb.HelloRequest, stream pb.Greeter_HelloServerStreamServer) error {
	resCount := 5
	for i := 0; i < resCount; i++ {
		if err := stream.Send(&pb.HelloReply{
			Message: fmt.Sprintf("[%d] Hello, %s!", i, req.GetName()),
		}); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}
	return nil
}

func (s *greeter_server) HelloClientStream(stream pb.Greeter_HelloClientStreamServer) error {
	nameList := make([]string, 0)
	for {
		req, err := stream.Recv()
		// END is err == io.EOF
		if errors.Is(err, io.EOF) {
			message := fmt.Sprintf("Hello, %v!", nameList)
			return stream.SendAndClose(&pb.HelloReply{
				Message: message,
			})
		}
		if err != nil {
			return err
		}
		nameList = append(nameList, req.GetName())
	}
}

func (s *greeter_server) HelloBiStreams(stream pb.Greeter_HelloBiStreamsServer) error {
	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		log.Println(md)
	}

	// metadata
	headerMD := metadata.New(map[string]string{"type": "stream", "from": "server", "in": "header"})
	if err := stream.SendHeader(headerMD); err != nil {
		return err
	}

	trailerMD := metadata.New(map[string]string{"type": "stream", "from": "server", "in": "trailer"})
	stream.SetTrailer(trailerMD)
	for {
		req, err := stream.Recv()
		// END is err == io.EOF
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		message := fmt.Sprintf("Hello, %v!", req.GetName())
		if err := stream.Send(&pb.HelloReply{
			Message: message,
		}); err != nil {
			return err
		}
	}
}
