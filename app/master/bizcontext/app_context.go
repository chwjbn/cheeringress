package bizcontext

import (
	"github.com/chwjbn/cheeringress/app/master/dbservice"
	"github.com/chwjbn/cheeringress/config"
	"errors"
	"fmt"
)

type AppContext struct {
	Config   *config.ConfigAppMaster
	AppDbSvc *dbservice.DbMongoSvc
}

func (this *AppContext) Init() error {

	var xError error

	var xDbError error
	xDbError, this.AppDbSvc = dbservice.NewDbMongoSvc(this.Config.DbAppMongodbUri)
	if xDbError != nil {
		xError = errors.New(fmt.Sprintf("AppContext.Init DbError=[%s]", xDbError.Error()))
		return xError
	}

	return xError

}
