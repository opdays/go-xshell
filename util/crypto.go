package util

import (
	"encoding/base64"
	"crypto/cipher"
	"bytes"
	"crypto/aes"
)

const AesKey  = "sfe023f_9fd&fwfl"

func UrlBase64Encrypt(input string) (output string) {
	output = base64.URLEncoding.EncodeToString([]byte(input))
	return
}
func UrlBase64Decrypt(input string) (string ,error) {
	outputBytes,err := base64.URLEncoding.DecodeString(string(input))
	if err != nil{
		return "",nil
	}
	return string(outputBytes),nil
}


//https://github.com/polaris1119/myblog_article_code/blob/master/aes/aes.go
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}



func SimpleDecrypt(text string) (result string) {
	resultBs, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return
	}
	resultBs, err = AesDecrypt([]byte(resultBs), []byte(AesKey))
	if err != nil {
		return
	}
	result = string(resultBs)
	return
}
func SimpleEncrypt(text string) (result string){
	resultBs, err := AesEncrypt([]byte(text), []byte(AesKey))
	if err != nil {
		return
	}
	result = base64.StdEncoding.EncodeToString(resultBs)
	return
}