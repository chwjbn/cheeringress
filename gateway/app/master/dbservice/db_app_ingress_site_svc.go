package dbservice

import (
	"context"
	"github.com/chwjbn/cheeringress/app/master/dbmodel"
)

func (this *DbMongoSvc) GetIngressSiteByTitle(ctx context.Context, title string) dbmodel.AppDataIngressSite {

	xData := dbmodel.AppDataIngressSite{}

	xWhere := make(map[string]interface{})
	xWhere["title"] = title

	xSort := make(map[string]interface{})

	xError := this.GetAppDataWithWhereAndOrder(ctx, &xData, xWhere, xSort)

	if xError != nil {
		xData = dbmodel.AppDataIngressSite{}
	}

	return xData

}

func (this *DbMongoSvc) GetFirstIngressSiteByActionValue(ctx context.Context, actionValue string) dbmodel.AppDataIngressSite {

	xData := dbmodel.AppDataIngressSite{}

	xWhere := make(map[string]interface{})
	xWhere["action_value"] = actionValue

	xSort := make(map[string]interface{})

	xError := this.GetAppDataWithWhereAndOrder(ctx, &xData, xWhere, xSort)

	if xError != nil {
		xData = dbmodel.AppDataIngressSite{}
	}

	return xData
}

func (this *DbMongoSvc) GetIngressSiteRuleByTitle(ctx context.Context, siteId string,title string) dbmodel.AppDataIngressSiteRule {

	xData := dbmodel.AppDataIngressSiteRule{}

	xWhere := make(map[string]interface{})
	xWhere["title"] = title
	xWhere["site_id"]=siteId

	xSort := make(map[string]interface{})

	xError := this.GetAppDataWithWhereAndOrder(ctx, &xData, xWhere, xSort)

	if xError != nil {
		xData = dbmodel.AppDataIngressSiteRule{}
	}

	return xData

}

func (this *DbMongoSvc) GetFirstIngressSiteRuleBySiteId(ctx context.Context, siteId string) dbmodel.AppDataIngressSiteRule {

	xData := dbmodel.AppDataIngressSiteRule{}

	xWhere := make(map[string]interface{})
	xWhere["site_id"] = siteId

	xSort := make(map[string]interface{})

	xError := this.GetAppDataWithWhereAndOrder(ctx, &xData, xWhere, xSort)

	if xError != nil {
		xData = dbmodel.AppDataIngressSiteRule{}
	}

	return xData
}

func (this *DbMongoSvc) GetFirstIngressSiteRuleByActionValue(ctx context.Context, actionValue string) dbmodel.AppDataIngressSiteRule {

	xData := dbmodel.AppDataIngressSiteRule{}

	xWhere := make(map[string]interface{})
	xWhere["action_value"] = actionValue

	xSort := make(map[string]interface{})

	xError := this.GetAppDataWithWhereAndOrder(ctx, &xData, xWhere, xSort)

	if xError != nil {
		xData = dbmodel.AppDataIngressSiteRule{}
	}

	return xData
}
