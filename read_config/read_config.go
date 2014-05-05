package read_config

import(
    "fmt"
    "strings"
    "github.com/jimlawless/cfg"
)

const (
    MAX_ALLOWED_HOSTS = 32
)

type ReadConfig struct {
    AllowedCommands map[string]string
    FileName string
    ServerPort uint16
    CommandPrefix string
    Server string
    AllowedHosts [MAX_ALLOWED_HOSTS]string
    Debug bool
    Nrpe_user string
    Nrpe_group string
    PidFile string
    ConfigMap map[string]string
    //TODO implement everything
}

//TODO
//design a constructor
func (rc *ReadConfig) Init(file_name string) {
    rc.AllowedCommands = make(map[string]string)
    rc.ConfigMap = make(map[string]string)
    rc.FileName = file_name
}

func (rc *ReadConfig) ReadConfigFile() error {
    if err := cfg.Load(rc.FileName,rc.ConfigMap); err != nil {
        return err
    }
    return nil
}
func (rc *ReadConfig) ReadCommands() {
    for key, value := range rc.ConfigMap {
        if strings.HasPrefix(key,"command[") {
            init_str := strings.Index(key,"[")
            end_str  := strings.Index(key,"]")
            fmt.Println(key[init_str+1:end_str])
            rc.AllowedCommands[key[init_str+1:end_str]] = value
        }
    }
}
