package main

import (
    "os"
    "fmt"
    "net"
    "github.com/vpereira/nrped/read_config"
    "github.com/vpereira/nrped/common"
    "github.com/droundy/goopt"
)


func main() {

    if len(os.Args) < 2 {
        fmt.Fprintf(os.Stderr, "Usage: %s -h for help\n", os.Args[0])
        os.Exit(1)
    }
    config_file := goopt.String([]string{"-c","--config"},"nrpe.cfg",
    "config file to use")
    //the first option, will be the default, if the -m isnt given
    run_mode := goopt.Alternatives([]string{"-m", "--mode"},
[]string{"foreground", "daemon", "systemd"},"operating mode")
    goopt.Parse(nil)

    //implement different run modes.. 
    fmt.Println(*run_mode)
    config_obj := new(read_config.ReadConfig)
    config_obj.Init(*config_file)
    err := config_obj.ReadConfigFile();
    common.CheckError(err)
    //extract the commands command[cmd_name] = "/bin/foobar"
    config_obj.ReadCommands()

    service := ":5666"
    tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
    common.CheckError(err)

    listener, err := net.ListenTCP("tcp", tcpAddr)
    common.CheckError(err)

    for {
        if conn, err := listener.Accept(); err != nil {
            continue
        }else{
            // run as a goroutine
            go handleClient(conn,config_obj)
        }
    }
}

func handleClient(conn net.Conn,config_obj *read_config.ReadConfig) {
	// close connection on exit
    defer conn.Close()
    pkt_rcv,_ := common.ReceivePacket(conn)
    cmd := string(pkt_rcv.CommandBuffer[:common.GetLen(pkt_rcv.CommandBuffer[:])])
    pkt_send := common.PrepareToSend(cmd,common.RESPONSE_PACKET)
    if pkt_send.ResultCode == common.STATE_UNKNOWN { //its a response, but not to the HELLO_COMMAND
        if config_obj.IsCommandAllowed(cmd) {
            str_cmd := config_obj.GetCommand(cmd)
            fmt.Println("executing:",str_cmd)
            return_id,return_stdout := common.ExecuteCommand(str_cmd)
            pkt_send.ResultCode = return_id
            copy(pkt_send.CommandBuffer[:],return_stdout)
        } else {
            pkt_send.ResultCode = common.STATE_CRITICAL
        }
    }
    err := common.SendPacket(conn,pkt_send)
	common.CheckError(err)
}
