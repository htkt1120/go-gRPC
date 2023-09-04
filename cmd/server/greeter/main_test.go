package greeter

import (
	"context"
	"errors"
	pb "grpc/pkg/grpc"
	"io"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, NewGreaterServer())
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
}

func bufDialer(ctx context.Context, address string) (net.Conn, error) {
	return lis.Dial()
}

func Test_greeter_server_UnaryHello(t *testing.T) {
	tests := []struct {
		name    string
		req     *pb.HelloRequest
		want    *pb.HelloReply
		wantErr bool
		err     error
	}{
		{
			name:    "correct",
			req:     &pb.HelloRequest{Name: "test"},
			want:    &pb.HelloReply{Message: "Hello, test!"},
			wantErr: false,
			err:     nil,
		},
		{
			name:    "name none error",
			req:     &pb.HelloRequest{Name: ""},
			want:    nil,
			wantErr: true,
			err:     status.Error(codes.InvalidArgument, "name is none"),
		},
		{
			name:    "name long error",
			req:     &pb.HelloRequest{Name: "AAAAAA"},
			want:    nil,
			wantErr: true,
			err:     status.Error(codes.InvalidArgument, "name is long"),
		},
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := client.UnaryHello(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("greeter_server.UnaryHello() error = %v, wantErr %v", err, tt.wantErr)
			}

			if (tt.wantErr) && (!errors.Is(err, tt.err)) {
				t.Errorf("greeter_server.UnaryHello() = %v, want %v", err, tt.err)
			}

			if (!tt.wantErr) && (res.GetMessage() != tt.want.Message) {
				t.Errorf("greeter_server.UnaryHello() = %v, want %v", res, tt.want)
			}
		})
	}
}

func Test_greeter_server_HelloServerStream(t *testing.T) {
	tests := []struct {
		name    string
		req     *pb.HelloRequest
		want    [5]*pb.HelloReply
		wantErr bool
	}{
		{
			name: "correct",
			req:  &pb.HelloRequest{Name: "test"},
			want: [5]*pb.HelloReply{
				{Message: "[0] Hello, test!"},
				{Message: "[1] Hello, test!"},
				{Message: "[2] Hello, test!"},
				{Message: "[3] Hello, test!"},
				{Message: "[4] Hello, test!"},
			},
			wantErr: false,
		},
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream, err := client.HelloServerStream(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("greeter_server.HelloServerStream() error = %v, wantErr %v", err, tt.wantErr)
			}
			count := 0
			for {
				res, err := stream.Recv()
				if errors.Is(err, io.EOF) {
					break
				}

				if (err != nil) != tt.wantErr {
					t.Errorf("greeter_server.UnaryHello() error = %v, wantErr %v", err, tt.wantErr)
				}

				if res.GetMessage() != tt.want[count].Message {
					t.Errorf("greeter_server.UnaryHello() = %v, want %v", res, tt.want)
				}
				count++
			}
			if count != 5 {
				t.Errorf("Expected 5 responses, got %d", count)
			}
		})
	}
}

func Test_greeter_server_HelloClientStream(t *testing.T) {
	tests := []struct {
		name    string
		req     [5]*pb.HelloRequest
		want    *pb.HelloReply
		wantErr bool
	}{
		{
			name: "correct",
			req: [5]*pb.HelloRequest{
				{Name: "test"},
				{Name: "test"},
				{Name: "test"},
				{Name: "test"},
				{Name: "test"},
			},
			want:    &pb.HelloReply{Message: "Hello, [test test test test test]!"},
			wantErr: false,
		},
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream, err := client.HelloClientStream(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("greeter_server.HelloClientStream() error = %v, wantErr %v", err, tt.wantErr)
			}

			sendCount := 5
			for i := 0; i < sendCount; i++ {
				if err := stream.Send(tt.req[i]); (err != nil) != tt.wantErr {
					t.Errorf("greeter_server.HelloClientStream() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			res, err := stream.CloseAndRecv()
			if (err != nil) != tt.wantErr {
				t.Errorf("greeter_server.HelloClientStream() error = %v, wantErr %v", err, tt.wantErr)
			}

			if res.GetMessage() != tt.want.Message {
				t.Errorf("greeter_server.HelloClientStream() = %v, want %v", res, tt.want)
			}
		})
	}
}

func Test_greeter_server_HelloBiStreams(t *testing.T) {
	tests := []struct {
		name    string
		req     []*pb.HelloRequest
		want    []*pb.HelloReply
		wantErr bool
	}{
		{
			name: "correct",
			req: []*pb.HelloRequest{
				{Name: "test"},
			},
			want: []*pb.HelloReply{
				{Message: "Hello, test!"},
			},
			wantErr: false,
		},
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream, err := client.HelloBiStreams(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("greeter_server.HelloBiStreams() error = %v, wantErr %v", err, tt.wantErr)
			}

			for i := 0; i < len(tt.req); i++ {
				stream.Send(tt.req[i])
				res, err := stream.Recv()
				if errors.Is(err, io.EOF) {
					break
				}

				if (err != nil) != tt.wantErr {
					t.Errorf("greeter_server.HelloBiStreams() error = %v, wantErr %v", err, tt.wantErr)
				}

				if res.GetMessage() != tt.want[i].Message {
					t.Errorf("greeter_server.HelloBiStreams() = %v, want %v", res, tt.want)
				}
			}

			if err := stream.CloseSend(); err != nil {
				t.Errorf("stream.CloseSend() faild")
			}
			if _, err := stream.Recv(); err != nil {
				if !errors.Is(err, io.EOF) {
					t.Errorf("stream.CloseSend() is not error")
				}
			}
		})
	}
}
