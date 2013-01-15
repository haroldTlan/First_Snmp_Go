package main

import (
	"database/sql"
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

type Disk struct {
	Uid        int `beedb:"PK"`
	Uuid       string
	Location   string
	MachineId  string
	Created    time.Time `orm:"index"`
}

type Machine struct {
	Uid        int `beedb:"PK"`
	Uuid       string
	Ip         string
	Slotnr     int
	Created    time.Time `orm:"index"`
}

func Initdb() {
	db, err := sql.Open("mymysql", "speediodb/root/passwd")
	if err != nil {
		panic(err)
	}
	orm = beedb.New(db)
	return
}

func InsertDisk(uuid string, location string, machineId string) error{
	var disk Disk
	disk.Uuid = uuid
	disk.Location = location
	disk.MachineId = machineId
	disk.Created = time.Now()
	if err := orm.Save(&disk); err != nil {
		return err
	}

	return nil
}

func SelectDisksOfMachine(uuid string) ([]Disk, error) {
	var ones []Disk
	if err := orm.Where("MachineId=?", uuid).FindAll(&ones); err != nil {
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

func SelectDisk(uuid string) (Disk, error){
	var one Disk
	if err := orm.Where("Uuid=?", uuid).Find(&one); err != nil {
		return one, err
	}
	return one, nil
}

func UpdateDisk(uuid string, location string, machineId string) error{
	// //update data
	saveone, _ := SelectDisk(uuid)
	saveone.Uuid = uuid
	saveone.Location = location
	saveone.MachineId = machineId
	saveone.Created = time.Now()
	if err := orm.Save(&saveone); err != nil {
		return err
	}
	return nil
}

func DeleteDisk(uuid string) error{
	// // //delete one data
	if _, err := orm.SetTable("disk").Where("uuid=?", uuid).DeleteRow(); err != nil {
		return err
	}

	return nil
}

func DeleteAllDisks() error{
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



func InsertMachine(uuid string, ip string, slotnr int) error{
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

func SelectAllMachines() ([]Machine, error){
	//get all data
	var ones []Machine
	if err := orm.FindAll(&ones); err != nil {
		return ones, err
	}
	return ones, nil
}

func SelectMachine(uuid string) (Machine, error){
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


