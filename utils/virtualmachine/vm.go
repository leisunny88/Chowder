package virtualmachine

import (
	"context"
	"fmt"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/units"
	"github.com/vmware/govmomi/vapi/library"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vapi/vcenter"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"net/url"
	"os"
	"strings"
)

type VmWare struct {
	IP     string
	User   string
	Pwd    string
	client *govmomi.Client
	ctx    context.Context
}

func NewVmWare(IP, User, Pwd string) *VmWare {
	u := &url.URL{
		Scheme: "https",
		Host:   IP,
		Path:   "/sdk",
	}
	ctx := context.Background()
	u.User = url.UserPassword(User, Pwd)
	client, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		panic(err)
	}
	return &VmWare{
		IP:     IP,
		User:   User,
		Pwd:    Pwd,
		client: client,
		ctx:    ctx,
	}
}

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func (vw *VmWare) getBase(tp string) (v *view.ContainerView, error error) {
	m := view.NewManager(vw.client.Client)

	v, err := m.CreateContainerView(vw.ctx, vw.client.Client.ServiceContent.RootFolder, []string{tp}, true)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (vw *VmWare) GetAllVmClient() (vmList []VirtualMachines, templateList []TemplateInfo, err error) {
	v, err := vw.getBase("VirtualMachine")
	if err != nil {
		return nil, nil, err
	}
	defer v.Destroy(vw.ctx)
	var vms []mo.VirtualMachine
	err = v.Retrieve(vw.ctx, []string{"VirtualMachine"}, []string{"summary"}, &vms)
	if err != nil {
		return nil, nil, err
	}
	for _, vm := range vms {
		//if vm.Summary.Config.Name == "测试机器" {
		//	v := object.NewVirtualMachine(vw.client.Client, vm.Self)
		//	vw.setIP(v)
		//}
		if vm.Summary.Config.Template {
			templateList = append(templateList, TemplateInfo{
				Name:   vm.Summary.Config.Name,
				System: vm.Summary.Config.GuestFullName,
				Self: Self{
					Type:  vm.Self.Type,
					Value: vm.Self.Value,
				},
				VM: vm.Self,
			})
		} else {
			vmList = append(vmList, VirtualMachines{
				Name:   vm.Summary.Config.Name,
				System: vm.Summary.Config.GuestFullName,
				Self: Self{
					Type:  vm.Self.Type,
					Value: vm.Self.Value,
				},
				VM: vm.Self,
			})
		}
	}
	fmt.Println(vmList)
	return vmList, templateList, nil
}

func (vw *VmWare) GetAllHost() (hostList []*HostSummary, err error) {
	v, err := vw.getBase("HostSystem")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.ctx)
	var hss []mo.HostSystem
	err = v.Retrieve(vw.ctx, []string{"HostSystem"}, []string{"summary"}, &hss)
	if err != nil {
		return nil, err
	}
	for _, hs := range hss {
		totalCPU := int64(hs.Summary.Hardware.CpuMhz) * int64(hs.Summary.Hardware.NumCpuCores)
		freeCPU := int64(totalCPU) - int64(hs.Summary.QuickStats.OverallCpuUsage)
		freeMemory := int64(hs.Summary.Hardware.MemorySize) - (int64(hs.Summary.QuickStats.OverallMemoryUsage) * 1024 * 1024)
		hostList = append(hostList, &HostSummary{
			Host: Host{
				Type:  hs.Summary.Host.Type,
				Value: hs.Summary.Host.Value,
			},
			Name:        hs.Summary.Config.Name,
			UsedCPU:     int64(hs.Summary.QuickStats.OverallCpuUsage),
			TotalCPU:    totalCPU,
			FreeCPU:     freeCPU,
			UsedMemory:  int64((units.ByteSize(hs.Summary.QuickStats.OverallMemoryUsage)) * 1024 * 1024),
			TotalMemory: int64(units.ByteSize(hs.Summary.Hardware.MemorySize)),
			FreeMemory:  freeMemory,
			HostSelf:    hs.Self,
		})
	}
	return hostList, err
}

func (vw *VmWare) GetAllNetwork() (networkList []map[string]string, err error) {
	v, err := vw.getBase("Network")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.ctx)
	var networks []mo.Network
	err = v.Retrieve(vw.ctx, []string{"Network"}, nil, &networks)
	if err != nil {
		return nil, err
	}
	for _, net := range networks {
		networkList = append(networkList, map[string]string{
			"Vlan":      net.Name,
			"NetworkID": strings.Split(net.Reference().String(), ":")[1],
		})
	}
	return networkList, nil
}

func (vw *VmWare) GetAllDatastore() (datastoreList []DatastoreSummary, err error) {
	v, err := vw.getBase("Datastore")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.ctx)
	var dss []mo.Datastore
	err = v.Retrieve(vw.ctx, []string{"Datastore"}, []string{"summary"}, &dss)
	if err != nil {
		return nil, err
	}
	for _, ds := range dss {
		datastoreList = append(datastoreList, DatastoreSummary{
			Name: ds.Summary.Name,
			Datastore: Datastore{
				Type:  ds.Summary.Datastore.Type,
				Value: ds.Summary.Datastore.Value,
			},
			Type:          ds.Summary.Type,
			Capacity:      int64(units.ByteSize(ds.Summary.Capacity)),
			FreeSpace:     int64(units.ByteSize(ds.Summary.FreeSpace)),
			DatastoreSelf: ds.Self,
		})
	}
	return
}

func (vw *VmWare) GetHostVm() (hostVm map[string][]VMS, err error) {
	hostList, err := vw.GetAllHost() //
	if err != nil {
		return
	}
	var hostIDList []string
	hostVm = make(map[string][]VMS)
	for _, host := range hostList {
		hostIDList = append(hostIDList, host.Host.Value)
		hostVm[host.Host.Value] = []VMS{}
	}
	v, err := vw.getBase("VirtualMachine")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.ctx)
	var vms []mo.VirtualMachine
	err = v.Retrieve(vw.ctx, []string{"VirtualMachine"}, []string{"summary"}, &vms)
	if err != nil {
		return nil, err
	}
	for _, vm := range vms {
		if IsContain(hostIDList, vm.Summary.Runtime.Host.Value) {
			hostVm[vm.Summary.Runtime.Host.Value] = append(hostVm[vm.Summary.Runtime.Host.Value], VMS{
				Name:  vm.Summary.Config.Name,
				Value: vm.Summary.Vm.Value,
			})
		}
		//s, _ := json.Marshal(vm.Summary)
		//fmt.Println(string(s))
	}
	//fmt.Println(hostVm)
	return
}

func (vw *VmWare) GetAllCluster() (clusterList []ClusterInfo, err error) {
	v, err := vw.getBase("ClusterComputeResource")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.ctx)
	var crs []mo.ClusterComputeResource
	err = v.Retrieve(vw.ctx, []string{"ClusterComputeResource"}, []string{}, &crs)
	if err != nil {
		return nil, err
	}
	for _, cr := range crs {
		clusterList = append(clusterList, ClusterInfo{
			Cluster: Self{
				Type:  cr.Self.Type,
				Value: cr.Self.Value,
			},
			Name: cr.Name,
			Parent: Self{
				Type:  cr.Parent.Type,
				Value: cr.Parent.Value,
			},
			ResourcePool: Self{
				Type:  cr.ResourcePool.Type,
				Value: cr.ResourcePool.Value,
			},
			Hosts:     cr.Host,
			Datastore: cr.Datastore,
		})
	}
	fmt.Println(clusterList)
	return
}

func (vw *VmWare) GetAllDatacenter() (dataCenterList []DataCenter, err error) {
	v, err := vw.getBase("Datacenter")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.ctx)
	var dcs []mo.Datacenter
	err = v.Retrieve(vw.ctx, []string{"Datacenter"}, []string{}, &dcs)
	if err != nil {
		return nil, err
	}
	for _, dc := range dcs {
		dataCenterList = append(dataCenterList, DataCenter{
			Datacenter: Self{
				Type:  dc.Self.Type,
				Value: dc.Self.Value,
			},
			Name: dc.Name,
			VmFolder: Self{
				Type:  dc.VmFolder.Type,
				Value: dc.VmFolder.Value,
			},
			HostFolder: Self{
				Type:  dc.HostFolder.Type,
				Value: dc.HostFolder.Value,
			},
			DatastoreFolder: Self{
				Type:  dc.DatastoreFolder.Type,
				Value: dc.DatastoreFolder.Value,
			},
		})
	}
	fmt.Println(dataCenterList)
	return
}

func (vw *VmWare) GetAllResourcePool() (resourceList []ResourcePoolInfo, err error) {
	v, err := vw.getBase("ResourcePool")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.ctx)
	var rps []mo.ResourcePool
	err = v.Retrieve(vw.ctx, []string{"ResourcePool"}, []string{}, &rps)
	for _, rp := range rps {
		//if rp.Name == "测试虚机" {
		//	s, _ := json.Marshal(rp)
		//	fmt.Println(string(s))
		//}
		resourceList = append(resourceList, ResourcePoolInfo{
			ResourcePool: Self{
				Type:  rp.Self.Type,
				Value: rp.Self.Value,
			},
			Name: rp.Name,
			Parent: Self{
				Type:  rp.Parent.Type,
				Value: rp.Parent.Value,
			},
			ResourcePoolList: rp.ResourcePool,
			Resource:         rp.Self,
		})
	}
	return
}

func (vw *VmWare) GetFolder() (folderList []FolderInfo, err error) {
	v, err := vw.getBase("Folder")
	if err != nil {
		return nil, err
	}
	defer v.Destroy(vw.ctx)
	var folders []mo.Folder
	err = v.Retrieve(vw.ctx, []string{"Folder"}, []string{}, &folders)
	for _, folder := range folders {
		//newFolder := object.NewFolder(vw.client.Client, folder.Self)
		//fmt.Println(newFolder)
		folderList = append(folderList, FolderInfo{
			Folder: Self{
				Type:  folder.Self.Type,
				Value: folder.Self.Value,
			},
			Name:        folder.Name,
			ChildEntity: folder.ChildEntity,
			Parent: Self{
				Type:  folder.Parent.Type,
				Value: folder.Parent.Value,
			},
			FolderSelf: folder.Self,
		})
		//break
	}
	return folderList, nil
}

func (vw *VmWare) getLibraryItem(ctx context.Context, rc *rest.Client) (*library.Item, error) {
	const (
		libraryName     = "模板"
		libraryItemName = "template-rehl7.7"
		libraryItemType = "ovf"
	)

	m := library.NewManager(rc)
	libraries, err := m.FindLibrary(ctx, library.Find{Name: libraryName})
	if err != nil {
		fmt.Printf("Find library by name %s failed, %v", libraryName, err)
		return nil, err
	}

	if len(libraries) == 0 {
		fmt.Printf("Library %s was not found", libraryName)
		return nil, fmt.Errorf("library %s was not found", libraryName)
	}

	if len(libraries) > 1 {
		fmt.Printf("There are multiple libraries with the name %s", libraryName)
		return nil, fmt.Errorf("there are multiple libraries with the name %s", libraryName)
	}

	items, err := m.FindLibraryItems(ctx, library.FindItem{Name: libraryItemName,
		Type: libraryItemType, LibraryID: libraries[0]})

	if err != nil {
		fmt.Printf("Find library item by name %s failed", libraryItemName)
		return nil, fmt.Errorf("find library item by name %s failed", libraryItemName)
	}

	if len(items) == 0 {
		fmt.Printf("Library item %s was not found", libraryItemName)
		return nil, fmt.Errorf("library item %s was not found", libraryItemName)
	}

	if len(items) > 1 {
		fmt.Printf("There are multiple library items with the name %s", libraryItemName)
		return nil, fmt.Errorf("there are multiple library items with the name %s", libraryItemName)
	}

	item, err := m.GetLibraryItem(ctx, items[0])
	if err != nil {
		fmt.Printf("Get library item by %s failed, %v", items[0], err)
		return nil, err
	}
	return item, nil
}

func (vw *VmWare) CreateVM() {
	createData := CreateMap{
		TempName:    "template-rehl7.7",
		Datacenter:  "Datacenter",
		Cluster:     "AsiaLink-Production",
		Host:        "192.168.100.201",
		Resources:   "测试虚机",
		Storage:     "local-esxi-201",
		VmName:      "测试机器one",
		SysHostName: "test",
		Network:     "vlan80",
	}
	_, templateList, err := vw.GetAllVmClient()
	if err != nil {
		panic(err)
	}
	var templateNameList []string
	for _, template := range templateList {
		templateNameList = append(templateNameList, template.Name)
	}
	if !IsContain(templateNameList, createData.TempName) {
		fmt.Fprintf(os.Stderr, "模版不存在，虚拟机创建失败")
		return
	}
	resourceList, err := vw.GetAllResourcePool()
	if err != nil {
		panic(err)
	}
	var resourceStr, resourceID string
	for _, resource := range resourceList {
		if resource.Name == createData.Resources {
			resourceStr = resource.Name
			resourceID = resource.ResourcePool.Value
		}
	}
	if resourceStr == "" {
		fmt.Fprintf(os.Stderr, "资源池不存在，虚拟机创建失败")
		return
	}
	fmt.Println("ResourceID", resourceID)
	datastoreList, err := vw.GetAllDatastore()
	if err != nil {
		panic(err)
	}
	var datastoreID, datastoreStr string
	for _, datastore := range datastoreList {
		if datastore.Name == createData.Storage {
			datastoreID = datastore.Datastore.Value
			datastoreStr = datastore.Name
		}
	}
	if datastoreStr == "" {
		fmt.Fprintf(os.Stderr, "存储中心不存在，虚拟机创建失败")
		return
	}
	fmt.Println("DatastoreID", datastoreID)
	networkList, err := vw.GetAllNetwork()
	if err != nil {
		panic(err)
	}
	var networkID, networkStr string
	for _, network := range networkList {
		if network["Vlan"] == createData.Network {
			networkStr = network["Vlan"]
			networkID = network["NetworkID"]
		}
	}

	if networkStr == "" {
		fmt.Fprintf(os.Stderr, "网络不存在，虚拟机创建失败")
		return
	}
	fmt.Println("NetworkID", networkID)
	finder := find.NewFinder(vw.client.Client)
	//resourcePools, err := finder.DatacenterList(vw.ctx, "*")
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Failed to list resource pool at vc %v", err)
	//	os.Exit(1)
	//}
	//fmt.Println(reflect.TypeOf(resourcePools[0].Reference().Value), resourcePools)
	folders, err := finder.FolderList(vw.ctx, "*")
	var folderID string
	for _, folder := range folders {
		if folder.InventoryPath == "/"+createData.Datacenter+"/vm" {
			folderID = folder.Reference().Value
		}
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list folder at vc  %v", err)
		return
	}
	rc := rest.NewClient(vw.client.Client)
	if err := rc.Login(vw.ctx, url.UserPassword(vw.User, vw.Pwd)); err != nil {
		fmt.Fprintf(os.Stderr, "rc Login filed, %v", err)
		return
	}
	item, err := vw.getLibraryItem(vw.ctx, rc)
	if err != nil {
		panic(err)
	}
	//cloneSpec := &types.VirtualMachineCloneSpec{
	//	PowerOn:  false,
	//	Template: cmd.template,
	//}
	// 7fa9e782-cba2-4061-95fc-4ebb08ec127a
	fmt.Println("Item", item.ID)
	m := vcenter.NewManager(rc)
	fr := vcenter.FilterRequest{
		Target: vcenter.Target{
			ResourcePoolID: resourceID,
			FolderID:       folderID,
		},
	}
	r, err := m.FilterLibraryItem(vw.ctx, item.ID, fr)
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
	fmt.Println(11111111111, r.Networks, r.StorageGroups)
	networkKey := r.Networks[0]
	//storageKey := r.StorageGroups[0]
	deploy := vcenter.Deploy{
		DeploymentSpec: vcenter.DeploymentSpec{
			Name:               createData.VmName,
			DefaultDatastoreID: datastoreID,
			AcceptAllEULA:      true,
			NetworkMappings: []vcenter.NetworkMapping{
				{
					Key:   networkKey,
					Value: networkID,
				},
			},
			StorageMappings: []vcenter.StorageMapping{{
				Key: "",
				Value: vcenter.StorageGroupMapping{
					Type:         "DATASTORE",
					DatastoreID:  datastoreID,
					Provisioning: "thin",
				},
			}},
			StorageProvisioning: "thin",
		},
		Target: vcenter.Target{
			ResourcePoolID: resourceID,
			FolderID:       folderID,
		},
	}
	ref, err := vcenter.NewManager(rc).DeployLibraryItem(vw.ctx, item.ID, deploy)
	if err != nil {
		fmt.Println(4444444444, err)
		panic(err)
	}
	f := find.NewFinder(vw.client.Client)
	obj, err := f.ObjectReference(vw.ctx, *ref)
	if err != nil {
		panic(err)
	}
	_ = obj.(*object.VirtualMachine)

	//datastores, err := finder.VirtualMachineList(vw.ctx, "*/group-v629")
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Failed to list datastore at vc %v", err)
	//	os.Exit(1)
	//}
	//fmt.Println(datastores)
}

func (vw *VmWare) CloneVM() {
	cloneData := CreateMap{
		TempName:    "template-rehl7.7",
		Datacenter:  "Datacenter",
		Cluster:     "AsiaLink-Production",
		Host:        "192.168.100.201",
		Resources:   "测试虚机",
		Storage:     "local-esxi-201",
		VmName:      "测试机器",
		SysHostName: "test",
		Network:     "vlan80",
	}
	vmList, templateList, err := vw.GetAllVmClient()
	if err != nil {
		panic(err)
	}
	var templateNameList []string
	var vmTemplate types.ManagedObjectReference
	for _, template := range templateList {
		templateNameList = append(templateNameList, template.Name)
		if template.Name == cloneData.TempName {
			vmTemplate = template.VM
		}
	}
	if !IsContain(templateNameList, cloneData.TempName) {
		fmt.Fprintf(os.Stderr, "模版不存在，虚拟机克隆失败")
		return
	}
	dataCenterList, err := vw.GetAllDatacenter()
	if err != nil {
		panic(err)
	}
	var datacenterID, datacenterName string
	for _, datacenter := range dataCenterList {
		if datacenter.Name == cloneData.Datacenter {
			datacenterID = datacenter.Datacenter.Value
			datacenterName = datacenter.Name
		}
	}
	if datacenterName == "" {
		fmt.Fprintf(os.Stderr, "数据中心不存在，虚拟机克隆失败")
		return
	}
	hostList, err := vw.GetAllHost()
	if err != nil {
		panic(err)
	}
	var hostName string
	var hostRef types.ManagedObjectReference
	for _, host := range hostList {
		if host.Name == cloneData.Host {
			hostName = host.Name
			hostRef = host.HostSelf
		}
	}
	if hostName == "" {
		fmt.Fprintf(os.Stderr, "主机不存在，虚拟机克隆失败")
		return
	}
	resourceList, err := vw.GetAllResourcePool()
	if err != nil {
		panic(err)
	}
	var resourceStr, resourceID string
	var poolRef types.ManagedObjectReference
	for _, resource := range resourceList {
		if resource.Name == cloneData.Resources {
			resourceStr = resource.Name
			resourceID = resource.ResourcePool.Value
			poolRef = resource.Resource
		}
	}
	if resourceStr == "" {
		fmt.Fprintf(os.Stderr, "资源池不存在，虚拟机克隆失败")
		return
	}
	fmt.Println("ResourceID", resourceID)
	datastoreList, err := vw.GetAllDatastore()
	if err != nil {
		panic(err)
	}
	var datastoreID, datastoreStr string
	var datastoreRef types.ManagedObjectReference
	for _, datastore := range datastoreList {
		if datastore.Name == cloneData.Storage {
			datastoreID = datastore.Datastore.Value
			datastoreStr = datastore.Name
			datastoreRef = datastore.DatastoreSelf
		}
	}
	if datastoreStr == "" {
		fmt.Fprintf(os.Stderr, "存储中心不存在，虚拟机克隆失败")
		return
	}
	fmt.Println("DatastoreID", datastoreID)
	networkList, err := vw.GetAllNetwork()
	if err != nil {
		panic(err)
	}
	var networkID, networkStr string
	for _, network := range networkList {
		if network["Vlan"] == cloneData.Network {
			networkStr = network["Vlan"]
			networkID = network["NetworkID"]
		}
	}

	if networkStr == "" {
		fmt.Fprintf(os.Stderr, "网络不存在，虚拟机克隆失败")
		return
	}
	fmt.Println("NetworkID", networkID)
	clusterList, err := vw.GetAllCluster()
	if err != nil {
		panic(err)
	}
	var clusterID, clusterName string
	for _, cluster := range clusterList {
		if cluster.Name == cloneData.Cluster {
			clusterID = cluster.Cluster.Value
			clusterName = cluster.Name
		}
	}
	if clusterName == "" {
		fmt.Fprintf(os.Stderr, "集群不存在，虚拟机克隆失败")
		return
	}
	configSpecs := []types.BaseVirtualDeviceConfigSpec{}
	fmt.Println("ClusterID", clusterID)
	for _, vms := range vmList {
		if vms.Name == cloneData.VmName {
			fmt.Fprintf(os.Stderr, "虚机已存在，虚拟机克隆失败")
			return
		}
	}
	finder := find.NewFinder(vw.client.Client)
	folders, err := finder.FolderList(vw.ctx, "*")
	var Folder *object.Folder
	for _, folder := range folders {
		if folder.InventoryPath == "/"+cloneData.Datacenter+"/vm" {
			Folder = folder
		}
	}
	fmt.Println(Folder)
	folderList, err := vw.GetFolder()
	if err != nil {
		panic(err)
	}
	var folderRef types.ManagedObjectReference
	for _, folder := range folderList {
		if folder.Parent.Value == datacenterID && folder.Name == "vm" {
			folderRef = folder.FolderSelf
		}
	}
	fmt.Println("poolRef", poolRef)
	relocateSpec := types.VirtualMachineRelocateSpec{
		DeviceChange: configSpecs,
		Folder:       &folderRef,
		Pool:         &poolRef,
		Host:         &hostRef,
		Datastore:    &datastoreRef,
	}
	vmConf := &types.VirtualMachineConfigSpec{
		NumCPUs:  4,
		MemoryMB: 16 * 1024,
	}
	cloneSpec := &types.VirtualMachineCloneSpec{
		PowerOn:  false,
		Template: false,
		Location: relocateSpec,
		Config:   vmConf,
	}
	t := object.NewVirtualMachine(vw.client.Client, vmTemplate)
	newFolder := object.NewFolder(vw.client.Client, folderRef)
	fmt.Println(newFolder)
	fmt.Println(cloneData.VmName)
	fmt.Println(cloneSpec.Location)
	task, err := t.Clone(vw.ctx, newFolder, cloneData.VmName, *cloneSpec)
	if err != nil {
		panic(err)
	}
	fmt.Println("克隆任务开始，", task.Wait(vw.ctx))
}

func (vw *VmWare) setIP(vm *object.VirtualMachine) error {
	ipAddr := IpAddr{
		ip:       "192.168.80.108",
		netmask:  "255.255.255.0",
		gateway:  "192.168.80.254",
		hostname: "test",
	}
	cam := types.CustomizationAdapterMapping{
		Adapter: types.CustomizationIPSettings{
			Ip:         &types.CustomizationFixedIp{IpAddress: ipAddr.ip},
			SubnetMask: ipAddr.netmask,
			Gateway:    []string{ipAddr.gateway},
		},
	}
	customSpec := types.CustomizationSpec{
		NicSettingMap: []types.CustomizationAdapterMapping{cam},
		Identity:      &types.CustomizationLinuxPrep{HostName: &types.CustomizationFixedName{Name: ipAddr.hostname}},
	}
	task, err := vm.Customize(vw.ctx, customSpec)
	if err != nil {
		return err
	}
	return task.Wait(vw.ctx)
}

type IpAddr struct {
	ip       string
	netmask  string
	gateway  string
	hostname string
}

func (vw *VmWare) MigrateVM() {
	migrateData := "测试虚机"
	v, err := vw.getBase("VirtualMachine")
	if err != nil {
		panic(err)
	}
	defer v.Destroy(vw.ctx)
	var vms []mo.VirtualMachine
	err = v.Retrieve(vw.ctx, []string{"VirtualMachine"}, []string{"summary"}, &vms)
	if err != nil {
		panic(err)
	}
	var vmTarget types.ManagedObjectReference
	for _, vm := range vms {
		if vm.Summary.Config.Name == migrateData {
			vmTarget = vm.Self
		}
	}
	resourceList, err := vw.GetAllResourcePool()
	if err != nil {
		panic(err)
	}
	var resourceStr, resourceID string
	var poolRef types.ManagedObjectReference
	for _, resource := range resourceList {
		if resource.Name == "" {
			resourceStr = resource.Name
			resourceID = resource.ResourcePool.Value
			poolRef = resource.Resource
		}
	}
	if resourceStr == "" {
		fmt.Fprintf(os.Stderr, "资源池不存在，虚拟机迁移失败")
		return
	}
	fmt.Println("ResourceID", resourceID)
	hostList, err := vw.GetAllHost()
	if err != nil {
		panic(err)
	}
	var hostName string
	var hostRef types.ManagedObjectReference
	for _, host := range hostList {
		if host.Name == "192.168.100.201" {
			hostName = host.Name
			hostRef = host.HostSelf
		}
	}
	if hostName == "" {
		fmt.Fprintf(os.Stderr, "主机不存在，虚拟机迁移失败")
		return
	}
	t := object.NewVirtualMachine(vw.client.Client, vmTarget)
	pool := object.NewResourcePool(vw.client.Client, poolRef)
	host := object.NewHostSystem(vw.client.Client, hostRef)
	//var priority types.VirtualMachineMovePriority
	//var state types.VirtualMachinePowerState
	task, err := t.Migrate(vw.ctx, pool, host, "defaultPriority", "poweredOff")
	if err != nil {
		panic(err)
	}
	fmt.Println("虚拟机迁移中......")
	_ = task.Wait(vw.ctx)
	fmt.Println("虚拟机迁移完成.....")
}
