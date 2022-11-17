package dbservice

import (
	"context"
	"github.com/chwjbn/cheeringress/app/master/dbmodel"
)

func (this *DbMongoSvc) GetIngressActionStaticByTitle(ctx context.Context, title string) dbmodel.AppDataIngressActionStatic {

	xData := dbmodel.AppDataIngressActionStatic{}

	xWhere := make(map[string]interface{})
	xWhere["title"] = title

	xSort := make(map[string]interface{})

	xError := this.GetAppDataWithWhereAndOrder(ctx, &xData, xWhere, xSort)

	if xError != nil {
		xData = dbmodel.AppDataIngressActionStatic{}
	}

	return xData

}
