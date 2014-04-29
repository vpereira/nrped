package main

import (
    "os"
    "fmt"
    "net"
    "bytes"
    "time"
    "math/rand"
    "encoding/binary"
    "github.com/vpereira/nrped/common"

)

//it will be read from the config file
//its for now its just a mockup
var allowedCommands = map[string]int16 {
    "check_ping":int16(0),
    "check_foo":int16(1),
}
func main() {

	service := ":5666"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	common.CheckError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	common.CheckError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		// run as a goroutine
		go handleClient(conn)
	}
}

func receivePackets(conn net.Conn) common.NrpePacket {
    pkt_rcv := new(common.NrpePacket)
	err := binary.Read(conn, binary.BigEndian, pkt_rcv)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
    return *pkt_rcv
}

func IsCommandAllowed(cmd string) bool {
    if _,ok := allowedCommands[cmd]; ok {
        return true
    }else{
        return false
    }
}

func getCommand(cmd string) int16 {
    return allowedCommands[cmd]
}

func fillRandomData() string {
    char := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
    rand.Seed(time.Now().UTC().UnixNano())
    buf := make([]byte, 1024)
    for i := 0; i < 1024; i++ {
        buf[i] = char[rand.Intn(len(char)-1)]
    }
    return string(buf)
}
func prepareToSend(cmd string) common.NrpePacket {
    pkt_send := common.NrpePacket{PacketVersion:common.NRPE_PACKET_VERSION_2,PacketType:common.RESPONSE_PACKET,
        Crc32Value:0,ResultCode:common.STATE_UNKNOWN}
    if cmd == common.HELLO_COMMAND {
       copy(pkt_send.CommandBuffer[:],common.PROGRAM_VERSION)
       pkt_send.ResultCode = common.STATE_OK
    } else if IsCommandAllowed(cmd) {
        pkt_send.ResultCode = getCommand(cmd)
        copy(pkt_send.CommandBuffer[:],fillRandomData())
    } else {
        pkt_send.ResultCode = common.STATE_CRITICAL
        copy(pkt_send.CommandBuffer[:],fillRandomData())
    }
    pkt_send.Crc32Value = common.DoCRC32(pkt_send)
    return pkt_send
}
func sendPacket(conn net.Conn, pkt_send common.NrpePacket) {
    buf := new(bytes.Buffer)
    if err := binary.Write(buf, binary.BigEndian, &pkt_send); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    _, err := conn.Write([]byte(buf.Bytes()))
	common.CheckError(err)

}
func handleClient(conn net.Conn) {
	// close connection on exit
    defer conn.Close()
    pkt_rcv := receivePackets(conn)
    pkt_send := prepareToSend(string(pkt_rcv.CommandBuffer[:common.GetLen(pkt_rcv.CommandBuffer[:])]))
    sendPacket(conn,pkt_send)
}
