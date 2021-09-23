package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"my_workspace/calculator/calculatorpb"
	"time"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")
	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	// doUnary(c)
	doClientStream(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.SumRequest{
		FirstNum:  3,
		SecondNum: 10,
	}

	res, err := c.Sum(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Calculator: %v", res.Sum)
}

func doServerStream(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")
	req := &calculatorpb.PrimeNumberDecomposeRequest{
		Number: 120,
	}

	resStream, err := c.PrimeNumberDecompose(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecompose RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from PrimeNumberDecompose: %v", msg.GetResult())
	}
}

func doClientStream(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")
	request := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Number: 1,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 2,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 3,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 4,
		},
	}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling ComputeAverage: %v", err)
	}

	for _, req := range request {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet: %v", err)
	}
	fmt.Printf("ComputeAverage Response: %v\n", res)
}
