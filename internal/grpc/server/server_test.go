package server

import (
	"context"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/exiffM/final-project/internal/config"
	rpcapi "github.com/exiffM/final-project/internal/grpc/pb"
	"github.com/exiffM/final-project/internal/monitoring"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var agent *monitoring.Agent

func init() {
	file, err := os.Open("/etc/system.monitor/config.yml")
	if err != nil {
		return
	}
	defer file.Close()

	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(file); err != nil {
		log.Fatal(err)
	}

	configuration := config.NewConfig()
	err = viper.Unmarshal(configuration)
	if err != nil {
		log.Fatalf("Can't convert config to struct %v", err.Error())
	}
	agent = monitoring.NewAgent(*configuration)
}

func TestLogic(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		defer wg.Done()
		if err := agent.AccumulateStats(ctx); err != nil {
			log.Fatal("Accumulation finished with error!", err.Error())
		}
	}()

	serv := NewServer(agent)

	go func() {
		defer wg.Done()
		if err := serv.Start("localhost:50051"); err != nil {
			log.Fatal("Grpc server didn't start cause of error!")
		}
	}()

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err) //nolint: gocritic
	}
	defer conn.Close()

	client := rpcapi.NewMonitorClient(conn)

	r := &rpcapi.Request{Timeout: 5, AverageInterval: 15}
	monitorClient, err := client.SendStatistic(context.Background(), r)
	if err != nil {
		log.Fatal("Invalid request!")
	}
	var stats *rpcapi.Statistic
MAINFOR:
	for {
		for {
			stats, err = monitorClient.Recv()
			if err != nil {
				log.Printf("response error: %v\n", err)
				return
			}
			require.NotEmpty(t, stats, "Statistic hasn't been received!")
			break MAINFOR
		}
	}
	serv.Shutdown()
	cancel()
	wg.Wait()
}
