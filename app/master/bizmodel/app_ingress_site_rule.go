package bizmodel

import validation "github.com/go-ozzo/ozzo-validation"

type IngressSiteRulePageRequest struct {
	PageRequest
	SiteId string `json:"site_id"`
	Title  string `json:"title"`
}

type IngressSiteRuleAddRequest struct {
	NamespaceId string `json:"namespace_id"`
	SiteId      string `json:"site_id"`

	Title   string `json:"title"`
	OrderNo int    `json:"order_no"`

	HttpMethod string `json:"http_method"`

	MatchTarget string `json:"match_target"`
	MatchOp     string `json:"match_op"`
	MatchValue  string `json:"match_value"`

	ActionType  string `json:"action_type"`
	ActionValue string `json:"action_value"`
}

func (this *IngressSiteRuleAddRequest) Validate() error {

	xError := validation.ValidateStruct(this,
		validation.Field(&this.NamespaceId, validation.Required.Error("请选择网关空间!")),
		validation.Field(&this.SiteId, validation.Required.Error("请选择流量站点!")),
		validation.Field(&this.Title, validation.Required.Error("请输入规则名称!"), validation.Length(4, 50).Error("规则名称有效长度4-50!")),
		validation.Field(&this.OrderNo, validation.Required.Error("请输入规则序号!")),
		validation.Field(&this.HttpMethod, validation.Required.Error("请选择HTTP方法!")),
		validation.Field(&this.MatchTarget, validation.Required.Error("请选择规则匹配对象!")),
		validation.Field(&this.MatchOp, validation.Required.Error("请选择规则匹配方式!")),
		validation.Field(&this.MatchValue, validation.Required.Error("请选择规则匹配内容!")),
		validation.Field(&this.ActionType, validation.Required.Error("请选择响应类型!")),
		validation.Field(&this.ActionValue, validation.Required.Error("请选择响应内容!")),
	)

	return GetValidateError(xError)

}
