package main

import (
    "net"
    "os"
    "fmt"
    "github.com/vpereira/nrped/common"
    "github.com/droundy/goopt"
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


func main() {
    if len(os.Args) < 2 {
        fmt.Printf("%s -h for help\n",os.Args[0])
		os.Exit(1)
	}

    var host = goopt.String([]string{"-H","--host"},"127.0.0.1","The remote host running NRPE-Server")
    var port = goopt.Int([]string{"-p","--port"},5666,"The remote port on which the NRPE-server listens")
    var command = goopt.String([]string{"-c","--command"},"version","The check command defined in the nrpe.cfg file you would like to trigger")
    goopt.Parse(nil)
    service := fmt.Sprintf("%s:%d",*host,*port)
    conn := prepareConnection(service)
    pkt_to_send := common.PrepareToSend(*command,common.QUERY_PACKET)
    err := common.SendPacket(conn,pkt_to_send)
	common.CheckError(err)
    response_from_command,_ := common.ReceivePacket(conn)
    fmt.Println(response_from_command.CommandBuffer)
    os.Exit(int(response_from_command.ResultCode))
}
