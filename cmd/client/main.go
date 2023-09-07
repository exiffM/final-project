package main

import (
	"context"
	rpcapi "final-project/internal/grpc/pb"
	"fmt"
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
			printStatistic(stats)
		}
	}
}

func printStatistic(s *rpcapi.Statistic) {
	// Average system load
	fmt.Printf("Average system load: One minute: %v, Five minutes: %v, Fiveteen minutes: %v;\n",
		s.SysLoad.One, s.SysLoad.Five, s.SysLoad.Quater)
	// Average cpu load
	fmt.Printf("Average CPU load: %vus, %vsy, %vid;\n", s.CpuLoad.Usr, s.CpuLoad.Sys, s.CpuLoad.Idle)
	// Average disk load
	fmt.Printf("Device\t\ttps\t\tKb_read/s\t\tKb_wrtn/s\n")
	for key, elem := range s.DiskInfo.Stats {
		fmt.Printf("%v\t\t%v\t\t%v\t\t%v\n", key, elem.Tps, elem.Kbrps, elem.Kbwps)
	}
	// Disks file system info
	fmt.Println("Disks file system info by blocks")
	fmt.Printf("Source\t\tFile system\t\tSize\t\tUsed\t\tPercent")
	for _, elem := range s.FsInfo.Fsdblocks {
		fmt.Printf("%v\t\t%v\t\t%v\t\t%v\t\t%v\n",
			elem.Source, elem.Fs, elem.Total, elem.Used, elem.Percent)
	}
	fmt.Println("Disks file system info by inodes")
	fmt.Println("Source\t\tFile system\t\tTotal\t\tiUsed\t\tiPercent")
	for _, elem := range s.FsInfo.Fsdinodes {
		fmt.Printf("%v\t\t%v\t\t%v\t\t%v\t\t%v\n",
			elem.Source, elem.Fs, elem.Total, elem.Used, elem.Percent)
	}
	// Net stats: tcp/udp listeners
	fmt.Println("Net statistic tcp/udp listeners")
	fmt.Println("Prog/PID\t\tUser\t\tProtocol\t\tPort")
	for _, elem := range s.Net.TuListeners {
		fmt.Printf("%v\t\t%v\t\t%v\t\t%v\n",
			elem.Pid, elem.User, elem.Protoc, elem.Port)
	}
	// Net stats: tcp sockets state
	fmt.Println("Averaged TCP sockets states:")
	for key := range s.Net.States {
		fmt.Printf("%v - %v\n", key, s.Net.States[key])
	}
}
