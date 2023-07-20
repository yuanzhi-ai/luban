package comm

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/yuanzhi-ai/luban/server/log"
)

// Md5Encode 对一个字符串进行md5加密
func Md5Encode(data string) string {
	has := md5.Sum([]byte(data))
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

// CBC模式的AES加密
func AesEncryptCBC(originData []byte, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()
	origData := pkcs5Padding(originData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encrypted := make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)
	return encrypted
}

// CBC模式解密
func AesDecryptCBC(encrypted []byte, key []byte) []byte {

	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	decrypted := make([]byte, len(encrypted))
	blockMode.CryptBlocks(decrypted, encrypted)
	decrypted = pkcs5UnPadding(decrypted)
	return decrypted
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// 电话号码aes对称加密
func AesEncrptyPhone(phone string) (string, error) {
	skeyInstance := GetSkeyInstance()
	phoneSkey, err := skeyInstance.GetSkey(PhoneSkey)
	if err != nil {
		log.Errorf("get phone skey err:%v", err)
		return "", err
	}
	aesPhone := AesEncryptCBC([]byte(phone), []byte(phoneSkey))
	return fmt.Sprintf("%x", aesPhone), nil
}

// 电话号码对称解密
func AesDecryptPhone(decPhone string) (string, error) {
	skeyInstance := GetSkeyInstance()
	phoneSkey, err := skeyInstance.GetSkey(PhoneSkey)
	if err != nil {
		log.Errorf("get phone skey err:%v", err)
		return "", err
	}
	bPhone, err := hex.DecodeString(decPhone)
	if err != nil {
		return "", fmt.Errorf("decode phone from hex fail. hex:%v err:%v", decPhone, err)
	}
	phone := AesDecryptCBC(bPhone, []byte(phoneSkey))
	return string(phone), nil
}

// 账号密码计算s2
func CalculateS2(phone string, pswd string) string {
	md5Pswd := Md5Encode(pswd)
	s2 := Md5Encode(phone + md5Pswd)
	return s2
}
