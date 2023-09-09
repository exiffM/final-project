package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/exiffM/final-project/internal/grpc/convert"
	rpcapi "github.com/exiffM/final-project/internal/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var port int

func init() {
	flag.IntVar(&port, "port", 50051, "Port of rpc server")
}

func main() {
	flag.Parse()

	host := os.Getenv("MONITOR_HOST")
	if host == "" {
		host = "localhost"
	}
	conn, err := grpc.Dial(net.JoinHostPort(host, strconv.Itoa(port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := rpcapi.NewMonitorClient(conn)

	r := &rpcapi.Request{Timeout: 5, AverageInterval: 15}
	monitorClient, err := client.SendStatistic(context.Background(), r)
	if err != nil {
		conn.Close()
		log.Fatal("Invalid request!")
	}
	for {
		for {
			stats, err := monitorClient.Recv()
			if err != nil {
				log.Printf("response error: %v\n", err)
				conn.Close()
				return
			}
			convert.PrintStatistic(stats)
		}
	}
}
