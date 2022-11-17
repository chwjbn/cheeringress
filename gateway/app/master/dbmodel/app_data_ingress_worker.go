package dbmodel

type AppDataIngressWorker struct {
	AppDataBase  `bson:",inline"`
	NamespaceId  string  `bson:"namespace_id" json:"namespace_id" index:"namespace_id"`
	NodeToken    string  `bson:"node_token" json:"node_token" index:"node_token"`
	NodeName     string  `bson:"node_name" json:"node_name"`
	NodeOs       string  `bson:"node_os" json:"node_os"`
	NodeArch     string  `bson:"node_arch" json:"node_arch"`
	NodeUserName string  `bson:"node_user_name" json:"node_user_name"`
	NodeCpuCore  int     `bson:"node_cpu_core" json:"node_cpu_core"`
	NodeCpuUsed  float64 `bson:"node_cpu_used" json:"node_cpu_used"`
	NodeMemTotal uint64  `bson:"node_mem_total" json:"node_mem_total"`
	NodeMemUsed  uint64  `bson:"node_mem_used" json:"node_mem_used"`

	NodeAddr     string `bson:"node_addr" json:"node_addr" index:"node_addr"`
	NodeLastTime string `bson:"node_last_time" json:"node_last_time" index:"node_last_time"`
}

func (this *AppDataIngressWorker) GetTableName() string {
	return "t_app_ingress_worker"
}
