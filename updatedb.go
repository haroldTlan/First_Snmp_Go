package main

func Rundb() {
	go func() {
		sub := trapTopic.Subscribe()
		defer trapTopic.Unsubscribe(sub)

		for {
			data := <-sub
			event := data.(DiskEvent)

			switch event.Name {
				case "DiskPlugged":
					InsertDisk(event.Uuid, event.Location, event.MachineId)
					
				case "DiskUnplugged":
					DeleteDisk(event.Uuid)
					
				case "DiskUpdate":
					UpdateDisk(event.Uuid, event.Location, event.MachineId)

				default:
			}
		}
	}()
}