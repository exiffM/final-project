package config

type AgentConfig struct {
	StorageTime int64 // Time of storaging system data
	AvgSysLoad  bool  // On/off average system load
	AvgCPULoad  bool  // On/off average CPU load
	DiskLoad    bool  // On/off disks load
	DiskFsInfo  bool  // On/off disks file system info
	TTNet       bool  // On/off top talkers net stats
	NetStats    bool  // On/off net stats
}

func NewConfig() *AgentConfig {
	return &AgentConfig{}
}
