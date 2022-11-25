package app

import (
	"fmt"
	"github.com/chwjbn/cheeringress/app/protocol"
	"github.com/chwjbn/cheeringress/app/worker/workerservice"
	"github.com/chwjbn/cheeringress/app/worker/workerutil"
	"github.com/chwjbn/cheeringress/cheerapp"
	"github.com/chwjbn/cheeringress/cheerlib"
	"github.com/chwjbn/cheeringress/config"
	v3 "github.com/chwjbn/go4sky/plugins/gin/v3"
	"github.com/gin-gonic/gin"
	"hash/crc32"
	"io/ioutil"
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

	xSpan := cheerapp.SpanBeginBizFunction(ctx.Request.Context(), "CheerWorkerApp.onHttpRequest")
	defer func() {
		cheerapp.SpanEnd(xSpan)
	}()

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

	xSpan := cheerapp.SpanBeginBizFunction(ctx.Request.Context(), "CheerWorkerApp.checkSiteAuth")
	defer func() {
		cheerapp.SpanEnd(xSpan)
	}()

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

	xSpan := cheerapp.SpanBeginBizFunction(ctx.Request.Context(), "CheerWorkerApp.processAction")
	defer func() {
		cheerapp.SpanEnd(xSpan)
	}()

	cheerapp.LogInfoWithContext(ctx.Request.Context(), "CheerWorkerApp.processAction actionType=[%s],actionData=[%v]", actionType, actionData)

	if ctx.Request == nil {
		workerutil.ActionShowErrorPage(ctx, 400, "500400", "非法请求!")
		return
	}

	if ctx.Request.URL == nil {
		workerutil.ActionShowErrorPage(ctx, 400, "500401", "非法请求!")
		return
	}

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

	xSpan := cheerapp.SpanBeginBizFunction(ctx.Request.Context(), "CheerWorkerApp.processActionStatic")
	defer func() {
		cheerapp.SpanEnd(xSpan)
	}()

	xStaticContent := staticInfo.Data

	if strings.EqualFold(staticInfo.DataType, "PlainText") {
		ctx.Data(200, staticInfo.ContentType, []byte(xStaticContent))
		return
	}

	//Base64编码内容
	if strings.EqualFold(staticInfo.DataType, "Base64Data") {
		xStaticContent = cheerlib.EncryptBase64Decode(xStaticContent)
		ctx.Data(200, staticInfo.ContentType, []byte(xStaticContent))
		return
	}

	//Http内容
	if strings.EqualFold(staticInfo.DataType, "HttpResContent") {

		xContentFilePath, xContentErr := workerutil.ActionFetchHttpResouce(ctx, staticInfo.Data)

		if xContentErr != nil {
			cheerapp.SpanError(xSpan, xContentErr.Error())
			workerutil.ActionShowErrorPage(ctx, 403, "400404", fmt.Sprintf("获取HTTP资源失败,URL=[%s],错误信息=[%s]", staticInfo.Data, xContentErr.Error()))
			return
		}

		xContentFileData, xContentErr := ioutil.ReadFile(xContentFilePath)

		if xContentErr != nil {
			cheerapp.SpanError(xSpan, xContentErr.Error())
			workerutil.ActionShowErrorPage(ctx, 403, "400405", fmt.Sprintf("获取HTTP资源失败,URL=[%s],错误信息=[%s]", staticInfo.Data, xContentErr.Error()))
			return
		}

		ctx.Data(200, staticInfo.ContentType, xContentFileData)

		return
	}

	//Http打包内容
	if strings.EqualFold(staticInfo.DataType, "HttpResZip") {

		xContentFilePath, xContentErr := workerutil.ActionFetchHttpResouce(ctx, staticInfo.Data)

		if xContentErr != nil {
			cheerapp.SpanError(xSpan, xContentErr.Error())
			workerutil.ActionShowErrorPage(ctx, 403, "400404", fmt.Sprintf("获取HTTP资源失败,URL=[%s],错误信息=[%s]", staticInfo.Data, xContentErr.Error()))
			return
		}

		xContentErr, xContentFileData := cheerlib.ZipReadStaticFile(xContentFilePath, "", ctx.Request.URL.Path)
		if xContentErr != nil {
			cheerapp.SpanError(xSpan, xContentErr.Error())
			workerutil.ActionShowErrorPage(ctx, 403, "400405", fmt.Sprintf("获取HTTP压缩资源失败,URL=[%s],错误信息=[%s]", staticInfo.Data, xContentErr.Error()))
			return
		}

		xContentType := cheerlib.WebGetContentType(ctx.Request.URL.Path)

		ctx.Data(200, xContentType, xContentFileData)

		return
	}

	if strings.EqualFold(staticInfo.DataType,"Http301Redirect"){

		if len(staticInfo.Data)<1{
			workerutil.ActionShowErrorPage(ctx, 404, "400404", fmt.Sprintf("不正确的跳转参数=[%s]", staticInfo.Data))
			return
		}

		ctx.Redirect(301,staticInfo.Data)
		return
	}

	workerutil.ActionShowErrorPage(ctx, 404, "600404", "当前请求的服务不存在.")

}

func (this *CheerWorkerApp) processActionBackend(ctx *gin.Context, backendInfo protocol.WorkerDataActionBackend, backendNodeInfoList []protocol.WorkerDataActionBackendNode) {

	xSpan := cheerapp.SpanBeginBizFunction(ctx.Request.Context(), "CheerWorkerApp.processActionBackend")

	xExitSpan:=cheerapp.SpanBeginHttpClient(ctx.Request.Context(),ctx.Request)

	defer func() {
		cheerapp.SpanEnd(xSpan)
		cheerapp.SpanEnd(xExitSpan)
	}()

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
		req.URL.Host=xBackendNodeAddr
	}

	xProxy := &httputil.ReverseProxy{Director: xDirector, ErrorHandler: this.processActionBackendErrorHandler}
	xProxy.ServeHTTP(ctx.Writer, ctx.Request)
}

func (this *CheerWorkerApp) processActionBackendErrorHandler(resp http.ResponseWriter, req *http.Request, err error) {

	xSpan := cheerapp.SpanBeginBizFunction(req.Context(), "CheerWorkerApp.processActionBackendErrorHandler")
	defer func() {
		cheerapp.SpanEnd(xSpan)
	}()

	cheerapp.SpanError(xSpan, err.Error())

	xPageContent := workerutil.GetShowErrorPageContent("500503", "后端服务节点当前无法正确响应.")

	resp.WriteHeader(http.StatusServiceUnavailable)

	resp.Write([]byte(xPageContent))
}

func (this *CheerWorkerApp) getMatchedSiteInfo(ctx *gin.Context) protocol.WorkerDataSite {

	xSpan := cheerapp.SpanBeginBizFunction(ctx.Request.Context(), "CheerWorkerApp.getMatchedSiteInfo")
	defer func() {
		cheerapp.SpanEnd(xSpan)
	}()

	xSiteInfo := protocol.WorkerDataSite{}

	xSiteHost := ctx.Request.Host

	cheerapp.LogInfoWithContext(ctx.Request.Context(), "CheerWorkerApp.getMatchedSiteInfo xSiteHost=[%s]", xSiteHost)

	xConfigSiteInfoList := this.mConfigService.GetSiteInfoList()
	if len(xConfigSiteInfoList) < 1 {
		return xSiteInfo
	}

	cheerapp.LogInfoWithContext(ctx.Request.Context(), "CheerWorkerApp.getMatchedSiteInfo xConfigSiteInfoList.Count=[%d]", len(xConfigSiteInfoList))

	for _, xConfigSiteInfo := range xConfigSiteInfoList {

		if workerutil.StringRuleMatched(xSiteHost, xConfigSiteInfo.MatchOp, xConfigSiteInfo.MatchValue) {
			xSiteInfo = xConfigSiteInfo
			break
		}
	}

	cheerapp.LogInfoWithContext(ctx.Request.Context(), "CheerWorkerApp.getMatchedSiteInfo xSiteInfo.DataId=[%s]", xSiteInfo.DataId)

	return xSiteInfo
}

func (this *CheerWorkerApp) getMatchedSiteRuleInfo(ctx *gin.Context, siteId string) protocol.WorkerDataSiteRule {

	xSpan := cheerapp.SpanBeginBizFunction(ctx.Request.Context(), "CheerWorkerApp.getMatchedSiteRuleInfo")
	defer func() {
		cheerapp.SpanEnd(xSpan)
	}()

	xSiteRuleInfo := protocol.WorkerDataSiteRule{}

	cheerapp.LogInfoWithContext(ctx.Request.Context(), "CheerWorkerApp.getMatchedSiteRuleInfo siteId=[%s]", siteId)

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

	cheerapp.LogInfoWithContext(ctx.Request.Context(), "CheerWorkerApp.xSiteRuleInfo xSiteRuleInfo.DataId=[%s]", xSiteRuleInfo.DataId)

	return xSiteRuleInfo
}
