//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/exiffM/final-project/internal/config"
	rpcapi "github.com/exiffM/final-project/internal/grpc/pb"
	"github.com/exiffM/final-project/internal/grpc/server"
	"github.com/exiffM/final-project/internal/monitoring"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MonitorSuite struct {
	suite.Suite
	agent         *monitoring.Agent
	ctx           context.Context
	cancel        context.CancelFunc
	serv          *server.Server
	configuration *config.AgentConfig
	wg            sync.WaitGroup
}

func (m *MonitorSuite) SetupSuite() {
	file, err := os.Open("/etc/system.monitor/config.yml")
	m.Require().NoError(err, "Open config file error occurred!")

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(file)
	m.Require().NoError(err, "Read config file error occurred!")

	m.configuration = config.NewConfig()
	err = viper.Unmarshal(m.configuration)
	m.Require().NoError(err, "Unmarshal config file error occurred!")

	m.ctx, m.cancel = context.WithCancel(context.Background())
	m.agent = monitoring.NewAgent(*m.configuration)
	m.serv = server.NewServer(m.agent)
}

func (m *MonitorSuite) SetupTest() {
	fmt.Print("Test set up")
	m.wg.Add(2)

	go func() {
		defer m.wg.Done()
		err := m.agent.AccumulateStats(m.ctx)
		m.Require().NoError(err)
	}()

	go func() {
		defer m.wg.Done()
		err := m.serv.Start("localhost:50051")
		m.Require().NoError(err)
	}()
}

func (m *MonitorSuite) TestService() {
	host := os.Getenv("MONITOR_HOST")
	if host == "" {
		host = "localhost:50051"
	}
	var err error
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	m.Require().NoError(err, "create client connection error")
	ctx := context.Background()
	client := rpcapi.NewMonitorClient(conn)
	r := &rpcapi.Request{
		Timeout:         5,
		AverageInterval: 15,
	}
	monitorClient, err := client.SendStatistic(ctx, r)
	m.Require().NoError(err)

	stats, err := monitorClient.Recv()
	m.Require().NoError(err)
	m.Require().NotEmpty(stats, "Empty statistics")
	m.Require().Greater(stats.CpuLoad.Usr, float64(0), "Cpu usr stat is less than zero")
	m.Require().Greater(stats.CpuLoad.Sys, float64(0), "Cpu sys stat is less than zero")
	m.Require().Greater(stats.CpuLoad.Idle, float64(0), "Cpu idle stat is less than zero")
	m.Require().Greater(stats.SysLoad.One, float64(0), "System load 1 is less than zero")
	m.Require().Greater(stats.SysLoad.Five, float64(0), "System load 5 is less than zero")
	m.Require().Greater(stats.SysLoad.Quater, float64(0), "System load 15 is less than zero")
}

func (m *MonitorSuite) TearDownSuite() {
	m.serv.Shutdown()
	m.cancel()
	m.wg.Wait()
}

func TestMonitorService(t *testing.T) {
	suite.Run(t, new(MonitorSuite))
}
