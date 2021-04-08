package model

import "gorm.io/gorm"

type Device struct {
	gorm.Model
	Name               string `json:"host_name"`
	ResourceChildName  string `json:"resource_child_name"`
	ContentLibraryName string `json:"content_library_name"`
	TemplateName       string `json:"template_name"`
	VmName             string `json:"vm_name"`
	VmIP               string `json:"vm_ip"`
	VmVlan             string `json:"vm_vlan"`
}

func (Device) TableName() string {
	return "device"
}

type CPU struct {
	gorm.Model
	CPU       string `json:"cpu"`
	CORES     string `json:"cores"`
	ModelName string `json:"model_name"`
	MHZ       string `json:"mhz"`
	CacheSize string `json:"cache_size"`
	DeviceID  uint
}

func (CPU) TableName() string {
	return "cpu"
}

type Memory struct {
	gorm.Model
	Total       string `json:"total"`
	Free        string `json:"free"`
	Used        string `json:"used"`
	UsedPercent string `json:"used_percent"`
	DeviceID    uint
}

func (Memory) TableName() string {
	return "memory"
}

type Disk struct {
	gorm.Model
	Name     string `json:"name"`
	IOTime   uint64 `json:"io_time"`
	DeviceID uint
}

func (Disk) TableName() string {
	return "disk"
}
