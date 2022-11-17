package dbmodel

type AppDataIngressConfig struct {
	AppDataBase `bson:",inline"`
	NamespaceId string `bson:"namespace_id" json:"namespace_id" index:"namespace_id"`
	Version     string `bson:"version" json:"version" index:"version"`
	Data        string `bson:"data" json:"data"`
}

func (this *AppDataIngressConfig) GetTableName() string {
	return "t_app_ingress_config"
}
