package common

import (
	"os"
	"fmt"
    "io"
    "hash/crc32"
    "bytes"
    "net"
    "math/rand"
    "time"
    "encoding/binary"
    "os/exec"
    "strings"
)

//define states
const (
    STATE_OK = 0
    STATE_WARNING = 1
    STATE_CRITICAL = 2
    STATE_UNKNOWN =  3
)

//packet type
const (
    QUERY_PACKET = 1
    RESPONSE_PACKET = 2
)

//packet version
const (
    NRPE_PACKET_VERSION_1 =  1
    NRPE_PACKET_VERSION_2 =  2
    NRPE_PACKET_VERSION_3 =  3               /* packet version identifier */
)

//max buffer size 
const MAX_PACKETBUFFER_LENGTH = 1024

const HELLO_COMMAND = "version"

const PROGRAM_VERSION = "0.02"

type NrpePacket struct {
    PacketVersion int16
    PacketType int16
    Crc32Value uint32
    ResultCode int16
    CommandBuffer [MAX_PACKETBUFFER_LENGTH]byte
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

//todo return error as well
func ReceivePacket(conn net.Conn) (NrpePacket,error) {
    pkt_rcv := new(NrpePacket)
	if err := binary.Read(conn, binary.BigEndian, pkt_rcv); err != nil {
        return *pkt_rcv,err
	}
    return *pkt_rcv,nil
}

func SendPacket(conn net.Conn, pkt_send NrpePacket) error {
    buf := new(bytes.Buffer)
    if err := binary.Write(buf, binary.BigEndian, &pkt_send); err != nil {
        fmt.Println(err)
    }
    _, err := conn.Write([]byte(buf.Bytes()))
    if err != nil {
        return err
    }
    return nil
}

func PrepareToSend(cmd string, pkt_type int16) NrpePacket {
    var pkt_send NrpePacket = NrpePacket{PacketVersion:NRPE_PACKET_VERSION_2,
            Crc32Value:0,ResultCode:STATE_UNKNOWN}
    if pkt_type == RESPONSE_PACKET {  //its a response
        pkt_send.PacketType = RESPONSE_PACKET
        if cmd == HELLO_COMMAND {
           copy(pkt_send.CommandBuffer[:],PROGRAM_VERSION)
           pkt_send.ResultCode = STATE_OK
        }
    } else {  // Query Packet
        pkt_send.ResultCode = STATE_OK
        pkt_send.PacketType = QUERY_PACKET
        copy(pkt_send.CommandBuffer[:],cmd)
    }
    pkt_send.Crc32Value,_ = DoCRC32(pkt_send)
    return pkt_send
}

func FillRandomData() string {
    char := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
    rand.Seed(time.Now().UTC().UnixNano())
    buf := make([]byte, 1024)
    for i := 0; i < 1024; i++ {
        buf[i] = char[rand.Intn(len(char)-1)]
    }
    return string(buf)
}


func DoCRC32(pkt NrpePacket) (uint32, error) {
    buf := new(bytes.Buffer)
    if err := binary.Write(buf, binary.LittleEndian, &pkt); err != nil {
        return uint32(0), err
    }
    return crc32.ChecksumIEEE(buf.Bytes()),nil
}

// count the numbers of bytes until 0 is found
func GetLen(b []byte) int {
    return bytes.Index(b, []byte{0})
}

func ExecCommand(cmd_in string) (uint16, io.Writer) {
    parts := strings.Fields(cmd_in)
	head := parts[0]
	parts = parts[1:len(parts)]
    cmd := exec.Command(head,parts...)
    err := cmd.Run()
    if err != nil {
        return uint16(2),cmd.Stdout
    }
    return uint16(0),cmd.Stdout
}
