package main

//import	"fmt"


func Rundb() {
	go func() {
		sub := trapTopic.Subscribe()
		defer trapTopic.Unsubscribe(sub)

		for {
			data := <-sub
			event := data.(DiskEvent)
			
			//fmt.Printf("From updatedb.go:%s\n",event)
						
			switch event.Name {
				case "DiskPlugged":
					InsertDisk(event.Uuid, event.Location, event.MachineId)
					
				case "DiskUnplugged":
					DeleteDisk(event.Uuid)
					
				case "DiskUpdate":
					UpdateDisk(event.Uuid, event.Location, event.MachineId)
					
				case "DiskAlarm":
					InsertSmartInfo(event.Uuid, event.Location, event.MachineId, event.RawReadErrorRate, event.SpinUpTime, event.StartStopCount, event.ReallocatedSectorCt, event.SeekErrorRate, event.PowerOnHours, event.SpinRetryCount, event.PowerCycleCount, event.PowerOffRetractCount, event.LoadCycleCount, event.CurrentPendingSector, event.OfflineUncorrectable, event.UDMACRCErrorCount )					
					//fmt.Printf("DiskAlarm\n",event.Uuid, event.Location, event.MachineId, event.RawReadErrorRate, event.SpinUpTime, event.StartStopCount, event.ReallocatedSectorCt, event.SeekErrorRate, event.PowerOnHours, event.SpinRetryCount, event.PowerCycleCount, event.PowerOffRetractCount, event.LoadCycleCount, event.CurrentPendingSector, event.OfflineUncorrectable, event.UDMACRCErrorCount )
				default:
			}
		}
	}()
}