package cheerlib

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/pprof"
)

func ApplicationBaseDirectory() string {

	sRet := ""

	xFilePath, xFilePathErr := filepath.Abs(os.Args[0])
	if xFilePathErr != nil {
		return sRet
	}

	sRet = filepath.Dir(xFilePath)

	return sRet

}

func ApplicationFileName() string {

	sRet := ""

	xFilePath, xFilePathErr := filepath.Abs(os.Args[0])
	if xFilePathErr != nil {
		return sRet
	}

	sRet = filepath.Base(xFilePath)

	return sRet

}

func ApplicationFullPath() string {

	sRet := ""

	xFilePath, xFilePathErr := filepath.Abs(os.Args[0])
	if xFilePathErr != nil {
		return sRet
	}

	sRet = xFilePath

	return sRet

}

func ApplicationWriteHeapFile(fileName string) {

	xHeapFileName := fmt.Sprintf("%s/%s", ApplicationBaseDirectory(), fileName)
	xHeapFile, xHeapErr := os.OpenFile(xHeapFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if xHeapErr == nil {
		pprof.Lookup("heap").WriteTo(xHeapFile, 0)
		xHeapFile.Close()
	}

}
