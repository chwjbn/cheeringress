package cheerlib

import (
	"errors"
	"fmt"
	"github.com/alexmullins/zip"
	"io"
	"io/ioutil"
	"os"
	"path"
)

func FileZip(desZipPath string, srcFilePath string, passwd string) error {

	var xError error = nil

	xZipFile, xZipErr := os.Create(desZipPath)
	if xZipErr != nil {
		xError = errors.New(fmt.Sprintf("FileZip ZipError=[%s]", xZipErr.Error()))
		return xError
	}

	xZipHandle := zip.NewWriter(xZipFile)

	defer func() {
		xZipHandle.Close()
	}()

	xZipWriter, xZipErr := xZipHandle.Encrypt(path.Base(srcFilePath), passwd)
	if xZipErr != nil {
		xError = errors.New(fmt.Sprintf("FileZip ZipError=[%s]", xZipErr.Error()))
		return xError
	}

	xSrcFileData, xZipErr := ioutil.ReadFile(srcFilePath)
	if xZipErr != nil {
		xError = errors.New(fmt.Sprintf("FileZip ZipError=[%s]", xZipErr.Error()))
		return xError
	}

	xZipWriter.Write(xSrcFileData)
	xZipHandle.Flush()

	return xError

}

func FileExists(path string) bool {

	bRet := false

	xFileInfo, xFileInfoErr := os.Stat(path)
	if xFileInfoErr != nil {
		return bRet
	}

	if !xFileInfo.IsDir() {
		bRet = true
	}

	return bRet

}

func FileDelete(path string) bool {

	bRet := false

	xErr := os.Remove(path)
	if xErr == nil {
		bRet = true
	}

	return bRet

}

func FileRename(oldPath string, newPath string) bool {

	bRet := false

	xErr := os.Rename(oldPath, newPath)
	if xErr == nil {
		bRet = true
	}

	return bRet

}

func FileCopy(srcPath string, desPath string) bool {

	bRet := false

	if !FileExists(srcPath) {
		LogError("FileCopy srcPath Not Exists")
		return bRet
	}

	if FileExists(desPath) {
		LogError("FileCopy desPath Exists")
		return bRet
	}

	xDesFile, xDesFileErr := os.Create(desPath)
	if xDesFileErr != nil {
		LogError("FileCopy Create DesFile Error:" + xDesFileErr.Error())
		return bRet
	}

	defer xDesFile.Close()

	xSrcFile, xSrcFileErr := os.Open(srcPath)

	if xSrcFileErr != nil {
		LogError("FileCopy Open SrcFile Error:" + xSrcFileErr.Error())
		return bRet
	}

	defer xSrcFile.Close()

	_, xCopyErr := io.Copy(xDesFile, xSrcFile)

	if xCopyErr != nil {
		LogError("FileCopy Error:" + xCopyErr.Error())
		return bRet
	}

	bRet = true

	return bRet

}

func FileReadAllText(path string) string {

	sData := ""

	xFileData, xFileDataErr := ioutil.ReadFile(path)
	if xFileDataErr != nil {
		return sData
	}

	sData = string(xFileData)

	return sData

}

func FileWriteAllText(path string, data string) bool {

	bRet := false

	xFileDataErr := ioutil.WriteFile(path, []byte(data), 0644)
	if xFileDataErr == nil {
		bRet = true
	}

	return bRet

}
