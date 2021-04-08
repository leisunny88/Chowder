package service

import (
	"errors"
	"fmt"
	"gin-project/global"
	"gin-project/model"
	"gin-project/model/response"
	"gin-project/utils"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"strconv"
)

// 处理错误
func ErrorInfo(c *gin.Context, err string)  {
	utils.ResponseNotFoundCode(c, err)
}

// 删除设备
func DeleteDevice(m model.Device, id int) (err error) {
	dbErr := global.DB.Where("id = ?", id).Delete(&m).Error
	//dbErr := global.DB.Delete(&m).Error
	if dbErr != nil {
		return errors.New("data deletion failed")
	} else {
		return nil
	}
}

// 入库设备的disk信息


// 处理CPU信息
func LogicDeviceCPU(c *gin.Context, data []cpu.InfoStat) {
	var cpuInfo response.CPUInfo
	for _, info := range data {
		cpuInfo.CPU = strconv.Itoa(int(info.CPU))
		cpuInfo.CORES = strconv.Itoa(int(info.Cores))
		cpuInfo.MHZ = strconv.Itoa(int(info.Mhz))
		cpuInfo.ModelName = info.ModelName
		cpuInfo.CacheSize = strconv.Itoa(int(info.CacheSize))
	}
	utils.ResponseJsonMsg(c, cpuInfo)
}

// 处理disk IO 信息
func LogicDeviceDiskIO(c *gin.Context, data map[string]disk.IOCountersStat, id uint) {
	var diskInfo response.Disk
	for _, info := range data {
		diskInfo.Name = info.Name
		diskInfo.IOTime = info.IoTime
		fmt.Println(info)
	}
}