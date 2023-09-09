package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/exiffM/final-project/internal/grpc/convert"
	rpcapi "github.com/exiffM/final-project/internal/grpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	host := os.Getenv("MONITOR_HOST")
	if host == "" {
		host = "localhost"
	}
	conn, err := grpc.Dial(net.JoinHostPort(host, "50051"),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
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
