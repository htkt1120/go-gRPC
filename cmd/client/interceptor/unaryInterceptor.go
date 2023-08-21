package interceptor

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

func UnaryClientInteceptor1(ctx context.Context, method string, req, res interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// preprocessing
	fmt.Println("[pre] my unary client interceptor 1", method, req)
	err := invoker(ctx, method, req, res, cc, opts...)
	// postprocessing
	fmt.Println("[post] my unary client interceptor 1", res)
	return err
}
