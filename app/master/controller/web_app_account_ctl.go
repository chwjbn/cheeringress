package controller

import (
	"github.com/chwjbn/cheeringress/app/master/bizmodel"
	"github.com/chwjbn/cheeringress/app/master/bizservice"
	"github.com/chwjbn/cheeringress/app/master/dbmodel"
	"github.com/chwjbn/cheeringress/cheerlib"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

func (this *WebAppCtl) CtlAccountCheckCodeImage(ctx *gin.Context) {

	xDriver := base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)
	xCaptcha := base64Captcha.NewCaptcha(xDriver, base64Captcha.DefaultMemStore)

	xCheckCodeImageId, xCheckCodeImageData, xCheckCodeImageErr := xCaptcha.Generate()

	if xCheckCodeImageErr != nil {
		this.ReturnIntenalError(ctx)
		return
	}

	xImageData := bizmodel.CheckCodeImageData{}
	xImageData.CodeId = xCheckCodeImageId
	xImageData.ImageData = xCheckCodeImageData

	this.ReturnAppSuccessData(ctx, xImageData)

}

func (this *WebAppCtl) CtlAccountLogin(ctx *gin.Context) {

	xRequest := bizmodel.AccountLoginRequest{}
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

	if !base64Captcha.DefaultMemStore.Verify(xRequest.ImageCodeId, xRequest.ImageCodeData, true) {
		this.ReturnMsg(ctx, "250", "app.server.msg.account.checkimgcode.error", nil)
		return
	}

	xAccountData := this.AppContext.AppDbSvc.GetAccountByUserName(ctx.Request.Context(),xRequest.Username)

	if len(xAccountData.DataId) < 1 {
		if this.AppContext.AppDbSvc.GetAccountCount(ctx) < 1 {
			this.AppContext.AppDbSvc.CreateAccount(ctx.Request.Context(),xRequest.Username, xRequest.Password, this.GetClientIp(ctx))
			xAccountData = this.AppContext.AppDbSvc.GetAccountByUserName(ctx.Request.Context(),xRequest.Username)
		}
	}

	if len(xAccountData.DataId) < 1 {
		this.ReturnAppError(ctx, "app.server.msg.account.account.error")
		return
	}

	xAppDataToken := dbmodel.AppDataToken{LastAliveTime: cheerlib.TimeGetNow(), LastAliveIp: this.GetClientIp(ctx)}

	xAppDataToken.TokenData = xAccountData.DataId
	xAppDataToken.AccountId = xAccountData.DataId
	xAppDataToken.TenantId = xAccountData.TenantId

	xAppDataToken.Username = xAccountData.Username
	xAppDataToken.Nickname = xAccountData.Nickname
	xAppDataToken.RealName = xAccountData.RealName
	xAppDataToken.Role = xAccountData.Role
	xAppDataToken.Avatar = "/user_avar.png"

	xAccountSvc := bizservice.AppAccountService{Context: this.AppContext}

	xAppDataToken.SetUpdateIp(this.GetClientIp(ctx))
	xAppDataToken.SetState("enable")

	xBizError, xTokenId := xAccountSvc.CreateAccountToken(ctx.Request.Context(),xAppDataToken)
	if xBizError != nil {
		cheerlib.LogError(xBizError.Error())
		this.ReturnIntenalError(ctx)
		return
	}

	xAppDataToken.SetDataId(xTokenId)
	this.ReturnMsg(ctx, "0", "app.server.msg.account.login.succ", xAppDataToken)

}

func (this *WebAppCtl) CtlAccountLogout(ctx *gin.Context) {

	xTokenData := this.getLoginToken(ctx)
	if len(xTokenData.AccountId) < 1 {
		this.ReturnMsg(ctx, "401", "app.server.msg.common.token.invalid", nil)
		return
	}

	xAccountSvc := bizservice.AppAccountService{Context: this.AppContext}
	xAccountSvc.RemoveAccountTokenById(ctx.Request.Context(),xTokenData.DataId)

	this.ReturnAppSuccess(ctx, "app.server.msg.common.action.succ")

}
