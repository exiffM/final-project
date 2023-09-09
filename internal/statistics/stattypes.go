package statistics

import "math"

type AvgSysLoadStat struct {
	OneMinLoad  float64
	FiveMinLoad float64
	QuaterLoad  float64
}

func NewASLS() *AvgSysLoadStat {
	return &AvgSysLoadStat{}
}

func (asls *AvgSysLoadStat) Sum(arg AvgSysLoadStat) {
	asls.OneMinLoad += arg.OneMinLoad
	asls.FiveMinLoad += arg.FiveMinLoad
	asls.QuaterLoad += arg.QuaterLoad
}

func (asls *AvgSysLoadStat) Avg(count float64) {
	asls.OneMinLoad /= count
	asls.FiveMinLoad /= count
	asls.QuaterLoad /= count
}

func (asls *AvgSysLoadStat) Ceil() {
	asls.OneMinLoad = math.Ceil(asls.OneMinLoad*100) / 100
	asls.FiveMinLoad = math.Ceil(asls.FiveMinLoad*100) / 100
	asls.QuaterLoad = math.Ceil(asls.QuaterLoad*100) / 100
}

type AvgCPULoadStat struct {
	Usr  float64
	Sys  float64
	Idle float64
}

func NewACLS() *AvgCPULoadStat {
	return &AvgCPULoadStat{}
}

func (acls *AvgCPULoadStat) Sum(arg AvgCPULoadStat) {
	acls.Usr += arg.Usr
	acls.Sys += arg.Sys
	acls.Idle += arg.Idle
}

func (acls *AvgCPULoadStat) Avg(count float64) {
	acls.Usr /= count
	acls.Sys /= count
	acls.Idle /= count
}

func (acls *AvgCPULoadStat) Ceil() {
	acls.Usr = math.Ceil(acls.Usr*100) / 100
	acls.Sys = math.Ceil(acls.Sys*100) / 100
	acls.Idle = math.Ceil(acls.Idle*100) / 100
}

type DiskInfo struct {
	Tps            float64
	KbReadPerSec   float64
	KbWritenPerSec float64
}

func (di *DiskInfo) Sum(arg DiskInfo) {
	di.Tps += arg.Tps
	di.KbReadPerSec += arg.KbReadPerSec
	di.KbWritenPerSec += arg.KbWritenPerSec
}

func (di *DiskInfo) Avg(count float64) {
	di.Tps /= count
	di.KbReadPerSec /= count
	di.KbWritenPerSec /= count
}

func (di *DiskInfo) Ceil() {
	di.Tps = math.Ceil(di.Tps*100) / 100
	di.KbReadPerSec = math.Ceil(di.KbReadPerSec*100) / 100
	di.KbWritenPerSec = math.Ceil(di.KbWritenPerSec*100) / 100
}

type DiskInfoStats map[string]DiskInfo

func (dis *DiskInfoStats) Sum(arg DiskInfoStats) {
	di := DiskInfo{}
	for key := range arg {
		di = (*dis)[key]
		di.Sum(arg[key])
		(*dis)[key] = di
	}
}

func (dis *DiskInfoStats) Avg(count float64) {
	di := DiskInfo{}
	for key := range *dis {
		di = (*dis)[key]
		di.Avg(count)
		(*dis)[key] = di
	}
}

func (dis *DiskInfoStats) Ceil() {
	di := DiskInfo{}
	for key := range *dis {
		di = (*dis)[key]
		di.Ceil()
		(*dis)[key] = di
	}
}

type FSD struct {
	Source  string
	FS      string
	Total   float64
	Used    float64
	Percent float64
}

func (fsd *FSD) Sum(arg FSD) {
	fsd.Total += arg.Total
	fsd.Used += arg.Used
	fsd.Percent += arg.Percent
}

func (fsd *FSD) Avg(count float64) {
	fsd.Total /= count
	fsd.Used /= count
	fsd.Percent /= count
}

func (fsd *FSD) Ceil() {
	fsd.Total = math.Ceil(fsd.Total)
	fsd.Used = math.Ceil(fsd.Used)
	fsd.Percent = math.Ceil(fsd.Percent)
}

type FSDInfoStat struct {
	FSDBlocks []FSD
	FSDInodes []FSD
}

func NewFSDIS(length int) *FSDInfoStat {
	return &FSDInfoStat{make([]FSD, length), make([]FSD, length)}
}

func (fis *FSDInfoStat) Sum(arg FSDInfoStat) {
	for idx := range arg.FSDBlocks {
		fis.FSDBlocks[idx].Source = arg.FSDBlocks[idx].Source
		fis.FSDBlocks[idx].FS = arg.FSDBlocks[idx].FS
		fis.FSDBlocks[idx].Sum(arg.FSDBlocks[idx])
	}
	for idx := range arg.FSDInodes {
		fis.FSDInodes[idx].Source = arg.FSDInodes[idx].Source
		fis.FSDInodes[idx].FS = arg.FSDInodes[idx].FS
		fis.FSDInodes[idx].Sum(arg.FSDInodes[idx])
	}
}

func (fis *FSDInfoStat) Avg(count float64) {
	for idx := range fis.FSDBlocks {
		fis.FSDBlocks[idx].Avg(count)
	}
	for idx := range fis.FSDInodes {
		fis.FSDInodes[idx].Avg(count)
	}
}

func (fis *FSDInfoStat) Ceil() {
	for idx := range fis.FSDBlocks {
		fis.FSDBlocks[idx].Ceil()
	}
	for idx := range fis.FSDInodes {
		fis.FSDInodes[idx].Ceil()
	}
}

// TODO:.
type TopTalkersStat struct {
	_ int
}

// No statistic for averaging send actual??
type Listeners struct {
	ProgPid string
	User    string
	Protoc  string
	Port    int
}

type NetStat struct {
	TUListeners    []Listeners
	TCPStatesCount map[string]float64
}

func NewNetStat(length int) *NetStat {
	return &NetStat{make([]Listeners, length), make(map[string]float64)}
}

func (ns *NetStat) Sum(arg NetStat) {
	ns.TUListeners = arg.TUListeners
	for key := range arg.TCPStatesCount {
		ns.TCPStatesCount[key] += arg.TCPStatesCount[key]
	}
}

func (ns *NetStat) Avg(count float64) {
	for key := range ns.TCPStatesCount {
		ns.TCPStatesCount[key] /= count
	}
}

func (ns *NetStat) Ceil() {
	for key := range ns.TCPStatesCount {
		ns.TCPStatesCount[key] = math.Ceil(ns.TCPStatesCount[key])
	}
}

type Statistic struct {
	ASLStat  *AvgSysLoadStat // Ave System Load statistic
	ACLStat  *AvgCPULoadStat // Ave Cpu Load statistic
	DIStat   DiskInfoStats   // Disk Load statistic
	FSDIStat *FSDInfoStat    // Fyle system info
	TTStat   *TopTalkersStat // Top talkers statistic
	NStat    *NetStat        // Net statistic
}
