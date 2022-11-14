package dbmodel

type AppDataIngressActionBackendNode struct {
	AppDataBase `bson:",inline"`

	NamespaceId string `bson:"namespace_id" json:"namespace_id" index:"namespace_id"`
	BackendId   string `bson:"backend_id" json:"backend_id" index:"backend_id"`

	Title       string `bson:"title" json:"title" index:"uniq"`
	ServerHost  string `bson:"server_host" json:"server_host" index:"server_host"`
	ServerPort  int    `bson:"server_port" json:"server_port" index:"server_port"`
	WeightScore int    `bson:"weight_score" json:"weight_score" index:"weight_score"`
}

func (this *AppDataIngressActionBackendNode) GetTableName() string {
	return "t_app_ingress_action_backend_node"
}
