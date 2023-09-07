package interfaces

type StatsAgent interface {
	AverageSystemLoad() (string, error)
	AverageCpuLoad() (string, error)
	LoadDiskInfo() (string, error)
	FileSystemDiskInfo() (string, error)
	TopTalkersInfo() (string, error)
	NetStatistics() (string, error)
}
