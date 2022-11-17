package cheerlib

import (
	"errors"
	"reflect"
)

func FlatStructCopy(srcData interface{}, desData interface{}) error {

	var xErr error = nil

	if srcData == nil || srcData == nil {
		xErr = errors.New("srcData or srcData can not be nil.")
		return xErr
	}

	if reflect.TypeOf(srcData).Kind() != reflect.Ptr {
		xErr = errors.New("srcData must be Ptr.")
		return xErr
	}

	if reflect.TypeOf(desData).Kind() != reflect.Ptr {
		xErr = errors.New("desData must be Ptr.")
		return xErr
	}

	xSrcJson := TextStructToJson(srcData)
	xErr = TextStructFromJson(desData, xSrcJson)

	return xErr

}

func FlatMapToStruct(srcDataMap map[string]interface{}, desData interface{}) error {

	var xErr error = nil

	if desData == nil {
		xErr = errors.New("desData can not be nil.")
		return xErr
	}

	if reflect.TypeOf(desData).Kind() != reflect.Ptr {
		xErr = errors.New("desData must be Ptr.")
		return xErr
	}

	xSrcJson := TextStructToJson(srcDataMap)
	xErr = TextStructFromJson(desData, xSrcJson)

	return xErr
}
