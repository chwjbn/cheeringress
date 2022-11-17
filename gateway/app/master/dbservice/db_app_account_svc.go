package dbservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/chwjbn/cheeringress/app/master/dbmodel"
	"github.com/chwjbn/cheeringress/cheerlib"
)

func (this *DbMongoSvc) GetAccountByUserName(ctx context.Context, username string) dbmodel.AppDataUser {

	xData := dbmodel.AppDataUser{}

	xWhere := make(map[string]interface{})
	xWhere["username"] = username

	xSort := make(map[string]interface{})

	xError := this.GetAppDataWithWhereAndOrder(ctx, &xData, xWhere, xSort)

	if xError != nil {
		xData = dbmodel.AppDataUser{}
	}

	return xData

}

func (this *DbMongoSvc) GetAccountCount(ctx context.Context) int64 {

	var nRet int64

	xData := dbmodel.AppDataUser{}
	xWhere := make(map[string]interface{})

	nRet = this.GetAppDataCount(ctx, &xData, xWhere)

	return nRet

}

func (this *DbMongoSvc) CreateAccount(ctx context.Context, username string, password string, clientIp string) error {

	var xError error = nil

	xAccountData := this.GetAccountByUserName(ctx, username)
	if len(xAccountData.Username) > 0 {
		xError = errors.New("app.server.msg.account.username.exists")
		return xError
	}

	xAccountData.TenantId = cheerlib.EncryptMd5("cheeradmin")
	xAccountData.SetUpdateTime(cheerlib.TimeGetNow())
	xAccountData.SetUpdateIp(clientIp)
	xAccountData.SetCreateTime(cheerlib.TimeGetNow())
	xAccountData.SetState("enable")

	xAccountData.Username = username
	xAccountData.Nickname = username
	xAccountData.RealName = username
	xAccountData.Role = "admin"

	xAccountData.PwdSalt = cheerlib.EncryptNewId()
	xAccountData.Password = xAccountData.GetEncryptPassword(password)

	xAccountData.Avatar = "/user_avar.png"

	xAccountData.InitDataId()

	xDbError, _ := this.AddAppData(ctx, &xAccountData)

	if xDbError != nil {
		xError = errors.New(fmt.Sprintf("Internal Error=[%s]", xDbError.Error()))
		return xError
	}

	return xError

}
