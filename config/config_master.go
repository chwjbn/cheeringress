package config

import "errors"

type ConfigAppMaster struct {
	ServerAddr string `yaml:"server_addr"`
	ServerPort int    `yaml:"server_port"`

	DbAppMongodbUri string `yaml:"db_app_mongodb_uri"`
}

func (this *ConfigAppMaster) Check() error {

	var xError error

	if len(this.ServerAddr) < 1 {
		xError = errors.New("invalid config node=[server_addr]")
		return xError
	}

	if this.ServerPort < 1 {
		xError = errors.New("invalid config node=[server_port]")
		return xError
	}

	if len(this.DbAppMongodbUri) < 1 {
		xError = errors.New("invalid config node=[db_app_mongodb_uri]")
		return xError
	}

	return xError

}
