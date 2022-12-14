package cheerlib

import (
	"fmt"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type CheerLog struct {
	mLogger *zap.Logger
}

func NewCheerLogOnFile() *CheerLog {

	var xError error = nil
	var xThis *CheerLog = nil

	xThis = new(CheerLog)

	xError, xInfoLogWriter := xThis.getLogFileWriter(fmt.Sprintf("%s_info", ApplicationFileName()))
	if xError != nil {
		StdError(fmt.Sprintf("NewCheerLog CheerLog.getLogWriter Error=[%s]", xError.Error()))
	}

	xError, xWarnLogWriter := xThis.getLogFileWriter(fmt.Sprintf("%s_warn", ApplicationFileName()))
	if xError != nil {
		StdError(fmt.Sprintf("NewCheerLog CheerLog.getLogWriter Error=[%s]", xError.Error()))
	}

	xError, xErrorLogWriter := xThis.getLogFileWriter(fmt.Sprintf("%s_error", ApplicationFileName()))
	if xError != nil {
		StdError(fmt.Sprintf("NewCheerLog CheerLog.getLogWriter Error=[%s]", xError.Error()))
	}

	xError = xThis.InitWithLogWriter(xInfoLogWriter, xWarnLogWriter, xErrorLogWriter)

	if xError != nil {
		StdError(fmt.Sprintf("NewCheerLog CheerLog.InitWithLogWriter Error=[%s]", xError.Error()))
	}

	return xThis

}

func (this *CheerLog) LogInfo(logContent string) {

	if this.mLogger == nil {
		return
	}

	this.mLogger.Info(logContent)
}

func (this *CheerLog) LogWarn(logContent string) {

	if this.mLogger == nil {
		return
	}

	this.mLogger.Warn(logContent)
}

func (this *CheerLog) LogError(logContent string) {

	if this.mLogger == nil {
		return
	}

	this.mLogger.Error(logContent)
}

func (this *CheerLog) InitWithLogWriter(logWriterForInfo zapcore.WriteSyncer, logWriterForWarn zapcore.WriteSyncer, logWriterForError zapcore.WriteSyncer) error {

	var xError error = nil

	xEncoder := this.getLogEncoder()

	xLogInfoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})

	xLogWarnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})

	xLogErrorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	xLogCore := zapcore.NewTee(
		zapcore.NewCore(xEncoder, logWriterForInfo, xLogInfoLevel),
		zapcore.NewCore(xEncoder, logWriterForWarn, xLogWarnLevel),
		zapcore.NewCore(xEncoder, logWriterForError, xLogErrorLevel),
	)

	this.mLogger = zap.New(xLogCore,
		zap.AddCaller(),
		zap.AddCallerSkip(2),
	)

	// ????????????ZAP
	if this.mLogger != nil {
		zap.ReplaceGlobals(this.mLogger)
	}

	return xError
}

func (this *CheerLog) getEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

func (this *CheerLog) getLogEncoder() zapcore.Encoder {

	xConfig := zap.NewProductionEncoderConfig()
	xConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	xConfig.EncodeLevel = this.getEncodeLevel

	xEncoder := zapcore.NewConsoleEncoder(xConfig)

	return xEncoder
}

func (this *CheerLog) getLogFileWriter(logFileName string) (error, zapcore.WriteSyncer) {

	var xError error = nil

	xLogFilePath := fmt.Sprintf("%s/log", ApplicationBaseDirectory())
	if !DirectoryExists(xLogFilePath) {
		DirectoryCreateDirectory(xLogFilePath)
	}

	xLogFilePath = fmt.Sprintf("%s/%s", xLogFilePath, logFileName)

	xLogFileWriter, xError := rotatelogs.New(
		xLogFilePath+"_%Y%m%d%H.log",                 //??????????????????
		rotatelogs.WithLinkName(xLogFilePath+".log"), //????????????,????????????????????????
		rotatelogs.WithRotationTime(time.Minute),     //?????????1?????????????????????60s  ??????1????????????1?????????
		rotatelogs.WithMaxAge(gCheerLogFileMaxAge),   //????????????
		rotatelogs.WithClock(rotatelogs.Local),       //????????????
	)

	if xError != nil {
		return xError, nil
	}

	xWriteSyncer := zapcore.AddSync(xLogFileWriter)

	return xError, xWriteSyncer

}

func (this *CheerLog) getStdErrorWriter() (error, zapcore.WriteSyncer) {

	var xError error = nil

	xWriteSyncer := zapcore.AddSync(os.Stderr)

	return xError, xWriteSyncer
}

func (this *CheerLog) getStdOutWriter() (error, zapcore.WriteSyncer) {

	var xError error = nil

	xWriteSyncer := zapcore.AddSync(os.Stdout)

	return xError, xWriteSyncer
}
