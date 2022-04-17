package encryption

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

//Create encryption service by passing a 32-bit key as parameter
func NewEncryptionService(key []byte) EncryptionService {
	return &encryptionService{
		key: key,
	}
}

//Helper functions to remove code duplication
func (e *encryptionService) initGCM() (cipher.AEAD, error) {
	//Create cipher with given key
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher, external aes package returned the following error: %s", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM, external cipher package returned the following error: %s", err)
	}

	return aesGCM, nil
}

func (e *encryptionService) sealWithNonce(aesGCM cipher.AEAD, str string) ([]byte, error) {
	//Create random nonce so that the same input value changes when it is encrypted
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	val := aesGCM.Seal(nonce, nonce, []byte(str), nil) //Encrypt using all values

	return val, nil
}

func (e *encryptionService) findFieldTag(object reflect.StructTag, fieldTagNames []string) (string, error) {
	for i := 0; i < len(fieldTagNames); i++ {
		field := object.Get(fieldTagNames[i])

		if field != "" && len(field) != 0 {
			return field, nil
		}
	}

	return "", errors.New("none of the necessary field tags were found in the structure could not map to interface")
}

//END

//Get encrypted []byte by inputting string
func (e *encryptionService) EncryptStr(str string) ([]byte, error) {
	aesGCM, err := e.initGCM()
	if err != nil {
		return nil, err
	}

	val, err := e.sealWithNonce(aesGCM, str)
	if err != nil {
		return nil, err
	}

	return val, nil
}

//Encrypt any 1-dimensional struct and get an interface ready for MongoDB storage
func (e *encryptionService) EncryptToInterface(eData interface{}) (map[string]interface{}, error) {
	//Get underlying object of interface
	object := reflect.ValueOf(eData)
	returnObj := map[string]interface{}{}

	aesGCM, err := e.initGCM()
	if err != nil {
		return nil, err
	}

	//For each field in object
	for i := 0; i < object.NumField(); i++ {
		fieldName, err := e.findFieldTag(object.Type().Field(i).Tag, []string{"bson", "ename"})
		if err != nil {
			return nil, err
		}

		//Get encrypted tag
		encrypt := object.Type().Field(i).Tag.Get("encrypted")
		//If encrypted == false don't encrypt otherwise encrypt
		if encrypt != "false" {
			val, err := e.sealWithNonce(aesGCM, Encode(object.Field(i)))
			if err != nil {
				return nil, err
			}

			returnObj[fieldName] = val
		} else {
			val := fmt.Sprint(object.Field(i))
			returnObj[fieldName] = val
		}

	}

	return returnObj, nil
}

//Input any 1-dimensional struct and get JSON bytes back representing the encrypted structure
func (e *encryptionService) EncryptToJSON(eData interface{}) ([]byte, error) {
	//Get underlying object of interface
	object := reflect.ValueOf(eData)
	returnObj := map[string]interface{}{}

	aesGCM, err := e.initGCM()
	if err != nil {
		return nil, err
	}

	//For each field in object
	for i := 0; i < object.NumField(); i++ {
		fieldName, err := e.findFieldTag(object.Type().Field(i).Tag, []string{"bson", "json"})
		if err != nil {
			return nil, err
		}

		encrypt := object.Type().Field(i).Tag.Get("encrypted")

		if encrypt != "false" {
			val, err := e.sealWithNonce(aesGCM, Encode(object.Field(i)))
			if err != nil {
				return nil, err
			}
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

//Decrypt a encrypted struct by passing the encrypted data and an empty object of the desired response type
func (e *encryptionService) Decrypt(encryptedData interface{}, desiredOutput interface{}) (interface{}, error) {
	//Get underlying object from interface
	object := reflect.ValueOf(encryptedData)
	returnObj := reflect.New(reflect.ValueOf(desiredOutput).Type())

	aesGCM, err := e.initGCM()
	if err != nil {
		return nil, err
	}

	//For each field in object
	for i := 0; i < object.NumField(); i++ {
		//Get encrypted tag
		encrypted := reflect.Indirect(returnObj).Type().Field(i).Tag.Get("encrypted")

		var decryptedStr string
		//If encrypted == false don't encrypt otherwise encrypt
		if encrypted != "false" {
			decryptedStr, err = e.getPlainText(object.Field(i).Bytes(), aesGCM)
			if err != nil {
				return nil, fmt.Errorf("failed to get text out of encrypted value, the following error occured: %s", err)
			}
		} else {
			decryptedStr = object.Field(i).String()
		}

		//Convert string to desired type
		field := reflect.Indirect(returnObj).Field(i)
		if field.IsValid() {
			val, err := Decode(decryptedStr, field)
			if err != nil {
				return nil, err
			}

			field.Set(reflect.ValueOf(val))
		}
	}

	return returnObj.Interface(), nil
}

//Decrypt encrypted []byte of a string
func (e *encryptionService) DecryptStr(b []byte) (string, error) {
	//Create cipher using given key
	aesGCM, err := e.initGCM()
	if err != nil {
		return "", err
	}

	//Decrypt string
	decryptedStr, err := e.getPlainText(b, aesGCM)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt string, the following error occured: %s", err)
	}

	return decryptedStr, nil
}

//Get decrypted []byte of encrypted value
func (e encryptionService) getPlainText(val []byte, aesGCM cipher.AEAD) (string, error) {
	//Get nonce size
	nonceSize := aesGCM.NonceSize()

	//Remove nonce from bytes
	nonce, ciphertext := val[:nonceSize], val[nonceSize:]
	//Get decrypted bytes
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

//Encrypt []byte
func (e *encryptionService) EncryptByt(b []byte) ([]byte, error) {
	//Create new cipher using given key
	aesGCM, err := e.initGCM()
	if err != nil {
		return nil, err
	}

	val, err := e.sealWithNonce(aesGCM, string(b))
	if err != nil {
		return nil, err
	}

	return val, nil
}

//Decrypt byte
func (e *encryptionService) DecryptByt(b []byte) ([]byte, error) {
	//Create new cipher using given key
	aesGCM, err := e.initGCM()
	if err != nil {
		return nil, err
	}

	//Decrypt bytes
	decryptedBytes, err := e.getPlainBytes(b, aesGCM)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt string, the following error occured: %s", err)
	}

	return decryptedBytes, nil
}

//Get decrypted bytes from encrypted []byte
func (e encryptionService) getPlainBytes(val []byte, aesGCM cipher.AEAD) ([]byte, error) {
	//Get nonce size
	nonceSize := aesGCM.NonceSize()

	nonce, ciphertext := val[:nonceSize], val[nonceSize:]
	plainbytes, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plainbytes, nil
}
