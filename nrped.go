package main

import (
    "os"
    "log"
    "fmt"
    "net"
    "strings"
    "github.com/vpereira/nrped/common"
    "github.com/jimlawless/cfg"
)

//it will hold the allowed commands specificed in the config file
var allowedCommands map[string]string

func read_commands(config map[string] string) map[string] string {
    allowedCommands = make(map[string]string)
    for key, value := range config {
        if strings.HasPrefix(key,"command[") {
            init_str := strings.Index(key,"[")
            end_str  := strings.Index(key,"]")
            fmt.Println(key[init_str+1:end_str])
            allowedCommands[key[init_str+1:end_str]] = value
        }
    }
    return allowedCommands
}
func main() {

    if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "Usage: %s config-file\n", os.Args[0])
        os.Exit(1)
    }

    config_file := os.Args[1]
	mymap := make(map[string]string)
    err := cfg.Load(config_file, mymap)

    if err != nil {
        log.Fatal(err)
    }

    //extract the commands command[cmd_name] = "/bin/foobar"
    read_commands(mymap)

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

func IsCommandAllowed(cmd string) bool {
    if _,ok := allowedCommands[cmd]; ok {
        return true
    }else{
        return false
    }
}

func getCommand(cmd string) string {
    return allowedCommands[cmd]
}

func prepareToSend(cmd string) common.NrpePacket {
    pkt_send := common.NrpePacket{PacketVersion:common.NRPE_PACKET_VERSION_2,PacketType:common.RESPONSE_PACKET,
        Crc32Value:0,ResultCode:common.STATE_UNKNOWN}
    if cmd == common.HELLO_COMMAND {
       copy(pkt_send.CommandBuffer[:],common.PROGRAM_VERSION)
       pkt_send.ResultCode = common.STATE_OK
    } else if IsCommandAllowed(cmd) {
        str_cmd := getCommand(cmd) //TODO it should be executed, and the result should be added to ResultCode
        fmt.Println("executing:",str_cmd)
        pkt_send.ResultCode = 0
        copy(pkt_send.CommandBuffer[:],common.FillRandomData())
    } else {
        pkt_send.ResultCode = common.STATE_CRITICAL
        copy(pkt_send.CommandBuffer[:],common.FillRandomData())
    }
    pkt_send.Crc32Value = common.DoCRC32(pkt_send)
    return pkt_send
}

func handleClient(conn net.Conn) {
	// close connection on exit
    defer conn.Close()
    pkt_rcv := common.ReceivePacket(conn)
    pkt_send := prepareToSend(string(pkt_rcv.CommandBuffer[:common.GetLen(pkt_rcv.CommandBuffer[:])]))
    err := common.SendPacket(conn,pkt_send)
	common.CheckError(err)
}
