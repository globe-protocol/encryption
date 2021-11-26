package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
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

func (e *encryptionService) EncryptStr(str string) ([]byte, error) {
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

	val := aesGCM.Seal(nonce, nonce, []byte(str), nil)

	return val, nil
}

func (e *encryptionService) EncryptToInterface(eData interface{}) (map[string]interface{}, error) {
	object := reflect.ValueOf(eData)
	returnObj := map[string]interface{}{}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher, external aes package returned the following error: %s", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM, external cipher package returned the following error: %s", err)
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
			fieldName = object.Type().Field(i).Tag.Get("json")
		}

		if len(fieldName) == 0 {
			return nil, fmt.Errorf("the provided struct does not have either the json or bson tag, this is necessary for creating return map")
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
		return nil, fmt.Errorf("failed to convert to json, external json package returned the following error: %s", err)
	}

	return jsonBytes, nil
}

func (e *encryptionService) Decrypt(encryptedData interface{}, desiredOutput interface{}) (interface{}, error) {
	object := reflect.ValueOf(encryptedData)
	returnObj := reflect.New(reflect.ValueOf(desiredOutput).Type())

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher, external aes package returned the following error: %s", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM, external cipher package returned the following error: %s", err)
	}

	for i := 0; i < object.NumField(); i++ {
		encrypted := reflect.Indirect(returnObj).Type().Field(i).Tag.Get("encrypted")

		var decryptedStr string
		if encrypted != "false" {
			decryptedStr, err = e.getPlainText(object.Field(i).Bytes(), aesGCM)
			if err != nil {
				return nil, fmt.Errorf("failed to get text out of encrypted value, the following error occured: %s", err)
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

func (e *encryptionService) DecryptStr(b []byte) (string, error) {

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher, external aes package returned the following error: %s", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM, external cipher package returned the following error: %s", err)
	}

	decryptedStr, err := e.getPlainText(b, aesGCM)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt string, the following error occured: %s", err)
	}

	return decryptedStr, nil
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

func (e *encryptionService) EncryptByt(b []byte) ([]byte, error) {
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

	val := aesGCM.Seal(nonce, nonce, b, nil)

	return val, nil
}

func (e *encryptionService) DecryptByt(b []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher, external aes package returned the following error: %s", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM, external cipher package returned the following error: %s", err)
	}

	decryptedBytes, err := e.getPlainBytes(b, aesGCM)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt string, the following error occured: %s", err)
	}

	return decryptedBytes, nil
}

func (e encryptionService) getPlainBytes(val []byte, aesGCM cipher.AEAD) ([]byte, error) {
	nonceSize := aesGCM.NonceSize()

	nonce, ciphertext := val[:nonceSize], val[nonceSize:]
	plainbytes, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plainbytes, nil
}
