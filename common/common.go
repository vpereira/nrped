package common

import (
	"os"
	"fmt"
    "hash/crc32"
    "bytes"
    "net"
    "math/rand"
    "time"
    "encoding/binary"
)

//define states
const STATE_UNKNOWN =  3
const STATE_CRITICAL = 2
const STATE_WARNING = 1
const STATE_OK = 0

//packet type
const QUERY_PACKET = 1
const RESPONSE_PACKET = 2

//max buffer size 
const MAX_PACKETBUFFER_LENGTH = 1024

const NRPE_PACKET_VERSION_3 =  3               /* packet version identifier */
const NRPE_PACKET_VERSION_2 =  2
const NRPE_PACKET_VERSION_1 =  1

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

func ReceivePacket(conn net.Conn) NrpePacket {
    pkt_rcv := new(NrpePacket)
	err := binary.Read(conn, binary.BigEndian, pkt_rcv)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
    return *pkt_rcv
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


func DoCRC32(pkt NrpePacket) uint32 {
    buf := new(bytes.Buffer)
    if err := binary.Write(buf, binary.LittleEndian, &pkt); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    return crc32.ChecksumIEEE(buf.Bytes())
}

// count the numbers of bytes until 0 is found
func GetLen(b []byte) int {
    return bytes.Index(b, []byte{0})
}
