package cheerlib

import (
	"encoding/json"
	"golang.org/x/text/encoding"
)

func TextStructFromJson(data interface{}, dataJson string) error {

	var xErr error = nil
	xErr = json.Unmarshal([]byte(dataJson), data)
	return xErr

}

func TextStructToJson(data interface{}) string {

	xData := "{}"

	jonData, jsonErr := json.Marshal(data)

	if jsonErr != nil {
		return xData
	}

	xData = string(jonData)

	return xData
}

func TextGetString(data []byte, encoding encoding.Encoding) string {

	sData := ""

	dataBuffer, xErr := encoding.NewDecoder().Bytes(data)

	if xErr != nil {
		return sData
	}

	sData = string(dataBuffer)

	return sData

}

func TextGetMapColumn(dataMapList []map[string]interface{}, colName string) []interface{} {

	dataList := []interface{}{}

	for _, dataMapItem := range dataMapList {

		xDataItem, bFind := dataMapItem[colName]
		if !bFind {
			continue
		}

		dataList = append(dataList, xDataItem)
	}

	return dataList
}
