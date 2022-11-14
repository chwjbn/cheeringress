package dbservice

import (
	"github.com/chwjbn/cheeringress/app/master/dbmodel"
	"context"
)

func (this *DbMongoSvc) GetIngressWorkerByToken(ctx context.Context,nodeToken string) dbmodel.AppDataIngressWorker {

	xData := dbmodel.AppDataIngressWorker{}

	xWhere := make(map[string]interface{})
	xWhere["node_token"] = nodeToken

	xSort := make(map[string]interface{})

	xError := this.GetAppDataWithWhereAndOrder(ctx,&xData, xWhere, xSort)

	if xError != nil {
		xData = dbmodel.AppDataIngressWorker{}
	}

	return xData

}

func (this *DbMongoSvc) RemoveIngressWorkerByNamespaceId(ctx context.Context,namespaceId string) {

	xData := dbmodel.AppDataIngressWorker{}

	xWhere := make(map[string]interface{})
	xWhere["namespace_id"] = namespaceId

	this.DeleteAppData(ctx,&xData, xWhere)

}
