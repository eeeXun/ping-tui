package main

import (
	"bytes"
	"github.com/rivo/tview"
	"net"
	"time"
)

var (
	// ICMP
	send_pkt    = ICPMPacket{Type: 8, SequenceNum: 1}
	recv_pkt    ICPMPacket
	laddr       = net.IPAddr{IP: net.ParseIP("0.0.0.0")}
	send_buffer bytes.Buffer
	recv_buffer = make([]byte, 1024)
	// UI
	app        = tview.NewApplication()
	input_box  = tview.NewInputField()
	output_box = tview.NewTextView()
	screen     = Screen{title: "PING"}
	// Control
	stop_ping bool
	sec       = time.Second
)

func main() {
	UIInit()

	if err := app.SetRoot(
		tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(input_box, 0, 1, true).
			AddItem(output_box, 0, 6, false),
		true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
