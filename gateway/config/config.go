package config

import (
	"fmt"
	"github.com/chwjbn/cheeringress/cheerlib"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type Config interface {
	Check() error
}

var configFileName string

func ParseConfigData(configData []byte, cfg Config) error {
	var xError error = nil
	xError = yaml.Unmarshal(configData, cfg)
	return xError

}

func PrintConfig(cfg Config) error {

	var xError error = nil

	xDataVal := ""

	xType := reflect.TypeOf(cfg).Elem()
	xVal := reflect.ValueOf(cfg).Elem()

	for i := 0; i < xType.NumField(); i++ {

		xTypeField := xType.Field(i)
		xTypeTag := xTypeField.Tag.Get("yaml")
		xValField := xVal.FieldByName(xTypeField.Name)

		xDataVal = fmt.Sprintf("%s\r\n[%s@%s]=[%v]", xDataVal, xTypeField.Name, xTypeTag, xValField)
	}

	cheerlib.LogInfo(fmt.Sprintf("\r\n+============================================+\r\nPrintConfig [%s]:%s\r\n+============================================+", xType.Name(), xDataVal))

	return xError

}

func ReadConfigFromEnv(cfg Config) error {

	var xError error = nil

	xType := reflect.TypeOf(cfg).Elem()
	xVal := reflect.ValueOf(cfg).Elem()

	for i := 0; i < xType.NumField(); i++ {

		xTypeField := xType.Field(i)
		xTypeTag := xTypeField.Tag.Get("yaml")

		if len(xTypeTag) < 1 {
			continue
		}

		xEnvKey := fmt.Sprintf("cheerenv_%s", xTypeTag)

		xEnvVal, xEnvHave := os.LookupEnv(xEnvKey)
		if !xEnvHave {
			continue
		}

		if !xVal.FieldByName(xTypeField.Name).CanSet() {
			continue
		}

		xValKind := xVal.FieldByName(xTypeField.Name).Kind().String()

		if strings.EqualFold(xValKind, "int") {
			xEnvValInt, xEnvValErr := strconv.ParseInt(xEnvVal, 10, 64)

			if xEnvValErr != nil {
				continue
			}

			xVal.FieldByName(xTypeField.Name).SetInt(xEnvValInt)

			continue
		}

		if strings.EqualFold(xValKind, "float") {

			xEnvValFloat, xEnvValErr := strconv.ParseFloat(xEnvVal, 64)

			if xEnvValErr != nil {
				continue
			}

			xVal.FieldByName(xTypeField.Name).SetFloat(xEnvValFloat)

			continue
		}

		if strings.EqualFold(xValKind, "bool") {

			xEnvValBool, xEnvValErr := strconv.ParseBool(xEnvVal)

			if xEnvValErr != nil {
				continue
			}

			xVal.FieldByName(xTypeField.Name).SetBool(xEnvValBool)

			continue
		}

		xVal.FieldByName(xTypeField.Name).SetString(xEnvVal)
	}

	return xError

}

func ReadConfigFromFile(filePath string, cfg Config) error {

	var xError error = nil
	var xFilePath = ""

	var xConfigData []byte

	xFilePath, xError = filepath.Abs(filePath)
	if xError != nil {
		return xError
	}

	xConfigData, xError = ioutil.ReadFile(xFilePath)

	if xError != nil {
		return xError
	}

	configFileName = xFilePath

	xError = ParseConfigData(xConfigData, cfg)

	if xError != nil {
		return xError
	}

	cheerlib.LogInfo(fmt.Sprintf("ReadConfigFile From [%s] Values(%v)", xFilePath, cfg))

	return xError

}

func SaveConfigFile(cfg Config) error {

	xConfigData, xError := yaml.Marshal(cfg)
	if xError != nil {
		return xError
	}

	xError = ioutil.WriteFile(configFileName, xConfigData, 0755)
	if xError != nil {
		return xError
	}

	cheerlib.LogInfo(fmt.Sprintf("SaveConfigFile To [%s] Values(%s)", configFileName, cfg))

	return nil
}
