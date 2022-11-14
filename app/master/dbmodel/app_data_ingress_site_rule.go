package dbmodel

type AppDataIngressSiteRule struct {
	AppDataBase `bson:",inline"`

	NamespaceId string `bson:"namespace_id" json:"namespace_id" index:"namespace_id"`
	SiteId      string `bson:"site_id" json:"site_id" index:"site_id"`

	Title   string `bson:"title" json:"title" index:"uniq"`
	OrderNo int    `bson:"order_no" json:"order_no" index:"order_no"`

	HttpMethod string `bson:"http_method" json:"http_method"`

	MatchTarget string `bson:"match_target" json:"match_target"`
	MatchOp     string `bson:"match_op" json:"match_op"`
	MatchValue  string `bson:"match_value" json:"match_value"`

	ActionType  string `bson:"action_type" json:"action_type" index:"action_type"`
	ActionValue string `bson:"action_value" json:"action_value" index:"action_value"`
}

func (this *AppDataIngressSiteRule) GetTableName() string {
	return "t_app_ingress_site_rule"
}
