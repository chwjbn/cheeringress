package bizservice

import (
	"context"
	"github.com/chwjbn/cheeringress/app/master/bizcontext"
	"github.com/chwjbn/cheeringress/app/master/dbmodel"
	"github.com/chwjbn/cheeringress/cheerlib"
)

type AppAccountService struct {
	Context *bizcontext.AppContext
}

func (this *AppAccountService) GetAccountCount(ctx context.Context) int64 {

	xWhere := make(map[string]interface{})

	return this.Context.AppDbSvc.GetAppDataCount(ctx, &dbmodel.AppDataUser{}, xWhere)

}

func (this *AppAccountService) GetAccountInfoByUsername(ctx context.Context, username string) dbmodel.AppDataUser {

	xData := dbmodel.AppDataUser{}

	xWhere := make(map[string]interface{})
	xSort := make(map[string]interface{})

	xError := this.Context.AppDbSvc.GetAppDataWithWhereAndOrder(ctx, &xData, xWhere, xSort)

	if xError != nil {
		xData = dbmodel.AppDataUser{}
		cheerlib.LogError(xError.Error())
	}

	return xData

}

func (this *AppAccountService) GetAccountInfoById(ctx context.Context, dataId string) dbmodel.AppDataUser {

	xData := dbmodel.AppDataUser{}
	xData.SetDataId(dataId)

	xError := this.Context.AppDbSvc.GetAppDataById(ctx, &xData)

	if xError != nil {
		xData = dbmodel.AppDataUser{}
		cheerlib.LogError(xError.Error())
	}

	return xData

}

func (this *AppAccountService) GetAccountTokenById(ctx context.Context, dataId string) dbmodel.AppDataToken {

	xData := dbmodel.AppDataToken{}
	xData.SetDataId(dataId)

	xError := this.Context.AppDbSvc.GetAppDataById(ctx, &xData)

	if xError != nil {
		xData = dbmodel.AppDataToken{}
		cheerlib.LogError(xError.Error())
	}

	return xData

}

func (this *AppAccountService) RemoveAccountTokenById(ctx context.Context, dataId string) {

	xData := dbmodel.AppDataToken{}
	xData.SetDataId(dataId)

	xError := this.Context.AppDbSvc.DeleteAppDataById(ctx, &xData)

	if xError != nil {
		cheerlib.LogError(xError.Error())
	}
}

func (this *AppAccountService) CreateAccountInfo(ctx context.Context, data dbmodel.AppDataUser) error {

	var xError error = nil

	data.PwdSalt = cheerlib.EncryptNewId()

	data.Password = data.GetEncryptPassword(data.Password)

	data.UpdateTime = cheerlib.TimeGetNow()

	data.State = "active"

	xError, _ = this.Context.AppDbSvc.AddAppData(ctx, &data)

	return xError

}

func (this *AppAccountService) CreateAccountToken(ctx context.Context, data dbmodel.AppDataToken) (error, string) {

	xError, xTokenId := this.Context.AppDbSvc.AddAppData(ctx, &data)

	return xError, xTokenId

}
