package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	hellopb "grpc-practice/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	scanner *bufio.Scanner
	client  hellopb.GreetingServiceClient
)

func main() {
	fmt.Println("start gRPC client")

	scanner = bufio.NewScanner(os.Stdin)

	address := "localhost:8080"
	conn, err := grpc.Dial(
		address,
		grpc.WithUnaryInterceptor(myUnaryClientInteceptor1),
		grpc.WithStreamInterceptor(myStreamClientInteceptor1),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal("connection error: ", err)
		return
	}
	defer conn.Close()

	client = hellopb.NewGreetingServiceClient(conn)

	for {
		fmt.Println("1: send request")
		fmt.Println("2: HelloServerStream")
		fmt.Println("3: HelloClientStream")
		fmt.Println("4: HelloBiStreams")
		fmt.Println("5: exit")
		fmt.Print("please enter > ")

		scanner.Scan()
		input := scanner.Text()

		switch input {
		case "1":
			Hello()
		case "2":
			HelloServerStream()
		case "3":
			HelloClientStream()
		case "4":
			HelloBiStreams()
		case "5":
			fmt.Println("exit")
			goto M
		}
	}
M:
}

func Hello() {
	var header, trailer metadata.MD
	fmt.Println("please enter your name")

	ctx := context.Background()
	md := metadata.New(map[string]string{
		"type": "unary",
		"from": "client",
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	scanner.Scan()
	name := scanner.Text()

	req := &hellopb.HelloRequest{
		Name: name,
	}
	res, err := client.Hello(ctx, req, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Println("header: ", header)
		fmt.Println("trailer: ", trailer)
		fmt.Println("response: ", res.Message)
	}
}

func HelloServerStream() {
	fmt.Println("please enter your name")
	scanner.Scan()
	name := scanner.Text()

	req := &hellopb.HelloRequest{
		Name: name,
	}
	stream, err := client.HelloServerStream(context.Background(), req)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	for {
		res, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("all the responses have been received")
			break
		}

		if err != nil {
			fmt.Println("error: ", err)
		}
		fmt.Println("response: ", res.Message)
	}
}

func HelloClientStream() {
	stream, err := client.HelloClientStream(context.Background())
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	sendCount := 5
	fmt.Printf("please enter %d names", sendCount)
	for i := 0; i < sendCount; i++ {
		scanner.Scan()
		name := scanner.Text()

		if err := stream.Send(&hellopb.HelloRequest{
			Name: name,
		}); err != nil {
			fmt.Println("error: ", err)
			return
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Println("response: ", res.GetMessage())
	}
}

func HelloBiStreams() {
	ctx := context.Background()
	md := metadata.New(map[string]string{
		"type": "bistream",
		"from": "client",
	})
	ctx = metadata.NewOutgoingContext(ctx, md)
	stream, err := client.HelloBiStreams(ctx)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	sendNum := 5
	fmt.Printf("please enter %d names", sendNum)

	var sendEnd, recvEnd bool
	sendCount := 0
	for !(sendEnd && recvEnd) {
		if !sendEnd {
			scanner.Scan()
			name := scanner.Text()

			sendCount++
			if err := stream.Send(&hellopb.HelloRequest{
				Name: name,
			}); err != nil {
				fmt.Println("error: ", err)
				sendEnd = true
			}

			if sendCount == sendNum {
				sendEnd = true
				if err := stream.CloseSend(); err != nil {
					fmt.Println("error: ", err)
				}
			}
		}

		var headerMD metadata.MD
		if !recvEnd {
			if headerMD == nil {
				headerMD, err = stream.Header()
				if err != nil {
					fmt.Println("error: ", err)
				}
				fmt.Println("header: ", headerMD)
			}

			if res, err := stream.Recv(); err != nil {
				if !errors.Is(err, io.EOF) {
					fmt.Println("error: ", err)
				}
				recvEnd = true
			} else {
				fmt.Println("response: ", res.Message)
			}
		}
	}

	trailerMD := stream.Trailer()
	fmt.Println("trailer: ", trailerMD)
}
