package response

type FailMsg struct {
	Fail string
}

type Error struct {
	Error string
}

type Memory struct {
	Total       string
	Free        string
	Used        string
	UsedPercent string
}

type CPUInfo struct {
	CPU       string
	CORES     string
	ModelName string
	MHZ       string
	CacheSize string
}

type Disk struct {
	Name     string
	IOTime   uint64
	DeviceID uint
}
