package main

import (
	"database/sql"
	"fmt"
	"github.com/astaxie/beedb"
	_ "github.com/ziutek/mymysql/godrv"
	"time"
)

/*
CREATE TABLE `disk` (
 `uid` INT(10) AUTO_INCREMENT,
 `uuid` VARCHAR(64),
 `location` VARCHAR(64),
 `machineId` VARCHAR(64),
 `created` DATETIME DEFAULT NULL,
 PRIMARY KEY (`uid`)
);
CREATE TABLE `machine` (
	`uid` INT(10) AUTO_INCREMENT,
	`uuid` VARCHAR(64),
	`ip` VARCHAR(64),
	`slotnr` INT(10),
	`created` DATETIME DEFAULT NULL,
	PRIMARY KEY (`uid`)
);
*/

var orm beedb.Model

/*
type Disk struct {
	Uid        int `beedb:"PK"`
	Uuid       string
	Location   string
	MachineId  string
	Created    time.Time `orm:"index"`
}*/

type Disk struct { //the table'name is disk
	Uid                  int `beedb:"PK"`
	Uuid                 string
	Location             string
	MachineId            string
	Created              time.Time `orm:"index"`
	Status               string
	Role                 string
	Raid                 string
	Size                 string
	RawReadErrorRate     string
	SpinUpTime           string
	StartStopCount       string
	ReallocatedSectorCt  string
	SeekErrorRate        string
	PowerOnHours         string
	SpinRetryCount       string
	PowerCycleCount      string
	PowerOffRetractCount string
	LoadCycleCount       string
	CurrentPendingSector string
	OfflineUncorrectable string
	UDMACRCErrorCount    string
}

type Machine struct {
	Uid     int `beedb:"PK"`
	Uuid    string
	Ip      string
	Slotnr  int
	Created time.Time `orm:"index"`
}

func Initdb() {
	db, err := sql.Open("mymysql", "speediodb/root/passwd")
	if err != nil {
		panic(err)
	}
	orm = beedb.New(db)
	//fmt.Printf("open database\n")
	return
}

//func InsertdiskInfo(uuid string, location string, machineId string, RawReadErrorRate string, SpinUpTime string, StartStopCount string, ReallocatedSectorCt string, SeekErrorRate string, PowerOnHours string, SpinRetryCount string, PowerCycleCount string, PowerOffRetractCount string, LoadCycleCount string, CurrentPendingSector string, OfflineUncorrectable string, UDMACRCErrorCount string) error{
func InsertSmartInfo(uuid, location, machineId, RawReadErrorRate, SpinUpTime, StartStopCount, ReallocatedSectorCt, SeekErrorRate, PowerOnHours, SpinRetryCount, PowerCycleCount, PowerOffRetractCount, LoadCycleCount, CurrentPendingSector, OfflineUncorrectable, UDMACRCErrorCount string) error {

	fmt.Printf("Insert smartInfo start\n")
	//fmt.Printf("@@@",uuid, location, machineId, RawReadErrorRate, SpinUpTime, StartStopCount, ReallocatedSectorCt, SeekErrorRate, PowerOnHours, SpinRetryCount, PowerCycleCount, PowerOffRetractCount, LoadCycleCount, CurrentPendingSector, OfflineUncorrectable, UDMACRCErrorCount )
	var disk Disk
	disk.Uuid = uuid
	disk.Location = location
	disk.MachineId = machineId
	disk.Created = time.Now()
	//disk.Status	= ""
	//disk.Role	= ""
	//disk.Raid	= ""
	//disk.Size	= ""
	disk.RawReadErrorRate = RawReadErrorRate
	disk.SpinUpTime = SpinUpTime
	disk.StartStopCount = StartStopCount
	disk.ReallocatedSectorCt = ReallocatedSectorCt
	disk.SeekErrorRate = SeekErrorRate
	disk.PowerOnHours = PowerOnHours
	disk.SpinRetryCount = SpinRetryCount
	disk.PowerCycleCount = PowerCycleCount
	disk.PowerOffRetractCount = PowerOffRetractCount
	disk.LoadCycleCount = LoadCycleCount
	disk.CurrentPendingSector = CurrentPendingSector
	disk.OfflineUncorrectable = OfflineUncorrectable
	disk.UDMACRCErrorCount = UDMACRCErrorCount

	if err := orm.Save(&disk); err != nil {
		return err
	}
	fmt.Printf("Insert smartInfo finished\n")
	return nil
}

func InsertDisk(uuid string, location string, machineId string) error {
	var disk Disk
	disk.Uuid = uuid
	disk.Location = location
	disk.MachineId = machineId
	disk.Created = time.Now()
	if err := orm.Save(&disk); err != nil {
		return err
	}
	fmt.Printf("Insert Disk finished\n")
	return nil
}

func SelectDisksOfMachine(uuid string) ([]Disk, error) {
	var ones []Disk
	if err := orm.Where("MachineId=?", uuid).FindAll(&ones); err != nil {
		//select *,count(distinct uuid) from disk group by uuid;
		return ones, err
	}
	return ones, nil
}

func SelectAllDisks() ([]Disk, error) {
	//get all data
	var ones []Disk
	if err := orm.FindAll(&ones); err != nil {
		return ones, err
	}
	return ones, nil
}

func SelectDisk(uuid string) (Disk, error) {
	var one Disk
	if err := orm.Where("Uuid=?", uuid).Find(&one); err != nil {
		return one, err
	}
	return one, nil
}

func UpdateDisk(uuid string, location string, machineId string, status string, role string, raid string, size string) error {
	// //update data
	saveone, _ := SelectDisk(uuid)
	saveone.Uuid = uuid
	saveone.Location = location
	saveone.MachineId = machineId
	saveone.Created = time.Now()
	saveone.Status = status
	saveone.Role = role
	saveone.Raid = raid
	saveone.Size = size
	if err := orm.Save(&saveone); err != nil {
		return err
	}
	return nil
}

func DeleteDisk(uuid string) error {
	// // //delete one data
	if _, err := orm.SetTable("disk").Where("uuid=?", uuid).DeleteRow(); err != nil {
		return err
	}
	fmt.Printf("delete disk finished\n")
	return nil
}

func DeleteAllDisks() error {
	// //delete all data
	alldisks, err := SelectAllDisks()
	if err != nil {
		return err
	}

	if _, err = orm.DeleteAll(&alldisks); err != nil {
		return err
	}
	return nil
}

func InsertMachine(uuid string, ip string, slotnr int) error {
	var one Machine
	one.Uuid = uuid
	one.Ip = ip
	one.Slotnr = slotnr
	one.Created = time.Now()
	if err := orm.Save(&one); err != nil {
		return err
	}

	return nil
}

func SelectAllMachines() ([]Machine, error) {
	//get all data
	var ones []Machine
	if err := orm.FindAll(&ones); err != nil {
		return ones, err
	}
	return ones, nil
}

func SelectMachine(uuid string) (Machine, error) {
	var one Machine
	if err := orm.Where("Uuid=?", uuid).Find(&one); err != nil {
		return one, err
	}
	return one, nil
}

func DeleteMachine(uuid string) error {
	// // //delete one data
	if _, err := orm.SetTable("machines").Where("uuid=?", uuid).DeleteRow(); err != nil {
		return err
	}

	return nil
}

func UpdateMachine(uuid string, ip string, slotnr int) error {
	// //update data
	saveone, _ := SelectMachine(uuid)
	saveone.Uuid = uuid
	saveone.Ip = ip
	saveone.Slotnr = slotnr
	saveone.Created = time.Now()
	if err := orm.Save(&saveone); err != nil {
		return err
	}
	return nil
}
