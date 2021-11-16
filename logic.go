package AES

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
)

type encryptionService struct {
	key []byte
}

func NewEncryptionService(key []byte) EncryptionService {
	return &encryptionService{
		key: key,
	}
}

func (e *encryptionService) EncryptToInterface(eData interface{}) (map[string]interface{}, error) {
	object := reflect.ValueOf(eData)
	returnObj := map[string]interface{}{}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	for i := 0; i < object.NumField(); i++ {
		var fieldName string

		fieldName = object.Type().Field(i).Tag.Get("bson")
		if fieldName == "" || len(fieldName) == 0 {
			fieldName = object.Type().Field(i).Tag.Get("ename")
		}

		if len(fieldName) == 0 {
			return nil, fmt.Errorf("the provided struct does not have either the bson or ename tag, this is necessary for creating return map")
		}
		encrypt := object.Type().Field(i).Tag.Get("encrypted")

		if encrypt != "false" {
			val := aesGCM.Seal(nonce, nonce, []byte(fmt.Sprint(object.Field(i))), nil)
			returnObj[fieldName] = val
		} else {
			val := fmt.Sprint(object.Field(i))
			returnObj[fieldName] = val
		}

	}

	return returnObj, nil
}

func (e *encryptionService) EncryptToJSON(eData interface{}) ([]byte, error) {
	object := reflect.ValueOf(eData)
	returnObj := map[string]interface{}{}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	for i := 0; i < object.NumField(); i++ {
		var fieldName string

		fieldName = object.Type().Field(i).Tag.Get("bson")
		if fieldName == "" || len(fieldName) == 0 {
			fieldName = object.Type().Field(i).Tag.Get("ename")
		}

		if len(fieldName) == 0 {
			return nil, fmt.Errorf("the provided struct does not have either the bson or ename tag, this is necessary for creating return map")
		}
		encrypt := object.Type().Field(i).Tag.Get("encrypted")

		if encrypt != "false" {
			val := aesGCM.Seal(nonce, nonce, []byte(fmt.Sprint(object.Field(i))), nil)
			returnObj[fieldName] = val
		} else {
			val := fmt.Sprint(object.Field(i))
			returnObj[fieldName] = val
		}

	}

	jsonBytes, err := json.Marshal(returnObj)
	if err != nil {
		return nil, errors.New("could not convert interface to json")
	}

	return jsonBytes, nil
}

func (e *encryptionService) Decrypt(encryptedData interface{}, desiredOutput interface{}) (interface{}, error) {
	object := reflect.ValueOf(encryptedData)
	returnObj := reflect.New(reflect.ValueOf(desiredOutput).Type())

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, errors.New("failed to create cipher")
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("failed to create GCM")
	}

	for i := 0; i < object.NumField(); i++ {
		encrypted := reflect.Indirect(returnObj).Type().Field(i).Tag.Get("encrypted")

		var decryptedStr string
		if encrypted != "false" {
			decryptedStr, err = e.getPlainText(object.Field(i).Bytes(), aesGCM)
			if err != nil {
				return nil, err
			}
		} else {
			decryptedStr = object.Field(i).String()
		}

		field := reflect.Indirect(returnObj).Field(i)
		if field.IsValid() {
			val, err := Convert(decryptedStr, field)
			if err != nil {
				return nil, err
			}

			field.Set(reflect.ValueOf(val))
		}
	}

	return returnObj.Interface(), nil
}

func (e encryptionService) getPlainText(val []byte, aesGCM cipher.AEAD) (string, error) {
	nonceSize := aesGCM.NonceSize()

	nonce, ciphertext := val[:nonceSize], val[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", plaintext), nil
}
