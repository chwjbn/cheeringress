package dbmodel

type AppDataToken struct {
	AppDataBase   `bson:",inline"`
	TokenData     string `json:"token_data" bson:"token_data"`
	AccountId     string `json:"account_id" bson:"account_id" index:"account_id"`
	Username      string `json:"username" bson:"username"`
	Nickname      string `json:"nickname" bson:"nickname"`
	RealName      string `json:"real_name" bson:"real_name"`
	Avatar        string `json:"avatar" bson:"avatar"`
	Role          string `json:"role" bson:"role"`
	LastAliveTime string `bson:"last_alive_time" json:"last_alive_time" index:"last_alive_time"`
	LastAliveIp   string `bson:"last_alive_ip" json:"last_alive_ip" index:"last_alive_ip"`
}

func (this *AppDataToken) GetTableName() string {
	return "t_app_token"
}
