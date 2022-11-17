package bizmodel

import validation "github.com/go-ozzo/ozzo-validation"

type IngressSitePageRequest struct {
	PageRequest
	NamespaceId string `json:"namespace_id"`
	Title       string `json:"title"`
	State       string `bson:"state"`
}

type IngressSiteAddRequest struct {
	NamespaceId string `json:"namespace_id"`
	Title       string `json:"title"`
	OrderNo     int    `json:"order_no"`

	AuthNeed     string `json:"auth_need"`
	AuthUserName string `json:"auth_user_name"`
	AuthPassword string `json:"auth_password"`

	MatchOp    string `json:"match_op"`
	MatchValue string `json:"match_value"`

	ActionType  string `json:"action_type"`
	ActionValue string `json:"action_value"`
}

type IngressSiteSaveRequest struct {
	AppDataIdRequest
	NamespaceId string `json:"namespace_id"`
	Title       string `json:"title"`
	OrderNo     int    `json:"order_no"`

	AuthNeed     string `json:"auth_need"`
	AuthUserName string `json:"auth_user_name"`
	AuthPassword string `json:"auth_password"`

	MatchOp    string `json:"match_op"`
	MatchValue string `json:"match_value"`

	ActionType  string `json:"action_type"`
	ActionValue string `json:"action_value"`
	State       string `bson:"state"`
}

func (this *IngressSiteAddRequest) Validate() error {

	xError := validation.ValidateStruct(this,
		validation.Field(&this.NamespaceId, validation.Required.Error("请选择网关空间!")),
		validation.Field(&this.Title, validation.Required.Error("请输入站点名称!"), validation.Length(4, 50).Error("站点名称有效长度4-50!")),
		validation.Field(&this.OrderNo, validation.Required.Error("请输入站点序号!")),
		validation.Field(&this.AuthNeed, validation.Required.Error("请选择认证状态!")),
		validation.Field(&this.MatchOp, validation.Required.Error("请选择域名匹配方式!")),
		validation.Field(&this.MatchValue, validation.Required.Error("请输入域名匹配规则!")),
		validation.Field(&this.ActionType, validation.Required.Error("请选择响应类型!")),
		validation.Field(&this.ActionValue, validation.Required.Error("请选择响应内容!")),
	)

	return GetValidateError(xError)

}

func (this *IngressSiteSaveRequest) Validate() error {

	xError := validation.ValidateStruct(this,
		validation.Field(&this.NamespaceId, validation.Required.Error("请选择网关空间!")),
		validation.Field(&this.Title, validation.Required.Error("请输入站点名称!"), validation.Length(4, 50).Error("站点名称有效长度4-50!")),
		validation.Field(&this.OrderNo, validation.Required.Error("请输入站点序号!")),
		validation.Field(&this.AuthNeed, validation.Required.Error("请选择认证状态!")),
		validation.Field(&this.MatchOp, validation.Required.Error("请选择域名匹配方式!")),
		validation.Field(&this.MatchValue, validation.Required.Error("请输入域名匹配内容!")),
		validation.Field(&this.ActionType, validation.Required.Error("请选择响应类型!")),
		validation.Field(&this.ActionValue, validation.Required.Error("请选择响应内容!")),
		validation.Field(&this.State, validation.Required.Error("请选择状态!"), validation.In("enable", "disable").Error("请选择正确的状态!")),
	)

	return GetValidateError(xError)

}
