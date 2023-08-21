package interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func UnaryServerInterceptor1(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// preprocessing
	log.Println("[pre] my unary server interceptor 1: ", info.FullMethod, req)
	res, err := handler(ctx, req)
	// postprocessing
	log.Println("[post] my unary server interceptor 1: ", res)
	return res, err
}
