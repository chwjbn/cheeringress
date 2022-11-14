package controller

import (
	"github.com/chwjbn/cheeringress/app/master/bizmodel"
	"github.com/chwjbn/cheeringress/app/master/dbmodel"
	"github.com/chwjbn/cheeringress/cheerlib"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (this *WebAppCtl) CtlIngressActionStaticMapData(ctx *gin.Context) {

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

	xData := dbmodel.AppDataIngressActionStatic{}

	xMapData := []bizmodel.DataMapNode{}

	xError, xDataList := this.AppContext.AppDbSvc.GetAppDataListWithWhereAndOrder(ctx.Request.Context(),&xData, xWhere, xSort, -1, -1)

	if xError == nil {

		for _, xDataItem := range xDataList {
			xDataItemVal := xDataItem.(*dbmodel.AppDataIngressActionStatic)
			xMapDataNode := bizmodel.DataMapNode{DataId: xDataItemVal.DataId, DataName: fmt.Sprintf("%s", xDataItemVal.Title)}
			xMapData = append(xMapData, xMapDataNode)
		}

	}

	this.ReturnAppSuccessData(ctx, xMapData)
}

func (this *WebAppCtl) CtlIngressActionStaticPageData(ctx *gin.Context) {

	if this.checkLogin(ctx) != nil {
		return
	}

	xRequest := bizmodel.IngressActionStaticPageRequest{}
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

	xData := dbmodel.AppDataIngressActionStatic{}

	xPageData := this.AppContext.AppDbSvc.GetDataPageList(ctx.Request.Context(),&xData, xWhere, xSort, xRequest.PageNo, xRequest.PageSize)

	this.ReturnPageData(ctx, xPageData)
}

func (this *WebAppCtl) CtlIngressActionStaticAdd(ctx *gin.Context) {

	if this.checkLogin(ctx) != nil {
		return
	}

	xRequest := bizmodel.IngressActionStaticAddRequest{}
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

	xData := this.AppContext.AppDbSvc.GetIngressActionStaticByTitle(ctx.Request.Context(),xRequest.Title)

	if len(xData.DataId) > 0 {
		this.ReturnAppError(ctx, "输入的资源名称已经存在!")
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
	xData.ContentType = xRequest.ContentType
	xData.DataType = xRequest.DataType
	xData.Data = xRequest.Data

	xData.InitDataIdWithRand(xData.NamespaceId)

	xError, _ = this.AppContext.AppDbSvc.AddAppData(ctx.Request.Context(),&xData)

	if xError != nil {
		cheerlib.LogError(xError.Error())
		this.ReturnIntenalError(ctx)
	}

	this.AppContext.AppDbSvc.UpdateNamespaceLastVer(ctx.Request.Context(),xData.NamespaceId)

	this.ReturnAppSuccess(ctx, "app.server.msg.common.op.succ")
}

func (this *WebAppCtl) CtlIngressActionStaticInfo(ctx *gin.Context) {

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

	xData := dbmodel.AppDataIngressActionStatic{}
	xData.SetDataId(xRequest.DataId)

	this.AppContext.AppDbSvc.GetAppDataById(ctx.Request.Context(),&xData)

	this.ReturnAppSuccessData(ctx, xData)
}

func (this *WebAppCtl) CtlIngressActionStaticSave(ctx *gin.Context) {

	if this.checkLogin(ctx) != nil {
		return
	}

	xRequest := bizmodel.IngressActionStaticSaveRequest{}
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

	xData := dbmodel.AppDataIngressActionStatic{}
	xData.SetDataId(xRequest.DataId)

	this.AppContext.AppDbSvc.GetAppDataById(ctx.Request.Context(),&xData)

	if len(xData.State) < 1 {
		this.ReturnAppError(ctx, "操作失败,提交的数据不存在!")
		return
	}

	xData.NamespaceId = xRequest.NamespaceId
	xData.Title = xRequest.Title
	xData.ContentType = xRequest.ContentType
	xData.DataType = xRequest.DataType
	xData.Data = xRequest.Data

	xClientIp := this.GetClientIp(ctx)
	xData.SetUpdateTime(cheerlib.TimeGetNow())
	xData.SetUpdateIp(xClientIp)

	xError = this.AppContext.AppDbSvc.UpdateAppDataById(ctx.Request.Context(),&xData)

	if xError != nil {
		cheerlib.LogError(xError.Error())
		this.ReturnIntenalError(ctx)
	}

	this.AppContext.AppDbSvc.UpdateNamespaceLastVer(ctx.Request.Context(),xData.NamespaceId)

	this.ReturnAppSuccess(ctx, "app.server.msg.common.op.succ")
}

func (this *WebAppCtl) CtlIngressActionStaticRemove(ctx *gin.Context) {

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

	xSiteInfo := this.AppContext.AppDbSvc.GetFirstIngressSiteByActionValue(ctx.Request.Context(),xRequest.DataId)
	if len(xSiteInfo.DataId) > 0 {
		this.ReturnAppError(ctx, fmt.Sprintf("操作失败,此静态资源被站点[%s]引用!", xSiteInfo.Title))
		return
	}

	xSiteRuleInfo := this.AppContext.AppDbSvc.GetFirstIngressSiteRuleByActionValue(ctx.Request.Context(),xRequest.DataId)
	if len(xSiteRuleInfo.DataId) > 0 {
		this.ReturnAppError(ctx, fmt.Sprintf("操作失败,此静态资源被站点路由规则[%s]引用!", xSiteRuleInfo.Title))
		return
	}

	xData := dbmodel.AppDataIngressActionStatic{}
	xData.SetDataId(xRequest.DataId)

	this.AppContext.AppDbSvc.GetAppDataById(ctx.Request.Context(),&xData)

	xWhere := make(map[string]interface{})
	xWhere["data_id"] = xRequest.DataId

	this.AppContext.AppDbSvc.DeleteAppData(ctx.Request.Context(),&xData, xWhere)

	this.AppContext.AppDbSvc.UpdateNamespaceLastVer(ctx.Request.Context(),xData.NamespaceId)

	this.ReturnAppSuccess(ctx, "app.server.msg.common.op.succ")
}
