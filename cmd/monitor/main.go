package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/exiffM/final-project/internal/config"
	"github.com/exiffM/final-project/internal/grpc/server"
	"github.com/exiffM/final-project/internal/monitoring"
	"github.com/spf13/viper"
)

var (
	configFilePath string
	port           int
)

func init() {
	flag.StringVar(&configFilePath, "config", "/etc/system.monitor/config.yml", "Path to configuration file")
	flag.IntVar(&port, "port", 50051, "Port of rpc server")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		fmt.Println("System monitor v0.0.1")
		return
	}

	file, err := os.Open(configFilePath)
	if err != nil {
		fmt.Println(err.Error() + "No such file o directory")
		return
	}
	defer file.Close()

	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(file); err != nil {
		log.Fatal("Reading config error!") //nolint: gocritic
	}

	configuration := config.NewConfig()
	err = viper.Unmarshal(configuration)
	if err != nil {
		log.Fatalf("Can't convert config to struct %v", err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	ctx, cancel := context.WithCancel(context.Background())
	_ = ctx
	defer cancel()

	agent := monitoring.NewAgent(*configuration)
	go func() {
		defer wg.Done()
		if err := agent.AccumulateStats(ctx); err != nil {
			log.Fatal("Accumulation finished with error!")
		}
	}()

	serv := server.NewServer(agent)

	contained := os.Getenv("IS_IN_CONTAINER")
	var host string
	if contained == "1" {
		host = "0.0.0.0"
	} else {
		host = "localhost"
	}

	go func() {
		defer wg.Done()
		if err := serv.Start(net.JoinHostPort(host, strconv.Itoa(port))); err != nil {
			log.Fatal("Grpc server didn't start cause of error!")
		}
	}()

	log.Println("Daemon started!")

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, syscall.SIGINT, syscall.SIGTERM)
	<-exitCh

	serv.Shutdown()
	cancel()
	wg.Wait()
	log.Println("Daemon finished!")
}
