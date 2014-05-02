package main

import (
    "os"
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
    common.CheckError(err)

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

func handleClient(conn net.Conn) {
	// close connection on exit
    defer conn.Close()
    pkt_rcv := common.ReceivePacket(conn)
    cmd := string(pkt_rcv.CommandBuffer[:common.GetLen(pkt_rcv.CommandBuffer[:])])
    pkt_send := common.PrepareToSend(cmd,common.RESPONSE_PACKET)
    if pkt_send.ResultCode == common.STATE_UNKNOWN { //its a response, but not to the HELLO_COMMAND
        if IsCommandAllowed(cmd) {
            str_cmd := getCommand(cmd)
            fmt.Println("executing:",str_cmd)
            pkt_send.ResultCode = common.STATE_OK //it will be updated with the return code from the executed command
            copy(pkt_send.CommandBuffer[:],common.FillRandomData())
        } else {
            pkt_send.ResultCode = common.STATE_CRITICAL
        }
        copy(pkt_send.CommandBuffer[:],common.FillRandomData())
    }
    err := common.SendPacket(conn,pkt_send)
	common.CheckError(err)
}
