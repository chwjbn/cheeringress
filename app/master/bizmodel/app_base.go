package bizmodel

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"strings"
)

type PageRequest struct {
	PageNo   int `json:"current"`
	PageSize int `json:"pageSize"`
}

type AppDataIdRequest struct {
	DataId string `json:"data_id"`
}

type DataMapNode struct {
	DataId   string `json:"data_id"`
	DataName string `json:"data_name"`
}

func GetValidateError(rawErr error) error {

	if rawErr == nil {
		return nil
	}

	xErrMsg := rawErr.Error()

	xIdex := strings.Index(xErrMsg, ";")
	if xIdex > 0 {
		xErrMsg = xErrMsg[0:xIdex]
	}

	xIdex = strings.Index(xErrMsg, ":")
	if xIdex > 0 {
		xErrMsg = xErrMsg[xIdex+1:]
	}

	return errors.New(xErrMsg)
}

func (this *AppDataIdRequest) Validate() error {

	xError := validation.ValidateStruct(this, validation.Field(&this.DataId, validation.Required.Error("缺失数据ID!")))

	return GetValidateError(xError)

}
