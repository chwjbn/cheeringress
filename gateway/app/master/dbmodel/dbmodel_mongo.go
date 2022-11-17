package dbmodel

import (
	"fmt"
	"github.com/chwjbn/cheeringress/cheerlib"
)

type IDbModelMongo interface {
	GetDbName() string

	GetTableName() string

	InitDataId() error
	SetDataId(data string)
	GetDataId() string

	SetTenantId(data string)
	GetTenantId() string

	SetCreateTime(data string)
	GetCreateTime() string

	SetUpdateTime(data string)
	GetUpdateTime() string

	SetUpdateIp(data string)
	GetUpdateIp() string

	SetState(data string)
	GetState() string
}

type DbModelMongo struct {
	DataId string `bson:"data_id" json:"data_id" index:"data_id"`

	TenantId string `bson:"tenant_id" json:"tenant_id" index:"tenant_id"`

	CreateTime string `bson:"create_time" json:"create_time" index:"create_time"`

	UpdateTime string `bson:"update_time" json:"update_time" index:"update_time"`

	UpdateIp string `bson:"update_ip" json:"update_ip" index:"update_ip"`

	State string `bson:"state" json:"state" index:"state"`
}

func (this *DbModelMongo) InitDataIdWithRand(rnd string) error {

	var xError error = nil

	this.SetDataId(cheerlib.EncryptMd5(fmt.Sprintf("%s-%s-%s", this.GetTenantId(), cheerlib.EncryptNewId(), rnd)))

	return xError

}

func (this *DbModelMongo) InitDataId() error {

	return this.InitDataIdWithRand("")

}

func (this *DbModelMongo) SetDataId(data string) {
	this.DataId = data
}

func (this *DbModelMongo) GetDataId() string {
	return this.DataId
}

func (this *DbModelMongo) SetTenantId(data string) {
	this.TenantId = data
}

func (this *DbModelMongo) GetTenantId() string {
	return this.TenantId
}

func (this *DbModelMongo) SetCreateTime(data string) {
	this.CreateTime = data
}

func (this *DbModelMongo) GetCreateTime() string {
	return this.CreateTime
}

func (this *DbModelMongo) SetUpdateTime(data string) {
	this.UpdateTime = data
}

func (this *DbModelMongo) GetUpdateTime() string {
	return this.UpdateTime
}

func (this *DbModelMongo) SetUpdateIp(data string) {
	this.UpdateIp = data
}

func (this *DbModelMongo) GetUpdateIp() string {
	return this.UpdateIp
}

func (this *DbModelMongo) SetState(data string) {
	this.State = data
}

func (this *DbModelMongo) GetState() string {
	return this.State
}
