package cheerlib

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"
)

func EncryptMd5(data string) string {

	xHash := md5.New()
	xHash.Write([]byte(data))

	return hex.EncodeToString(xHash.Sum(nil))
}

func EncryptFileMd5(filePath string) string {

	xData := ""

	if !FileExists(filePath) {
		return xData
	}

	xFileData, xFileErr := ioutil.ReadFile(filePath)
	if xFileErr != nil {
		return xData
	}

	xHash := md5.New()
	xHash.Write(xFileData)
	xData = hex.EncodeToString(xHash.Sum(nil))

	return xData
}

func EncryptBase64Encode(data string) string {

	xData := base64.StdEncoding.EncodeToString([]byte(data))
	return xData

}

func EncryptBase64Decode(data string) string {

	xData := ""

	xDataBuf, xDataErr := base64.StdEncoding.DecodeString(data)

	if xDataErr != nil {
		return xData
	}

	xData = string(xDataBuf)

	return xData

}

func EncryptUrlEncode(data string) string {

	xData := ""

	xData = url.QueryEscape(data)

	return xData

}

func EncryptUrlDecode(data string) string {

	xData := ""

	xDataTemp, xDataTempErr := url.QueryUnescape(data)

	if xDataTempErr != nil {
		return xData
	}

	xData = xDataTemp

	return xData

}

func EncryptNewId() string {

	xEnvStr := strings.Join(os.Environ(), "|")
	xEnvStr = EncryptMd5(xEnvStr)

	xEnvStr = fmt.Sprintf("%s-%s-%s-%s", xEnvStr, OSName(), OsHostName(), OsIPV4())
	xEnvStr = EncryptMd5(xEnvStr)

	xNewId := fmt.Sprintf("%d", time.Now().Nanosecond())

	rand.Seed(time.Now().Unix())
	xNewId = fmt.Sprintf("id-%s-%d", xNewId, rand.Int63())

	xNewId = EncryptMd5(xNewId + xEnvStr)

	return xNewId

}
