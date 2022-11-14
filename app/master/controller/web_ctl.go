package controller

import (
	"github.com/chwjbn/cheeringress/app/master/dbmodel"
	"github.com/chwjbn/cheeringress/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WebCtl struct {
	mConfig *config.ConfigApp
}

func (this *WebCtl) Init(cfg *config.ConfigApp) error {

	var xError error

	this.mConfig = cfg

	return xError
}

func (this *WebCtl) GetClientIp(ctx *gin.Context) string {

	xData := ctx.GetHeader("X-Cheer-Client-IP")
	if len(xData) > 0 {
		return xData
	}

	xData = ctx.GetHeader("x-cheer-client-ip")
	if len(xData) > 0 {
		return xData
	}

	xData = ctx.GetHeader("X-Real-IP")
	if len(xData) > 0 {
		return xData
	}

	xData = ctx.GetHeader("x-real-ip")
	if len(xData) > 0 {
		return xData
	}

	xData = ctx.ClientIP()

	return xData

}

func (this *WebCtl) ReturnPageData(ctx *gin.Context, pageData dbmodel.PageData) {

	xRetData := map[string]interface{}{
		"total":   pageData.TotalCount,
		"success": true,
		"data":    pageData.DataList,
	}

	ctx.IndentedJSON(http.StatusOK, xRetData)
}

func (this *WebCtl) ReturnMsg(ctx *gin.Context, errorCode string, errorMessage string, data interface{}) {

	xRetData := map[string]interface{}{
		"error_code":    errorCode,
		"error_message": errorMessage,
		"data":          data,
	}

	ctx.IndentedJSON(http.StatusOK, xRetData)
}

func (this *WebCtl) ReturnIntenalError(ctx *gin.Context) {
	this.ReturnMsg(ctx, "500", "app.server.msg.system.upgrade", nil)
}

func (this *WebCtl) ReturnAppError(ctx *gin.Context, errorMessage string) {
	this.ReturnMsg(ctx, "250", errorMessage, nil)
}

func (this *WebCtl) ReturnAppSuccess(ctx *gin.Context, errorMessage string) {
	this.ReturnMsg(ctx, "0", errorMessage, nil)
}

func (this *WebCtl) ReturnAppSuccessData(ctx *gin.Context, data interface{}) {
	this.ReturnMsg(ctx, "0", "app.server.msg.common.request.succ", data)
}
