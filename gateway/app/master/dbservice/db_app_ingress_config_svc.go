package dbservice

import (
	"context"
	"github.com/chwjbn/cheeringress/app/master/dbmodel"
)

func (this *DbMongoSvc) GetIngressConfigByNamespaceAndVersion(ctx context.Context, namespaceId string, version string) dbmodel.AppDataIngressConfig {

	xData := dbmodel.AppDataIngressConfig{}

	xWhere := make(map[string]interface{})
	xWhere["namespace_id"] = namespaceId
	xWhere["version"] = version

	xSort := make(map[string]interface{})

	xError := this.GetAppDataWithWhereAndOrder(ctx, &xData, xWhere, xSort)

	if xError != nil {
		xData = dbmodel.AppDataIngressConfig{}
	}

	return xData

}

func (this *DbMongoSvc) RemoveIngressConfigByNamespaceId(ctx context.Context, namespaceId string) {

	xData := dbmodel.AppDataIngressConfig{}

	xWhere := make(map[string]interface{})
	xWhere["namespace_id"] = namespaceId

	this.DeleteAppData(ctx, &xData, xWhere)

}
