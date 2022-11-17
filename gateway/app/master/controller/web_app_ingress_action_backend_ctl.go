package controller

import (
	"fmt"
	"github.com/chwjbn/cheeringress/app/master/bizmodel"
	"github.com/chwjbn/cheeringress/app/master/dbmodel"
	"github.com/chwjbn/cheeringress/cheerlib"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

func (this *WebAppCtl) CtlIngressActionBackendMapData(ctx *gin.Context) {

	if this.checkLogin(ctx) != nil {
		return
	}

	xRequest := bizmodel.AppDataIdRequest{}
	xError := ctx.BindJSON(&xRequest)
	if xError != nil {
		this.ReturnAppError(ctx, "app.server.msg.common.request.invalid")
		return
	}

	xWhere := make(map[string]interface{})
	xSort := make(map[string]interface{})

	xWhere["namespace_id"] = xRequest.DataId

	xSort["update_time"] = -1

	xData := dbmodel.AppDataIngressActionBackend{}

	xMapData := []bizmodel.DataMapNode{}

	xError, xDataList := this.AppContext.AppDbSvc.GetAppDataListWithWhereAndOrder(ctx.Request.Context(), &xData, xWhere, xSort, -1, -1)

	if xError == nil {

		for _, xDataItem := range xDataList {
			xDataItemVal := xDataItem.(*dbmodel.AppDataIngressActionBackend)
			xMapDataNode := bizmodel.DataMapNode{DataId: xDataItemVal.DataId, DataName: fmt.Sprintf("%s", xDataItemVal.Title)}
			xMapData = append(xMapData, xMapDataNode)
		}

	}

	this.ReturnAppSuccessData(ctx, xMapData)
}

func (this *WebAppCtl) CtlIngressActionBackendPageData(ctx *gin.Context) {

	if this.checkLogin(ctx) != nil {
		return
	}

	xRequest := bizmodel.IngressActionBackendPageRequest{}
	xError := ctx.BindJSON(&xRequest)
	if xError != nil {
		this.ReturnAppError(ctx, "app.server.msg.common.request.invalid")
		return
	}

	xWhere := make(map[string]interface{})
	xSort := make(map[string]interface{})

	if len(xRequest.NamespaceId) > 0 {
		xWhere["namespace_id"] = xRequest.NamespaceId
	}

	if len(xRequest.Title) > 0 {
		xWhere["title"] = bson.M{"$regex": xRequest.Title}
	}

	if len(xRequest.State) > 0 {
		xWhere["state"] = xRequest.State
	}

	xSort["update_time"] = -1

	xData := dbmodel.AppDataIngressActionBackend{}

	xPageData := this.AppContext.AppDbSvc.GetDataPageList(ctx.Request.Context(), &xData, xWhere, xSort, xRequest.PageNo, xRequest.PageSize)

	this.ReturnPageData(ctx, xPageData)
}

func (this *WebAppCtl) CtlIngressActionBackendAdd(ctx *gin.Context) {

	if this.checkLogin(ctx) != nil {
		return
	}

	xRequest := bizmodel.IngressActionBackendAddRequest{}
	xError := ctx.BindJSON(&xRequest)
	if xError != nil {
		this.ReturnAppError(ctx, "app.server.msg.common.request.invalid")
		return
	}

	xError = xRequest.Validate()
	if xError != nil {
		this.ReturnAppError(ctx, xError.Error())
		return
	}

	xData := this.AppContext.AppDbSvc.GetIngressActionBackendByTitle(ctx.Request.Context(), xRequest.Title)

	if len(xData.DataId) > 0 {
		this.ReturnAppError(ctx, "输入的代理名称已经存在!")
		return
	}

	xClientIp := this.GetClientIp(ctx)

	xData.TenantId = cheerlib.EncryptMd5("cheeradmin")
	xData.SetUpdateTime(cheerlib.TimeGetNow())
	xData.SetUpdateIp(xClientIp)
	xData.SetCreateTime(cheerlib.TimeGetNow())
	xData.SetState("enable")

	xData.NamespaceId = xRequest.NamespaceId
	xData.Title = xRequest.Title
	xData.BalanceType = xRequest.BalanceType
	xData.NodeCount = 0

	xData.InitDataIdWithRand(xData.NamespaceId)

	xError, _ = this.AppContext.AppDbSvc.AddAppData(ctx.Request.Context(), &xData)

	if xError != nil {
		cheerlib.LogError(xError.Error())
		this.ReturnIntenalError(ctx)
	}

	this.AppContext.AppDbSvc.UpdateNamespaceLastVer(ctx.Request.Context(), xData.NamespaceId)

	this.ReturnAppSuccess(ctx, "app.server.msg.common.op.succ")
}

func (this *WebAppCtl) CtlIngressActionBackendInfo(ctx *gin.Context) {

	if this.checkLogin(ctx) != nil {
		return
	}

	xRequest := bizmodel.AppDataIdRequest{}
	xError := ctx.BindJSON(&xRequest)
	if xError != nil {
		this.ReturnAppError(ctx, "app.server.msg.common.request.invalid")
		return
	}

	xError = xRequest.Validate()
	if xError != nil {
		this.ReturnAppError(ctx, xError.Error())
		return
	}

	xData := dbmodel.AppDataIngressActionBackend{}
	xData.SetDataId(xRequest.DataId)

	this.AppContext.AppDbSvc.GetAppDataById(ctx.Request.Context(), &xData)

	this.ReturnAppSuccessData(ctx, xData)
}

func (this *WebAppCtl) CtlIngressActionBackendSave(ctx *gin.Context) {

	if this.checkLogin(ctx) != nil {
		return
	}

	xRequest := bizmodel.IngressActionBackendSaveRequest{}
	xError := ctx.BindJSON(&xRequest)
	if xError != nil {
		this.ReturnAppError(ctx, "app.server.msg.common.request.invalid")
		return
	}

	xError = xRequest.Validate()
	if xError != nil {
		this.ReturnAppError(ctx, xError.Error())
		return
	}

	xData := dbmodel.AppDataIngressActionBackend{}
	xData.SetDataId(xRequest.DataId)

	this.AppContext.AppDbSvc.GetAppDataById(ctx.Request.Context(), &xData)

	if len(xData.State) < 1 {
		this.ReturnAppError(ctx, "操作失败,提交的数据不存在!")
		return
	}

	xData.NamespaceId = xRequest.NamespaceId
	xData.Title = xRequest.Title
	xData.BalanceType = xRequest.BalanceType
	xData.State = xRequest.State

	xClientIp := this.GetClientIp(ctx)
	xData.SetUpdateTime(cheerlib.TimeGetNow())
	xData.SetUpdateIp(xClientIp)

	xError = this.AppContext.AppDbSvc.UpdateAppDataById(ctx.Request.Context(), &xData)

	if xError != nil {
		cheerlib.LogError(xError.Error())
		this.ReturnIntenalError(ctx)
	}

	this.AppContext.AppDbSvc.UpdateNamespaceLastVer(ctx.Request.Context(), xData.NamespaceId)

	this.ReturnAppSuccess(ctx, "app.server.msg.common.op.succ")
}

func (this *WebAppCtl) CtlIngressActionBackendRemove(ctx *gin.Context) {

	if this.checkLogin(ctx) != nil {
		return
	}

	xRequest := bizmodel.AppDataIdRequest{}
	xError := ctx.BindJSON(&xRequest)
	if xError != nil {
		this.ReturnAppError(ctx, "app.server.msg.common.request.invalid")
		return
	}

	xError = xRequest.Validate()
	if xError != nil {
		this.ReturnAppError(ctx, xError.Error())
		return
	}

	xSiteInfo := this.AppContext.AppDbSvc.GetFirstIngressSiteByActionValue(ctx.Request.Context(), xRequest.DataId)
	if len(xSiteInfo.DataId) > 0 {
		this.ReturnAppError(ctx, fmt.Sprintf("操作失败,此反向代理被站点[%s]引用!", xSiteInfo.Title))
		return
	}

	xSiteRuleInfo := this.AppContext.AppDbSvc.GetFirstIngressSiteRuleByActionValue(ctx.Request.Context(), xRequest.DataId)
	if len(xSiteRuleInfo.DataId) > 0 {
		this.ReturnAppError(ctx, fmt.Sprintf("操作失败,此反向代理被站点路由规则[%s]引用!", xSiteRuleInfo.Title))
		return
	}

	xBackendNodeInfo := this.AppContext.AppDbSvc.GetFirstIngressActionBackendNodeByBackendId(ctx.Request.Context(), xRequest.DataId)
	if len(xBackendNodeInfo.DataId) > 0 {
		this.ReturnAppError(ctx, fmt.Sprintf("操作失败,此反向代理下包含节点[%s]!", xBackendNodeInfo.Title))
		return
	}

	xData := dbmodel.AppDataIngressActionBackend{}
	xData.SetDataId(xRequest.DataId)

	this.AppContext.AppDbSvc.GetAppDataById(ctx.Request.Context(), &xData)

	xWhere := make(map[string]interface{})
	xWhere["data_id"] = xRequest.DataId

	this.AppContext.AppDbSvc.DeleteAppData(ctx.Request.Context(), &xData, xWhere)

	this.AppContext.AppDbSvc.UpdateNamespaceLastVer(ctx.Request.Context(), xData.NamespaceId)

	this.ReturnAppSuccess(ctx, "app.server.msg.common.op.succ")
}

func (this *WebAppCtl) CtlIngressActionBackendNodePageData(ctx *gin.Context) {

	if this.checkLogin(ctx) != nil {
		return
	}

	xRequest := bizmodel.IngressActionBackendNodePageRequest{}
	xError := ctx.BindJSON(&xRequest)
	if xError != nil {
		this.ReturnAppError(ctx, "app.server.msg.common.request.invalid")
		return
	}

	xRequest.PageSize = 10000
	xRequest.PageNo = 1

	xWhere := make(map[string]interface{})
	xSort := make(map[string]interface{})

	xWhere["backend_id"] = xRequest.BackendId

	xSort["update_time"] = -1

	xData := dbmodel.AppDataIngressActionBackendNode{}

	xPageData := this.AppContext.AppDbSvc.GetDataPageList(ctx.Request.Context(), &xData, xWhere, xSort, xRequest.PageNo, xRequest.PageSize)

	this.ReturnPageData(ctx, xPageData)
}

func (this *WebAppCtl) CtlIngressActionBackendNodeAdd(ctx *gin.Context) {

	if this.checkLogin(ctx) != nil {
		return
	}

	xRequest := bizmodel.IngressActionBackendNodeAddRequest{}
	xError := ctx.BindJSON(&xRequest)
	if xError != nil {
		this.ReturnAppError(ctx, "app.server.msg.common.request.invalid")
		return
	}

	xError = xRequest.Validate()
	if xError != nil {
		this.ReturnAppError(ctx, xError.Error())
		return
	}

	xData := this.AppContext.AppDbSvc.GetIngressActionBackendNodeByTitle(ctx.Request.Context(), xRequest.Title)

	if len(xData.DataId) > 0 {
		this.ReturnAppError(ctx, "输入的节点名称已经存在!")
		return
	}

	xClientIp := this.GetClientIp(ctx)

	xData.TenantId = cheerlib.EncryptMd5("cheeradmin")
	xData.SetUpdateTime(cheerlib.TimeGetNow())
	xData.SetUpdateIp(xClientIp)
	xData.SetCreateTime(cheerlib.TimeGetNow())
	xData.SetState("enable")

	xData.NamespaceId = xRequest.NamespaceId
	xData.BackendId = xRequest.BackendId
	xData.Title = xRequest.Title
	xData.ServerHost = strings.TrimSpace(xRequest.ServerHost)
	xData.ServerPort = xRequest.ServerPort
	xData.WeightScore = xRequest.WeightScore

	xData.InitDataIdWithRand(xData.BackendId)

	xError, _ = this.AppContext.AppDbSvc.AddAppData(ctx.Request.Context(), &xData)

	if xError != nil {
		cheerlib.LogError(xError.Error())
		this.ReturnIntenalError(ctx)
	}

	this.AppContext.AppDbSvc.UpdateIngressActionBackendNodeCount(ctx.Request.Context(), xData.BackendId)

	this.AppContext.AppDbSvc.UpdateNamespaceLastVer(ctx.Request.Context(), xData.NamespaceId)

	this.ReturnAppSuccess(ctx, "app.server.msg.common.op.succ")
}

func (this *WebAppCtl) CtlIngressActionBackendNodeRemove(ctx *gin.Context) {

	if this.checkLogin(ctx) != nil {
		return
	}

	xRequest := bizmodel.AppDataIdRequest{}
	xError := ctx.BindJSON(&xRequest)
	if xError != nil {
		this.ReturnAppError(ctx, "app.server.msg.common.request.invalid")
		return
	}

	xError = xRequest.Validate()
	if xError != nil {
		this.ReturnAppError(ctx, xError.Error())
		return
	}

	xData := dbmodel.AppDataIngressActionBackendNode{}
	xData.SetDataId(xRequest.DataId)

	this.AppContext.AppDbSvc.GetAppDataById(ctx.Request.Context(), &xData)

	xWhere := make(map[string]interface{})
	xWhere["data_id"] = xRequest.DataId

	this.AppContext.AppDbSvc.DeleteAppData(ctx.Request.Context(), &xData, xWhere)

	if len(xData.BackendId) > 0 {
		this.AppContext.AppDbSvc.UpdateIngressActionBackendNodeCount(ctx.Request.Context(), xData.BackendId)
	}

	this.AppContext.AppDbSvc.UpdateNamespaceLastVer(ctx.Request.Context(), xData.NamespaceId)

	this.ReturnAppSuccess(ctx, "app.server.msg.common.op.succ")
}
