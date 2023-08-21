package interceptor

import (
	"context"
	"errors"
	"io"
	"log"

	"google.golang.org/grpc"
)

func StreamClientInteceptor1(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	// preprocessing
	log.Println("[pre] my stream client interceptor 1", method)

	stream, err := streamer(ctx, desc, cc, method, opts...)
	// postprocessing
	return &clientStreamWrapper1{stream}, err
}

type clientStreamWrapper1 struct {
	grpc.ClientStream
}

// preprocessing SendMsg
func (s *clientStreamWrapper1) SendMsg(m interface{}) error {
	log.Println("[pre message] my stream client interceptor 1: ", m)
	return s.ClientStream.SendMsg(m)
}

// postprocessing RecvMsg
func (s *clientStreamWrapper1) RecvMsg(m interface{}) error {
	err := s.ClientStream.RecvMsg(m)
	if !errors.Is(err, io.EOF) {
		log.Println("[post message] my stream client interceptor 1: ", m)
	}
	return err
}

// postprocessing CloseSend
func (s *clientStreamWrapper1) CloseSend() error {
	err := s.ClientStream.CloseSend()
	log.Println("[post] my stream client interceptor 1")
	return err
}
