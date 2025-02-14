package toolbox

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
)

type encrypt struct {
}

func (u *utils) NewEncrypt() *encrypt {
	return &encrypt{}
}

func (g *encrypt) Md5StringInt(num int) string {
	str := strconv.Itoa(num)
	return g.Md5String(str)
}

func (g *encrypt) Md5String(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

func (g *encrypt) Sha1(str string) string {
	s := sha1.New()
	s.Write([]byte(str))
	return hex.EncodeToString(s.Sum(nil))
}

func (g *encrypt) Sha256(str string) string {
	s := sha256.New()
	s.Write([]byte(str))
	return hex.EncodeToString(s.Sum(nil))
}

func (g *encrypt) HmacSha256(key string, data string) []byte {
	mac := hmac.New(sha256.New, []byte(key))
	_, _ = mac.Write([]byte(data))

	return mac.Sum(nil)
}

// HmacSha256ToHex 将加密后的二进制转16进制字符串
func (g *encrypt) HmacSha256ToHex(key string, data string) string {
	return hex.EncodeToString(g.HmacSha256(key, data))
}

// HmacSha256ToHex 将加密后的二进制转Base64字符串
func (g *encrypt) HmacSha256ToBase64(key string, data string) string {
	return base64.URLEncoding.EncodeToString(g.HmacSha256(key, data))
}

func (g *encrypt) EncodePwd(str string) string {
	return g.Md5String(g.Sha1(str))
}

/*
DES CBC加密
key的长度为8个字节， iv必须相同长度
*/
func (e *encrypt) EncryptDES_CBC_PKCS5_HEX(src, key, iv string) (str string, err error) {
	data := []byte(src)
	keyByte := []byte(key)
	block, err := des.NewCipher(keyByte)
	if err != nil {
		return
	}
	data = e.pkcs5Padding(data, block.BlockSize())
	//获取CBC加密模式
	//iv := keyByte //用密钥作为向量(不建议这样使用)
	ivByte := []byte(iv)
	mode := cipher.NewCBCEncrypter(block, ivByte)
	out := make([]byte, len(data))
	mode.CryptBlocks(out, data)
	str = fmt.Sprintf("%X", out)
	return
}

// DESC CBC解密
func (e *encrypt) DecryptDES_CBC_PKCS5_HEX(src, key, iv string) (str string, err error) {
	keyByte := []byte(key)
	data, err := hex.DecodeString(src)
	if err != nil {
		return
	}
	block, err := des.NewCipher(keyByte)
	if err != nil {
		return
	}
	ivBye := []byte(iv)
	mode := cipher.NewCBCDecrypter(block, ivBye)
	plaintext := make([]byte, len(data))
	mode.CryptBlocks(plaintext, data)
	plaintext = e.pkcs5UnPadding(plaintext)
	str = string(plaintext)
	return
}

// ECB加密
func (e *encrypt) EncryptDES_ECB_PKCS5_HEX(src, key string) (str string, err error) {
	data := []byte(src)
	keyByte := []byte(key)
	block, err := des.NewCipher(keyByte)
	if err != nil {
		return
	}
	bs := block.BlockSize()
	//对明文数据进行补码
	data = e.pkcs5Padding(data, bs)
	if len(data)%bs != 0 {
		panic("Need a multiple of the blocksize")
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		//对明文按照blocksize进行分块加密
		//必要时可以使用go关键字进行并行加密
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	str = fmt.Sprintf("%X", out)
	return
}

// ECB解密
func (e *encrypt) DecryptDES_ECB_PKCS5_HEX(src, key string) (str string, err error) {
	data, err := hex.DecodeString(src)
	if err != nil {
		return
	}
	keyByte := []byte(key)
	block, err := des.NewCipher(keyByte)
	if err != nil {
		return
	}
	bs := block.BlockSize()
	if len(data)%bs != 0 {
		err = errors.New("crypto/cipher: input not full blocks")
		return
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Decrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	out = e.pkcs5UnPadding(out)
	str = string(out)
	return
}

/*
key参数的长度 iv必须相同长度
16 字节 - AES-128
24 字节 - AES-192
32 字节 - AES-256
*/
func (e *encrypt) EncryptAES_CBC_PKCS5_HEX(src, key, iv string) (str string, err error) {
	data := []byte(src)
	keyByte := []byte(key)
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return
	}

	data = e.pkcs5Padding(data, block.BlockSize())
	//获取CBC加密模式
	//iv := keyByte //用密钥作为向量(不建议这样使用)
	ivByte := []byte(iv)
	mode := cipher.NewCBCEncrypter(block, ivByte)
	out := make([]byte, len(data))
	mode.CryptBlocks(out, data)
	str = fmt.Sprintf("%X", out)
	return
}

// AES CBC解密
func (e *encrypt) DecryptAES_CBC_PKCS5_HEXAndBase64(src, key, iv string) (str string, err error) {
	keyByte := []byte(key)
	data, err := hex.DecodeString(src)
	if err != nil {
		data, err = base64.StdEncoding.DecodeString(src)
		if err != nil {
			return
		}
	}

	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return
	}
	//iv := keyByte //用密钥作为向量(不建议这样使用)
	ivBye := []byte(iv)
	mode := cipher.NewCBCDecrypter(block, ivBye)
	plaintext := make([]byte, len(data))
	mode.CryptBlocks(plaintext, data)
	plaintext = e.pkcs5UnPadding(plaintext)
	str = string(plaintext)
	return
}

// 明文补码算法
func (e *encrypt) pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 明文减码算法
func (e *encrypt) pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func (e *encrypt) pkcs7Padding(data []byte, blockSize int) []byte {
	//判断缺少几位长度。最少1，最多 blockSize
	padding := blockSize - len(data)%blockSize
	//补足位数。把切片[]byte{byte(padding)}复制padding个
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// 明文减码算法
func (e *encrypt) pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	}
	//获取填充的个数
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

func (e *encrypt) EncryptAES_CBC_PKCS7_Base64(data, key, iv string) (str string, err error) {
	//创建加密实例
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return
	}
	//判断加密快的大小
	blockSize := block.BlockSize()
	//填充
	encryptBytes := e.pkcs7Padding([]byte(data), blockSize)
	//初始化加密数据接收切片
	crypted := make([]byte, len(encryptBytes))
	//使用cbc加密模式
	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	//执行加密
	blockMode.CryptBlocks(crypted, encryptBytes)
	str = string(crypted)
	return
}

func (e *encrypt) DecryptAES_CBC_PKCS7_Base64(data []byte, key, iv []byte) (crypted []byte, err error) {
	//创建实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	//获取块的大小
	//blockSize := block.BlockSize()
	//使用cbc
	blockMode := cipher.NewCBCDecrypter(block, iv)
	//初始化解密数据接收切片
	crypted = make([]byte, len(data))
	//执行解密
	blockMode.CryptBlocks(crypted, data)
	//去除填充
	crypted, err = e.pkcs7UnPadding(crypted)
	return
}

func (e *encrypt) DecryptAES_CBC_PKCS7_HEXAndBase64(dataByte string, key, iv []byte) (crypted []byte, err error) {
	data, err := hex.DecodeString(dataByte)
	if err != nil {
		data, err = base64.StdEncoding.DecodeString(dataByte)
		if err != nil {
			return
		}
	}
	//创建实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	//获取块的大小
	//blockSize := block.BlockSize()
	//使用cbc
	blockMode := cipher.NewCBCDecrypter(block, iv)
	//初始化解密数据接收切片
	crypted = make([]byte, len(data))
	//执行解密
	blockMode.CryptBlocks(crypted, data)
	//去除填充
	crypted, err = e.pkcs7UnPadding(crypted)
	return
}
