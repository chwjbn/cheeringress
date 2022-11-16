package protocol

import (
	"errors"
	"fmt"
	"github.com/chwjbn/cheeringress/cheerlib"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"runtime"
	"time"
)

type Exchanger struct {
	mMasterHost  string
	mNamespaceId string
}

func NewExchanger(masterHost string, namespaceId string) *Exchanger {

	xThis := new(Exchanger)
	xThis.mMasterHost = masterHost
	xThis.mNamespaceId = namespaceId

	return xThis
}

func (this *Exchanger) GetToken() WorkerGetTokenResponse {

	xResp := WorkerGetTokenResponse{}
	xResp.Code = "500"
	xResp.Message = "Internal Error"

	xApiPath := "/xapi/worker/get-token"

	xReq := WorkerGetTokenRequest{}
	xReq.NodeName = cheerlib.OsHostName()
	xReq.NodeOs = cheerlib.OSName()
	xReq.NodeArch = runtime.GOARCH
	xReq.NodeUserName = cheerlib.OSUserName()

	xRespCall := WorkerGetTokenResponse{}

	xApiCallErr := this.callApi(this.mNamespaceId, xApiPath, xReq, &xRespCall)
	if xApiCallErr != nil {
		xResp.Code = "555"
		xResp.Message = xApiCallErr.Error()
		return xResp
	}

	if len(xRespCall.Code) < 1 {
		xResp.Code = "505"
		xResp.Message = fmt.Sprintf("call [%s%s] faild", this.mMasterHost, xApiPath)
		return xResp
	}

	xResp.Code = xRespCall.Code
	xResp.Message = xRespCall.Message
	xResp.Data = xRespCall.Data

	return xResp
}

func (this *Exchanger) QueryConfig(token string) WorkerQueryConfigResponse {

	xResp := WorkerQueryConfigResponse{}
	xResp.Code = "500"
	xResp.Message = "Internal Error"

	xApiPath := "/xapi/worker/query-config"

	xReq := WorkerQueryConfigRequest{}
	xReq.NodeCpuCore, _ = cpu.Counts(true)

	xReq.NodeCpuUsed = 0
	xCpuPercents, xDeviceErr := cpu.Percent(time.Second, false)
	if xDeviceErr == nil {
		for _, xCpuPercentItem := range xCpuPercents {
			xReq.NodeCpuUsed = xCpuPercentItem
		}
	}

	xReq.NodeMemTotal = 0
	xReq.NodeMemUsed = 0
	xMemInfo, xDeviceErr := mem.VirtualMemory()
	if xDeviceErr == nil {
		xReq.NodeMemTotal = xMemInfo.Total
		xReq.NodeMemUsed = xMemInfo.Used
	}

	xRespCall := WorkerQueryConfigResponse{}

	xApiCallErr := this.callApi(token, xApiPath, xReq, &xRespCall)
	if xApiCallErr != nil {
		xResp.Code = "555"
		xResp.Message = xApiCallErr.Error()
		return xResp
	}

	if len(xRespCall.Code) < 1 {
		xResp.Code = "505"
		xResp.Message = fmt.Sprintf("call [%s%s] faild", this.mMasterHost, xApiPath)
		return xResp
	}

	xResp.Code = xRespCall.Code
	xResp.Message = xRespCall.Message
	xResp.Data = xRespCall.Data

	return xResp
}

func (this *Exchanger) FetchConfig(token string, configVer string) WorkerFetchConfigResponse {

	xResp := WorkerFetchConfigResponse{}
	xResp.Code = "500"
	xResp.Message = "Internal Error"

	xApiPath := "/xapi/worker/fetch-config"

	xReq := WorkerFetchConfigRequest{}
	xReq.ConfigVersion = configVer

	xRespCall := WorkerFetchConfigResponse{}

	xApiCallErr := this.callApi(token, xApiPath, xReq, &xRespCall)
	if xApiCallErr != nil {
		xResp.Code = "555"
		xResp.Message = xApiCallErr.Error()
		return xResp
	}

	if len(xRespCall.Code) < 1 {
		xResp.Code = "505"
		xResp.Message = fmt.Sprintf("call [%s%s] faild", this.mMasterHost, xApiPath)
		return xResp
	}

	xResp.Code = xRespCall.Code
	xResp.Message = xRespCall.Message
	xResp.Data = xRespCall.Data

	return xResp
}

func (this *Exchanger) GetResourceData(token string, resType string, resId string) WorkerGetResourceDataResponse {

	xResp := WorkerGetResourceDataResponse{}
	xResp.Code = "500"
	xResp.Message = "Internal Error"

	xApiPath := "/xapi/worker/get-resource-data"

	xReq := WorkerGetResourceDataRequest{}
	xReq.ResType = resType
	xReq.ResId = resId

	xRespCall := WorkerGetResourceDataResponse{}

	xApiCallErr := this.callApi(token, xApiPath, xReq, &xRespCall)
	if xApiCallErr != nil {
		xResp.Code = "555"
		xResp.Message = xApiCallErr.Error()
		return xResp
	}

	if len(xRespCall.Code) < 1 {
		xResp.Code = "505"
		xResp.Message = fmt.Sprintf("call [%s%s] faild", this.mMasterHost, xApiPath)
		return xResp
	}

	xResp.Code = xRespCall.Code
	xResp.Message = xRespCall.Message
	xResp.Data = xRespCall.Data

	return xResp
}

func (this *Exchanger) callApi(authData string, apiPath string, reqData interface{}, respData interface{}) error {

	var xError error

	xReqUrl := fmt.Sprintf("%s%s", this.mMasterHost, apiPath)

	xReqJson := cheerlib.TextStructToJson(reqData)

	xRespErr, xRespData := cheerlib.NetHttpPostJson(xReqUrl, xReqJson, authData)
	if xRespErr != nil {
		xError = errors.New(fmt.Sprintf("NetHttpPostJson Return Error=[%s]", xRespErr.Error()))
		return xError
	}

	xDecodeErr := cheerlib.TextStructFromJson(respData, xRespData)
	if xDecodeErr != nil {
		xError = errors.New(fmt.Sprintf("TextStructFromJson Response Data Error=[%s]", xDecodeErr.Error()))
		return xError
	}

	return xError

}
