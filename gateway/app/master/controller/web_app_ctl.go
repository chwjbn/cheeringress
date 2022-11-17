package controller

import (
	"errors"
	"github.com/chwjbn/cheeringress/app/master/bizcontext"
	"github.com/chwjbn/cheeringress/app/master/dbmodel"
	"github.com/chwjbn/cheeringress/cheerlib"
	"github.com/gin-gonic/gin"
	_ "strings"
)

type WebAppCtl struct {
	WebCtl
	AppContext *bizcontext.AppContext
}

func NewWebAppCtl(appCtx *bizcontext.AppContext) *WebAppCtl {

	xWebAppCtl := WebAppCtl{AppContext: appCtx}

	return &xWebAppCtl
}

func (this *WebAppCtl) getLoginToken(ctx *gin.Context) dbmodel.AppDataToken {

	xTokenData := dbmodel.AppDataToken{}

	xAuthorization := ctx.GetHeader("Authorization")

	if len(xAuthorization) < 1 {
		return xTokenData
	}

	xTokenData.SetDataId(xAuthorization)

	xError := this.AppContext.AppDbSvc.GetAppDataById(ctx.Request.Context(), &xTokenData)
	if xError != nil {
		xTokenData = dbmodel.AppDataToken{}
		return xTokenData
	}

	xTokenData.LastAliveTime = this.GetClientIp(ctx)
	xTokenData.LastAliveIp = cheerlib.TimeGetNow()

	this.AppContext.AppDbSvc.UpdateAppDataById(ctx.Request.Context(), &xTokenData)

	return xTokenData

}

func (this *WebAppCtl) checkLogin(ctx *gin.Context) error {

	xError := errors.New("app.server.msg.common.token.invalid")

	xTokenData := this.getLoginToken(ctx)

	if len(xTokenData.Username) < 1 {
		this.ReturnMsg(ctx, "401", "app.server.msg.common.token.invalid", nil)
		return xError
	}

	xError = nil

	return xError
}
