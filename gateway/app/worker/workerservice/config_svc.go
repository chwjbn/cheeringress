package workerservice

import (
	"errors"
	"fmt"
	"github.com/chwjbn/cheeringress/app/protocol"
	"github.com/chwjbn/cheeringress/app/worker/workermodel"
	"github.com/chwjbn/cheeringress/cheerlib"
	"github.com/chwjbn/cheeringress/config"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type ConfigService struct {
	mConfig    *config.ConfigAppWorker
	mCache     *workermodel.LocalCache
	mExchanger *protocol.Exchanger

	mBackendNodeMap      map[string][]string
	mBackendNodeIndexMap map[string]int
	mBackendNodeLocker   sync.RWMutex
}

func NewConfigService(config *config.ConfigAppWorker) (error, *ConfigService) {

	var xError error

	xThis := new(ConfigService)
	xThis.mConfig = config

	xThis.mBackendNodeLocker = sync.RWMutex{}
	xThis.mBackendNodeIndexMap = make(map[string]int)
	xThis.mBackendNodeMap = make(map[string][]string)

	xError, xThis.mCache = workermodel.NewLocalCache()
	if xError != nil {
		return errors.New(fmt.Sprintf("NewLocalCache Error=[%s]", xError.Error())), nil
	}

	xThis.mExchanger = protocol.NewExchanger(xThis.mConfig.MasterHost, xThis.mConfig.NamespaceId)

	go func() {
		xThis.runCheckWorkerConfigTask()
	}()

	return nil, xThis
}

func (this *ConfigService) runCheckWorkerConfigTask() {

	for {

		xCheckErr := this.checkAndLoadConfigFromMaster()
		if xCheckErr != nil {
			cheerlib.LogError(fmt.Sprintf("ConfigService.checkAndLoadConfigFromMaster Error=[%s]", xCheckErr.Error()))
			time.Sleep(10 * time.Second)
			continue
		}

		time.Sleep(20 * time.Second)

	}

}

func (this *ConfigService) getToken() string {

	var xToken string

	xTokenCacheKey := "api_worker_token"

	_, xToken = this.mCache.GetData(xTokenCacheKey)

	if len(xToken) < 1 {

		xTokenResp := this.mExchanger.GetToken()
		if xTokenResp.Code != "0" {
			cheerlib.LogError(fmt.Sprintf("Exchanger.GetToken Error=[%s]", xTokenResp.Message))
			return xToken
		}

		xToken = xTokenResp.Data

		if len(xToken) > 0 {
			this.mCache.SetData(xTokenCacheKey, xToken)
		}
	}

	return xToken

}

func (this *ConfigService) removeToken() {
	xTokenCacheKey := "api_worker_token"
	this.mCache.DelData(xTokenCacheKey)
}

func (this *ConfigService) checkAndLoadConfigFromMaster() error {

	var xError error

	xToken := this.getToken()
	if len(xToken) < 1 {
		xError = errors.New("ConfigService.checkAndLoadConfigFromMaster faild for invalid token")
		return xError
	}

	xQueryConfigResp := this.mExchanger.QueryConfig(xToken)

	if xQueryConfigResp.Code != "0" {

		if xQueryConfigResp.Code == "401" {
			this.removeToken()
		}

		xError = errors.New(fmt.Sprintf("ConfigService.checkAndLoadConfigFromMaster Exchanger.QueryConfig Error=[%s]", xQueryConfigResp.Message))
		return xError
	}

	xConfigVerCacheKey := "api_worker_config_current_ver"
	_, xConfigCurrentVer := this.mCache.GetData(xConfigVerCacheKey)
	xConfigMasterVer := xQueryConfigResp.Data

	//版本一致,不需要更新
	if strings.EqualFold(xConfigCurrentVer, xConfigMasterVer) {
		return xError
	}

	//获取对应的配置
	xFetchConfigResp := this.mExchanger.FetchConfig(xToken, xConfigMasterVer)
	if xQueryConfigResp.Code != "0" {
		xError = errors.New(fmt.Sprintf("ConfigService.checkAndLoadConfigFromMaster Exchanger.FetchConfig Error=[%s]", xQueryConfigResp.Message))
		return xError
	}

	xDataContent := xFetchConfigResp.Data

	//反向代理
	for _, xDataItem := range xDataContent.ActionBackendInfo {
		xDataCacheKey := fmt.Sprintf("worker_config_action_backend_%s", xDataItem.DataId)
		this.mCache.SetData(xDataCacheKey, cheerlib.TextStructToJson(xDataItem))
	}

	//反向代理的节点
	xActionBackendNodeMap := make(map[string][]protocol.WorkerDataActionBackendNode)
	for _, xDataItem := range xDataContent.ActionBackendNodeInfo {

		_, xIsHave := xActionBackendNodeMap[xDataItem.BackendId]
		if !xIsHave {
			xActionBackendNodeMap[xDataItem.BackendId] = []protocol.WorkerDataActionBackendNode{}
		}

		xActionBackendNodeMap[xDataItem.BackendId] = append(xActionBackendNodeMap[xDataItem.BackendId], xDataItem)
	}

	for xDataKey, xDataVal := range xActionBackendNodeMap {

		xDataCacheKey := fmt.Sprintf("worker_config_action_backend_node_%s", xDataKey)

		this.mCache.SetData(xDataCacheKey, cheerlib.TextStructToJson(xDataVal))

		//初始化加权轮询
		this.SetRoundRobinBackendNodeList(xDataKey, xDataVal)
	}

	//静态资源
	for _, xDataItem := range xDataContent.ActionStaticInfo {
		xDataCacheKey := fmt.Sprintf("worker_config_action_static_%s", xDataItem.DataId)
		this.mCache.SetData(xDataCacheKey, cheerlib.TextStructToJson(xDataItem))
	}

	//站点信息
	this.mCache.SetData("worker_config_site", cheerlib.TextStructToJson(xDataContent.SiteInfo))

	//站点规则
	xSiteRuleMap := make(map[string][]protocol.WorkerDataSiteRule)
	for _, xDataItem := range xDataContent.SiteRuleInfo {
		_, xIsHave := xSiteRuleMap[xDataItem.SiteId]
		if !xIsHave {
			xSiteRuleMap[xDataItem.SiteId] = []protocol.WorkerDataSiteRule{}
		}

		xSiteRuleMap[xDataItem.SiteId] = append(xSiteRuleMap[xDataItem.SiteId], xDataItem)
	}

	for xDataKey, xDataVal := range xSiteRuleMap {
		xDataCacheKey := fmt.Sprintf("worker_config_site_rule_%s", xDataKey)
		this.mCache.SetData(xDataCacheKey, cheerlib.TextStructToJson(xDataVal))
	}

	//设置本地规则版本
	this.mCache.SetData(xConfigVerCacheKey, xConfigMasterVer)

	return xError

}

//设置加权轮询的所有节点
func (this *ConfigService) SetRoundRobinBackendNodeList(backendId string, backendNodeList []protocol.WorkerDataActionBackendNode) {

	this.mBackendNodeLocker.Lock()
	defer func() {
		this.mBackendNodeLocker.Unlock()
	}()

	this.mBackendNodeMap[backendId] = []string{}
	this.mBackendNodeIndexMap[backendId] = 0

	xTotalScore := 0
	for _, xNodeItem := range backendNodeList {
		xTotalScore = xTotalScore + xNodeItem.WeightScore
	}

	if xTotalScore < 1 {
		return
	}

	xBackendIds := []string{}

	for _, xNodeItem := range backendNodeList {

		xCurrentCount := int(xNodeItem.WeightScore / xTotalScore)

		if xCurrentCount < 1 {
			xCurrentCount = 1
		}

		if xCurrentCount > 100 {
			xCurrentCount = 100
		}

		for i := 0; i < xCurrentCount; i++ {
			xBackendIds = append(xBackendIds, xNodeItem.DataId)
		}
	}

	//随机打乱这个数组
	rand.Seed(time.Now().UnixNano())
	xLen := len(xBackendIds)
	for i := xLen - 1; i > 0; i-- {
		randNum := rand.Intn(i)
		xBackendIds[i], xBackendIds[randNum] = xBackendIds[randNum], xBackendIds[i]
	}

	this.mBackendNodeMap[backendId] = xBackendIds
}

//移动指定反向代理后端节点的游标
func (this *ConfigService) SetRoundRobinBackendNodeMoveNext(backendId string) {

	this.mBackendNodeLocker.Lock()
	defer func() {
		this.mBackendNodeLocker.Unlock()
	}()

	xBackendIndex, xBackendIndexHave := this.mBackendNodeIndexMap[backendId]
	if !xBackendIndexHave {
		this.mBackendNodeIndexMap[backendId] = 0
		xBackendIndex = 0
		return
	}

	xBackendIndex = xBackendIndex + 1

	xBackends, xBackendsHave := this.mBackendNodeMap[backendId]
	if !xBackendsHave {
		xBackends = []string{}
	}

	if xBackendIndex < 0 {
		xBackendIndex = 0
	}

	if xBackendIndex >= len(xBackends) {
		xBackendIndex = 0
	}

	this.mBackendNodeIndexMap[backendId] = xBackendIndex
}

//获取当前加权轮询的节点ID
func (this *ConfigService) GetCurrentRoundRobinBackendNodeId(backendId string) string {

	this.mBackendNodeLocker.RLock()
	defer func() {
		this.mBackendNodeLocker.RUnlock()
	}()

	xNodeId := ""

	xBackends, xBackendsHave := this.mBackendNodeMap[backendId]
	if !xBackendsHave {
		return xNodeId
	}

	if len(xBackends) < 1 {
		return xNodeId
	}

	xBackendIndex, xBackendIndexHave := this.mBackendNodeIndexMap[backendId]
	if !xBackendIndexHave {
		xBackendIndex = 0
	}

	if xBackendIndex >= len(xBackends) {
		xBackendIndex = 0
	}

	xNodeId = xBackends[xBackendIndex]

	return xNodeId

}

func (this *ConfigService) GetCurrentRoundRobinBackendNodeIdWithMoveNext(backendId string) string {
	xNodeId := this.GetCurrentRoundRobinBackendNodeId(backendId)
	this.SetRoundRobinBackendNodeMoveNext(backendId)
	return xNodeId
}

func (this *ConfigService) GetActionBackendInfo(backendId string) protocol.WorkerDataActionBackend {

	xData := protocol.WorkerDataActionBackend{}

	xDataCacheKey := fmt.Sprintf("worker_config_action_backend_%s", backendId)
	_, xDataCacheVal := this.mCache.GetData(xDataCacheKey)

	if len(xDataCacheVal) < 1 {
		return xData
	}

	if cheerlib.TextStructFromJson(&xData, xDataCacheVal) != nil {
		xData = protocol.WorkerDataActionBackend{}
	}

	return xData

}

func (this *ConfigService) GetActionBackendNodeInfo(backendId string) []protocol.WorkerDataActionBackendNode {

	xDataList := []protocol.WorkerDataActionBackendNode{}

	xDataCacheKey := fmt.Sprintf("worker_config_action_backend_node_%s", backendId)
	_, xDataCacheVal := this.mCache.GetData(xDataCacheKey)

	if len(xDataCacheVal) < 1 {
		return xDataList
	}

	if cheerlib.TextStructFromJson(&xDataList, xDataCacheVal) != nil {
		xDataList = []protocol.WorkerDataActionBackendNode{}
	}

	return xDataList

}

func (this *ConfigService) GetActionStaticInfo(staticId string) protocol.WorkerDataActionStatic {

	xData := protocol.WorkerDataActionStatic{}

	xDataCacheKey := fmt.Sprintf("worker_config_action_static_%s", staticId)
	_, xDataCacheVal := this.mCache.GetData(xDataCacheKey)

	if len(xDataCacheVal) < 1 {
		return xData
	}

	if cheerlib.TextStructFromJson(&xData, xDataCacheVal) != nil {
		xData = protocol.WorkerDataActionStatic{}
	}

	return xData

}

func (this *ConfigService) GetSiteInfoList() []protocol.WorkerDataSite {

	xDataList := []protocol.WorkerDataSite{}

	xDataCacheKey := "worker_config_site"
	_, xDataCacheVal := this.mCache.GetData(xDataCacheKey)

	if len(xDataCacheVal) < 1 {
		return xDataList
	}

	if cheerlib.TextStructFromJson(&xDataList, xDataCacheVal) != nil {
		xDataList = []protocol.WorkerDataSite{}
	}

	return xDataList
}

func (this *ConfigService) GetSiteRuleInfo(siteId string) []protocol.WorkerDataSiteRule {

	xDataList := []protocol.WorkerDataSiteRule{}

	xDataCacheKey := fmt.Sprintf("worker_config_site_rule_%s", siteId)
	_, xDataCacheVal := this.mCache.GetData(xDataCacheKey)

	if len(xDataCacheVal) < 1 {
		return xDataList
	}

	if cheerlib.TextStructFromJson(&xDataList, xDataCacheVal) != nil {
		xDataList = []protocol.WorkerDataSiteRule{}
	}

	return xDataList

}
