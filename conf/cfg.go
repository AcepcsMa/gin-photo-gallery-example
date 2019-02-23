package conf

import (
	"encoding/json"
	"log"
	"os"
)

var ServerCfg Cfg

type Cfg struct {
	ConfigMap map[string]string
}

func init() {
	confFile, err := os.Open("conf/server.conf")
	defer confFile.Close()
	if err != nil {
		log.Fatalln(err)
	}

	ServerCfg.ConfigMap = make(map[string]string)
	err = json.NewDecoder(confFile).Decode(&ServerCfg.ConfigMap)
	if err != nil {
		log.Fatalln(err)
	}
}

func (cfg *Cfg) Get(key string) string {
	if val, ok := cfg.ConfigMap[key]; ok {
		return val
	}
	log.Fatalf("No such config term: %s!\n", key)
	return ""
}