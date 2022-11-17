package dbservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/chwjbn/cheeringress/app/master/dbmodel"
	"github.com/chwjbn/cheeringress/cheerapp"
	"github.com/chwjbn/cheeringress/cheerlib"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"reflect"
	"strings"
	"time"
)

type DbMongoSvc struct {
	mDbUri    string
	mDbClient *mongo.Client
}

func NewDbMongoSvc(dbUri string) (error, *DbMongoSvc) {

	var xError error = nil

	if len(dbUri) < 1 {
		return xError, nil
	}

	xThis := &DbMongoSvc{}
	xThis.mDbUri = dbUri
	xError = xThis.initClient()

	if xError != nil {
		return xError, nil
	}

	return xError, xThis
}

func (this *DbMongoSvc) initClient() error {

	var xError error = nil
	var xDbError error = nil

	//创建连接到数据库
	clientOptions := options.Client().
		ApplyURI(this.mDbUri).
		SetHeartbeatInterval(3 * time.Second).
		SetConnectTimeout(10 * time.Second).
		SetMaxPoolSize(50).
		SetMaxConnIdleTime(30 * time.Second)

	cheerlib.LogInfo(fmt.Sprintf("DbMongoSvc Begin Connect To [%s]", this.mDbUri))

	this.mDbClient, xDbError = mongo.Connect(context.Background(), clientOptions)
	if xDbError != nil {
		xError = errors.New(fmt.Sprintf("DbMongoSvc.CheckClient mongo.Connect Error=[%s]", xDbError.Error()))
		cheerlib.LogError(xError.Error())
		return xError
	}

	cheerlib.LogInfo(fmt.Sprintf("DbMongoSvc Successfully Connected To [%s]", this.mDbUri))

	return xError

}

func (this *DbMongoSvc) GetDbHandle(ctx context.Context, dbName string) (error, *mongo.Database) {

	xDbHandle := this.mDbClient.Database(dbName)

	return nil, xDbHandle
}

func (this *DbMongoSvc) GetTableHandle(ctx context.Context, dbName string, tableName string) (error, *mongo.Collection) {

	xErr, xDbHandle := this.GetDbHandle(ctx, dbName)

	if xErr != nil {
		return xErr, nil
	}

	return xErr, xDbHandle.Collection(tableName)
}

//确保数据索引是创建的
func (this *DbMongoSvc) CheckAppDataIndexs(ctx context.Context, data dbmodel.IDbModelMongo) error {

	var xError error = nil

	xIndexList := this.getAppDataIndexs(data)

	//没有索引
	if len(xIndexList) < 1 {
		return xError
	}

	xError, xTableHandle := this.GetTableHandle(ctx, data.GetDbName(), data.GetTableName())
	if xError != nil {
		return xError
	}

	xTableIndexView := xTableHandle.Indexes()
	xIndexCursor, xIndexError := xTableIndexView.List(context.TODO())
	if xIndexError != nil {
		xError = errors.New(fmt.Sprintf("DbMongoSvc.CheckAppDataIndexs List Index Error=[%s]", xIndexError.Error()))
		return xError
	}

	defer func() {
		xIndexCursor.Close(context.TODO())
	}()

	//获取存在的索引列表
	var xIndexResult []bson.M
	xIndexError = xIndexCursor.All(context.TODO(), &xIndexResult)
	if xIndexError != nil {
		xError = errors.New(fmt.Sprintf("DbMongoSvc.CheckAppDataIndexs IndexCursor.All Error=[%s]", xIndexError.Error()))
		return xError
	}

	xExistsIndexList := []string{}
	for _, xIndexDataItem := range xIndexResult {
		xIndexName := xIndexDataItem["name"].(string)

		if len(xIndexName) < 1 {
			continue
		}

		xExistsIndexList = append(xExistsIndexList, xIndexName)
	}

	xExistsIndexStr := fmt.Sprintf("|%s|", strings.Join(xExistsIndexList, "|"))

	//把索引拆解为数组
	xNeedAddIndexList := [][]string{}
	for _, xIndexDataItem := range xIndexList {

		xIndexDataSpit := strings.Split(xIndexDataItem, "|")
		if len(xIndexDataSpit) < 2 {
			continue
		}

		xIndexName := fmt.Sprintf("idex_%s", xIndexDataSpit[0])
		if strings.Contains(xExistsIndexStr, fmt.Sprintf("|%s|", xIndexName)) {
			continue
		}

		xNeedAddIndexList = append(xNeedAddIndexList, xIndexDataSpit)
	}

	if len(xNeedAddIndexList) < 1 {
		return xError
	}

	//对于需要创建的索引执行创建
	for _, xIndexDataItem := range xNeedAddIndexList {

		if len(xIndexDataItem) < 2 {
			continue
		}

		xIndexName := fmt.Sprintf("idex_%s", xIndexDataItem[0])

		xIndexOpt := new(options.IndexOptions)
		xIndexOpt.SetName(xIndexName)
		xIndexOpt.SetBackground(true)

		if strings.EqualFold(xIndexDataItem[1], "uniq") {
			xIndexOpt.SetUnique(true)
		}

		xIndexData := mongo.IndexModel{Options: xIndexOpt, Keys: bsonx.Doc{}}
		xIndexData.Keys = xIndexData.Keys.(bsonx.Doc).Append(xIndexDataItem[0], bsonx.Int32(1))

		xCreateIndexResult, xCreateIndexError := xTableIndexView.CreateOne(context.TODO(), xIndexData)
		if xCreateIndexError != nil {
			cheerlib.LogError(fmt.Sprintf("DbMongoSvc.CheckAppDataIndexs TableIndexView.CreateOne IndexName=[%s.%s.%s] Error=[%s]", data.GetDbName(), data.GetTableName(), xIndexName, xCreateIndexError.Error()))
			continue
		}

		cheerlib.LogInfo(fmt.Sprintf("DbMongoSvc.CheckAppDataIndexs TableIndexView.CreateOne IndexName=[%s.%s.%s] CreateIndexResult=[%s]", data.GetDbName(), data.GetTableName(), xIndexName, xCreateIndexResult))
	}

	return xError

}

//获取一个数据对象，带条件和排序
func (this *DbMongoSvc) GetAppDataWithWhereAndOrder(ctx context.Context, data dbmodel.IDbModelMongo, whereMap map[string]interface{}, sortMap map[string]interface{}) error {

	var xError error = nil

	xSpan := cheerapp.SpanBeginDbService(ctx, "DbMongoSvc.GetAppDataWithWhereAndOrder")
	defer func() {

		if xError != nil {
			cheerapp.SpanError(xSpan, xError.Error())
		}

		cheerapp.SpanEnd(xSpan)
	}()

	cheerapp.SpanTag(xSpan, "_ARG_data", cheerlib.TextStructToJson(data))
	cheerapp.SpanTag(xSpan, "_ARG_whereMap", cheerlib.TextStructToJson(whereMap))
	cheerapp.SpanTag(xSpan, "_ARG_sortMap", cheerlib.TextStructToJson(sortMap))

	xError, xTableHandle := this.GetTableHandle(ctx, data.GetDbName(), data.GetTableName())
	if xError != nil {
		return xError
	}

	xWhere := bson.M{}
	for dKey, dVal := range whereMap {
		xWhere[dKey] = dVal
	}

	xFindOneOpt := options.FindOne()

	xSort := bson.M{}
	for dKey, dVal := range sortMap {
		xSort[dKey] = dVal
	}

	xFindOneOpt.SetSort(xSort)

	xFindResult := xTableHandle.FindOne(ctx, xWhere, xFindOneOpt)
	if xFindResult.Err() != nil {
		return xError
	}

	xDbDataDecodeError := xFindResult.Decode(data)
	if xDbDataDecodeError != nil {
		xError = errors.New(fmt.Sprintf("DbMongoSvc.GetAppDataWithWhereAndOrder FindOne Decode Error=[%s]", xDbDataDecodeError.Error()))
		return xError
	}

	return xError

}

//根据ID查找数据
func (this *DbMongoSvc) GetAppDataById(ctx context.Context, data dbmodel.IDbModelMongo) error {

	xWhere := make(map[string]interface{})
	xWhere["data_id"] = data.GetDataId()

	xSort := make(map[string]interface{})
	xSort["data_id"] = 1

	return this.GetAppDataWithWhereAndOrder(ctx, data, xWhere, xSort)

}

//查询对应条件的数量
func (this *DbMongoSvc) GetAppDataCount(ctx context.Context, data dbmodel.IDbModelMongo, whereMap map[string]interface{}) int64 {

	var xCount int64 = 0

	xSpan := cheerapp.SpanBeginDbService(ctx, "DbMongoSvc.GetAppDataCount")
	defer func() {
		cheerapp.SpanEnd(xSpan)
	}()

	cheerapp.SpanTag(xSpan, "_ARG_data", cheerlib.TextStructToJson(data))
	cheerapp.SpanTag(xSpan, "_ARG_whereMap", cheerlib.TextStructToJson(whereMap))

	xError, xTableHandle := this.GetTableHandle(ctx, data.GetDbName(), data.GetTableName())
	if xError != nil {
		cheerlib.LogError(fmt.Sprintf("DbMongoSvc.GetAppDataCount Error=[%s]", xError.Error()))
		return xCount
	}

	xWhere := bson.M{}
	for dataK, dataV := range whereMap {
		xWhere[dataK] = dataV
	}

	xCount, xError = xTableHandle.CountDocuments(ctx, xWhere)
	if xError != nil {
		cheerlib.LogError(fmt.Sprintf("DbMongoSvc.GetAppDataCount Error=[%s]", xError.Error()))
		return xCount
	}

	return xCount

}

//获取一个数据列表，带条件和排序
func (this *DbMongoSvc) GetAppDataListWithWhereAndOrder(ctx context.Context, data dbmodel.IDbModelMongo, whereMap map[string]interface{}, sortMap map[string]interface{}, fromIndex int64, limitCount int64) (error, []dbmodel.IDbModelMongo) {

	var xError error = nil

	xSpan := cheerapp.SpanBeginDbService(ctx, "DbMongoSvc.GetAppDataListWithWhereAndOrder")
	defer func() {
		if xError != nil {
			cheerapp.SpanError(xSpan, xError.Error())
		}
		cheerapp.SpanEnd(xSpan)
	}()

	cheerapp.SpanTag(xSpan, "_ARG_data", cheerlib.TextStructToJson(data))
	cheerapp.SpanTag(xSpan, "_ARG_whereMap", cheerlib.TextStructToJson(whereMap))
	cheerapp.SpanTag(xSpan, "_ARG_sortMap", cheerlib.TextStructToJson(sortMap))
	cheerapp.SpanTag(xSpan, "_ARG_fromIndex", cheerlib.TextStructToJson(fromIndex))
	cheerapp.SpanTag(xSpan, "_ARG_limitCount", cheerlib.TextStructToJson(limitCount))

	xDataList := []dbmodel.IDbModelMongo{}

	xError, xTableHandle := this.GetTableHandle(ctx, data.GetDbName(), data.GetTableName())
	if xError != nil {
		return xError, xDataList
	}

	xWhere := bson.M{}
	for dKey, dVal := range whereMap {
		xWhere[dKey] = dVal
	}

	xFindOpt := options.Find()

	if fromIndex > 0 {
		xFindOpt.SetSkip(fromIndex)
	}

	if limitCount > 0 {
		xFindOpt.SetLimit(limitCount)
	}

	xSort := bson.M{}
	for dKey, dVal := range sortMap {
		xSort[dKey] = dVal
	}

	xFindOpt.SetSort(xSort)

	xFindCur, xDbError := xTableHandle.Find(ctx, xWhere, xFindOpt)
	if xDbError != nil {
		xError = errors.New(fmt.Sprintf("DbMongoSvc.GetAppDataListWithWhereAndOrder Find DbError=[%s]", xDbError.Error()))
		return xError, xDataList
	}

	defer func() {
		xFindCur.Close(context.TODO())
	}()

	xDataType := reflect.TypeOf(data)
	if xDataType.Kind() == reflect.Ptr {
		xDataType = xDataType.Elem()
	}

	for xFindCur.Next(context.TODO()) {

		xDataHandle := reflect.New(xDataType)
		xDataItem := xDataHandle.Interface()

		xDbDecodeError := xFindCur.Decode(xDataItem)
		if xDbDecodeError != nil {
			cheerlib.LogError(fmt.Sprintf("DbMongoSvc.GetAppDataListWithWhereAndOrder Find DbDecodeError=[%s]", xDbDecodeError.Error()))
			continue
		}

		xDataList = append(xDataList, xDataItem.(dbmodel.IDbModelMongo))
	}

	return xError, xDataList
}

//分页获取数据
func (this *DbMongoSvc) GetDataPageList(ctx context.Context, data dbmodel.IDbModelMongo, whereMap map[string]interface{}, sortMap map[string]interface{}, pageNo int, pageSize int) dbmodel.PageData {

	xPageData := dbmodel.PageData{}

	xPageData.TotalCount = this.GetAppDataCount(ctx, data, whereMap)
	xPageData.PageNo = int64(pageNo)
	xPageData.PageSize = int64(pageSize)
	xPageData.Calc()

	xFromIndex := (xPageData.PageNo - 1) * xPageData.PageSize
	xLimitCount := xPageData.PageSize

	xError, xDataList := this.GetAppDataListWithWhereAndOrder(ctx, data, whereMap, sortMap, xFromIndex, xLimitCount)

	if xError != nil {
		cheerlib.LogError(fmt.Sprintf("DbMongoSvc.GetDataPageList Error=[%s]", xError.Error()))
		return xPageData
	}

	xPageData.DataList = xDataList

	return xPageData
}

//添加数据
func (this *DbMongoSvc) AddAppData(ctx context.Context, data dbmodel.IDbModelMongo) (error, string) {

	var xError error = nil
	xDataId := ""

	xSpan := cheerapp.SpanBeginDbService(ctx, "DbMongoSvc.AddAppData")
	defer func() {
		if xError != nil {
			cheerapp.SpanError(xSpan, xError.Error())
		}
		cheerapp.SpanEnd(xSpan)
	}()

	cheerapp.SpanTag(xSpan, "_ARG_data", cheerlib.TextStructToJson(data))

	//检查索引
	xError = this.CheckAppDataIndexs(ctx, data)
	if xError != nil {
		return xError, xDataId
	}

	xError, xTableHandle := this.GetTableHandle(ctx, data.GetDbName(), data.GetTableName())
	if xError != nil {
		return xError, xDataId
	}

	if len(data.GetDataId()) < 1 {
		data.InitDataId()
	}

	if len(data.GetCreateTime()) < 1 {
		data.SetCreateTime(cheerlib.TimeGetNow())
	}

	if len(data.GetUpdateTime()) < 1 {
		data.SetUpdateTime(cheerlib.TimeGetNow())
	}

	if len(data.GetUpdateIp()) < 1 {
		data.SetUpdateIp("127.0.0.1")
	}

	_, xActInsertError := xTableHandle.InsertOne(ctx, data)
	if xActInsertError != nil {
		xError = errors.New(fmt.Sprintf("DbMongoSvc.AddAppData Insert Error=[%s]", xActInsertError.Error()))
		return xError, xDataId
	}

	xDataId = data.GetDataId()

	return xError, xDataId

}

func (this *DbMongoSvc) UpdateAppData(ctx context.Context, data dbmodel.IDbModelMongo, whereMap map[string]interface{}) error {

	var xError error = nil

	xSpan := cheerapp.SpanBeginDbService(ctx, "DbMongoSvc.UpdateAppData")
	defer func() {
		if xError != nil {
			cheerapp.SpanError(xSpan, xError.Error())
		}
		cheerapp.SpanEnd(xSpan)
	}()

	cheerapp.SpanTag(xSpan, "_ARG_data", cheerlib.TextStructToJson(data))
	cheerapp.SpanTag(xSpan, "_ARG_whereMap", cheerlib.TextStructToJson(whereMap))

	xError, xTableHandle := this.GetTableHandle(ctx, data.GetDbName(), data.GetTableName())
	if xError != nil {
		return xError
	}

	xWhere := bson.M{}
	for dKey, dVal := range whereMap {
		xWhere[dKey] = dVal
	}

	_, xDbError := xTableHandle.UpdateMany(ctx, xWhere, bson.M{"$set": data})

	if xDbError != nil {
		xError = errors.New(fmt.Sprintf("DbMongoSvc.SaveAppData DbError=[%s]", xDbError.Error()))
	}

	return xError

}

func (this *DbMongoSvc) UpdateAppDataById(ctx context.Context, data dbmodel.IDbModelMongo) error {

	var xError error = nil

	xWhereMap := make(map[string]interface{})
	xWhereMap["data_id"] = data.GetDataId()

	xError = this.UpdateAppData(ctx, data, xWhereMap)

	return xError

}

func (this *DbMongoSvc) DeleteAppData(ctx context.Context, data dbmodel.IDbModelMongo, whereMap map[string]interface{}) error {

	var xError error = nil

	xSpan := cheerapp.SpanBeginDbService(ctx, "DbMongoSvc.DeleteAppData")
	defer func() {
		if xError != nil {
			cheerapp.SpanError(xSpan, xError.Error())
		}
		cheerapp.SpanEnd(xSpan)
	}()

	cheerapp.SpanTag(xSpan, "_ARG_data", cheerlib.TextStructToJson(data))
	cheerapp.SpanTag(xSpan, "_ARG_whereMap", cheerlib.TextStructToJson(whereMap))

	xError, xTableHandle := this.GetTableHandle(ctx, data.GetDbName(), data.GetTableName())
	if xError != nil {
		return xError
	}

	xWhere := bson.M{}
	for dKey, dVal := range whereMap {
		xWhere[dKey] = dVal
	}

	_, xDbError := xTableHandle.DeleteMany(ctx, xWhere)

	if xDbError != nil {
		xError = errors.New(fmt.Sprintf("DbMongoSvc.DeleteAppData DbError=[%s]", xDbError.Error()))
	}

	return xError

}

func (this *DbMongoSvc) DeleteAppDataById(ctx context.Context, data dbmodel.IDbModelMongo) error {

	var xError error = nil

	xWhereMap := make(map[string]interface{})
	xWhereMap["data_id"] = data.GetDataId()

	xError = this.DeleteAppData(ctx, data, xWhereMap)

	return xError

}

//根据结构体获取索引列表
func (this *DbMongoSvc) getAppDataIndexs(data interface{}) []string {

	xIndexList := []string{}

	xDataType := reflect.TypeOf(data)
	xDataVal := reflect.ValueOf(data)

	//如果传过来的是指针
	if xDataType.Kind() == reflect.Ptr {
		xDataType = xDataType.Elem()
		xDataVal = xDataVal.Elem()
	}

	xFieldCount := xDataType.NumField()

	for i := 0; i < xFieldCount; i++ {

		xDataField := xDataType.Field(i)

		xDataFieldName, xDataFieldIsBson := xDataField.Tag.Lookup("bson")
		if !xDataFieldIsBson {
			continue
		}

		if strings.EqualFold(xDataFieldName, ",inline") {

			xDataFieldVal := xDataVal.Field(i)
			xSubIndexList := this.getAppDataIndexs(xDataFieldVal.Interface())
			if len(xSubIndexList) < 1 {
				continue
			}

			for _, xSubIndex := range xSubIndexList {
				xIndexList = append(xIndexList, xSubIndex)
			}

			continue
		}

		xDataIndexName, xDataIndexExists := xDataField.Tag.Lookup("index")
		if !xDataIndexExists {
			continue
		}

		if len(xDataIndexName) < 1 {
			continue
		}

		if strings.EqualFold(xDataFieldName, "data_id") || strings.EqualFold(xDataIndexName, "uniq") {
			xIndexItem := fmt.Sprintf("%s|%s", xDataFieldName, "uniq")
			xIndexList = append(xIndexList, xIndexItem)
			continue
		}

		xIndexItem := fmt.Sprintf("%s|%s", xDataFieldName, "normal")
		xIndexList = append(xIndexList, xIndexItem)
	}

	return xIndexList

}
