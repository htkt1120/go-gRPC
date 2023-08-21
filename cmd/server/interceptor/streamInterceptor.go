package interceptor

import (
	"errors"
	"io"
	"log"

	"google.golang.org/grpc"
)

func StreamServerInterceptor1(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// preprocessing
	log.Println("[pre stream] my stream server interceptor 1: ", info.FullMethod)

	err := handler(srv, &serverStreamWrapper1{ss})

	// postprocessing
	log.Println("[post stream] my stream server interceptor 1: ")
	return err
}

type serverStreamWrapper1 struct {
	grpc.ServerStream
}

// preprocessing handler
func (s *serverStreamWrapper1) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	if !errors.Is(err, io.EOF) {
		log.Println("[pre message] my stream server interceptor 1: ", m)
	}
	return err
}

// preprocessing SendMsg
func (s *serverStreamWrapper1) SendMsg(m interface{}) error {
	log.Println("[post message] my stream server interceptor 1: ", m)
	return s.ServerStream.SendMsg(m)
}
