package api

import (
	"fmt"
	"gin-project/global"
	"gin-project/model"
	"gin-project/model/response"
	"gin-project/service"
	"gin-project/utils"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"reflect"
	"strconv"
	"time"
)

const defaultFailMsg = "device information acquisition failed"

// 从Excel文件批量导入设备信息
func BatchCreateDevice(c *gin.Context) {
	err, dataInfo := utils.UploadExcel(c)
	if err != nil {
		utils.ResponseCodeMsg(c, err.Error())
		return
	}
	var device []model.Device
	for _, line := range dataInfo {
		hostName := line["host_name"]
		//fmt.Println(reflect.TypeOf(hostName))
		resourceChildName := line["resource_child_name"]
		contentLibrary := line["content_library_name"]
		templateName := line["template_name"]
		vmIP := line["vm_ip"]
		vmName := line["vm_name"]
		vmVlan := line["vm_vlan"]
		devices := model.Device{Name: hostName, ResourceChildName: resourceChildName,
			ContentLibraryName: contentLibrary, TemplateName: templateName, VmIP: vmIP,
			VmName: vmName, VmVlan: vmVlan}
		device = append(device, devices)
	}
	//fmt.Println(device)
	fmt.Println(reflect.TypeOf(device))
	global.DB.Create(&device)
	utils.ResponseCodeMsg(c, "device create success")
	return

}

// 更新设备信息

type DeleteID struct {
	//ID int `json:"id"`
	IDList []int `json:"id"`
}

// 删除设备
func BatchDeleteDevice(c *gin.Context) {
	req := DeleteID{}
	err := c.BindJSON(&req)
	device := model.Device{}
	var devices []model.Device
	global.DB.Find(&devices)
	fmt.Println(devices)
	//fmt.Println(reflect.TypeOf(device))
	if err != nil {
		utils.ResponseCodeMsg(c, "JSON data error")
		return
	}
	if len(req.IDList) == 1 {
		return
	}
	fmt.Println(req.IDList)
	dbErr := service.DeleteDevice(device, req.IDList[0])
	if dbErr != nil {
		utils.ResponseCodeMsg(c, dbErr.Error())
		return
	}
	utils.ResponseCodeMsg(c, "data deletion successful")

}

// 删除单条设备
func DeleteDevice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if id <= 0 || err != nil {
		return
	}
	fmt.Println(id)

}

// 添加用户
//func AddUser(user *models.User) (error error) {
//	user.Password = utils.ScryptPassword(user.Password)
//	err := global.DB.Create(&user).Error
//	if err != nil {
//		return err
//	}
//	return nil
//}

// 获取设备内存信息
func GetDeviceMem(c *gin.Context) {
	Num := uint64(1024 * 1024 * 1024)
	NumFree := uint64(1024 * 1024)
	var memJson response.Memory
	v, err := mem.VirtualMemory()
	if err != nil {
		utils.ResponseCodeMsg(c, defaultFailMsg)
		return
	}
	memJson.Total = strconv.Itoa(int(v.Total/Num)) + "GB"
	memJson.Free = strconv.Itoa(int(v.Free/NumFree)) + "MB"
	memJson.Used = strconv.Itoa(int(v.Used/Num)) + "GB"
	memJson.UsedPercent = strconv.Itoa(int(v.UsedPercent)) + "%"
	utils.ResponseJsonMsg(c, memJson)
}

// 获取设备CPU信息
func GetDeviceCPU(c *gin.Context) {
	var errMsg response.Error
	v, err := cpu.Info()
	if err == nil {
		errMsg.Error = defaultFailMsg
		service.ErrorInfo(c, errMsg.Error)
		return
	}
	service.LogicDeviceCPU(c, v)
}

// 获取设备硬盘信息
func GetDeviceDisk(c *gin.Context) {
	// IO 统计信息
	var errMsg response.Error
	id, errID := strconv.Atoi(c.Param("id"))
	if id <= 0 || errID != nil {
		return
	}
	v, err := disk.IOCounters()
	infos, _ := disk.Usage("/")
	fmt.Println(infos)
	//vPart, err1 := disk.Partitions(true)
	if err != nil {
		errMsg.Error = defaultFailMsg
		service.ErrorInfo(c, errMsg.Error)
		return
	}
	service.LogicDeviceDiskIO(c, v, uint(id))
}

// 获取内核版本和平台信息 开机时间
func GetHostInfo(c *gin.Context) {
	id, errID := strconv.Atoi(c.Param("id"))
	if id <= 0 || errID != nil {
		return
	}
	// 开机时间
	timestamp, _ := host.BootTime()
	t := time.Unix(int64(timestamp), 0)
	onTime := t.Local().Format("2006-01-02 15:04:05")
	fmt.Println(onTime)
	// 内核版本
	version, _ := host.KernelVersion()
	fmt.Println(version)
	// 平台信息
	platform, family, version, _ := host.PlatformInformation()
	fmt.Println("platform:", platform)
	fmt.Println("family:", family)
	fmt.Println("version:", version)
	users, _ := host.Users()
	fmt.Println("user", users)
}

// 获取机器上所有的进程
func GetDeviceProcess(c *gin.Context) {
	id, errID := strconv.Atoi(c.Param("id"))
	if id <= 0 || errID != nil {
		return
	}
	processes, _ := process.Processes()
	for _, p := range processes {
		name, _ := p.Name()
		fmt.Println(p.Pid, name)
	}
}