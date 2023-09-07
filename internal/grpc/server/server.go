package server

import (
	"errors"
	"final-project/internal/grpc/convert"
	rpcapi "final-project/internal/grpc/pb"
	"final-project/internal/monitoring"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

var (
	errInvalidRequest = errors.New("invalid time intervals requested")
)

type Server struct {
	rpcapi.UnimplementedMonitorServer
	grpcServer *grpc.Server
	monitor    *monitoring.Agent
}

func NewServer(m *monitoring.Agent) *Server {
	return &Server{monitor: m}
}

func (s *Server) SendStatistic(r *rpcapi.Request, stream rpcapi.Monitor_SendStatisticServer) error {
	if r.Timeout < 1 || r.AverageInterval < 1 {
		return errInvalidRequest
	}
	for {
		select {
		case <-stream.Context().Done():
			// finish stream due to client disconnect
			return nil
		case <-time.After(time.Duration(r.AverageInterval) * time.Second):
			sendTicker := time.NewTicker(time.Duration(r.Timeout) * time.Second)
			for {
				select {
				case <-stream.Context().Done():
					return nil
				case <-sendTicker.C:
					stats := s.monitor.Average(int(r.AverageInterval))
					if err := stream.Send(convert.ConvertStatistic(stats)); err != nil {
						return err
					}
				}
			}
		}
	}
}

func (s *Server) Start(adress string) error {
	log.Print("Start GRPC server")
	lis, err := net.Listen("tcp", adress)
	if err != nil {
		return err
	}

	s.grpcServer = grpc.NewServer()
	rpcapi.RegisterMonitorServer(s.grpcServer, s)
	err = s.grpcServer.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown() {
	s.grpcServer.Stop()
}
