package cheerlib

import "os"

func DirectoryExists(path string) bool {

	bRet := false

	xFileInfo, xFileInfoErr := os.Stat(path)
	if xFileInfoErr != nil {
		return bRet
	}

	if xFileInfo.IsDir() {
		bRet = true
	}

	return bRet

}

func DirectoryCreateDirectory(path string) bool {

	bRet := false

	xMkdirErr := os.MkdirAll(path, 0755)
	if xMkdirErr == nil {
		bRet = true
	}

	return bRet
}

func DirectoryDeleteDirectory(path string) bool {

	bRet := false

	xErr := os.Remove(path)
	if xErr == nil {
		bRet = true
	}

	return bRet
}
