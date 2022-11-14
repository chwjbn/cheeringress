package config

import "errors"

type ConfigApp struct {
	AppMode string `yaml:"app_mode"`

	SkyapmAppName string `yaml:"skyapm_app_name"`
	SkyapmOapGrpcAddr string `yaml:"skyapm_oap_grpc_addr"`
}

func (this *ConfigApp) Check() error {

	var xError error

	if len(this.AppMode) < 1 {
		xError = errors.New("invalid config node=[app_mode]")
		return xError
	}

	return xError

}
