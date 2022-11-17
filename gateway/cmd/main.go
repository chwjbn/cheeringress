package main

import (
	"errors"
	"fmt"
	"github.com/chwjbn/cheeringress/app"
	"github.com/chwjbn/cheeringress/cheerapp"
	"github.com/chwjbn/cheeringress/cheerlib"
	"github.com/chwjbn/cheeringress/config"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

func RunMasterApp() error {

	var xError error

	configFilePath := path.Join(cheerlib.ApplicationBaseDirectory(), "config", "app_master.yml")
	if !cheerlib.FileExists(configFilePath) {
		xError = errors.New("Lost Config File:" + configFilePath)
		return xError
	}

	var cfg config.ConfigAppMaster

	xError = config.ReadConfigFromFile(configFilePath, &cfg)
	if xError != nil {
		return xError
	}

	config.PrintConfig(&cfg)

	xError = config.ReadConfigFromEnv(&cfg)
	if xError != nil {
		return xError
	}

	config.PrintConfig(&cfg)

	xError = cfg.Check()
	if xError != nil {
		return xError
	}

	xError = app.RunMaster(&cfg)
	if xError != nil {
		return xError
	}

	return xError
}

func RunWorkerApp() error {
	var xError error

	configFilePath := path.Join(cheerlib.ApplicationBaseDirectory(), "config", "app_worker.yml")
	if !cheerlib.FileExists(configFilePath) {
		xError = errors.New("Lost Config File:" + configFilePath)
		return xError
	}

	var cfg config.ConfigAppWorker

	xError = config.ReadConfigFromFile(configFilePath, &cfg)
	if xError != nil {
		return xError
	}

	config.PrintConfig(&cfg)

	xError = config.ReadConfigFromEnv(&cfg)
	if xError != nil {
		return xError
	}

	config.PrintConfig(&cfg)

	xError = cfg.Check()
	if xError != nil {
		return xError
	}

	xError = app.RunWorker(&cfg)
	if xError != nil {
		return xError
	}

	return xError
}

func RunMixApp() error {

	var xError error

	go func() {
		RunMasterApp()
	}()

	xError = RunWorkerApp()

	return xError

}

func AppWork() error {

	var xError error

	defer func() {

		if xError != nil {
			cheerlib.LogError(fmt.Sprintf("AppWork Error=[%s]", xError.Error()))
		}

	}()

	configFilePath := path.Join(cheerlib.ApplicationBaseDirectory(), "config", "app.yml")
	if !cheerlib.FileExists(configFilePath) {
		xError = errors.New("Lost Config File:" + configFilePath)
		return xError
	}

	var cfg config.ConfigApp

	xError = config.ReadConfigFromFile(configFilePath, &cfg)
	if xError != nil {
		return xError
	}

	config.PrintConfig(&cfg)

	xError = config.ReadConfigFromEnv(&cfg)
	if xError != nil {
		return xError
	}

	config.PrintConfig(&cfg)

	xError = cfg.Check()
	if xError != nil {
		return xError
	}

	//启用了SkyApm
	os.Setenv("SKYAPM_APP_NAME", fmt.Sprintf("arch::cheeringress_%s", cfg.AppMode))
	if len(cfg.SkyapmAppName) > 0 {
		os.Setenv("SKYAPM_APP_NAME", cfg.SkyapmAppName)
	}

	if len(cfg.SkyapmOapGrpcAddr) > 0 {
		os.Setenv("SKYAPM_OAP_GRPC_ADDR", cfg.SkyapmOapGrpcAddr)
	}

	cheerapp.StartSkyapm()

	if strings.EqualFold(cfg.AppMode, "master") {
		xError = RunMasterApp()
		return xError
	}

	if strings.EqualFold(cfg.AppMode, "worker") {
		xError = RunWorkerApp()
		return xError
	}

	if strings.EqualFold(cfg.AppMode, "mix") {
		xError = RunMixApp()
		return xError
	}

	return xError
}

func RunApp() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	xServiceMgr, xServiceMgrErr := cheerapp.CreateServiceMgr(cheerlib.GetGlobalAppName(), cheerlib.GetGlobalAppDescription(), cheerlib.GetGlobalAppDescription(), AppWork)
	if xServiceMgrErr != nil {
		cheerlib.LogError("xServiceMgrErr=" + xServiceMgrErr.Error())
	}

	xRunArgs := os.Args

	if len(xRunArgs) > 1 {
		xRunArg := xRunArgs[1]
		if strings.Contains(xRunArg, "install") {
			xServiceMgr.Install()
			return
		}

		if strings.Contains(xRunArg, "remove") {
			xServiceMgr.StopService()
			xServiceMgr.Uninstall()
			return
		}
	}

	xServiceMgr.RunService()

}

func main() {

	//设置应用信息
	cheerlib.SetGlobalAppInfo("cheeringress", "CheerIngress Service")

	// 最长日志保留时间
	cheerlib.SetGlobalCheerLogFileMaxAge(48 * time.Hour)

	cheerlib.StdInfo(fmt.Sprintf("==================================================%s Begin==================================================", cheerlib.GetGlobalAppDescription()))
	cheerlib.LogInfo(fmt.Sprintf("==================================================%s Begin==================================================", cheerlib.GetGlobalAppDescription()))

	RunApp()

	cheerlib.LogInfo(fmt.Sprintf("==================================================%s End==================================================", cheerlib.GetGlobalAppDescription()))
	cheerlib.StdInfo(fmt.Sprintf("==================================================%s End==================================================", cheerlib.GetGlobalAppDescription()))

}
