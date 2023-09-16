package convert

import (
	"fmt"

	rpcapi "github.com/exiffM/final-project/internal/grpc/pb"
	types "github.com/exiffM/final-project/internal/statistics"
)

type ClientRequest struct {
	Timeout           int
	AveragingInterval int
}

func NewRPCDiskInfo() *rpcapi.DiskStats {
	return &rpcapi.DiskStats{Stats: make(map[string]*rpcapi.DiskInfo)}
}

func Statistic(s types.Statistic) *rpcapi.Statistic {
	var pbStat rpcapi.Statistic
	if s.ASLStat == nil {
		pbStat.SysLoad = nil
	} else {
		pbStat.SysLoad = &rpcapi.AvgSysLoad{
			One:    s.ASLStat.OneMinLoad,
			Five:   s.ASLStat.FiveMinLoad,
			Quater: s.ASLStat.QuaterLoad,
		}
	}
	if s.ACLStat == nil {
		pbStat.CpuLoad = nil
	} else {
		pbStat.CpuLoad = &rpcapi.AvgCpuLoad{
			Usr:  s.ACLStat.Usr,
			Sys:  s.ACLStat.Sys,
			Idle: s.ACLStat.Idle,
		}
	}
	if s.DIStat == nil {
		pbStat.DiskInfo = nil
	} else {
		pbStat.DiskInfo = &rpcapi.DiskStats{Stats: make(map[string]*rpcapi.DiskInfo)}
		for key := range s.DIStat {
			pbDi := &rpcapi.DiskInfo{}
			pbDi.Tps = s.DIStat[key].Tps
			pbDi.Kbrps = s.DIStat[key].KbReadPerSec
			pbDi.Kbwps = s.DIStat[key].KbWritenPerSec
			pbStat.DiskInfo.Stats[key] = pbDi
		}
	}
	if s.FSDIStat == nil {
		pbStat.FsInfo = nil
	} else {
		pbStat.FsInfo = &rpcapi.FSDStats{
			Fsdblocks: make([]*rpcapi.FSD, 0),
			Fsdinodes: make([]*rpcapi.FSD, 0),
		}
		for idx := range s.FSDIStat.FSDBlocks {
			block := rpcapi.FSD{}
			block.Source = s.FSDIStat.FSDBlocks[idx].Source
			block.Fs = s.FSDIStat.FSDBlocks[idx].FS
			block.Total = s.FSDIStat.FSDBlocks[idx].Total
			block.Used = s.FSDIStat.FSDBlocks[idx].Used
			block.Percent = s.FSDIStat.FSDBlocks[idx].Percent
			pbStat.FsInfo.Fsdblocks = append(pbStat.FsInfo.Fsdblocks, &block)
		}
		for idx := range s.FSDIStat.FSDInodes {
			block := rpcapi.FSD{}
			block.Source = s.FSDIStat.FSDInodes[idx].Source
			block.Fs = s.FSDIStat.FSDInodes[idx].FS
			block.Total = s.FSDIStat.FSDInodes[idx].Total
			block.Used = s.FSDIStat.FSDInodes[idx].Used
			block.Percent = s.FSDIStat.FSDInodes[idx].Percent
			pbStat.FsInfo.Fsdinodes = append(pbStat.FsInfo.Fsdinodes, &block)
		}
	}
	if s.NStat == nil {
		pbStat.Net = nil
	} else {
		pbStat.Net = &rpcapi.NetStats{
			TuListeners: make([]*rpcapi.Listeners, 0),
			States:      make(map[string]float64),
		}
		for idx := range s.NStat.TUListeners {
			lis := rpcapi.Listeners{}
			lis.Pid = s.NStat.TUListeners[idx].ProgPid
			lis.User = s.NStat.TUListeners[idx].User
			lis.Protoc = s.NStat.TUListeners[idx].Protoc
			lis.Port = int32(s.NStat.TUListeners[idx].Port)
			pbStat.Net.TuListeners = append(pbStat.Net.TuListeners, &lis)
		}
		pbStat.Net.States = s.NStat.TCPStatesCount
	}
	return &pbStat
}

func PrintStatistic(s *rpcapi.Statistic) {
	// Average system load
	if s.SysLoad != nil {
		fmt.Printf("Average system load: One minute: %v, Five minutes: %v, Fiveteen minutes: %v;\n",
			s.SysLoad.One, s.SysLoad.Five, s.SysLoad.Quater)
	}
	// Average cpu load
	if s.CpuLoad != nil {
		fmt.Printf("Average CPU load: %vus, %vsy, %vid;\n", s.CpuLoad.Usr, s.CpuLoad.Sys, s.CpuLoad.Idle)
	}
	// Average disk load
	if s.DiskInfo != nil {
		fmt.Println("Disks load information:")
		fmt.Printf("Device\t\ttps\t\tKb_read/s\t\tKb_wrtn/s\n")
		for key, elem := range s.DiskInfo.Stats {
			fmt.Printf("%v\t\t%v\t\t%v\t\t%v\n", key, elem.Tps, elem.Kbrps, elem.Kbwps)
		}
	}
	// Disks file system info
	if s.FsInfo != nil {
		fmt.Println("Disks file system info by blocks:")
		fmt.Println("Source\t\tFile system\t\tSize\t\tUsed\t\tPercent")
		for _, elem := range s.FsInfo.Fsdblocks {
			fmt.Printf("%v\t\t%v\t\t%v\t\t%v\t\t%v\n",
				elem.Source, elem.Fs, elem.Total, elem.Used, elem.Percent)
		}
		fmt.Println("Disks file system info by inodes:")
		fmt.Println("Source\t\tFile system\t\tTotal\t\tiUsed\t\tiPercent")
		for _, elem := range s.FsInfo.Fsdinodes {
			fmt.Printf("%v\t\t%v\t\t%v\t\t%v\t\t%v\n",
				elem.Source, elem.Fs, elem.Total, elem.Used, elem.Percent)
		}
	}
	// Net stats: tcp/udp listeners
	if s.Net != nil {
		fmt.Println("Net statistic tcp/udp listeners:")
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
}
