package dbservice

import (
	"context"
	"github.com/chwjbn/cheeringress/app/master/dbmodel"
)

func (this *DbMongoSvc) GetIngressActionBackendByTitle(ctx context.Context, title string) dbmodel.AppDataIngressActionBackend {

	xData := dbmodel.AppDataIngressActionBackend{}

	xWhere := make(map[string]interface{})
	xWhere["title"] = title

	xSort := make(map[string]interface{})

	xError := this.GetAppDataWithWhereAndOrder(ctx, &xData, xWhere, xSort)

	if xError != nil {
		xData = dbmodel.AppDataIngressActionBackend{}
	}

	return xData

}

func (this *DbMongoSvc) GetIngressActionBackendNodeByTitle(ctx context.Context, title string) dbmodel.AppDataIngressActionBackendNode {

	xData := dbmodel.AppDataIngressActionBackendNode{}

	xWhere := make(map[string]interface{})
	xWhere["title"] = title

	xSort := make(map[string]interface{})

	xError := this.GetAppDataWithWhereAndOrder(ctx, &xData, xWhere, xSort)

	if xError != nil {
		xData = dbmodel.AppDataIngressActionBackendNode{}
	}

	return xData

}

func (this *DbMongoSvc) GetFirstIngressActionBackendNodeByBackendId(ctx context.Context, backendId string) dbmodel.AppDataIngressActionBackendNode {

	xData := dbmodel.AppDataIngressActionBackendNode{}

	xWhere := make(map[string]interface{})
	xWhere["backend_id"] = backendId

	xSort := make(map[string]interface{})

	xError := this.GetAppDataWithWhereAndOrder(ctx, &xData, xWhere, xSort)

	if xError != nil {
		xData = dbmodel.AppDataIngressActionBackendNode{}
	}

	return xData
}

func (this *DbMongoSvc) UpdateIngressActionBackendNodeCount(ctx context.Context, backendId string) {

	xNodeWhere := make(map[string]interface{})
	xNodeWhere["backend_id"] = backendId

	xNodeData := dbmodel.AppDataIngressActionBackendNode{}

	xNodeCount := this.GetAppDataCount(ctx, &xNodeData, xNodeWhere)

	xData := dbmodel.AppDataIngressActionBackend{}
	xData.SetDataId(backendId)

	this.GetAppDataById(ctx, &xData)

	if len(xData.NamespaceId) < 1 {
		return
	}

	xData.NodeCount = int(xNodeCount)

	this.UpdateAppDataById(ctx, &xData)

}
