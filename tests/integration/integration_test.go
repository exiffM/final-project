//go:build integration

package integration_test

import (
	"context"
	rpcapi "final-project/internal/grpc/pb"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MonitorSuite struct {
	suite.Suite
	conn   *grpc.ClientConn
	client rpcapi.MonitorClient
	ctx    context.Context
}

func (m *MonitorSuite) SetupTest() {
	// Create grpc conn for client
	fmt.Print("Test set up")
	host := os.Getenv("MONITOR_HOST")
	if host == "" {
		host = "localhost:50051"
	}
	var err error
	m.conn, err = grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	m.Require().NoError(err, "create client connection error")
	m.ctx = context.Background()
	m.client = rpcapi.NewMonitorClient(m.conn)
}

func (m *MonitorSuite) TestService() {
	r := &rpcapi.Request{
		Timeout:         5,
		AverageInterval: 15,
	}
	monitorClient, err := m.client.SendStatistic(m.ctx, r)
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

func (m *MonitorSuite) TearDownTest() {
	m.conn.Close()
}

func TestMonitorService(t *testing.T) {
	suite.Run(t, new(MonitorSuite))
}
