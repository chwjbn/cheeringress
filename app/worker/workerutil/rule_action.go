package workerutil

import (
	"errors"
	"fmt"
	"github.com/chwjbn/cheeringress/cheerapp"
	"github.com/chwjbn/cheeringress/cheerlib"
	"github.com/gin-gonic/gin"
	"path"
	"strings"
)

func ActionShowErrorPage(ctx *gin.Context, httpStatus int, errorCode string, errorMessage string) {

	xPageContent := "PCFET0NUWVBFIGh0bWwgUFVCTElDICItLy9XM0MvL0RURCBYSFRNTCAxLjAgVHJhbnNpdGlvbmFsLy9FTiIgImh0dHA6Ly93d3cudzMub3JnL1RSL3hodG1sMS9EVEQveGh0bWwxLXRyYW5zaXRpb25hbC5kdGQiPgo8aHRtbCB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMTk5OS94aHRtbCI+CjxoZWFkPgo8bWV0YSBodHRwLWVxdWl2PSJDb250ZW50LVR5cGUiIGNvbnRlbnQ9InRleHQvaHRtbDsgY2hhcnNldD11dGYtOCIgLz4KPHRpdGxlPkNoZWVySW5ncmVzczwvdGl0bGU+CjxzdHlsZSB0eXBlPSJ0ZXh0L2NzcyI+CmJvZHkgewoJbWFyZ2luOiAwcHg7CglwYWRkaW5nOiAwcHg7Cglmb250LWZhbWlseTogIuW+rui9r+mbhem7kSIsIEFyaWFsLCAiVHJlYnVjaGV0IE1TIiwgVmVyZGFuYSwgR2VvcmdpYSwgQmFza2VydmlsbGUsIFBhbGF0aW5vLCBUaW1lczsKCWZvbnQtc2l6ZTogMTZweDsKfQpkaXYgewoJbWFyZ2luLWxlZnQ6IGF1dG87CgltYXJnaW4tcmlnaHQ6IGF1dG87Cn0KaDEsIGgyLCBoMywgaDQgewoJbWFyZ2luOiAwOwoJZm9udC13ZWlnaHQ6IG5vcm1hbDsKCWZvbnQtZmFtaWx5OiAi5b6u6L2v6ZuF6buRIiwgQXJpYWwsICJUcmVidWNoZXQgTVMiLCBIZWx2ZXRpY2EsIFZlcmRhbmE7Cn0KaDEgewoJZm9udC1zaXplOiAyMnB4OwoJY29sb3I6ICMwMTg4REU7CglwYWRkaW5nOiAyMHB4IDBweCAxMHB4IDBweDsKfQpoMiB7Cgljb2xvcjogIzAxODhERTsKCWZvbnQtc2l6ZTogMTZweDsKCXBhZGRpbmc6IDEwcHggMHB4IDQwcHggMHB4Owp9CiNwYWdlIHsKCXdpZHRoOiA4MCU7CglwYWRkaW5nOiAyMHB4IDIwcHggMjBweCAyMHB4OwoJbWFyZ2luLXRvcDogMjBweDsKCWJvcmRlci1zdHlsZTogZGFzaGVkOwoJYm9yZGVyLWNvbG9yOiAjZTRlNGU0OwoJbGluZS1oZWlnaHQ6IDMwcHg7Cn0KPC9zdHlsZT4KPC9oZWFkPgo8Ym9keT4KPGRpdiBpZD0icGFnZSI+CiAgPGgxPuW9k+WJjeiuv+mXrueahOacjeWKoeS4jeWPr+eUqDwvaDE+CiAgPHA+PGZvbnQgY29sb3I9IiM2NjY2NjYiPumUmeivr+egge+8mltbRUNPREVdXTwvZm9udD48L3A+CiAgPHA+PGZvbnQgY29sb3I9IiM2NjY2NjYiPumUmeivr+S/oeaBr++8mltbRU1TR11dPC9mb250PjwvcD4KPC9kaXY+CjwvYm9keT4KPC9odG1sPg=="

	xPageContent = cheerlib.EncryptBase64Decode(xPageContent)

	xPageContent = strings.ReplaceAll(xPageContent, "[[ECODE]]", errorCode)
	xPageContent = strings.ReplaceAll(xPageContent, "[[EMSG]]", errorMessage)

	ctx.Header("Content-Type", "text/html")
	ctx.String(httpStatus, xPageContent)
}

func ActionFetchHttpResouce(ctx *gin.Context,url string) (string,error) {

	xSpan := cheerapp.SpanBeginBizFunction(ctx.Request.Context(), "workerutil.ActionFetchHttpResouce")
	defer func() {
		cheerapp.SpanEnd(xSpan)
	}()

	cheerapp.SpanTag(xSpan,"url",url)

	xFileRoot:=path.Join(cheerlib.ApplicationBaseDirectory(),"data","res")
	if !cheerlib.DirectoryExists(xFileRoot){
		cheerlib.DirectoryCreateDirectory(xFileRoot)
	}

	xFilePath:=path.Join(xFileRoot,cheerlib.EncryptMd5(url))

	cheerapp.SpanTag(xSpan,"filepath",xFilePath)

	if cheerlib.FileExists(xFilePath){
		return xFilePath,nil
	}

	xError,_:=cheerlib.NetHttpDownloadFile(url,xFilePath)
	if xError!=nil{
		return "", errors.New(fmt.Sprintf("Download Url=[%s] To FilePath=[%s] Error=[%s]",url,xFilePath,xError.Error()))
	}

	return xFilePath,nil
}