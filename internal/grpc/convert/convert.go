package convert

import (
	rpcapi "final-project/internal/grpc/pb"
	types "final-project/internal/statistics"
)

type ClientRequest struct {
	Timeout           int
	AveragingInterval int
}

func NewRpcDiskInfo() *rpcapi.DiskStats {
	return &rpcapi.DiskStats{Stats: make(map[string]*rpcapi.DiskInfo)}
}

func ConvertRequest(r *rpcapi.Request) ClientRequest {
	return ClientRequest{int(r.GetTimeout()), int(r.GetAverageInterval())}
}

func ConvertStatistic(s types.Statistic) *rpcapi.Statistic {
	var pbStat rpcapi.Statistic
	if s.ASLStat == nil {
		pbStat.SysLoad = &rpcapi.AvgSysLoad{}
	} else {
		pbStat.SysLoad = &rpcapi.AvgSysLoad{
			One:    s.ASLStat.OneMinLoad,
			Five:   s.ASLStat.FiveMinLoad,
			Quater: s.ASLStat.QuaterLoad,
		}
	}
	if s.ACLStat == nil {
		pbStat.CpuLoad = &rpcapi.AvgCpuLoad{}
	} else {
		pbStat.CpuLoad = &rpcapi.AvgCpuLoad{
			Usr:  s.ACLStat.Usr,
			Sys:  s.ACLStat.Sys,
			Idle: s.ACLStat.Idle,
		}
	}
	if s.DIStat == nil {
		pbStat.DiskInfo = &rpcapi.DiskStats{}
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
		pbStat.FsInfo = &rpcapi.FSDStats{}
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
	if s.TTStat == nil {
		pbStat.Tts = &rpcapi.TopTalkersStats{}
	} else {
		// TODO:
		_ = 5
	}
	if s.NStat == nil {
		pbStat.Net = &rpcapi.NetStats{}
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
