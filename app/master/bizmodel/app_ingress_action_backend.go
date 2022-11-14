package bizmodel

import validation "github.com/go-ozzo/ozzo-validation"

type IngressActionBackendPageRequest struct {
	PageRequest
	NamespaceId string `json:"namespace_id"`
	Title       string `json:"title"`
	State       string `bson:"state"`
}

type IngressActionBackendAddRequest struct {
	NamespaceId string `json:"namespace_id"`
	Title       string `json:"title"`
	BalanceType string `json:"balance_type"`
}

type IngressActionBackendSaveRequest struct {
	AppDataIdRequest
	NamespaceId string `json:"namespace_id"`
	Title       string `json:"title"`
	BalanceType string `json:"balance_type"`
	State       string `bson:"state"`
}

type IngressActionBackendNodePageRequest struct {
	PageRequest
	BackendId string `json:"backend_id"`
}

type IngressActionBackendNodeAddRequest struct {
	NamespaceId string `json:"namespace_id"`
	BackendId   string `json:"backend_id"`

	Title       string `json:"title"`
	ServerHost  string `json:"server_host"`
	ServerPort  int    `json:"server_port"`
	WeightScore int    `json:"weight_score"`
}

func (this *IngressActionBackendAddRequest) Validate() error {

	xError := validation.ValidateStruct(this,
		validation.Field(&this.NamespaceId, validation.Required.Error("请选择网关空间!")),
		validation.Field(&this.Title, validation.Required.Error("请输入代理名称!"), validation.Length(4, 50).Error("代理名称有效长度4-50!")),
		validation.Field(&this.BalanceType, validation.Required.Error("请选择负载均衡策略!"), validation.In("IPHash", "RoundRobin", "LeastConnection").Error("请选择正确的负载均衡策略!")),
	)

	return GetValidateError(xError)

}

func (this *IngressActionBackendNodeAddRequest) Validate() error {

	xError := validation.ValidateStruct(this,
		validation.Field(&this.NamespaceId, validation.Required.Error("请选择网关空间!")),
		validation.Field(&this.BackendId, validation.Required.Error("请选择反向代理!")),
		validation.Field(&this.Title, validation.Required.Error("请输入节点名称!"), validation.Length(4, 50).Error("节点名称有效长度4-50!")),
		validation.Field(&this.ServerHost, validation.Required.Error("请输入服务器地址!")),
		validation.Field(&this.ServerPort, validation.Required.Error("请输入服务器端口!")),
		validation.Field(&this.WeightScore, validation.Required.Error("请输入服务器权重!")),
	)

	return GetValidateError(xError)

}

func (this *IngressActionBackendSaveRequest) Validate() error {

	xError := validation.ValidateStruct(this,
		validation.Field(&this.NamespaceId, validation.Required.Error("请选择网关空间!")),
		validation.Field(&this.Title, validation.Required.Error("请输入代理名称!"), validation.Length(4, 50).Error("代理名称有效长度4-50!")),
		validation.Field(&this.BalanceType, validation.Required.Error("请选择负载均衡策略!"), validation.In("IPHash", "RoundRobin", "LeastConnection").Error("请选择正确的负载均衡策略!")),
		validation.Field(&this.State, validation.Required.Error("请选择状态!"), validation.In("enable", "disable").Error("请选择正确的状态!")),
	)

	return GetValidateError(xError)

}
