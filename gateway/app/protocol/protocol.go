package protocol

type WorkerDataBase struct {
	DataId      string `json:"data_id"`
	NamespaceId string `json:"namespace_id"`
}

type WorkerDataActionBackend struct {
	WorkerDataBase
	BalanceType string `json:"balance_type"`
}

type WorkerDataActionBackendNode struct {
	WorkerDataBase
	BackendId   string `json:"backend_id"`
	ServerHost  string `json:"server_host"`
	ServerPort  int    `json:"server_port"`
	WeightScore int    `json:"weight_score"`
}

type WorkerDataActionStatic struct {
	WorkerDataBase
	ContentType string `json:"content_type"`
	DataType    string `json:"data_type"`
	Data        string `json:"data"`
}

type WorkerDataSite struct {
	WorkerDataBase

	OrderNo int `bson:"order_no" json:"order_no"`

	AuthNeed     string `json:"auth_need"`
	AuthUserName string `json:"auth_user_name"`
	AuthPassword string `json:"auth_password"`

	MatchOp    string `json:"match_op"`
	MatchValue string `json:"match_value"`

	ActionType  string `json:"action_type"`
	ActionValue string `json:"action_value"`
}

type WorkerDataSiteRule struct {
	WorkerDataBase

	SiteId string `json:"site_id" index:"site_id"`

	OrderNo int `json:"order_no" index:"order_no"`

	HttpMethod string `json:"http_method"`

	MatchTarget string `json:"match_target"`
	MatchOp     string `json:"match_op"`
	MatchValue  string `json:"match_value"`

	ActionType  string `json:"action_type"`
	ActionValue string `json:"action_value"`
}

type WorkerDataNamespaceDataContent struct {
	ActionBackendInfo     []WorkerDataActionBackend     `json:"action_backend_info"`
	ActionBackendNodeInfo []WorkerDataActionBackendNode `json:"action_backend_node_info"`
	ActionStaticInfo      []WorkerDataActionStatic      `json:"action_static_info"`
	SiteInfo              []WorkerDataSite              `json:"site_info"`
	SiteRuleInfo          []WorkerDataSiteRule          `json:"site_rule_info"`
}

type WorkerBaseRequest struct {
	Token string `json:"token"`
}

type WorkerBaseResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

//获取token
type WorkerGetTokenRequest struct {
	NodeName     string `json:"node_name"`
	NodeOs       string `json:"node_os"`
	NodeArch     string `json:"node_arch"`
	NodeUserName string `json:"node_user_name"`
}

type WorkerGetTokenResponse struct {
	WorkerBaseResponse
	Data string `json:"data"`
}

//查询配置
type WorkerQueryConfigRequest struct {
	WorkerBaseRequest
	NodeCpuCore  int     `json:"node_cpu_core"`
	NodeCpuUsed  float64 `json:"node_cpu_used"`
	NodeMemTotal uint64  `json:"node_mem_total"`
	NodeMemUsed  uint64  `json:"node_mem_used"`
}

type WorkerQueryConfigResponse struct {
	WorkerBaseResponse
	Data string `json:"data"`
}

//获取配置
type WorkerFetchConfigRequest struct {
	WorkerBaseRequest
	ConfigVersion string `json:"config_version"`
}

type WorkerFetchConfigResponse struct {
	WorkerBaseResponse
	Data WorkerDataNamespaceDataContent `json:"data"`
}

//获取资源内容
type WorkerGetResourceDataRequest struct {
	WorkerBaseRequest
	ResType string `json:"res_type"`
	ResId   string `json:"res_id"`
}

type WorkerGetResourceDataResponse struct {
	WorkerBaseResponse
	Data string `json:"data"`
}
