package monitoring

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/exiffM/final-project/internal/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

var agent *Agent

func init() {
	file, err := os.Open("/etc/system.monitor/config.yml")
	if err != nil {
		return
	}

	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(file); err != nil {
		file.Close()
		log.Fatal(err)
	}

	configuration := config.NewConfig()
	err = viper.Unmarshal(configuration)
	if err != nil {
		file.Close()
		log.Fatalf("Can't convert config to struct %v", err.Error())
	}
	file.Close()
	agent = NewAgent(*configuration)
}

func TestCreate(t *testing.T) {
	require.Equal(t, time.Duration(60)*time.Second, agent.StorageTime, "Storaging time in config file has been modyfied")
	// require.True(t, agent.AllowAvgSysLoad, "Average system load stat is off")
	// require.True(t, agent.AllowAvgCpuLoad, "Average cpu load stat is off")
	// require.True(t, agent.AllowDiskLoad, "Disk load stat is off")
	// require.True(t, agent.AllowDiskFsInfo, "Disk file system info is off")
	// require.False(t, agent.AllowTTNet, "Top talkers net stat is off")
	// require.True(t, agent.AllowNetStats, "Net stat is off")
}

func TestAvgSysLoad(t *testing.T) {
	out, err := agent.averageSystemLoad()
	require.Nil(t, err, "Error occurred!")
	require.NotEmpty(t, out, "Average system load not parsed")
}

func TestAvgCpuLoad(t *testing.T) {
	out, err := agent.averageCPULoad()
	require.Nil(t, err, "Error occurred!")
	require.NotEmpty(t, out, "Average cpu load not parsed")
}

func TestLoadDiskInfo(t *testing.T) {
	out, err := agent.loadDiskInfo()
	require.Nil(t, err, "Error occurred!")
	require.NotEmpty(t, out, "Disks information not parsed")
}

func TestFSDInfo(t *testing.T) {
	out, err := agent.fileSystemDiskInfo()
	require.Nil(t, err, "Error occurred!")
	require.NotEmpty(t, out, "Disks file system information not parsed")
}

func TestNetStatistics(t *testing.T) {
	out, err := agent.netStatistics()
	require.Nil(t, err, "Error occurred!")
	require.NotEmpty(t, out, "Net statistics not parsed")
}

func TestAvgStats(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := agent.AccumulateStats(ctx); err != nil {
		require.Fail(t, "Accumulation error occurred")
	}
	st := agent.Average(5)
	require.NotEmpty(t, st, "Net statistics not parsed")
}
