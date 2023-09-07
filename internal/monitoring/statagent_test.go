package monitoring

import (
	"context"
	"final-project/internal/config"
	"log"
	"os"
	"testing"
	"time"

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
	viper.ReadConfig(file)

	configuration := config.NewConfig()
	err = viper.Unmarshal(configuration)
	if err != nil {
		log.Fatalf("Can't convert config to struct %v", err.Error())
	}
	agent = NewAgent(*configuration)
}

func TestCreate(t *testing.T) {

	require.Equal(t, time.Duration(60)*time.Second, agent.StorageTime, "Storaging time in config file has been modyfied")
	require.True(t, agent.AllowAvgSysLoad, "Average system load stat is off")
	require.True(t, agent.AllowAvgCpuLoad, "Average cpu load stat is off")
	require.True(t, agent.AllowDiskLoad, "Disk load stat is off")
	require.True(t, agent.AllowDiskFsInfo, "Disk file system info is off")
	require.True(t, agent.AllowTTNet, "Top talkers net stat is off")
	require.True(t, agent.AllowNetStats, "Net stat is off")
}

func TestAvgSysLoad(t *testing.T) {
	out, err := agent.averageSystemLoad()
	require.Nil(t, err, "Error occured!")
	require.Contains(t, out, "System load average:")
}

func TestAvgCpuLoad(t *testing.T) {
	out, err := agent.averageCpuLoad()
	require.Nil(t, err, "Error occured!")
	require.Contains(t, out, "Cpu", "Invalid output: %q", out)
	require.Contains(t, out, "us", "Invalid output: %q", out)
	require.Contains(t, out, "sy", "Invalid output: %q", out)
	require.Contains(t, out, "id", "Invalid output: %q", out)
}

func TestLoadDiskInfo(t *testing.T) {
	out, err := agent.loadDiskInfo()
	require.Nil(t, err, "Error occured!")
	require.Contains(t, out, "sda", "Invalid output: %q", out)
}

func TestFSDInfo(t *testing.T) {
	out, err := agent.fileSystemDiskInfo()
	require.Errorf(t, err, "Error: %q", err.Error())
	require.Contains(t, out, "sda", "Invalid output: %q", out)
}

func TestNetStatistics(t *testing.T) {
	out, err := agent.netStatistics()
	require.Errorf(t, err, "Error: %q", err.Error())
	require.Contains(t, out, "tcp", "Invalid output: %q", out)
}

func TestAvgStats(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	agent.AccumulateStats(ctx)
	st := agent.Average(5)
	_ = st
	// TODO:
}
