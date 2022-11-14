package dbmodel

type AppDataIngressSite struct {
	AppDataBase `bson:",inline"`

	NamespaceId string `bson:"namespace_id" json:"namespace_id" index:"namespace_id"`

	Title   string `bson:"title" json:"title" index:"uniq"`
	OrderNo int    `bson:"order_no" json:"order_no" index:"order_no"`

	AuthNeed     string `bson:"auth_need" json:"auth_need" index:"auth_need"`
	AuthUserName string `bson:"auth_user_name" json:"auth_user_name"`
	AuthPassword string `bson:"auth_password" json:"auth_password"`

	MatchOp    string `bson:"match_op" json:"match_op"`
	MatchValue string `bson:"match_value" json:"match_value"`

	ActionType  string `bson:"action_type" json:"action_type" index:"action_type"`
	ActionValue string `bson:"action_value" json:"action_value" index:"action_value"`
}

func (this *AppDataIngressSite) GetTableName() string {
	return "t_app_ingress_site"
}
