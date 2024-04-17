package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"log"
	"math/rand"
)

const MaxMsgSize = 256

var aes_cipher cipher.Block;
var key string = "962B4D8D0F032F2948D84590557373AB"

func init(){
    var err error
    aes_cipher, err = aes.NewCipher([]byte(key))
    if err != nil {
        log.Fatal(err, aes_cipher)
    }
}

func generateNewRandomIV() []byte{
    iv := make([]byte, aes.BlockSize)
    // melhorar isso aqui
    for i := 0; i < aes.BlockSize; i++ {
        iv[i] = byte(rand.Int()) 
    }
    return iv
}

func Encrypt(src string) (string, error) {
    if len(src) > MaxMsgSize {
        return "", fmt.Errorf("Input não pode ser maior que %d caracteres!", MaxMsgSize)
    }
    msg := make([]byte, MaxMsgSize)
    copy(msg, []byte(src))

    dec := make([]byte, MaxMsgSize + aes.BlockSize)

    iv := generateNewRandomIV()
    encrypter := cipher.NewCBCEncrypter(aes_cipher, iv)
    encrypter.CryptBlocks(dec[16:], msg)
    // copia o iv para o inicio do valor encriptado, os primeiros 16 bytes
    copy(dec[:16], iv)
    return string(dec), nil
}

func Decrypt(src []byte) string {
    decrypter := cipher.NewCBCDecrypter(aes_cipher, src[:16])
    dec := make([]byte, MaxMsgSize)
    decrypter.CryptBlocks(dec, src[16:])
    return string(dec)
}
