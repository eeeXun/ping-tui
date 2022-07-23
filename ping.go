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
	sum := (^((uint16(pkt.Type) << 8) +
		uint16(pkt.Code) +
		uint16(pkt.Identifier) +
		uint16(pkt.SequenceNum)))

	return sum
}

func Ping(dest string) {
	var (
		// ICMP
		raddr, _        = net.ResolveIPAddr("ip", dest)
		seq      uint16 = 1
		TTL      uint8
		// Time
		startTime time.Time
		duration  float64
	)

	for ; ; seq++ {
		if stop_ping {
			return
		}

		// Init
		send_pkt.SequenceNum = seq
		send_buffer.Reset()
		recv_buffer = make([]byte, 1024)

		// Generate Checksum
		send_pkt.Checksum = send_pkt.CalcChecksum()
		binary.Write(&send_buffer, binary.BigEndian, send_pkt)

		conn, err := net.DialIP("ip4:icmp", &laddr, raddr)
		CheckErr(err)
		defer conn.Close()

		// Send packet
		screen.title = fmt.Sprintf("PING %s (%s)", dest, raddr.String())
		screen.UpdateTitle()
		_, err = conn.Write(send_buffer.Bytes())
		CheckErr(err)
		conn.SetReadDeadline(time.Now().Add(3 * sec))

		// Recv packet
		startTime = time.Now()
		_, err = conn.Read(recv_buffer)
		CheckErr(err)
		duration = float64(time.Since(startTime).Nanoseconds()) / 1000000

		TTL = uint8(recv_buffer[8])
		// Check if Checksum is correct
		recv_pkt = ICPMPacket{Type: uint8(recv_buffer[20]),
			Code:        uint8(recv_buffer[21]),
			Checksum:    (uint16(recv_buffer[22]) << 8) + uint16(recv_buffer[23]),
			Identifier:  (uint16(recv_buffer[24]) << 8) + uint16(recv_buffer[25]),
			SequenceNum: (uint16(recv_buffer[26]) << 8) + uint16(recv_buffer[27]),
		}

		if recv_pkt.CalcChecksum() != recv_pkt.Checksum {
			app.Stop()
			fmt.Println(recv_buffer)
			fmt.Println(recv_pkt.Checksum, recv_pkt.CalcChecksum())
			panic("The checksum of the reply is incorrect")
		}

		if len(screen.content) == 0 {
			screen.content = fmt.Sprintf("Reply from (%s): imcp_seq=%d ttl=%d time=%.2fms", raddr.String(), recv_pkt.SequenceNum, TTL, duration)
		} else {
			screen.content = fmt.Sprintf("Reply from (%s): imcp_seq=%d ttl=%d time=%.2fms\n%s", raddr.String(), recv_pkt.SequenceNum, TTL, duration, screen.content)
		}

		screen.UpdateContent()

		time.Sleep(1 * sec)
	}
}
