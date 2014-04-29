package main

import (
	"net"
	"os"
	"fmt"
    "flag"
    "bytes"
    "encoding/binary"
    "github.com/vpereira/nrped/common"
)

func prepareConnection(endpoint string) net.Conn {
    tcpAddr, err := net.ResolveTCPAddr("tcp4", endpoint)
	common.CheckError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	common.CheckError(err)
    if conn != nil {
        return conn
    }
    return nil
}

func prepareBufToSend(command string) *bytes.Buffer {
    var pkt_send common.NrpePacket
    pkt_send = common.NrpePacket{PacketVersion:common.NRPE_PACKET_VERSION_2,PacketType:common.QUERY_PACKET,Crc32Value:0,ResultCode:0}
    copy(pkt_send.CommandBuffer[:],command)
    pkt_send.Crc32Value = common.DoCRC32(pkt_send)

    buf := new(bytes.Buffer)
    if err := binary.Write(buf, binary.BigEndian, &pkt_send); err != nil {
        fmt.Println(err)
        return nil
    }
  return buf
}
func handleResponse(conn net.Conn) common.NrpePacket {
    var pkt_rcv common.NrpePacket
    err := binary.Read(conn, binary.BigEndian, &pkt_rcv)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
    return pkt_rcv
}
func main() {
    if len(os.Args) < 2 {
        fmt.Printf("%s -h for help\n",os.Args[0])
		os.Exit(1)
	}

    var host = flag.String("host","127.0.0.1","The remote host running NRPE-Server")
    var port = flag.Int("port",5666,"The remote port on which the NRPE-server listens")
    var command = flag.String("command","version","The check command defined in the nrpe.cfg file you would like to trigger")
    flag.Parse()

    service := fmt.Sprintf("%s:%d",*host,*port)
    buf := prepareBufToSend(*command)
    conn := prepareConnection(service)
    _, err := conn.Write([]byte(buf.Bytes()))
	common.CheckError(err)
    response_from_command := handleResponse(conn)
    fmt.Println(response_from_command.CommandBuffer)
    os.Exit(int(response_from_command.ResultCode))
}
