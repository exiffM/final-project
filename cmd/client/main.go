package main

import (
	"context"
	"final-project/internal/grpc/convert"
	rpcapi "final-project/internal/grpc/pb"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := rpcapi.NewMonitorClient(conn)

	r := &rpcapi.Request{Timeout: 5, AverageInterval: 15}
	monitorClient, err := client.SendStatistic(context.Background(), r)
	if err != nil {
		log.Fatal("Invalid request!")
	}
	for {
		for {
			stats, err := monitorClient.Recv()
			if err != nil {
				log.Printf("response error: %v\n", err)
				return
			}
			convert.PrintStatistic(stats)
		}
	}
}
