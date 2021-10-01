package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"io"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

// getKey gets a 32 byte key from user provided password using SHA-256 hash.
func getKey(password, salt []byte) ([]byte, []byte) {

	// if salt is not provided then produce a random salt
	// Note: This should be nil only when encrypting, when
	// decrypting salt should already be provided
	if salt == nil {
		salt = make([]byte, 8)
		rand.Read(salt)
	}

	key := pbkdf2.Key(password, salt, 1000, 32, sha256.New)

	return key, salt
}

// Encrypt will encrypt plainData by producing a 32 byte key using AES-256
// and return a string with the following format:
// <plainDataLength>-<salt>-<iv>-<cipherData>.
func Encrypt(password, plainData []byte) (string, error) {

	var err error

	plainDataLength := make([]byte, 8)
	// Get the length of plainData and save it into a byte slice in little endian
	// format
	binary.LittleEndian.PutUint64(plainDataLength, uint64(len(plainData)))

	// If the plainData is not a multiple of aes.BlockSize then the plainData needs
	// to be padded
	if len(plainData)%aes.BlockSize != 0 {

		// calculate number of bytes to pad
		numBytesToPad := aes.BlockSize - (len(plainData) % aes.BlockSize)
		pad := make([]byte, numBytesToPad)

		// fill the pad with random data
		if _, err = io.ReadFull(rand.Reader, pad); err != nil {
			return "", err
		}

		// append the pad to the plainData so that it is a multiple of aes.BlockSize
		plainData = append(plainData, pad...)
	}

	key, salt := getKey(password, nil)

	var block cipher.Block
	if block, err = aes.NewCipher(key); err != nil {
		return "", err
	}

	iv := make([]byte, aes.BlockSize)
	// fill iv with random data
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cipherData := make([]byte, len(plainData))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherData, plainData)

	// return data by concatenating the encoding of all the fields as hex strings
	return hex.EncodeToString(plainDataLength) + "-" + hex.EncodeToString(salt) + "-" + hex.EncodeToString(iv) + "-" + hex.EncodeToString(cipherData), err
}

// Decrpyt will decrypt the ciphertext and return the original text that was decrypted.
func Decrypt(password, ciphertext []byte) ([]byte, error) {

	var err error

	// split on "-" which was meant to be used as a delimiter to extract the
	// indiviual fields
	savedData := strings.Split(string(ciphertext), "-")

	// data is saved in the following order:
	// plainDataLength, salt, iv and cipherData
	var plainDataLengthTmp, salt, iv, cipherData []byte

	if plainDataLengthTmp, err = hex.DecodeString(savedData[0]); err != nil {
		return nil, err
	}
	plainDataLength := binary.LittleEndian.Uint64(plainDataLengthTmp)

	if salt, err = hex.DecodeString(savedData[1]); err != nil {
		return nil, err
	}

	if iv, err = hex.DecodeString(savedData[2]); err != nil {
		return nil, err
	}

	if cipherData, err = hex.DecodeString(savedData[3]); err != nil {
		return nil, err
	}

	// have to use the extracted salt otherwise the produced key will not be
	// valid
	key, _ := getKey(password, salt)

	var block cipher.Block
	if block, err = aes.NewCipher(key); err != nil {
		return nil, err
	}

	plainData := make([]byte, len(cipherData))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plainData, cipherData)

	// when returning the plainData it needs to be sliced to remove any padding
	// that may have been added
	return plainData[:plainDataLength], err
}
