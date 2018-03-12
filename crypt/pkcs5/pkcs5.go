package pkcs5

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"

	"github.com/pkg/errors"
)

//Cipher implements the crypt.Cipher interface
type Cipher struct {
	Key []byte
}

//Encrypt implements the encrypt func of crypt.Cipher
func (pk Cipher) Encrypt(message string) (string, error) {
	crypt, err := encrypt([]byte(message), pk.Key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(crypt), nil
}

//Decrypt implements the decrypt func of crypt.Cipher
func (pk Cipher) Decrypt(message string) (string, error) {
	cmessage, err := hex.DecodeString(message)
	if err != nil {
		return "", err
	}
	crypt, err := decrypt(cmessage, pk.Key)
	if err != nil {
		return "", err
	}
	return string(crypt), nil
}

func encrypt(src []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, "encrypt key error1")
	}
	if len(src) == 0 {
		return nil, errors.New("encrypt plain content empty")
	}
	ecb := cipher.NewCBCEncrypter(block, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	src = pkcs5padding(src, block.BlockSize())
	crypted := make([]byte, len(src))
	ecb.CryptBlocks(crypted, src)

	return crypted, nil
}

func decrypt(crypt []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, "decrypt key cipher error")
	}
	if len(crypt) == 0 {
		return nil, errors.New("decrypt plain content empty")
	}
	ecb := cipher.NewCBCDecrypter(block, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	decrypted := make([]byte, len(crypt))
	ecb.CryptBlocks(decrypted, crypt)

	return pkcs5trim(decrypted), nil
}

func pkcs5padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func pkcs5trim(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
