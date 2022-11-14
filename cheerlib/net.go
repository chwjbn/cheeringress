package cheerlib

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func NetHttpDo(method string, url string, data string, contentType string, authData string) (error, string) {

	var xErr error = nil
	var xData = ""

	dataReader := bytes.NewReader([]byte(data))

	xReq, xReqErr := http.NewRequest(method, url, dataReader)
	if xReqErr != nil {
		xErr = errors.New(fmt.Sprintf("Net_HttpDo url=%s xReqErr=%s", url, xReqErr.Error()))
		return xErr, xData
	}

	xReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36")

	if len(contentType) > 0 {
		xReq.Header.Set("Content-Type", contentType)
	}

	if len(authData) > 0 {
		xReq.Header.Set("Authorization", authData)
	}

	xHttpClient := http.Client{}
	xResp, xRespErr := xHttpClient.Do(xReq)

	if xRespErr != nil {
		xErr = errors.New(fmt.Sprintf("Net_HttpDo url=%s xRespErr=%s", url, xRespErr.Error()))
		return xErr, xData
	}

	xRespBody, xRespBodyErr := ioutil.ReadAll(xResp.Body)

	if xRespBodyErr != nil {
		xErr = errors.New(fmt.Sprintf("Net_HttpDo url=%s xBodyErr=%s", url, xRespBodyErr.Error()))
		return xErr, xData
	}

	xData = TextGetString(xRespBody, unicode.UTF8)

	return xErr, xData
}

func NetHttpGet(url string) (error, string) {

	var xErr error = nil
	var xData = ""

	xErr, xData = NetHttpDo("GET", url, "", "", "")

	return xErr, xData
}

func NetHttpGetWithAuth(url string, authData string) (error, string) {

	var xErr error = nil
	var xData = ""

	xErr, xData = NetHttpDo("GET", url, "", "", authData)

	return xErr, xData
}

func NetHttpPostJson(url string, jsonData string, authData string) (error, string) {

	var xErr error = nil
	var xData = ""

	xErr, xData = NetHttpDo("POST", url, jsonData, "application/json;charset=UTF-8", authData)

	return xErr, xData

}

func NetHttpPostForm(url string, formData string, authData string) (error, string) {

	var xErr error = nil
	var xData = ""

	xErr, xData = NetHttpDo("POST", url, formData, "application/x-www-form-urlencoded", authData)

	return xErr, xData

}

func NetHttpDownloadFile(url string, filePath string) (error, int64) {

	var xErr error = nil
	var xData int64 = 0

	xResp, xRespErr := http.Get(url)

	if xRespErr != nil {
		xErr = xRespErr
		return xErr, xData
	}

	xFile, xFileCreateErr := os.Create(filePath)

	if xFileCreateErr != nil {
		xErr = xFileCreateErr
		return xErr, xData
	}

	defer xFile.Close()

	xFileLen, xFileCopyErr := io.Copy(xFile, xResp.Body)

	if xFileCopyErr != nil {
		xErr = xFileCopyErr
		return xErr, xData
	}

	xData = xFileLen

	return xErr, xData

}
