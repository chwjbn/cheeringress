package dbmodel

type AppDataIngressNamespace struct {
	AppDataBase `bson:",inline"`
	Title       string `bson:"title" json:"title" index:"uniq"`
	WorkerCount int    `bson:"worker_count" json:"worker_count" index:"worker_count"`
	LastVer     string `bson:"last_ver" json:"last_ver" index:"last_ver"`
	LastPubVer  string `bson:"last_pub_ver" json:"last_pub_ver" index:"last_pub_ver"`
}

func (this *AppDataIngressNamespace) GetTableName() string {
	return "t_app_ingress_namespace"
}
