package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
)

func main() {
	UIInit()
	dst := "google.com"
	raddr, _ := net.ResolveIPAddr("ip", dst)

	// Generate Checksum
	send_pkt.Checksum = send_pkt.CalcChecksum()
	binary.Write(&send_buffer, binary.BigEndian, send_pkt)

	conn, err := net.DialIP("ip4:icmp", &laddr, raddr)
	CheckErr(err)
	defer conn.Close()

	// Send packet
	fmt.Printf("PING %s (%s)\n", dst, raddr.String())
	_, err = conn.Write(send_buffer.Bytes())
	CheckErr(err)
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))

	// Recv packet
	startTime := time.Now()
	_, err = conn.Read(recv_buffer)
	CheckErr(err)
	duration := float64(time.Since(startTime).Nanoseconds()) / 1000000

	// Check Checksum is correct
	recv_pkt = ICPMPacket{Type: uint8(recv_buffer[20]),
		Code:        uint8(recv_buffer[21]),
		Checksum:    (uint16(recv_buffer[22]) << 8) + uint16(recv_buffer[23]),
		Identifier:  (uint16(recv_buffer[24]) << 8) + uint16(recv_buffer[25]),
		SequenceNum: (uint16(recv_buffer[26]) << 8) + uint16(recv_buffer[27]),
	}

	if recv_pkt.CalcChecksum() != recv_pkt.Checksum {
		panic("The checksum of the reply is incorrect")
	}

	fmt.Printf("Reply from (%s): imcp_seq=%d time=%.2fms\n", raddr.String(), recv_pkt.SequenceNum, duration)

	if err := app.SetRoot(
		tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(input_box, 0, 1, true).
			AddItem(output_box, 0, 5, false),
		true).Run(); err != nil {
		panic(err)
	}
}
