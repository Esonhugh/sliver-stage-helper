package sliverClient

import (
	"encoding/json"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

type ClientConfig struct {
	Operator      string `json:"operator"` // This value is actually ignored for the most part (cert CN is used instead)
	LHost         string `json:"lhost"`
	LPort         int    `json:"lport"`
	Token         string `json:"token"`
	CACertificate string `json:"ca_certificate"`
	PrivateKey    string `json:"private_key"`
	Certificate   string `json:"certificate"`
}

func ReadConfig(confFilePath string) (*ClientConfig, error) {
	confFile, err := os.Open(confFilePath)
	if err != nil {
		log.Debugf("Open failed %v", err)
		return nil, err
	}
	defer confFile.Close()
	data, err := io.ReadAll(confFile)
	if err != nil {
		log.Debugf("Read failed %v", err)
		return nil, err
	}
	conf := &ClientConfig{}
	err = json.Unmarshal(data, conf)
	if err != nil {
		log.Debugf("Parse failed %v", err)
		return nil, err
	}
	return conf, nil
}
