package cheerapp

import (
	"fmt"
	"github.com/chwjbn/cheeringress/cheerlib"
	"github.com/kardianos/service"
	"os"
	"os/signal"
	"syscall"
)

type WorkCallBack func() error

type ServiceMgr struct {
	mName         string
	mDisplayName  string
	mDescription  string
	mWorkCallBack WorkCallBack
	mNotifyChan   chan os.Signal
	mService      service.Service
}

func CreateServiceMgr(name string, displayName string, description string, workCallBack WorkCallBack) (*ServiceMgr, error) {

	var xCreateErr error

	pThis := new(ServiceMgr)
	pThis.mName = name
	pThis.mDisplayName = displayName
	pThis.mDescription = description
	pThis.mWorkCallBack = workCallBack
	pThis.mNotifyChan = make(chan os.Signal)

	xSvcConfig := &service.Config{
		Name:        name,
		DisplayName: displayName,
		Description: description,
	}

	pThis.mService, xCreateErr = service.New(pThis, xSvcConfig)
	if xCreateErr != nil {
		return nil, xCreateErr
	}

	return pThis, nil

}

func (this *ServiceMgr) run() {

	if this.mWorkCallBack == nil {
		return
	}

	go func() {
		this.mWorkCallBack()
	}()

	signal.Notify(this.mNotifyChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		select {
		case xSignal := <-this.mNotifyChan:

			cheerlib.LogInfo(fmt.Sprintf("ServiceMgr Notify Exit Signal=[%v]", xSignal.String()))
			os.Exit(0)

		}
	}

}

func (this *ServiceMgr) Install() {

	cheerlib.LogInfo("ServiceMgr.Install Begin")

	this.mService.Install()

	cheerlib.LogInfo("ServiceMgr.Install End")

}

func (this *ServiceMgr) Uninstall() {

	cheerlib.LogInfo("ServiceMgr.Uninstall Begin")

	this.mService.Uninstall()

	cheerlib.LogInfo("ServiceMgr.Uninstall End")

}

func (this *ServiceMgr) Start(svc service.Service) error {

	var xErr error = nil

	cheerlib.LogInfo("ServiceMgr.Start Begin")

	go this.run()

	cheerlib.LogInfo("ServiceMgr.Start End")

	return xErr

}

func (this *ServiceMgr) Stop(svc service.Service) error {

	var xErr error = nil

	cheerlib.LogInfo("ServiceMgr.Stop Begin")

	this.mNotifyChan <- syscall.SIGQUIT

	cheerlib.LogInfo("ServiceMgr.Stop End")

	return xErr

}

func (this *ServiceMgr) RunService() {

	cheerlib.LogInfo("ServiceMgr.RunService Begin")

	this.mService.Run()

	cheerlib.LogInfo("ServiceMgr.RunService Begin")

}

func (this *ServiceMgr) StartService() {

	this.mService.Start()
}

func (this *ServiceMgr) StopService() {

	this.mService.Stop()

}
