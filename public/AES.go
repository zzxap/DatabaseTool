package public

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	//"unicode/utf8"

	//"encoding/hex"

	//"fmt" xsbnal666
	"fmt"
	"strings"
)

//https://studygolang.com/articles/4752
//http://www.cnblogs.com/lavin/p/5373188.html
//var aeskey = "Ap@12Sd#32Cv$12Z"

var aeskey = "A1234567A9abc8ef" // "WE3456@#8765432WR234567887654abc"

/*
 *src 要加密的字符串
 *key 用来加密的密钥 密钥长度可以是128bit、192bit、256bit中的任意一个 128bit(byte[16])、192bit(byte[24])、256bit(byte[32])
 *16位key对应128bit
 */
// src := "0.56"
//key := "A1234567A9abc8ef"

func AesEncrypt(content []byte) (byter []byte, err error) {
	block, err := aes.NewCipher([]byte(aeskey))
	if err != nil {
		fmt.Println("key error1", err)
		return nil, err
	}

	ecb := NewECBEncrypter(block)
	//content := []byte(src)
	content = PKCS7Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	// 普通base64编码加密 区别于urlsafe base64
	//fmt.Println("base64 result:", base64.StdEncoding.EncodeToString(crypted))

	//fmt.Println("base64UrlSafe result:", Base64UrlSafeEncode(crypted))
	return crypted, nil
}
func AesEncryptStr(str string) string {

	content := []byte(str)
	block, err := aes.NewCipher([]byte(aeskey))
	if err != nil {
		fmt.Println("key error1", err)
		return ""
	}

	ecb := NewECBEncrypter(block)
	//content := []byte(src)
	content = PKCS7Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	// 普通base64编码加密 区别于urlsafe base64
	//fmt.Println("base64 result:", base64.StdEncoding.EncodeToString(crypted))

	//fmt.Println("base64UrlSafe result:", Base64UrlSafeEncode(crypted))
	return base64.StdEncoding.EncodeToString(crypted)
}
func AesDecryptStr(src string) string {
	decdata, errr := base64.StdEncoding.DecodeString(src)
	if errr != nil {

		Log(errr.Error())
		return ""
	}
	data, err := AesDecrypt(decdata)

	if err == nil {
		return string(data[:])
	} else {
		Log(err.Error())
	}
	return ""
}
func AesDecryptStr2(str string) []byte {
	//key只能是 16 24 32长度

	crypted := []byte(str)
	block, err := aes.NewCipher([]byte(aeskey))

	//fmt.Println("block.BlockSize() is:", block.BlockSize())
	if err != nil {
		fmt.Println("err is:", err)
		return nil
	}
	// 长度不能小于aes.Blocksize
	if len(str) < block.BlockSize() {
		fmt.Println("crypto/cipher: ciphertext too short")

		return nil
	}
	blockMode := NewECBDecrypter(block)

	origData := make([]byte, len(crypted))

	blockMode.CryptBlocks(origData, crypted)
	fmt.Println("base64.StdEncoding1 is:", base64.StdEncoding.EncodeToString(origData))
	origData = PKCS7UnPadding(origData, block.BlockSize())

	fmt.Println("base64.StdEncoding2 is:", base64.StdEncoding.EncodeToString(origData))

	return origData //ConverToStr(origData)string(origData[:]) //
}
func base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func base64Decode(str string) (string, bool) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", true
	}
	return string(data), false
}
func Base64URLDecode(data string) ([]byte, error) {
	var missing = (4 - len(data)%4) % 4
	data += strings.Repeat("=", missing)
	//res, err := base64.URLEncoding.DecodeString(data)
	//fmt.Println("  decodebase64urlsafe is :", string(res), err)
	return base64.URLEncoding.DecodeString(data)

}

func Base64UrlSafeEncode(source []byte) string {
	// Base64 Url Safe is the same as Base64 but does not contain '/' and '+' (replaced by '_' and '-') and trailing '=' are removed.
	bytearr := base64.StdEncoding.EncodeToString(source)
	safeurl := strings.Replace(string(bytearr), "/", "_", -1)
	safeurl = strings.Replace(safeurl, "+", "-", -1)
	safeurl = strings.Replace(safeurl, "=", "", -1)
	return safeurl
}

func AesDecrypt(crypted []byte) ([]byte, error) {

	block, err := aes.NewCipher([]byte(aeskey))
	if err != nil {
		//fmt.Println("err is:", err)
		return nil, err
	}
	blockMode := NewECBDecrypter(block)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData, block.BlockSize())
	//fmt.Println("source is :", origData, string(origData))
	return origData, nil
}
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(plantText []byte, blockSize int) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])

	if unpadding > blockSize {
		fmt.Println("PKCS7UnPadding length error")
		return plantText
	} else {
		return plantText[:(length - unpadding)]

	}

}
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}
func (x *ecbEncrypter) BlockSize() int { return x.blockSize }
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		Log("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		Log("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

// NewECBDecrypter returns a BlockMode which decrypts in electronic code book
// mode, using the given Block.
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}
func (x *ecbDecrypter) BlockSize() int { return x.blockSize }
func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		Log("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		Log("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}
