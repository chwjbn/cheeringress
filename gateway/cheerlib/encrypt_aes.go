package cheerlib

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

func pkcs7Padding(data []byte, blockSize int) []byte {

	padding := blockSize - len(data)%blockSize

	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(data, padText...)
}

func pkcs7UnPadding(data []byte) ([]byte, error) {

	length := len(data)

	if length == 0 {
		return nil, errors.New("invalid encrypt data!")
	}

	unPadding := int(data[length-1])

	return data[:(length - unPadding)], nil
}

func EncryptAesEncryptData(data []byte, key []byte) ([]byte, error) {
	//创建加密实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//判断加密快的大小
	blockSize := block.BlockSize()
	//填充
	encryptBytes := pkcs7Padding(data, blockSize)
	//初始化加密数据接收切片
	crypted := make([]byte, len(encryptBytes))
	//使用cbc加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	//执行加密
	blockMode.CryptBlocks(crypted, encryptBytes)
	return crypted, nil
}

func EncryptAesDecryptData(data []byte, key []byte) ([]byte, error) {
	//创建实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//使用cbc
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	//初始化解密数据接收切片
	crypted := make([]byte, len(data))
	//执行解密
	blockMode.CryptBlocks(crypted, data)
	//去除填充
	crypted, err = pkcs7UnPadding(crypted)
	if err != nil {
		return nil, err
	}
	return crypted, nil
}

func EncryptAesEncrypt(data []byte, aesKey string) (string, error) {

	aesKey = EncryptMd5(aesKey)

	res, err := EncryptAesEncryptData(data, []byte(aesKey))

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(res), nil
}

func EncryptAesDecrypt(data string, aesKey string) ([]byte, error) {

	aesKey = EncryptMd5(aesKey)

	dataByte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	return EncryptAesDecryptData(dataByte, []byte(aesKey))
}
