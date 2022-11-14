package controller

import (
	"github.com/chwjbn/cheeringress/app/master/dbmodel"
	"github.com/chwjbn/cheeringress/app/protocol"
	"github.com/chwjbn/cheeringress/cheerlib"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (this *WebAppCtl) CtlWorkerApiGetToken(ctx *gin.Context) {

	xRequest := protocol.WorkerGetTokenRequest{}
	xError := ctx.ShouldBindJSON(&xRequest)
	if xError != nil {
		this.ReturnWrokerApiMsg(ctx, protocol.WorkerBaseResponse{Code: "255", Message: "invalid request."})
		return
	}

	xAuthorization := ctx.GetHeader("Authorization")

	if len(xAuthorization) < 1 {
		this.ReturnWrokerApiMsg(ctx, protocol.WorkerBaseResponse{Code: "250", Message: "invalid args."})
		return
	}

	xNamespaceInfo := dbmodel.AppDataIngressNamespace{}
	xNamespaceInfo.SetDataId(xAuthorization)
	this.AppContext.AppDbSvc.GetAppDataById(ctx.Request.Context(),&xNamespaceInfo)

	if len(xNamespaceInfo.State) < 1 {
		this.ReturnWrokerApiMsg(ctx, protocol.WorkerBaseResponse{Code: "251", Message: "namespace not found."})
		return
	}

	xClientIp := this.GetClientIp(ctx)

	xData := dbmodel.AppDataIngressWorker{}
	xData.NamespaceId = xNamespaceInfo.DataId

	xData.TenantId = cheerlib.EncryptMd5("cheeradmin")
	xData.SetUpdateTime(cheerlib.TimeGetNow())
	xData.SetUpdateIp(xClientIp)
	xData.SetCreateTime(cheerlib.TimeGetNow())
	xData.SetState("enable")

	xData.NamespaceId = xNamespaceInfo.DataId
	xData.NodeName = xRequest.NodeName
	xData.NodeOs = xRequest.NodeOs
	xData.NodeArch = xRequest.NodeArch
	xData.NodeUserName = xRequest.NodeUserName

	xData.NodeCpuCore = 0
	xData.NodeCpuUsed = 0

	xData.NodeMemTotal = 0
	xData.NodeMemUsed = 0

	xData.NodeLastTime = cheerlib.TimeGetNow()
	xData.NodeAddr = xClientIp

	xData.InitDataIdWithRand(xData.NamespaceId)

	xData.NodeToken = cheerlib.EncryptMd5(fmt.Sprintf("%s_%s", xData.GetDataId(), cheerlib.EncryptNewId()))

	xError, _ = this.AppContext.AppDbSvc.AddAppData(ctx.Request.Context(),&xData)
	if xError != nil {
		this.ReturnWrokerApiMsg(ctx, protocol.WorkerBaseResponse{Code: "501", Message: fmt.Sprintf("create token with system error=[%s]", xError.Error())})
		return
	}

	xResp := protocol.WorkerGetTokenResponse{}
	xResp.Code = "0"
	xResp.Message = "sucess"
	xResp.Data = xData.NodeToken

	this.ReturnWrokerApiMsg(ctx, xResp)
}

func (this *WebAppCtl) CtlWorkerApiQueryConfig(ctx *gin.Context) {

	xRequest := protocol.WorkerQueryConfigRequest{}
	xError := ctx.ShouldBindJSON(&xRequest)
	if xError != nil {
		this.ReturnWrokerApiMsg(ctx, protocol.WorkerBaseResponse{Code: "255", Message: "invalid request."})
		return
	}

	xAuthorization := ctx.GetHeader("Authorization")
	if len(xAuthorization) < 1 {
		this.ReturnWrokerApiMsg(ctx, protocol.WorkerBaseResponse{Code: "250", Message: "invalid args."})
		return
	}

	xIngressWorkerData := this.AppContext.AppDbSvc.GetIngressWorkerByToken(ctx.Request.Context(),xAuthorization)
	if len(xIngressWorkerData.DataId) < 1 {
		this.ReturnWrokerApiMsg(ctx, protocol.WorkerBaseResponse{Code: "401", Message: "invalid token."})
		return
	}

	xClientIp := this.GetClientIp(ctx)

	xIngressWorkerData.NodeCpuCore = xRequest.NodeCpuCore
	xIngressWorkerData.NodeCpuUsed = xRequest.NodeCpuUsed
	xIngressWorkerData.NodeMemTotal = xRequest.NodeMemTotal
	xIngressWorkerData.NodeMemUsed = xRequest.NodeMemUsed

	xIngressWorkerData.NodeLastTime = cheerlib.TimeGetNow()
	xIngressWorkerData.NodeAddr = xClientIp

	this.AppContext.AppDbSvc.UpdateAppDataById(ctx.Request.Context(),&xIngressWorkerData)

	xNamespaceInfo := dbmodel.AppDataIngressNamespace{}
	xNamespaceInfo.SetDataId(xIngressWorkerData.NamespaceId)

	xError = this.AppContext.AppDbSvc.GetAppDataById(ctx.Request.Context(),&xNamespaceInfo)
	if xError != nil {
		this.ReturnWrokerApiMsg(ctx, protocol.WorkerBaseResponse{Code: "258", Message: "invalid namespace data."})
		return
	}

	xResp := protocol.WorkerQueryConfigResponse{}
	xResp.Code = "0"
	xResp.Message = "sucess"
	xResp.Data = xNamespaceInfo.LastPubVer

	this.ReturnWrokerApiMsg(ctx, xResp)

}

func (this *WebAppCtl) CtlWorkerApiFetchConfig(ctx *gin.Context) {

	xRequest := protocol.WorkerFetchConfigRequest{}
	xError := ctx.ShouldBindJSON(&xRequest)
	if xError != nil {
		this.ReturnWrokerApiMsg(ctx, protocol.WorkerBaseResponse{Code: "255", Message: "invalid request."})
		return
	}

	xAuthorization := ctx.GetHeader("Authorization")
	if len(xAuthorization) < 1 {
		this.ReturnWrokerApiMsg(ctx, protocol.WorkerBaseResponse{Code: "250", Message: "invalid args."})
		return
	}

	xIngressWorkerData := this.AppContext.AppDbSvc.GetIngressWorkerByToken(ctx.Request.Context(),xAuthorization)
	if len(xIngressWorkerData.DataId) < 1 {
		this.ReturnWrokerApiMsg(ctx, protocol.WorkerBaseResponse{Code: "401", Message: "invalid token."})
		return
	}

	xClientIp := this.GetClientIp(ctx)
	xIngressWorkerData.NodeLastTime = cheerlib.TimeGetNow()
	xIngressWorkerData.NodeAddr = xClientIp

	this.AppContext.AppDbSvc.UpdateAppDataById(ctx.Request.Context(),&xIngressWorkerData)

	xConfigData := this.AppContext.AppDbSvc.GetIngressConfigByNamespaceAndVersion(ctx.Request.Context(),xIngressWorkerData.NamespaceId, xRequest.ConfigVersion)

	xDataContent := protocol.WorkerDataNamespaceDataContent{}

	xError = cheerlib.TextStructFromJson(&xDataContent, xConfigData.Data)

	xResp := protocol.WorkerFetchConfigResponse{}
	xResp.Code = "0"
	xResp.Message = "sucess"
	xResp.Data = xDataContent

	this.ReturnWrokerApiMsg(ctx, xResp)

}

func (this *WebAppCtl) ReturnWrokerApiMsg(ctx *gin.Context, msg interface{}) {
	ctx.IndentedJSON(http.StatusOK, msg)
}
