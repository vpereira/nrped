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
