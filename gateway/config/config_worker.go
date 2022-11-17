package config

import "errors"

type ConfigAppWorker struct {
	ServerAddr string `yaml:"server_addr"`
	ServerPort int    `yaml:"server_port"`

	MasterHost  string `yaml:"master_host"`
	NamespaceId string `yaml:"namespace_id"`
}

func (this *ConfigAppWorker) Check() error {

	var xError error

	if len(this.ServerAddr) < 1 {
		xError = errors.New("invalid config node=[server_addr]")
		return xError
	}

	if this.ServerPort < 1 {
		xError = errors.New("invalid config node=[server_port]")
		return xError
	}

	if len(this.MasterHost) < 1 {
		xError = errors.New("invalid config node=[master_host]")
		return xError
	}

	if len(this.NamespaceId) < 1 {
		xError = errors.New("invalid config node=[namespace_id]")
		return xError
	}

	return xError

}
