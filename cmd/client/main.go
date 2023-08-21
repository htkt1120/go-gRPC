package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"grpc/cmd/client/interceptor"
	pb "grpc/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	defaultName = "world"
)

var (
	scanner *bufio.Scanner
	client  pb.GreeterClient
	addr    = flag.String("addr", "localhost:50051", "the address to connect to")
	name    = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(
		*addr,
		grpc.WithUnaryInterceptor(interceptor.UnaryClientInteceptor1),
		grpc.WithStreamInterceptor(interceptor.StreamClientInteceptor1),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client = pb.NewGreeterClient(conn)
	scanner = bufio.NewScanner(os.Stdin)

	for {
		log.Println("1: send Request")
		log.Println("2: HelloServerStream")
		log.Println("3: HelloClientStream")
		log.Println("4: HelloBiStream")
		log.Println("5: exit")
		log.Print("please enter >")

		scanner.Scan()
		in := scanner.Text()

		switch in {
		case "1":
			Hello()
		case "2":
			HelloServerStream()
		case "3":
			HelloClientStream()
		case "4":
			HelloBiStream()
		case "5":
			log.Println("bye.")
			return
		}
	}
}

func Hello() {

	log.Println("Please enter your name.")
	scanner.Scan()
	name := scanner.Text()

	req := &pb.HelloRequest{
		Name: name,
	}

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// metadata
	md := metadata.New(map[string]string{"type": "unary", "from": "client"})
	ctx = metadata.NewOutgoingContext(ctx, md)

	var header, trailer metadata.MD
	res, err := client.UnaryHello(ctx, req, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		return
	}
	log.Println(header)
	log.Println(trailer)
	log.Printf("Greeting: %s", res.GetMessage())
	return
}

func HelloServerStream() {

	log.Println("Please enter your name.")
	scanner.Scan()
	name := scanner.Text()

	req := &pb.HelloRequest{
		Name: name,
	}

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// metadata
	md := metadata.New(map[string]string{"type": "stream", "from": "client"})
	ctx = metadata.NewOutgoingContext(ctx, md)

	stream, err := client.HelloServerStream(ctx, req)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		return
	}
	for {
		res, err := stream.Recv()
		// END is err == io.EOF
		if errors.Is(err, io.EOF) {
			log.Printf("all the responses have already received.")
			break
		}

		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Println(res)
	}
	return
}

func HelloClientStream() {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// metadata
	md := metadata.New(map[string]string{"type": "stream", "from": "client"})
	ctx = metadata.NewOutgoingContext(ctx, md)

	stream, err := client.HelloClientStream(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	sendCount := 5
	log.Printf("Please enter %d names.\n", sendCount)
	for i := 0; i < sendCount; i++ {
		scanner.Scan()
		name := scanner.Text()

		req := &pb.HelloRequest{
			Name: name,
		}

		if err := stream.Send(req); err != nil {
			log.Println(err)
			return
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Println(err)
	} else {
		log.Println(res.GetMessage())
	}
}

func HelloBiStream() {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// metadata
	md := metadata.New(map[string]string{"type": "stream", "from": "client"})
	ctx = metadata.NewOutgoingContext(ctx, md)

	stream, err := client.HelloBiStreams(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	sendNum := 5
	log.Printf("Please enter %d names.\n", sendNum)

	var sendEnd, recvEnd bool
	sendCount := 0
	for !(sendEnd && recvEnd) {
		// send
		if !sendEnd {
			scanner.Scan()
			name := scanner.Text()

			sendCount++
			req := &pb.HelloRequest{
				Name: name,
			}

			if err := stream.Send(req); err != nil {
				fmt.Println(err)
				sendEnd = true
			}

			if sendCount == sendNum {
				sendEnd = true
				if err := stream.CloseSend(); err != nil {
					fmt.Println(err)
				}
			}
		}

		// recv
		// recv metadata
		var headerMD metadata.MD
		if !recvEnd {
			if headerMD == nil {
				headerMD, err = stream.Header()
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(headerMD)
				}
			}
			if res, err := stream.Recv(); err != nil {
				if !errors.Is(err, io.EOF) {
					fmt.Println(err)
				}
				recvEnd = true
			} else {
				fmt.Println(res.GetMessage())
			}
		}
	}
	trailerMD := stream.Trailer()
	fmt.Println(trailerMD)
}
