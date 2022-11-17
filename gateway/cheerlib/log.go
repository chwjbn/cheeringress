package cheerlib

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	gCheerAppName        string = ApplicationFileName()
	gCheerAppDescription string = ApplicationFileName()

	gCheerLogFileMaxAge time.Duration = 24 * time.Hour
	gCheerLogger        *CheerLog     = NewCheerLogOnFile()
)

func SetGlobalAppInfo(appName string, appDescription string) {
	gCheerAppName = appName
	gCheerAppDescription = appDescription
}

func GetGlobalAppName() string {
	return gCheerAppName
}

func GetGlobalAppDescription() string {
	return gCheerAppDescription
}

func SetGlobalCheerLogFileMaxAge(maxAge time.Duration) {
	gCheerLogFileMaxAge = maxAge
}

func SetGlobalCheerLogger(cheerLog *CheerLog) {
	gCheerLogger = cheerLog
}

func LogInfo(logContent string) {

	if gCheerLogger == nil {
		StdInfo(logContent)
		return
	}

	gCheerLogger.LogInfo(logContent)

}

func LogWarn(logContent string) {

	if gCheerLogger == nil {
		StdError(logContent)
		return
	}

	gCheerLogger.LogWarn(logContent)
}

func LogError(logContent string) {

	if gCheerLogger == nil {
		StdError(logContent)
		return
	}

	gCheerLogger.LogError(logContent)

}

func StdError(logContent string) {
	logContent = strings.TrimSpace(logContent)
	os.Stderr.WriteString(fmt.Sprintf("[%s]%s\r\n", TimeGetNow(), logContent))
}

func StdInfo(logContent string) {
	logContent = strings.TrimSpace(logContent)
	os.Stdout.WriteString(fmt.Sprintf("[%s]%s\r\n", TimeGetNow(), logContent))
}
