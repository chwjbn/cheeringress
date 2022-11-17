package dbmodel

type AppDataIngressActionStatic struct {
	AppDataBase `bson:",inline"`

	NamespaceId string `bson:"namespace_id" json:"namespace_id" index:"namespace_id"`

	Title string `bson:"title" json:"title" index:"uniq"`

	ContentType string `bson:"content_type" json:"content_type" index:"content_type"`

	DataType string `bson:"data_type" json:"data_type"`
	Data     string `bson:"data" json:"data"`
}

func (this *AppDataIngressActionStatic) GetTableName() string {
	return "t_app_ingress_action_static"
}
