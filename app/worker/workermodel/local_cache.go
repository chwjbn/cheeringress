package workermodel

import (
	"github.com/chwjbn/cheeringress/cheerlib"
	"errors"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"path"
)

type LocalCache struct {
	mDbHandle *leveldb.DB
}

func NewLocalCache() (error, *LocalCache) {

	xThis := new(LocalCache)
	xError := xThis.intiCache()

	if xError != nil {
		return xError, nil
	}

	return xError, xThis
}

func (this *LocalCache) SetData(dataKey string, dataVal string) error {

	var xError error

	xOpt := opt.WriteOptions{}
	xOpt.Sync = true
	xOpt.NoWriteMerge = false

	xDbError := this.mDbHandle.Put([]byte(dataKey), []byte(dataVal), &xOpt)

	if xDbError != nil {
		xError = errors.New(fmt.Sprintf("leveldb.DbHandle.Put Error=[%s]", xDbError.Error()))
		return xError
	}

	return xError
}

func (this *LocalCache) GetData(dataKey string) (error, string) {

	var xError error
	var xData string

	xDataBuf, xDbError := this.mDbHandle.Get([]byte(dataKey), nil)
	if xDbError != nil {
		xError = errors.New(fmt.Sprintf("leveldb.DbHandle.Get Error=[%s]", xDbError.Error()))
		return xError, xData
	}

	xData = string(xDataBuf)

	return xError, xData
}

func (this *LocalCache) DelData(dataKey string) error {

	var xError error

	xDbError := this.mDbHandle.Delete([]byte(dataKey), nil)
	if xDbError != nil {
		xError = errors.New(fmt.Sprintf("leveldb.DbHandle.Delete Error=[%s]", xDbError.Error()))
		return xError
	}

	return xError
}

func (this *LocalCache) intiCache() error {

	var xError error

	xDbPath := path.Join(cheerlib.ApplicationBaseDirectory(), "data")
	if !cheerlib.DirectoryExists(xDbPath) {
		cheerlib.DirectoryCreateDirectory(xDbPath)
	}

	xDbPath = path.Join(xDbPath, "cheer.dat")

	var xDbError error
	this.mDbHandle, xDbError = leveldb.OpenFile(xDbPath, nil)

	if xDbError != nil {
		xError = errors.New(fmt.Sprintf("leveldb.OpenFile Error=[%s]", xDbError.Error()))
		return xError
	}

	return xError

}
