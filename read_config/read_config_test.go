package read_config

import (
    "testing"
)

func TestReadConfigInit(t *testing.T) {
    obj := new(ReadConfig)
    obj.Init("nrpe-test.cfg")
    if obj.FileName != "nrpe-test.cfg" {
        t.Error("Init failed to set FileName")
    }
}

func TestReadConfigReadFileConfig(t *testing.T) {
    obj := new(ReadConfig)
    obj.Init("nrpe-test.cfg")
    if err := obj.ReadConfigFile(); err != nil {
        t.Error("ReadConfigFile failed to read config file")
    }
}

func TestReadConfigReadCommands(t *testing.T) {
    obj := new(ReadConfig)
    obj.Init("nrpe-test.cfg")
    obj.ReadConfigFile()
    obj.ReadCommands()
    if len(obj.AllowedCommands) == 0 {
        t.Error("ReadCommands failed to parse nrpe commands")
    }
    if _,ok := obj.AllowedCommands["check_iostat"]; ok == false {
        t.Error("ReadCommands failed to parse nrpe commands")
    }
}

func TestReadConfigIsCommandAllowed( t *testing.T) {
    obj := new(ReadConfig)
    obj.Init("nrpe-test.cfg")
    obj.ReadConfigFile()
    obj.ReadCommands()
    if obj.IsCommandAllowed("check_iostat") == false {
        t.Error("IsCommandAllowed failed with check_iostat")
    }
    if obj.IsCommandAllowed("check_foobar") == true {
        t.Error("IsCommandAllowed failed with check_foobar")
    }
}

func TestReadConfigGetCommand( t *testing.T) {
    obj := new(ReadConfig)
    obj.Init("nrpe-test.cfg")
    obj.ReadConfigFile()
    obj.ReadCommands()
    if obj.GetCommand("check_iostat") == "" {
        t.Error("GetCommand failed with check_iostat")
    }
}
