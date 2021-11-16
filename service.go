package AES

//go:generate mockgen -source=service.go -destination=mock/mock_service.go -package=mock_encryption
type EncryptionService interface {
	Encrypt(eData interface{}) (map[string]interface{}, error)
	Decrypt(eData interface{}, eData2 interface{}) (interface{}, error)
}
