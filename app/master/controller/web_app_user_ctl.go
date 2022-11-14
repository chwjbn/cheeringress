package controller

import (
	"github.com/chwjbn/cheeringress/app/master/bizmodel"
	"github.com/chwjbn/cheeringress/app/master/dbmodel"
	"github.com/chwjbn/cheeringress/cheerlib"
	"github.com/gin-gonic/gin"
	"strings"
)

func (this *WebAppCtl) CtlUserGetCurrent(ctx *gin.Context) {

	xTokenData := this.getLoginToken(ctx)
	if len(xTokenData.AccountId) < 1 {
		this.ReturnMsg(ctx, "401", "app.server.msg.common.token.invalid", nil)
		return
	}

	xUserData := dbmodel.AppDataUser{}
	xUserData.SetDataId(xTokenData.AccountId)

	xError := this.AppContext.AppDbSvc.GetAppDataById(ctx.Request.Context(),&xUserData)

	if xError != nil {
		this.ReturnMsg(ctx, "401", "app.server.msg.common.token.invalid", nil)
		return
	}

	this.ReturnAppSuccessData(ctx, xUserData)

}

func (this *WebAppCtl) CtlUserUpdateInfo(ctx *gin.Context) {

	xTokenData := this.getLoginToken(ctx)
	if len(xTokenData.AccountId) < 1 {
		this.ReturnMsg(ctx, "401", "app.server.msg.common.token.invalid", nil)
		return
	}

	xRequest := bizmodel.UserInfoUpdateRequest{}
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

	xUserInfo := dbmodel.AppDataUser{}
	xUserInfo.SetDataId(xRequest.DataId)

	this.AppContext.AppDbSvc.GetAppDataById(ctx.Request.Context(),&xUserInfo)

	if len(xUserInfo.GetState()) < 1 {
		this.ReturnAppError(ctx, "当前用户信息不存在!")
		return
	}

	xUserInfo.Nickname = xRequest.Nickname
	xUserInfo.RealName = xRequest.RealName

	xUserInfo.UpdateTime = cheerlib.TimeGetNow()
	xUserInfo.UpdateIp = this.GetClientIp(ctx)

	this.AppContext.AppDbSvc.UpdateAppDataById(ctx.Request.Context(),&xUserInfo)

	this.ReturnAppSuccess(ctx, "app.server.msg.common.op.succ")
}

func (this *WebAppCtl) CtlUserUpdateSecurity(ctx *gin.Context) {

	xTokenData := this.getLoginToken(ctx)
	if len(xTokenData.AccountId) < 1 {
		this.ReturnMsg(ctx, "401", "app.server.msg.common.token.invalid", nil)
		return
	}

	xRequest := bizmodel.UserSecurityUpdateRequest{}
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

	xUserInfo := dbmodel.AppDataUser{}
	xUserInfo.SetDataId(xRequest.DataId)

	this.AppContext.AppDbSvc.GetAppDataById(ctx.Request.Context(),&xUserInfo)

	if len(xUserInfo.GetState()) < 1 {
		this.ReturnAppError(ctx, "当前用户信息不存在!")
		return
	}

	xPwdOld := xUserInfo.GetEncryptPassword(xRequest.PasswordOld)
	if !strings.EqualFold(xPwdOld, xUserInfo.Password) {
		this.ReturnAppError(ctx, "当前登录密码不正确!")
		return
	}

	xUserInfo.Password = xUserInfo.GetEncryptPassword(xRequest.PasswordNew)

	xUserInfo.UpdateTime = cheerlib.TimeGetNow()
	xUserInfo.UpdateIp = this.GetClientIp(ctx)

	this.AppContext.AppDbSvc.UpdateAppDataById(ctx.Request.Context(),&xUserInfo)

	this.ReturnAppSuccess(ctx, "app.server.msg.common.op.succ")

}
