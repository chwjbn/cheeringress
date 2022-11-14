package dbmodel

type AppDataIngressActionBackend struct {
	AppDataBase `bson:",inline"`

	NamespaceId string `bson:"namespace_id" json:"namespace_id" index:"namespace_id"`

	Title string `bson:"title" json:"title" index:"uniq"`

	BalanceType string `bson:"balance_type" json:"balance_type" index:"balance_type"`
	NodeCount   int    `bson:"node_count" json:"node_count" index:"node_count"`
}

func (this *AppDataIngressActionBackend) GetTableName() string {
	return "t_app_ingress_action_backend"
}
