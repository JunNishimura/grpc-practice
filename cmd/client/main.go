package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	hellopb "grpc-practice/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		fmt.Println("2: exit")
		fmt.Print("please enter > ")

		scanner.Scan()
		input := scanner.Text()

		switch input {
		case "1":
			Hello()
		case "2":
			fmt.Println("exit")
			goto M
		}
	}
M:
}

func Hello() {
	fmt.Println("please enter your name")
	scanner.Scan()
	name := scanner.Text()

	req := &hellopb.HelloRequest{
		Name: name,
	}
	res, err := client.Hello(context.Background(), req)
	if err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Println("response: ", res.Message)
	}
}
