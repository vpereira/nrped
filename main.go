package main

import (
	"./common"
	"./drop_privilege"
	"./read_config"
	"fmt"
	"github.com/droundy/goopt"
	"net"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s -h for help\n", os.Args[0])
		os.Exit(1)
	}
	config_file := goopt.String([]string{"-c", "--config"}, "nrpe.cfg",
		"config file to use")
	//the first option, will be the default, if the -m isnt given
	run_mode := goopt.Alternatives([]string{"-m", "--mode"},
		[]string{"foreground", "daemon", "systemd"}, "operating mode")
	goopt.Parse(nil)

	//implement different run modes..
	fmt.Println(*run_mode)
	config_obj := new(read_config.ReadConfig)
	config_obj.Init(*config_file)
	err := config_obj.ReadConfigFile()
	common.CheckError(err)
	//extract the commands command[cmd_name] = "/bin/foobar"
	config_obj.ReadCommands()
	config_obj.ReadPrivileges()
	//TODO check for errors
	//what we gonna do with the group?
	pwd := drop_privilege.Getpwnam(config_obj.Nrpe_user)
	drop_privilege.DropPrivileges(int(pwd.Uid), int(pwd.Gid))
	//we have to read it from config
	service := ":5666"
	err = setupSocket(4, service, config_obj)
	common.CheckError(err)
}

func setupSocket(socket_version int, service string, config_obj *read_config.ReadConfig) error {

	socket_type := "tcp4"

	if socket_version == 6 {
		socket_type = "tcp6"
	}

	tcpAddr, err := net.ResolveTCPAddr(socket_type, service)
	common.CheckError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		return err
	}

	for {
		if conn, err := listener.Accept(); err != nil {
			continue
		} else {
			// run as a goroutine
			go handleClient(conn, config_obj)
		}
	}
	return nil
}

func handleClient(conn net.Conn, config_obj *read_config.ReadConfig) {
	// close connection on exit
	defer conn.Close()
	pkt_rcv, _ := common.ReceivePacket(conn)

	cmd := string(pkt_rcv.CommandBuffer[:common.GetLen(pkt_rcv.CommandBuffer[:])])
	// TODO: CRC32 does NOT match even if DoCRC32p
	if crc32, _ := common.DoCRC32p(&pkt_rcv); crc32 != pkt_rcv.CRC32Value {
		fmt.Println("WARNING: CRC not matching", crc32, pkt_rcv.CRC32Value)
	}

	pkt_send := common.PrepareToSend(cmd, common.RESPONSE_PACKET)
	if pkt_send.ResultCode == common.STATE_UNKNOWN { //its a response, but not to the HELLO_COMMAND
		if config_obj.IsCommandAllowed(cmd) {
			str_cmd := config_obj.GetCommand(cmd)
			fmt.Println("executing:", str_cmd)
			return_id, return_stdout := common.ExecuteCommand(str_cmd)
			pkt_send.ResultCode = return_id
			copy(pkt_send.CommandBuffer[:], return_stdout)
			pkt_send.CRC32Value, _ = common.DoCRC32p(&pkt_send)
		} else {
			pkt_send.ResultCode = common.STATE_CRITICAL
		}
	}
	err := common.SendPacket(conn, pkt_send)
	common.CheckError(err)
}
