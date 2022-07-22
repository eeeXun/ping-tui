package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

//  0               1               2               3
//  0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 8 9 A B C D E F
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |     Type      |     Code      |           Checksum            |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |           Identifier          |        Sequence Number        |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |  	                          Data...                          |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
type ICPMPacket struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

// 0xFFFF - (Type * 256 + Code + Identifier + SequenceNum)
func (pkt ICPMPacket) CalcChecksum() uint16 {
	return (^((uint16(pkt.Type) << 8) +
		uint16(pkt.Code) +
		uint16(pkt.SequenceNum)))
}

func Ping(dest string) {
	var (
		raddr, _ = net.ResolveIPAddr("ip", dest)
	)

	// Generate Checksum
	send_pkt.Checksum = send_pkt.CalcChecksum()
	binary.Write(&send_buffer, binary.BigEndian, send_pkt)

	conn, err := net.DialIP("ip4:icmp", &laddr, raddr)
	CheckErr(err)
	defer conn.Close()

	// Send packet
	title = fmt.Sprintf("PING %s (%s)", dest, raddr.String())
	output_box.SetTitle(title)
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

	screen = fmt.Sprintf("Reply from (%s): imcp_seq=%d time=%.2fms\n%s", raddr.String(), recv_pkt.SequenceNum, duration, screen)
	fmt.Fprintf(output_box, screen)
}
