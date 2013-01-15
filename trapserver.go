// traserver
package main

import (
	"fmt"
	"net"
	"snmpserver/snmp"
	"snmpserver/topic"
)

var trapTopic = topic.New()

type DiskEvent struct {
	Name         string
	Uuid         string
	Location     string
	MachineId    string
}

const (
	DISKEVENTTYPE  = ".1.3.6.1.4.1.8888"
	EVENT          = ".1.3.6.1.4.1.8888.1.1"
	UUID           = ".1.3.6.1.4.1.8888.1.3"
	LOCATION       = ".1.3.6.1.4.1.8888.1.2"
	MACHINEID      = ".1.3.6.1.4.1.8888.1.4"
)

func newDiskEvent(values map[string]interface{}) DiskEvent {
	return DiskEvent{	Name:      values[EVENT].(string), 
						Uuid:      values[UUID].(string), 
						Location:  values[LOCATION].(string), 
						MachineId: "TODO"}
}

func TrapServer() {
	go func() {
		fmt.Println("Hello World!")
		socket,err := net.ListenUDP("udp4",&net.UDPAddr{IP:net.IPv4(0,0,0,0), Port:162})
		if err != nil {
			panic(err)
		}
		defer socket.Close()
	
		for {
			buf := make([]byte,2048)
			read,from,_:=socket.ReadFromUDP(buf)
			fmt.Println("Get msg from ",from.IP)
			HandleUdp(buf[:read])
		}
	}()
}

func HandleUdp(data []byte){
	trap,err := snmp.ParseUdp(data)
	if err !=nil{
		fmt.Println("Err",err.Error())
	}
	fmt.Println(trap.Version,trap.Community,trap.EnterpriseId,trap.Address)
	for k,v :=range trap.Values{
		fmt.Printf("%s = %s\n",k,v);
	}

	var event DiskEvent
	if trap.EnterpriseId == DISKEVENTTYPE {
		event = newDiskEvent(trap.Values)
		fmt.Println(event)
	}
	trapTopic.Publish(event)
}