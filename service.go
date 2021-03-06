package encryption

//go:generate mockgen -source=service.go -destination=mock/mock_service.go -package=mock_encryption
type EncryptionService interface {
	EncryptToInterface(eData interface{}) (map[string]interface{}, error)
	EncryptToJSON(eData interface{}) ([]byte, error)
	Decrypt(eData interface{}, eData2 interface{}) (interface{}, error)

	EncryptStr(str string) ([]byte, error)
	DecryptStr(b []byte) (string, error)

	EncryptByt(b []byte) ([]byte, error)
	DecryptByt(b []byte) ([]byte, error)
}
