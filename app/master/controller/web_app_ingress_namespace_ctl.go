package controller

import (
	"github.com/chwjbn/cheeringress/app/master/bizmodel"
	"github.com/chwjbn/cheeringress/app/master/dbmodel"
	"github.com/chwjbn/cheeringress/cheerapp"
	"github.com/chwjbn/cheeringress/cheerlib"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (this *WebAppCtl) CtlIngressNamespaceMapData(ctx *gin.Context) {

	if this.checkLogin(ctx) != nil {
		return
	}

    cheerapp.LogInfoWithContext(ctx.Request.Context(),"WebAppCtl.CtlIngressNamespaceMapData")

	xWhere := make(map[string]interface{})
	xSort := make(map[string]interface{})

	xSort["update_time"] = -1

	xData := dbmodel.AppDataIngressNamespace{}

	xMapData := []bizmodel.DataMapNode{}

	xError, xDataList := this.AppContext.AppDbSvc.GetAppDataListWithWhereAndOrder(ctx.Request.Context(),&xData, xWhere, xSort, -1, -1)

	if xError == nil {

		for _, xDataItem := range xDataList {
			xDataItemVal := xDataItem.(*dbmodel.AppDataIngressNamespace)
			xMapDataNode := bizmodel.DataMapNode{DataId: xDataItemVal.DataId, DataName: fmt.Sprintf("%s", xDataItemVal.Title)}
			xMapData = append(xMapData, xMapDataNode)
		}

	}

	this.ReturnAppSuccessData(ctx, xMapData)
}

func (this *WebAppCtl) CtlIngressNamespacePageData(ctx *gin.Context) {

	if this.checkLogin(ctx) != nil {
		return
	}

	xRequest := bizmodel.IngressNamespacePageRequest{}
	xError := ctx.BindJSON(&xRequest)
	if xError != nil {
		this.ReturnAppError(ctx, "app.server.msg.common.request.invalid")
		return
	}

	xWhere := make(map[string]interface{})
	xSort := make(map[string]interface{})

	if len(xRequest.Title) > 0 {
		xWhere["title"] = bson.M{"$regex": xRequest.Title}
	}

	if len(xRequest.State) > 0 {
		xWhere["state"] = xRequest.State
	}

	xSort["update_time"] = -1

	xData := dbmodel.AppDataIngressNamespace{}

	xPageData := this.AppContext.AppDbSvc.GetDataPageList(ctx.Request.Context(),&xData, xWhere, xSort, xRequest.PageNo, xRequest.PageSize)

	this.ReturnPageData(ctx, xPageData)
}

func (this *WebAppCtl) CtlIngressNamespaceAdd(ctx *gin.Context) {

	if this.checkLogin(ctx) != nil {
		return
	}

	xRequest := bizmodel.IngressNamespaceAddRequest{}
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

	xData := this.AppContext.AppDbSvc.GetIngressNamespaceByTitle(ctx.Request.Context(),xRequest.Title)

	if len(xData.DataId) > 0 {
		this.ReturnAppError(ctx, "输入的空间名称已经存在!")
		return
	}

	xClientIp := this.GetClientIp(ctx)

	xData.TenantId = cheerlib.EncryptMd5("cheeradmin")
	xData.SetUpdateTime(cheerlib.TimeGetNow())
	xData.SetUpdateIp(xClientIp)
	xData.SetCreateTime(cheerlib.TimeGetNow())
	xData.SetState("enable")

	xData.Title = xRequest.Title

	xData.InitDataId()

	xError, _ = this.AppContext.AppDbSvc.AddAppData(ctx.Request.Context(),&xData)

	if xError != nil {
		cheerlib.LogError(xError.Error())
		this.ReturnIntenalError(ctx)
	}

	this.ReturnAppSuccess(ctx, "app.server.msg.common.op.succ")
}

func (this *WebAppCtl) CtlIngressNamespaceInfo(ctx *gin.Context) {

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

	xData := dbmodel.AppDataIngressNamespace{}
	xData.SetDataId(xRequest.DataId)

	this.AppContext.AppDbSvc.GetAppDataById(ctx.Request.Context(),&xData)

	this.ReturnAppSuccessData(ctx, xData)
}

func (this *WebAppCtl) CtlIngressNamespaceSave(ctx *gin.Context) {

	if this.checkLogin(ctx) != nil {
		return
	}

	xRequest := bizmodel.IngressNamespaceSaveRequest{}
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

	xData := dbmodel.AppDataIngressNamespace{}
	xData.SetDataId(xRequest.DataId)

	this.AppContext.AppDbSvc.GetAppDataById(ctx.Request.Context(),&xData)

	if len(xData.State) < 1 {
		this.ReturnAppError(ctx, "操作失败,提交的数据不存在!")
		return
	}

	xData.Title = xRequest.Title
	xData.State = xRequest.State

	xClientIp := this.GetClientIp(ctx)
	xData.SetUpdateTime(cheerlib.TimeGetNow())
	xData.SetUpdateIp(xClientIp)

	xError = this.AppContext.AppDbSvc.UpdateAppDataById(ctx.Request.Context(),&xData)

	if xError != nil {
		cheerlib.LogError(xError.Error())
		this.ReturnIntenalError(ctx)
	}

	this.ReturnAppSuccess(ctx, "app.server.msg.common.op.succ")
}

func (this *WebAppCtl) CtlIngressNamespaceRemove(ctx *gin.Context) {

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

	xActionStaticData := dbmodel.AppDataIngressActionStatic{}
	this.AppContext.AppDbSvc.GetFirstIngressModelByNamespaceId(ctx.Request.Context(),xRequest.DataId, &xActionStaticData)
	if len(xActionStaticData.DataId) > 0 {
		this.ReturnAppError(ctx, fmt.Sprintf("操作失败,此空间下还存在静态资源[%s]!", xActionStaticData.Title))
		return
	}

	xActionBackendData := dbmodel.AppDataIngressActionBackend{}
	this.AppContext.AppDbSvc.GetFirstIngressModelByNamespaceId(ctx.Request.Context(),xRequest.DataId, &xActionBackendData)
	if len(xActionBackendData.DataId) > 0 {
		this.ReturnAppError(ctx, fmt.Sprintf("操作失败,此空间下还存在反向代理[%s]!", xActionBackendData.Title))
		return
	}

	xSiteData := dbmodel.AppDataIngressSite{}
	this.AppContext.AppDbSvc.GetFirstIngressModelByNamespaceId(ctx.Request.Context(),xRequest.DataId, &xSiteData)
	if len(xSiteData.DataId) > 0 {
		this.ReturnAppError(ctx, fmt.Sprintf("操作失败,此空间下还存在站点[%s]!", xSiteData.Title))
		return
	}

	//清理配置
	this.AppContext.AppDbSvc.RemoveIngressConfigByNamespaceId(ctx.Request.Context(),xRequest.DataId)

	//清理工作节点
	this.AppContext.AppDbSvc.RemoveIngressWorkerByNamespaceId(ctx.Request.Context(),xRequest.DataId)

	xData := dbmodel.AppDataIngressNamespace{}
	xData.SetDataId(xRequest.DataId)

	this.AppContext.AppDbSvc.GetAppDataById(ctx.Request.Context(),&xData)

	xWhere := make(map[string]interface{})
	xWhere["data_id"] = xRequest.DataId

	this.AppContext.AppDbSvc.DeleteAppData(ctx.Request.Context(),&xData, xWhere)

	this.ReturnAppSuccess(ctx, "app.server.msg.common.op.succ")
}

func (this *WebAppCtl) CtlIngressNamespacePublish(ctx *gin.Context) {

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

	xError = this.AppContext.AppDbSvc.PublishNamespaceConfig(ctx.Request.Context(),xRequest.DataId)

	if xError != nil {
		cheerlib.LogError(xError.Error())
		this.ReturnIntenalError(ctx)
	}

	this.ReturnAppSuccess(ctx, "app.server.msg.common.op.succ")
}
