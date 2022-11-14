package dbmodel

type AppDataBase struct {
	DbModelMongo `bson:",inline"`
}

func (a AppDataBase) GetDbName() string {
	return "db_cheer_ingress"
}
