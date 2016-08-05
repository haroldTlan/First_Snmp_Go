// traserver
package main

import (
	"fmt"
	"snmpserver/snmp"
	"regexp"
	"strings"
)

type DiskInfo struct {
	Uuid         string
	Location     string
	MachineId    string
}

func NewDiskInfo(machineId string) *DiskInfo {
	return &DiskInfo{MachineId: machineId}
}

func RefreshDisks(ip string, machineId string) error{
	oid := "1.3.6.1.4.1.8888.1.1.0"
	out, err := snmp.Get(ip, oid)
	if err != nil {
		return err
	}

	disks := extractDisks(out, machineId)
	for _, disk := range disks {
		fmt.Println(disk.Uuid, disk.Location, disk.MachineId)
		UpdateDisk(disk.Uuid, disk.Location, disk.MachineId)
	}

	return err
}

func extractDisks(out string, machineId string) []*DiskInfo {
	disks := make([]*DiskInfo, 0)

	disks_tmp := strings.Split(out, "[Disk_location]")

	for _, disk_tmp := range disks_tmp[1:] {
		disk := extractSingleDisk(disk_tmp, machineId)
		disks = append(disks, disk)
	}

	return disks
}

func extractSingleDisk(out string, machineId string) *DiskInfo {
	disk := NewDiskInfo(machineId)
	regLocation := regexp.MustCompile(`:(\d.\d.\d+)`)
	regUuid := regexp.MustCompile(`\[Disk_uuid\]:\s*(\S+)`)
	disk.Location = regLocation.FindStringSubmatch(out)[1]
	disk.Uuid = regUuid.FindStringSubmatch(out)[1]
	
	return disk
}