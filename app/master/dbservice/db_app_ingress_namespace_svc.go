package dbservice

import (
	"github.com/chwjbn/cheeringress/app/master/dbmodel"
	"github.com/chwjbn/cheeringress/app/protocol"
	"github.com/chwjbn/cheeringress/cheerlib"
	"context"
	"errors"
	"fmt"
)

func (this *DbMongoSvc) GetIngressNamespaceByTitle(ctx context.Context,title string) dbmodel.AppDataIngressNamespace {

	xData := dbmodel.AppDataIngressNamespace{}

	xWhere := make(map[string]interface{})
	xWhere["title"] = title

	xSort := make(map[string]interface{})

	xError := this.GetAppDataWithWhereAndOrder(ctx,&xData, xWhere, xSort)

	if xError != nil {
		xData = dbmodel.AppDataIngressNamespace{}
	}

	return xData

}

func (this *DbMongoSvc) GetFirstIngressModelByNamespaceId(ctx context.Context,namespaceId string, data dbmodel.IDbModelMongo) {

	xWhere := make(map[string]interface{})
	xWhere["namespace_id"] = namespaceId

	xSort := make(map[string]interface{})

	this.GetAppDataWithWhereAndOrder(ctx,data, xWhere, xSort)
}

func (this *DbMongoSvc) UpdateNamespaceLastVer(ctx context.Context,namespaceId string) {

	xData := dbmodel.AppDataIngressNamespace{}
	xData.SetDataId(namespaceId)
	this.GetAppDataById(ctx,&xData)
	if len(xData.State) < 1 {
		return
	}

	xData.UpdateTime = cheerlib.TimeGetNow()
	xData.LastVer = cheerlib.TimeGetNow()

	if len(xData.LastPubVer) < 1 {
		xData.LastPubVer = "0001-01-01 00:00:00"
	}

	this.UpdateAppDataById(ctx,&xData)
}

func (this *DbMongoSvc) PublishNamespaceConfig(ctx context.Context,namespaceId string) error {

	var xError error

	xNamespaceInfo := dbmodel.AppDataIngressNamespace{}
	xNamespaceInfo.SetDataId(namespaceId)
	this.GetAppDataById(ctx,&xNamespaceInfo)

	if len(xNamespaceInfo.State) < 1 {
		xError = errors.New("当前发布的网关空间数据不存在!")
		return xError
	}

	xIngressConfig := this.GetIngressConfigByNamespaceAndVersion(ctx,xNamespaceInfo.DataId, xNamespaceInfo.LastVer)

	//已经发布了
	if len(xIngressConfig.State) > 0 {
		return xError
	}

	xWhere := make(map[string]interface{})
	xWhere["namespace_id"] = namespaceId
	xWhere["state"] = "enable"

	xSort := make(map[string]interface{})
	xSort["create_time"] = 1

	xSortOrderNo := make(map[string]interface{})
	xSortOrderNo["order_no"] = 1

	//反向代理
	xModelActionBackend := dbmodel.AppDataIngressActionBackend{}
	_, xDataListActionBackend := this.GetAppDataListWithWhereAndOrder(ctx,&xModelActionBackend, xWhere, xSort, -1, -1)

	//反向代理节点
	xModelActionBackendNode := dbmodel.AppDataIngressActionBackendNode{}
	_, xDataListActionBackendNode := this.GetAppDataListWithWhereAndOrder(ctx,&xModelActionBackendNode, xWhere, xSort, -1, -1)

	//静态资源
	xModelActionStatic := dbmodel.AppDataIngressActionStatic{}
	_, xDataListActionStatic := this.GetAppDataListWithWhereAndOrder(ctx,&xModelActionStatic, xWhere, xSort, -1, -1)

	//静态资源跟后端节点不允许同时为空
	if len(xDataListActionBackend) < 1 && len(xDataListActionStatic) < 1 {
		xError = errors.New("发布失败,响应动作内容为空!")
		return xError
	}

	//不允许有空后端节点
	if len(xDataListActionBackend) > 0 {
		for _, xItemData := range xDataListActionBackend {
			xActionBackend := xItemData.(*dbmodel.AppDataIngressActionBackend)
			if xActionBackend.NodeCount < 1 {
				xError = errors.New(fmt.Sprintf("发布失败,反向代理[%s]后端节点为空!", xActionBackend.Title))
				return xError
			}
		}
	}

	//站点
	xModelSite := dbmodel.AppDataIngressSite{}
	_, xDataListSite := this.GetAppDataListWithWhereAndOrder(ctx,&xModelSite, xWhere, xSortOrderNo, -1, -1)

	//站点规则
	xModelSiteRule := dbmodel.AppDataIngressSiteRule{}
	_, xDataListSiteRule := this.GetAppDataListWithWhereAndOrder(ctx,&xModelSiteRule, xWhere, xSortOrderNo, -1, -1)

	//准备发布数据
	xPubData := protocol.WorkerDataNamespaceDataContent{}

	for _, xItemData := range xDataListActionBackend {

		xSrcData := xItemData.(*dbmodel.AppDataIngressActionBackend)
		xDesData := protocol.WorkerDataActionBackend{}

		xDesData.DataId = xSrcData.DataId
		xDesData.NamespaceId = xSrcData.NamespaceId
		xDesData.BalanceType = xSrcData.BalanceType

		xPubData.ActionBackendInfo = append(xPubData.ActionBackendInfo, xDesData)
	}

	for _, xItemData := range xDataListActionBackendNode {

		xSrcData := xItemData.(*dbmodel.AppDataIngressActionBackendNode)
		xDesData := protocol.WorkerDataActionBackendNode{}

		xDesData.DataId = xSrcData.DataId
		xDesData.NamespaceId = xSrcData.NamespaceId
		xDesData.BackendId = xSrcData.BackendId
		xDesData.ServerHost = xSrcData.ServerHost
		xDesData.ServerPort = xSrcData.ServerPort
		xDesData.WeightScore = xSrcData.WeightScore

		xPubData.ActionBackendNodeInfo = append(xPubData.ActionBackendNodeInfo, xDesData)
	}

	for _, xItemData := range xDataListActionStatic {

		xSrcData := xItemData.(*dbmodel.AppDataIngressActionStatic)
		xDesData := protocol.WorkerDataActionStatic{}

		xDesData.DataId = xSrcData.DataId
		xDesData.NamespaceId = xSrcData.NamespaceId
		xDesData.ContentType = xSrcData.ContentType
		xDesData.DataType = xSrcData.DataType
		xDesData.Data = xSrcData.Data

		xPubData.ActionStaticInfo = append(xPubData.ActionStaticInfo, xDesData)
	}

	for _, xItemData := range xDataListSite {

		xSrcData := xItemData.(*dbmodel.AppDataIngressSite)
		xDesData := protocol.WorkerDataSite{}

		xDesData.DataId = xSrcData.DataId
		xDesData.NamespaceId = xSrcData.NamespaceId
		xDesData.OrderNo = xSrcData.OrderNo

		xDesData.AuthNeed = xSrcData.AuthNeed
		xDesData.AuthUserName = xSrcData.AuthUserName
		xDesData.AuthPassword = xSrcData.AuthPassword

		xDesData.MatchOp = xSrcData.MatchOp
		xDesData.MatchValue = xSrcData.MatchValue

		xDesData.ActionType = xSrcData.ActionType
		xDesData.ActionValue = xSrcData.ActionValue

		xPubData.SiteInfo = append(xPubData.SiteInfo, xDesData)
	}

	for _, xItemData := range xDataListSiteRule {

		xSrcData := xItemData.(*dbmodel.AppDataIngressSiteRule)
		xDesData := protocol.WorkerDataSiteRule{}

		xDesData.DataId = xSrcData.DataId
		xDesData.NamespaceId = xSrcData.NamespaceId
		xDesData.SiteId = xSrcData.SiteId
		xDesData.OrderNo = xSrcData.OrderNo

		xDesData.HttpMethod = xSrcData.HttpMethod
		xDesData.MatchTarget = xSrcData.MatchTarget
		xDesData.MatchOp = xSrcData.MatchOp
		xDesData.MatchValue = xSrcData.MatchValue

		xDesData.ActionType = xSrcData.ActionType
		xDesData.ActionValue = xSrcData.ActionValue

		xPubData.SiteRuleInfo = append(xPubData.SiteRuleInfo, xDesData)
	}

	xPubDataJson := cheerlib.TextStructToJson(xPubData)

	xIngressConfig.TenantId = cheerlib.EncryptMd5("cheeradmin")
	xIngressConfig.SetUpdateTime(cheerlib.TimeGetNow())
	xIngressConfig.SetUpdateIp(xNamespaceInfo.GetUpdateIp())
	xIngressConfig.SetCreateTime(cheerlib.TimeGetNow())
	xIngressConfig.SetState("enable")

	xIngressConfig.NamespaceId = xNamespaceInfo.DataId
	xIngressConfig.Version = xNamespaceInfo.LastVer
	xIngressConfig.Data = xPubDataJson

	xIngressConfig.InitDataIdWithRand(xNamespaceInfo.DataId)

	xDbErr, _ := this.AddAppData(ctx,&xIngressConfig)

	if xDbErr != nil {
		xError = errors.New(fmt.Sprintf("发布失败,服务器错误=[%s]", xDbErr.Error()))
		return xError
	}

	xNamespaceInfo.LastPubVer = xNamespaceInfo.LastVer
	this.UpdateAppDataById(ctx,&xNamespaceInfo)

	return xError

}
