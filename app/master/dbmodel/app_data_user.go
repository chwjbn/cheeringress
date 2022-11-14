package dbmodel

import (
	"fmt"
	"github.com/chwjbn/cheeringress/cheerlib"
)

type AppDataUser struct {
	AppDataBase `bson:",inline"`
	Username    string `bson:"username" json:"username" index:"uniq"`
	Password    string `bson:"password" json:"password"`
	PwdSalt     string `bson:"pwd_salt" json:"pwd_salt"`

	Nickname string `json:"nickname" bson:"nickname"`
	RealName string `json:"real_name" bson:"real_name"`
	Avatar   string `json:"avatar" bson:"avatar"`
	Role     string `json:"role" bson:"role"`
}

func (this *AppDataUser) GetTableName() string {
	return "t_app_user"
}

func (this *AppDataUser) GetEncryptPassword(passwd string) string {

	return cheerlib.EncryptMd5(fmt.Sprintf("76==========%s=========77======%s", this.PwdSalt, passwd))

}
