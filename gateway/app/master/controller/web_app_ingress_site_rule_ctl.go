package controller

import (
	"github.com/chwjbn/cheeringress/app/master/bizmodel"
	"github.com/chwjbn/cheeringress/app/master/dbmodel"
	"github.com/chwjbn/cheeringress/cheerlib"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (this *WebAppCtl) CtlIngressSiteRulePageData(ctx *gin.Context) {

	if this.checkLogin(ctx) != nil {
		return
	}

	xRequest := bizmodel.IngressSiteRulePageRequest{}
	xError := ctx.BindJSON(&xRequest)
	if xError != nil {
		this.ReturnAppError(ctx, "app.server.msg.common.request.invalid")
		return
	}

	xWhere := make(map[string]interface{})
	xSort := make(map[string]interface{})

	xWhere["site_id"] = xRequest.SiteId

	if len(xRequest.Title) > 0 {
		xWhere["title"] = bson.M{"$regex": xRequest.Title}
	}

	xSort["order_no"] = 1

	xData := dbmodel.AppDataIngressSiteRule{}

	xPageData := this.AppContext.AppDbSvc.GetDataPageList(ctx.Request.Context(), &xData, xWhere, xSort, xRequest.PageNo, xRequest.PageSize)

	this.ReturnPageData(ctx, xPageData)
}

func (this *WebAppCtl) CtlIngressSiteRuleAdd(ctx *gin.Context) {

	if this.checkLogin(ctx) != nil {
		return
	}

	xRequest := bizmodel.IngressSiteRuleAddRequest{}
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

	xData := this.AppContext.AppDbSvc.GetIngressSiteRuleByTitle(ctx.Request.Context(),xRequest.SiteId, xRequest.Title)

	if len(xData.DataId) > 0 {
		this.ReturnAppError(ctx, "输入的规则名称已经存在!")
		return
	}

	xClientIp := this.GetClientIp(ctx)

	xData.TenantId = cheerlib.EncryptMd5("cheeradmin")
	xData.SetUpdateTime(cheerlib.TimeGetNow())
	xData.SetUpdateIp(xClientIp)
	xData.SetCreateTime(cheerlib.TimeGetNow())
	xData.SetState("enable")

	xData.NamespaceId = xRequest.NamespaceId
	xData.SiteId = xRequest.SiteId
	xData.Title = xRequest.Title
	xData.OrderNo = xRequest.OrderNo

	xData.HttpMethod = xRequest.HttpMethod
	xData.MatchTarget = xRequest.MatchTarget
	xData.MatchOp = xRequest.MatchOp
	xData.MatchValue = xRequest.MatchValue

	xData.ActionType = xRequest.ActionType
	xData.ActionValue = xRequest.ActionValue

	xData.InitDataIdWithRand(xData.NamespaceId)

	xError, _ = this.AppContext.AppDbSvc.AddAppData(ctx.Request.Context(), &xData)

	if xError != nil {
		cheerlib.LogError(xError.Error())
		this.ReturnIntenalError(ctx)
		return
	}

	this.AppContext.AppDbSvc.UpdateNamespaceLastVer(ctx.Request.Context(), xData.NamespaceId)

	this.ReturnAppSuccess(ctx, "app.server.msg.common.op.succ")
}

func (this *WebAppCtl) CtlIngressSiteRuleInfo(ctx *gin.Context) {

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

	xData := dbmodel.AppDataIngressSiteRule{}
	xData.SetDataId(xRequest.DataId)

	this.AppContext.AppDbSvc.GetAppDataById(ctx.Request.Context(), &xData)

	this.ReturnAppSuccessData(ctx, xData)
}

func (this *WebAppCtl) CtlIngressSiteRuleRemove(ctx *gin.Context) {

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

	xData := dbmodel.AppDataIngressSiteRule{}
	xData.SetDataId(xRequest.DataId)

	this.AppContext.AppDbSvc.GetAppDataById(ctx.Request.Context(), &xData)

	xWhere := make(map[string]interface{})
	xWhere["data_id"] = xRequest.DataId

	this.AppContext.AppDbSvc.DeleteAppData(ctx.Request.Context(), &xData, xWhere)

	this.AppContext.AppDbSvc.UpdateNamespaceLastVer(ctx.Request.Context(), xData.NamespaceId)

	this.ReturnAppSuccess(ctx, "app.server.msg.common.op.succ")
}
