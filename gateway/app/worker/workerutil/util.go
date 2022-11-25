package workerutil

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/chwjbn/go4sky"
	h "github.com/chwjbn/go4sky/plugins/http"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func HttpDoData(tracer *go4sky.Tracer,method string,url string,playload []byte,header map[string]string) ([]byte,error)  {

	respData:=[]byte{}
	var xError error

	xClient,xClientErr:=h.NewClient(tracer)

	if xError!=nil{
		xError=errors.New(fmt.Sprintf("go4sky NewClient Error=[%s]",xClientErr.Error()))
		return respData,xError
	}

	xClient.Timeout=60*time.Second

	xReqDataReader := bytes.NewReader(playload)
	xReq, xReqErr := http.NewRequest(method, url, xReqDataReader)
	if xReqErr!=nil{
		xError=errors.New(fmt.Sprintf("http.NewRequest Error=[%s]",xReqErr.Error()))
		return respData,xError
	}

	xReq.Header.Set("User-Agent", "Mozilla/5.5 CheerIngress HttpClient")
	for k,v:=range header{
		xReq.Header.Set(k,v)
	}

	xResp,xRespErr:=xClient.Do(xReq)
	if xRespErr!=nil{
		xError=errors.New(fmt.Sprintf("HttpClient Do Error=[%s]",xRespErr.Error()))
		return respData,xError
	}

	if xResp.StatusCode>=400{
		xError=errors.New(fmt.Sprintf("HttpClient Do Return Error StatusCode=[%d]",xResp.StatusCode))
		return respData,xError
	}

	xRespBody,xRespErr:=ioutil.ReadAll(xResp.Body)
	if xRespErr!=nil{
		xError=errors.New(fmt.Sprintf("HttpClient Response Read Body Error=[%s]",xRespErr.Error()))
		return respData,xError
	}

	respData=xRespBody

    return respData,xError
}

func HttpDownloadFile(tracer *go4sky.Tracer,method string,url string,playload []byte,header map[string]string,filePath string) (int64,error)  {

	var xDataSize int64
	var xError error

	xFileHandle, xFileHandleError := os.Create(filePath)
	if xFileHandleError != nil {
		xError=errors.New(fmt.Sprintf("create local file=[%s] error=[%s]",filePath,xFileHandleError.Error()))
		return xDataSize,xError
	}

	defer func() {

		if xFileHandle!=nil{
			xFileHandle.Close()
		}

	}()

	xClient,xClientErr:=h.NewClient(tracer)

	if xError!=nil{
		xError=errors.New(fmt.Sprintf("go4sky NewClient Error=[%s]",xClientErr.Error()))
		return xDataSize,xError
	}

	xClient.Timeout=600*time.Second

	xReqDataReader := bytes.NewReader(playload)
	xReq, xReqErr := http.NewRequest(method, url, xReqDataReader)
	if xReqErr!=nil{
		xError=errors.New(fmt.Sprintf("http.NewRequest Error=[%s]",xReqErr.Error()))
		return xDataSize,xError
	}

	xReq.Header.Set("User-Agent", "Mozilla/5.5 CheerIngress HttpClient")
	for k,v:=range header{
		xReq.Header.Set(k,v)
	}

	xResp,xRespErr:=xClient.Do(xReq)
	if xRespErr!=nil{
		xError=errors.New(fmt.Sprintf("HttpClient Do Error=[%s]",xRespErr.Error()))
		return xDataSize,xError
	}

	if xResp.StatusCode>=400{
		xError=errors.New(fmt.Sprintf("HttpClient Do Return Error StatusCode=[%d]",xResp.StatusCode))
		return xDataSize,xError
	}

	xFileCopySize,xFileCopyErr:=io.Copy(xFileHandle,xResp.Body)
	if xFileHandleError!=nil{
		xError=errors.New(fmt.Sprintf("Copy Net Data To Local File error=[%s]",xFileCopyErr.Error()))
		return xDataSize,xError
	}

	xDataSize=xFileCopySize

	return xDataSize,xError
}

func HttpPostJson(tracer *go4sky.Tracer,url string, jsonData string, authData string)(string,error)  {

	var xRespJson string
	var xError error

	var xPlayload=[]byte(jsonData)

	var xHeader=make(map[string]string)
	xHeader["Authorization"]=authData
	xHeader["Content-Type"]="application/json; charset=utf-8"

	xRespData,xError:=HttpDoData(tracer,"POST",url,xPlayload,xHeader)
	if xError!=nil{
		return xRespJson,xError
	}

	xRespJson=string(xRespData)

	return xRespJson,xError
}