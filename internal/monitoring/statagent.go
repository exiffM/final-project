package monitoring

import (
	"context"
	"errors"
	"final-project/internal/config"
	"final-project/internal/executor"
	types "final-project/internal/statistics"
	"final-project/internal/storage"
	"strconv"
	"strings"
	"time"
)

type Agent struct {
	data            *storage.Storage // Storage of accumulating statistics
	StorageTime     time.Duration    // Time of clearing storage
	AllowAvgSysLoad bool             // On/off average system load
	AllowAvgCpuLoad bool             // On/off average CPU load
	AllowDiskLoad   bool             // On/off disks load
	AllowDiskFsInfo bool             // On/off disks file system info
	AllowTTNet      bool             // On/off top talkers net stats
	AllowNetStats   bool             // On/off net stats
}

func NewAgent(cfg config.AgentConfig) *Agent {
	return &Agent{
		storage.NewStorage(),
		time.Second * time.Duration(cfg.StorageTime),
		cfg.AvgSysLoad,
		cfg.AvgCpuLoad,
		cfg.DiskLoad,
		cfg.DiskFsInfo,
		cfg.TTNet,
		cfg.NetStats,
	}
}

var (
	errAccumulationFSD = errors.New("accumulation of file system disk info error")
	errAccumulationLis = errors.New("accumulation of tcp/udp listeners error")
)

func accumulateFSDInfo(fields []string) ([]types.FSD, error) {
	result := make([]types.FSD, 0)
	info := types.FSD{}
	for i := 0; i < len(fields); i += 5 {
		info.Source = fields[i]
		info.FS = fields[i+1]
		sz, err := strconv.ParseFloat(fields[i+2], 64)
		if err != nil {
			return nil, errAccumulationFSD
		}
		info.Total = sz
		usd, err := strconv.ParseFloat(fields[i+3], 64)
		if err != nil {
			return nil, errAccumulationFSD
		}
		info.Used = usd
		var pcnt float64 // if % is unknown output is "-"
		if fields[i+4] == "-" {
			pcnt = 0
		} else {
			pcnt, err = strconv.ParseFloat(strings.TrimSuffix(fields[i+4], "%"), 64)
			if err != nil {
				return nil, errAccumulationFSD
			}
		}
		info.Percent = pcnt
		result = append(result, info)
	}
	return result, nil
}

func accumulateListeners(fields []string) ([]types.Listeners, error) {
	result := make([]types.Listeners, 0)
	info := types.Listeners{}
	for i := 0; i < len(fields); i += 9 {
		info.ProgPid = fields[i+8]
		info.User = fields[i+6]
		info.Protoc = fields[i]
		adress := strings.Split(fields[i+3], ":")
		value, err := strconv.Atoi(adress[len(adress)-1])
		if err != nil {
			return nil, errAccumulationLis
		}
		info.Port = value
		result = append(result, info)
	}
	return result, nil
}

func (a *Agent) averageSystemLoad() (types.AvgSysLoadStat, error) {
	// Execute top command and grepping load average statistics
	out, err := executor.RunCmd("top -b -n1 | grep average")
	if err != nil {
		return types.AvgSysLoadStat{}, err
	}
	// Prepare result
	slice := strings.SplitAfter(out, "load average: ")
	values := strings.Split(strings.TrimSpace(slice[1]), ", ")
	min, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return types.AvgSysLoadStat{}, err
	}
	five, err := strconv.ParseFloat(values[1], 64)
	if err != nil {
		return types.AvgSysLoadStat{}, err
	}
	quater, err := strconv.ParseFloat(values[2], 64)
	if err != nil {
		return types.AvgSysLoadStat{}, err
	}
	return types.AvgSysLoadStat{OneMinLoad: min, FiveMinLoad: five, QuaterLoad: quater}, nil
}

func (a *Agent) averageCpuLoad() (types.AvgCpuLoadStat, error) {
	// Execute top command and grepping CPU average load statistics
	out, err := executor.RunCmd("top -b -n1 | grep Cpu")
	if err != nil {
		return types.AvgCpuLoadStat{}, err
	}
	// Prepare result
	splited := strings.Split(out, ",")
	firstSplited := strings.Split(splited[0], " ")
	var strUs string
	if firstSplited[1] != "" {
		strUs = firstSplited[1]
	} else {
		strUs = firstSplited[2]
	}
	us, err := strconv.ParseFloat(strUs, 64)
	if err != nil {
		return types.AvgCpuLoadStat{}, err
	}
	tmp := strings.Split(strings.TrimSpace(splited[1]), " ")
	sy, err := strconv.ParseFloat(tmp[0], 64)
	if err != nil {
		return types.AvgCpuLoadStat{}, err
	}
	tmp = strings.Split(strings.TrimSpace(splited[3]), " ")
	id, err := strconv.ParseFloat(tmp[0], 64)
	if err != nil {
		return types.AvgCpuLoadStat{}, err
	}
	return types.AvgCpuLoadStat{Usr: us, Sys: sy, Idle: id}, nil
}

func (a *Agent) loadDiskInfo() (types.DiskInfoStats, error) {
	// Execute top command getting disk load statistics
	out, err := executor.RunCmd("iostat -d -k | tail -n+4 | head -n-2")
	if err != nil {
		return nil, err
	}
	// Prepare result
	table := strings.Split(out, "\n")
	result := make(types.DiskInfoStats)
	for _, elem := range table {
		if elem == "" {
			continue
		}
		columns := strings.Fields(elem)
		tps, err := strconv.ParseFloat(columns[1], 64)
		if err != nil {
			return nil, err
		}
		rd, err := strconv.ParseFloat(columns[2], 64)
		if err != nil {
			return nil, err
		}
		wr, err := strconv.ParseFloat(columns[3], 64)
		if err != nil {
			return nil, err
		}
		result[columns[0]] = types.DiskInfo{
			Tps:            tps,
			KbReadPerSec:   rd,
			KbWritenPerSec: wr,
		}
	}
	return result, nil
}

func (a *Agent) fileSystemDiskInfo() (types.FSDInfoStat, error) {
	// Execute top command disk file system info statistics
	// Intresting behaviour:
	// df returns error /run/user/1000/doc Operation not permited
	// This is harmless bug, but RunCmd returns error
	// This bug appears only if xdg-document-portal.service is working
	// To fix this issue xdg-document-portal.service should be switched off
	// by systemctl --user stop xdg-document-portal.service
	result := types.FSDInfoStat{}
	accumulate := func(cmd string) ([]types.FSD, error) {
		out, err := executor.RunCmd(cmd)
		if err != nil {
			return nil, err
		}
		fields := strings.Fields(out)
		info, err := accumulateFSDInfo(fields)
		if err != nil {
			return nil, err
		}
		return info, nil
	}
	out, err := accumulate("df --output='source','fstype','size','used','pcent' | tail -n+2")
	if err != nil {
		return types.FSDInfoStat{}, err
	}
	result.FSDBlocks = append(result.FSDBlocks, out...)
	out, err = accumulate("df --output='source','fstype','itotal','iused','ipcent' | tail -n+2")
	if err != nil {
		return types.FSDInfoStat{}, err
	}
	result.FSDInodes = append(result.FSDInodes, out...)

	return result, nil
}

// func (a *Agent) TopTalkersInfo() (string, error) {

// }

func (a *Agent) netStatistics() (types.NetStat, error) {
	result := types.NetStat{}
	out, err := executor.RunCmd("netstat -ltupeN --numeric-ports | grep LISTEN")
	if err != nil {
		return types.NetStat{}, err
	}
	fields := strings.Fields(out)
	listeners, err := accumulateListeners(fields)
	if err != nil {
		return types.NetStat{}, err
	}
	result.TUListeners = append(result.TUListeners, listeners...)
	out, err = executor.RunCmd("ss -ta | awk '{print $1}' | sort | uniq -c | sort -nr")
	if err != nil {
		return types.NetStat{}, err
	}
	result.TCPStatesCount = make(map[string]float64)
	splited := strings.Split(out, "\n")
	splited = splited[:len(splited)-1]
	for _, elem := range splited {
		if !strings.Contains(elem, "State") {
			pair := strings.Split(strings.TrimSpace(elem), " ")
			value, err := strconv.ParseFloat(pair[0], 64)
			if err != nil {
				return types.NetStat{}, err
			}
			result.TCPStatesCount[pair[1]] = value
		}
	}
	return result, nil
}

func (a *Agent) Statistics() (types.Statistic, error) {
	result := types.Statistic{}
	if a.AllowAvgSysLoad {
		avgsys, err := a.averageSystemLoad()
		if err != nil {
			return types.Statistic{}, err
		}
		result.ASLStat = &avgsys
	}
	if a.AllowAvgCpuLoad {
		avgcpu, err := a.averageCpuLoad()
		if err != nil {
			return types.Statistic{}, err
		}
		result.ACLStat = &avgcpu
	}
	if a.AllowDiskLoad {
		dl, err := a.loadDiskInfo()
		if err != nil {
			return types.Statistic{}, err
		}
		result.DIStat = dl
	}
	if a.AllowDiskFsInfo {
		fsdi, err := a.fileSystemDiskInfo()
		if err != nil {
			return types.Statistic{}, err
		}
		result.FSDIStat = &fsdi
	}
	if a.AllowTTNet {
		// TODO:
		_ = 5
	}
	if a.AllowNetStats {
		net, err := a.netStatistics()
		if err != nil {
			return types.Statistic{}, err
		}
		result.NStat = &net
	}
	return result, nil
}

func (a *Agent) AccumulateStats(ctx context.Context) error {
	clearTicker := time.NewTicker(a.StorageTime)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-clearTicker.C:
			a.data.Clear()
		case <-ticker.C:
			s, err := a.Statistics()
			if err != nil {
				return err
			}
			a.data.Append(s)
		}
	}
}

func (a *Agent) Average(m int) types.Statistic {
	stats := a.data.PullOut(m)
	avgStat := types.Statistic{}
	if a.AllowAvgSysLoad {
		avgStat.ASLStat = types.NewASLS()
	}
	if a.AllowAvgCpuLoad {
		avgStat.ACLStat = types.NewACLS()
	}
	if a.AllowDiskLoad {
		avgStat.DIStat = make(types.DiskInfoStats)
	}
	if a.AllowDiskFsInfo {
		avgStat.FSDIStat = types.NewFSDIS(len(stats[0].FSDIStat.FSDBlocks))
	}
	if a.AllowNetStats {
		avgStat.NStat = types.NewNetStat(len(stats[0].NStat.TUListeners))
	}
	// Sum up stats
	for _, elem := range stats {
		if a.AllowAvgSysLoad {
			avgStat.ASLStat.Sum(*elem.ASLStat)
		}
		if a.AllowAvgCpuLoad {
			avgStat.ACLStat.Sum(*elem.ACLStat)
		}
		if a.AllowDiskLoad {
			avgStat.DIStat.Sum(elem.DIStat)
		}
		if a.AllowDiskFsInfo {
			avgStat.FSDIStat.Sum(*elem.FSDIStat)
		}
		if a.AllowNetStats {
			avgStat.NStat.Sum(*elem.NStat)
		}
	}
	// Averaging and ceiling stats
	if a.AllowAvgSysLoad {
		avgStat.ASLStat.Avg(float64(len(stats)))
		avgStat.ASLStat.Ceil()
	}
	if a.AllowAvgCpuLoad {
		avgStat.ACLStat.Avg(float64(len(stats)))
		avgStat.ACLStat.Ceil()
	}
	if a.AllowDiskLoad {
		avgStat.DIStat.Avg(float64(len(stats)))
		avgStat.DIStat.Ceil()
	}
	if a.AllowDiskFsInfo {
		avgStat.FSDIStat.Avg(float64(len(stats)))
		avgStat.FSDIStat.Ceil()
	}
	if a.AllowNetStats {
		avgStat.NStat.Avg(float64(len(stats)))
		avgStat.NStat.Ceil()
	}

	return avgStat
}
