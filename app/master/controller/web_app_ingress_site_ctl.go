package controller

import (
	"fmt"
	"github.com/chwjbn/cheeringress/app/master/bizmodel"
	"github.com/chwjbn/cheeringress/app/master/dbmodel"
	"github.com/chwjbn/cheeringress/cheerlib"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (this *WebAppCtl) CtlIngressSitePageData(ctx *gin.Context) {

	if this.checkLogin(ctx) != nil {
		return
	}

	xRequest := bizmodel.IngressSitePageRequest{}
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

	xSort["order_no"] = 1

	xData := dbmodel.AppDataIngressSite{}

	xPageData := this.AppContext.AppDbSvc.GetDataPageList(ctx.Request.Context(), &xData, xWhere, xSort, xRequest.PageNo, xRequest.PageSize)

	this.ReturnPageData(ctx, xPageData)
}

func (this *WebAppCtl) CtlIngressSiteAdd(ctx *gin.Context) {

	if this.checkLogin(ctx) != nil {
		return
	}

	xRequest := bizmodel.IngressSiteAddRequest{}
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

	xData := this.AppContext.AppDbSvc.GetIngressSiteByTitle(ctx.Request.Context(), xRequest.Title)

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
	xData.OrderNo = xRequest.OrderNo

	xData.AuthNeed = xRequest.AuthNeed
	xData.AuthUserName = xRequest.AuthUserName
	xData.AuthPassword = xRequest.AuthPassword

	xData.MatchOp = xRequest.MatchOp
	xData.MatchValue = xRequest.MatchValue

	xData.ActionType = xRequest.ActionType
	xData.ActionValue = xRequest.ActionValue

	xData.InitDataIdWithRand(xData.NamespaceId)

	xError, _ = this.AppContext.AppDbSvc.AddAppData(ctx.Request.Context(), &xData)

	if xError != nil {
		cheerlib.LogError(xError.Error())
		this.ReturnIntenalError(ctx)
	}

	this.AppContext.AppDbSvc.UpdateNamespaceLastVer(ctx.Request.Context(), xData.NamespaceId)

	this.ReturnAppSuccess(ctx, "app.server.msg.common.op.succ")
}

func (this *WebAppCtl) CtlIngressSiteInfo(ctx *gin.Context) {

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

	xData := dbmodel.AppDataIngressSite{}
	xData.SetDataId(xRequest.DataId)

	this.AppContext.AppDbSvc.GetAppDataById(ctx.Request.Context(), &xData)

	this.ReturnAppSuccessData(ctx, xData)
}

func (this *WebAppCtl) CtlIngressSiteSave(ctx *gin.Context) {

	if this.checkLogin(ctx) != nil {
		return
	}

	xRequest := bizmodel.IngressSiteSaveRequest{}
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

	xData := dbmodel.AppDataIngressSite{}
	xData.SetDataId(xRequest.DataId)

	this.AppContext.AppDbSvc.GetAppDataById(ctx.Request.Context(), &xData)

	if len(xData.State) < 1 {
		this.ReturnAppError(ctx, "操作失败,提交的数据不存在!")
		return
	}

	xData.NamespaceId = xRequest.NamespaceId
	xData.Title = xRequest.Title
	xData.OrderNo = xRequest.OrderNo

	xData.AuthNeed = xRequest.AuthNeed
	xData.AuthUserName = xRequest.AuthUserName
	xData.AuthPassword = xRequest.AuthPassword

	xData.MatchOp = xRequest.MatchOp
	xData.MatchValue = xRequest.MatchValue

	xData.ActionType = xRequest.ActionType
	xData.ActionValue = xRequest.ActionValue

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

func (this *WebAppCtl) CtlIngressSiteRemove(ctx *gin.Context) {

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

	xSiteRuleData := this.AppContext.AppDbSvc.GetFirstIngressSiteRuleBySiteId(ctx.Request.Context(), xRequest.DataId)
	if len(xSiteRuleData.DataId) > 0 {
		this.ReturnAppError(ctx, fmt.Sprintf("操作失败,此站点下还存在规则[%s]!", xSiteRuleData.Title))
		return
	}

	xData := dbmodel.AppDataIngressSite{}
	xData.SetDataId(xRequest.DataId)

	this.AppContext.AppDbSvc.GetAppDataById(ctx.Request.Context(), &xData)

	xWhere := make(map[string]interface{})
	xWhere["data_id"] = xRequest.DataId

	this.AppContext.AppDbSvc.DeleteAppData(ctx.Request.Context(), &xData, xWhere)

	this.AppContext.AppDbSvc.UpdateNamespaceLastVer(ctx.Request.Context(), xData.NamespaceId)

	this.ReturnAppSuccess(ctx, "app.server.msg.common.op.succ")
}
