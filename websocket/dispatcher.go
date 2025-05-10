package websocket



func StartDispatcher() {
	go func() {
		for {
			noti := <-NotifyChan
			SendToClient(noti)

		}
	}()
}
