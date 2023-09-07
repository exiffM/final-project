package main

import (
	"context"
	"final-project/internal/config"
	"final-project/internal/grpc/server"
	"final-project/internal/monitoring"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spf13/viper"
)

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "/etc/system.monitor/config.yml", "Path to configuration file")
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

	viper.SetConfigType("yaml")
	viper.ReadConfig(file)

	configuration := config.NewConfig()
	err = viper.Unmarshal(configuration)
	if err != nil {
		log.Fatalf("Can't convert config to struct %v", err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	agent := monitoring.NewAgent(*configuration)
	go func() {
		defer wg.Done()
		if err := agent.AccumulateStats(ctx); err != nil {
			log.Fatal("Accumulation finished with error!")
		}
	}()

	serv := server.NewServer(agent)

	go func() {
		defer wg.Done()
		if err := serv.Start("localhost:50051"); err != nil {
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