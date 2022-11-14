package app

import (
	"github.com/chwjbn/cheeringress/app/master/bizcontext"
	"github.com/chwjbn/cheeringress/app/master/controller"
	"github.com/chwjbn/cheeringress/cheerlib"
	"github.com/chwjbn/cheeringress/config"
	"fmt"
	v3 "github.com/chwjbn/go4sky/plugins/gin/v3"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type CheerMasterApp struct {
	mAppContext *bizcontext.AppContext
	mWebCtl     *controller.WebAppCtl
}

func RunMaster(cfg *config.ConfigAppMaster) error {

	var xError error = nil

	xApp := new(CheerMasterApp)

	//初始化应用上下文
	xApp.mAppContext = &bizcontext.AppContext{Config: cfg}

	xError = xApp.mAppContext.Init()
	if xError != nil {
		return xError
	}

	xApp.mWebCtl = controller.NewWebAppCtl(xApp.mAppContext)
	xError = xApp.RunWebService()

	return xError

}

func (this *CheerMasterApp) RunWebService() error {

	var xError error = nil

	gin.SetMode(gin.DebugMode)

	xRouter := gin.Default()
	xRouter.Use(v3.Middleware(xRouter))

	xRouter.SetTrustedProxies([]string{"127.0.0.1"})

	xRouter.Any("/*path", this.onHttpRequest)

	xServerHostPort := fmt.Sprintf("%s:%d", this.mAppContext.Config.ServerAddr, this.mAppContext.Config.ServerPort)

	xError = xRouter.Run(xServerHostPort)

	return xError
}

func (this *CheerMasterApp) onHttpRequest(ctx *gin.Context) {

	xRouterMap := make(map[string]gin.HandlerFunc)

	//账号相关路由
	xRouterMap["/xapi/account/check-code-image"] = this.mWebCtl.CtlAccountCheckCodeImage
	xRouterMap["/xapi/account/login"] = this.mWebCtl.CtlAccountLogin
	xRouterMap["/xapi/account/logout"] = this.mWebCtl.CtlAccountLogout

	//用户相关路由
	xRouterMap["/xapi/user/get-current"] = this.mWebCtl.CtlUserGetCurrent
	xRouterMap["/xapi/user/update-info"] = this.mWebCtl.CtlUserUpdateInfo
	xRouterMap["/xapi/user/update-security"] = this.mWebCtl.CtlUserUpdateSecurity

	xRouterMap["/xapi/ingress/namespace-map"] = this.mWebCtl.CtlIngressNamespaceMapData
	xRouterMap["/xapi/ingress/namespace-page"] = this.mWebCtl.CtlIngressNamespacePageData
	xRouterMap["/xapi/ingress/namespace-add"] = this.mWebCtl.CtlIngressNamespaceAdd
	xRouterMap["/xapi/ingress/namespace-info"] = this.mWebCtl.CtlIngressNamespaceInfo
	xRouterMap["/xapi/ingress/namespace-save"] = this.mWebCtl.CtlIngressNamespaceSave
	xRouterMap["/xapi/ingress/namespace-remove"] = this.mWebCtl.CtlIngressNamespaceRemove
	xRouterMap["/xapi/ingress/namespace-publish"] = this.mWebCtl.CtlIngressNamespacePublish

	xRouterMap["/xapi/ingress/action-backend-map"] = this.mWebCtl.CtlIngressActionBackendMapData
	xRouterMap["/xapi/ingress/action-backend-page"] = this.mWebCtl.CtlIngressActionBackendPageData
	xRouterMap["/xapi/ingress/action-backend-add"] = this.mWebCtl.CtlIngressActionBackendAdd
	xRouterMap["/xapi/ingress/action-backend-info"] = this.mWebCtl.CtlIngressActionBackendInfo
	xRouterMap["/xapi/ingress/action-backend-save"] = this.mWebCtl.CtlIngressActionBackendSave
	xRouterMap["/xapi/ingress/action-backend-remove"] = this.mWebCtl.CtlIngressActionBackendRemove

	xRouterMap["/xapi/ingress/action-backend-node-page"] = this.mWebCtl.CtlIngressActionBackendNodePageData
	xRouterMap["/xapi/ingress/action-backend-node-add"] = this.mWebCtl.CtlIngressActionBackendNodeAdd
	xRouterMap["/xapi/ingress/action-backend-node-remove"] = this.mWebCtl.CtlIngressActionBackendNodeRemove

	xRouterMap["/xapi/ingress/action-static-map"] = this.mWebCtl.CtlIngressActionStaticMapData
	xRouterMap["/xapi/ingress/action-static-page"] = this.mWebCtl.CtlIngressActionStaticPageData
	xRouterMap["/xapi/ingress/action-static-add"] = this.mWebCtl.CtlIngressActionStaticAdd
	xRouterMap["/xapi/ingress/action-static-info"] = this.mWebCtl.CtlIngressActionStaticInfo
	xRouterMap["/xapi/ingress/action-static-save"] = this.mWebCtl.CtlIngressActionStaticSave
	xRouterMap["/xapi/ingress/action-static-remove"] = this.mWebCtl.CtlIngressActionStaticRemove

	xRouterMap["/xapi/ingress/site-page"] = this.mWebCtl.CtlIngressSitePageData
	xRouterMap["/xapi/ingress/site-add"] = this.mWebCtl.CtlIngressSiteAdd
	xRouterMap["/xapi/ingress/site-info"] = this.mWebCtl.CtlIngressSiteInfo
	xRouterMap["/xapi/ingress/site-save"] = this.mWebCtl.CtlIngressSiteSave
	xRouterMap["/xapi/ingress/site-remove"] = this.mWebCtl.CtlIngressSiteRemove

	xRouterMap["/xapi/ingress/site-rule-page"] = this.mWebCtl.CtlIngressSiteRulePageData
	xRouterMap["/xapi/ingress/site-rule-add"] = this.mWebCtl.CtlIngressSiteRuleAdd
	xRouterMap["/xapi/ingress/site-rule-info"] = this.mWebCtl.CtlIngressSiteRuleInfo
	xRouterMap["/xapi/ingress/site-rule-remove"] = this.mWebCtl.CtlIngressSiteRuleRemove

	xRouterMap["/xapi/worker/get-token"] = this.mWebCtl.CtlWorkerApiGetToken
	xRouterMap["/xapi/worker/query-config"] = this.mWebCtl.CtlWorkerApiQueryConfig
	xRouterMap["/xapi/worker/fetch-config"] = this.mWebCtl.CtlWorkerApiFetchConfig

	//先匹配api接口路由
	xReqPath := ctx.Request.URL.Path
	bIsMatch := false

	for xCallBackKey, xCallBackFunc := range xRouterMap {

		if strings.EqualFold(xReqPath, xCallBackKey) {
			bIsMatch = true
			xCallBackFunc(ctx)
			break
		}
	}

	if bIsMatch {
		return
	}

	//读取资源文件
	xFileErr, xFileContent := cheerlib.WebReadStaticFile(xReqPath)
	if xFileErr != nil {
		ctx.String(http.StatusNotFound, "Resource Error:%s", xFileErr.Error())
		return
	}

	//获取资源文件的MIME
	xContentType := cheerlib.WebGetContentType(xReqPath)

	ctx.Header("Content-Type", xContentType)

	ctx.Writer.Write(xFileContent)
}
