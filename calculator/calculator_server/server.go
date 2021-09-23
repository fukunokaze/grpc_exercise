package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"my_workspace/calculator/calculatorpb"
	"net"
	"time"

	"google.golang.org/grpc"
)

type server struct{}

func main() {
	fmt.Println("Hello world")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	result := &calculatorpb.SumResponse{
		Sum: req.GetFirstNum() + req.GetSecondNum(),
	}

	return result, nil
}

func (*server) PrimeNumberDecompose(req *calculatorpb.PrimeNumberDecomposeRequest, stream calculatorpb.CalculatorService_PrimeNumberDecomposeServer) error {
	fmt.Printf("Prime Number Decomposition function was invoked with %v\n", req)
	var k int32 = 2
	num := req.GetNumber()

	for {
		if num <= 1 {
			break
		}
		if num%k == 0 {
			num = num / k
			res := &calculatorpb.PrimeNumberDecomposeResponse{
				Result: k,
			}
			stream.Send(res)
			time.Sleep(1000 * time.Millisecond)
		} else {
			k += 1
		}
	}

	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage function was invoked with a streaming request\n")
	var cnt float64 = 0
	var sum float64 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			if cnt != 0 {
				return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
					Result: sum / cnt,
				})
			} else {
				return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
					Result: 0,
				})
			}
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		cnt += 1
		sum += float64(req.GetNumber())
	}
}
