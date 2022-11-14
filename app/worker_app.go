package app

import (
	"github.com/chwjbn/cheeringress/app/protocol"
	"github.com/chwjbn/cheeringress/app/worker/workerservice"
	"github.com/chwjbn/cheeringress/app/worker/workerutil"
	"github.com/chwjbn/cheeringress/cheerlib"
	"github.com/chwjbn/cheeringress/config"
	"fmt"
	v3 "github.com/chwjbn/go4sky/plugins/gin/v3"
	"github.com/gin-gonic/gin"
	"hash/crc32"
	"net/http"
	"net/http/httputil"
	"strings"
)

type CheerWorkerApp struct {
	mConfig        *config.ConfigAppWorker
	mConfigService *workerservice.ConfigService
}

func RunWorker(cfg *config.ConfigAppWorker) error {

	var xError error = nil

	xApp := new(CheerWorkerApp)
	xApp.mConfig = cfg

	xError, xApp.mConfigService = workerservice.NewConfigService(cfg)
	if xError != nil {
		return xError
	}

	xError = xApp.runService()

	return xError

}

func (this *CheerWorkerApp) runService() error {

	var xError error = nil

	gin.SetMode(gin.DebugMode)

	xRouter := gin.Default()
	xRouter.Use(v3.Middleware(xRouter))
	xRouter.SetTrustedProxies([]string{"127.0.0.1"})
	xRouter.Any("/*path", this.onHttpRequest)

	xServerHostPort := fmt.Sprintf("%s:%d", this.mConfig.ServerAddr, this.mConfig.ServerPort)

	xError = xRouter.Run(xServerHostPort)

	return xError

}

func (this *CheerWorkerApp) onHttpRequest(ctx *gin.Context) {

	xSiteInfo := this.getMatchedSiteInfo(ctx)
	if len(xSiteInfo.DataId) < 1 {
		workerutil.ActionShowErrorPage(ctx, 406, "200406", "系统升级或者访问用户超出系统限制,服务暂时不可用.")
		return
	}

	//检查认证
	if !this.checkSiteAuth(ctx, xSiteInfo) {
		ctx.Header("WWW-Authenticate", "Basic realm=\"Application Auth\"")
		ctx.String(401, "应用访问认证失败")
		return
	}

	xSiteRuleInfo := this.getMatchedSiteRuleInfo(ctx, xSiteInfo.DataId)
	if len(xSiteRuleInfo.DataId) > 0 {
		this.processAction(ctx, xSiteRuleInfo.ActionType, xSiteRuleInfo.ActionValue)
		return
	}

	this.processAction(ctx, xSiteInfo.ActionType, xSiteInfo.ActionValue)
}

func (this *CheerWorkerApp) checkSiteAuth(ctx *gin.Context, siteInfo protocol.WorkerDataSite) bool {

	bRet := false

	if !strings.EqualFold(siteInfo.AuthNeed, "yes") {
		bRet = true
		return bRet
	}

	xUserName, xUserPwd, xHaveAuth := ctx.Request.BasicAuth()
	if !xHaveAuth {
		return bRet
	}

	if strings.EqualFold(xUserName, siteInfo.AuthUserName) && strings.EqualFold(xUserPwd, siteInfo.AuthPassword) {
		bRet = true
		return bRet
	}

	return bRet

}

func (this *CheerWorkerApp) processAction(ctx *gin.Context, actionType string, actionData string) {

	if strings.EqualFold(actionType, "backend") {

		xBackendInfo := this.mConfigService.GetActionBackendInfo(actionData)
		if len(xBackendInfo.DataId) < 1 {
			workerutil.ActionShowErrorPage(ctx, 404, "300502", "当前请求的服务不存在.")
			return
		}

		xBackendNodeInfoList := this.mConfigService.GetActionBackendNodeInfo(actionData)
		if len(xBackendNodeInfoList) < 1 {
			workerutil.ActionShowErrorPage(ctx, 404, "300503", "当前请求的服务节点不存在.")
			return
		}

		this.processActionBackend(ctx, xBackendInfo, xBackendNodeInfoList)
		return
	}

	if strings.EqualFold(actionType, "static") {

		xStaticInfo := this.mConfigService.GetActionStaticInfo(actionData)

		if len(xStaticInfo.DataId) < 1 {
			workerutil.ActionShowErrorPage(ctx, 404, "300404", "当前请求的服务不存在.")
			return
		}

		this.processActionStatic(ctx, xStaticInfo)
		return
	}

}

func (this *CheerWorkerApp) processActionStatic(ctx *gin.Context, staticInfo protocol.WorkerDataActionStatic) {

	xStaticContent := staticInfo.Data
	if strings.EqualFold(staticInfo.DataType, "Base64Data") {
		xStaticContent = cheerlib.EncryptBase64Decode(xStaticContent)
	}

	ctx.Header("Content-Type", staticInfo.ContentType)
	ctx.String(200, xStaticContent)
}

func (this *CheerWorkerApp) processActionBackend(ctx *gin.Context, backendInfo protocol.WorkerDataActionBackend, backendNodeInfoList []protocol.WorkerDataActionBackendNode) {

	xBackendNodeInfo := protocol.WorkerDataActionBackendNode{}

	if strings.EqualFold(backendInfo.BalanceType, "IPHash") {
		xTarget := cheerlib.EncryptMd5(ctx.ClientIP())
		xIndex := int(crc32.ChecksumIEEE([]byte(xTarget))) % len(backendNodeInfoList)
		xBackendNodeInfo = backendNodeInfoList[xIndex]
	}

	if strings.EqualFold(backendInfo.BalanceType, "URIHash") {
		xTarget := cheerlib.EncryptMd5(ctx.Request.RequestURI)
		xIndex := int(crc32.ChecksumIEEE([]byte(xTarget))) % len(backendNodeInfoList)
		xBackendNodeInfo = backendNodeInfoList[xIndex]
	}

	if strings.EqualFold(backendInfo.BalanceType, "RoundRobin") {

		xNodeId := this.mConfigService.GetCurrentRoundRobinBackendNodeIdWithMoveNext(backendInfo.DataId)

		for _, xNodeInfo := range backendNodeInfoList {
			if strings.EqualFold(xNodeId, xNodeInfo.DataId) {
				xBackendNodeInfo = xNodeInfo
				break
			}
		}

	}

	xBackendNodeAddr := fmt.Sprintf("%s:%d", xBackendNodeInfo.ServerHost, xBackendNodeInfo.ServerPort)

	xDirector := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = xBackendNodeAddr
		req.Host = xBackendNodeAddr
	}

	xProxy := &httputil.ReverseProxy{Director: xDirector}
	xProxy.ServeHTTP(ctx.Writer, ctx.Request)
}

func (this *CheerWorkerApp) getMatchedSiteInfo(ctx *gin.Context) protocol.WorkerDataSite {

	xSiteInfo := protocol.WorkerDataSite{}

	xSiteHost := ctx.Request.Host

	xConfigSiteInfoList := this.mConfigService.GetSiteInfoList()
	if len(xConfigSiteInfoList) < 1 {
		return xSiteInfo
	}

	for _, xConfigSiteInfo := range xConfigSiteInfoList {

		if workerutil.StringRuleMatched(xSiteHost, xConfigSiteInfo.MatchOp, xConfigSiteInfo.MatchValue) {
			xSiteInfo = xConfigSiteInfo
			break
		}
	}

	return xSiteInfo
}

func (this *CheerWorkerApp) getMatchedSiteRuleInfo(ctx *gin.Context, siteId string) protocol.WorkerDataSiteRule {

	xSiteRuleInfo := protocol.WorkerDataSiteRule{}

	xSiteRuleInfoList := this.mConfigService.GetSiteRuleInfo(siteId)
	if len(xSiteRuleInfoList) < 1 {
		return xSiteRuleInfo
	}

	xMatchContext := workerutil.GetMatchContext(ctx)

	XMatchHttpMethod := strings.ToUpper(ctx.Request.Method)

	for _, xRuleItem := range xSiteRuleInfoList {

		if !strings.EqualFold(xRuleItem.HttpMethod, "ALL") {
			if !strings.EqualFold(XMatchHttpMethod, xRuleItem.HttpMethod) {
				continue
			}
		}

		xMatchTargetVal, xMatchTargetHave := xMatchContext[xRuleItem.MatchTarget]
		if !xMatchTargetHave {
			continue
		}

		xIsRuleMatched := workerutil.StringRuleMatched(xMatchTargetVal.(string), xRuleItem.MatchOp, xRuleItem.MatchValue)
		if xIsRuleMatched {
			xSiteRuleInfo = xRuleItem
			break
		}
	}

	return xSiteRuleInfo
}
