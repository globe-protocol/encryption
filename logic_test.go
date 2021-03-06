package encryption

import (
	"fmt"
	"reflect"
	"testing"
)

type stringfloatbool struct {
	String    string   `bson:"String" encrypted:"false"`
	Float64   float64  `bson:"Float64"`
	Bool      bool     `bson:"Bool"`
	StringArr []string `bson:"StringArr"`
	EmptyVal  string   `bson:"EmptyVal"`
}

type stringfloatboolEnc struct {
	String    string `bson:"String" encrypted:"false"`
	Float64   []byte `bson:"Float64"`
	Bool      []byte `bson:"Bool"`
	StringArr []byte `bson:"StringArr"`
	EmptyVal  []byte `bson:"EmptyVal"`
}

func Test_Encrypt_and_Decrypt(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "Successful data upload (string, float64, bool)",
			args: args{
				data: stringfloatbool{
					String:  "123Test",
					Float64: 64.64,
					Bool:    true,
				},
			},
			wantErr: false,
		},
		{
			name: "Successful data upload (string[])",
			args: args{
				data: stringfloatbool{
					StringArr: []string{"test value", "test, value 2"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := encryptionService{
				key: []byte{176, 55, 108, 116, 181, 15, 21, 190, 134, 27, 183, 18, 48, 179, 221, 123, 225, 172, 55, 54, 142, 158, 173, 59, 77, 239, 116, 99, 248, 15, 228, 254},
			}

			encryptedData, err := e.EncryptToInterface(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("encryptionService.Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			encryptedObj := stringfloatboolEnc{
				String:    fmt.Sprint(encryptedData["String"]),
				Float64:   encryptedData["Float64"].([]byte),
				Bool:      encryptedData["Bool"].([]byte),
				StringArr: encryptedData["StringArr"].([]byte),
				EmptyVal:  nil,
			}

			var desiredType stringfloatbool

			decryptedData, err := e.Decrypt(encryptedObj, desiredType)
			if err != nil {
				t.Errorf("Did not expect error: %s while decrypting object", err.Error())
			}

			decryptedObj, ok := reflect.ValueOf(decryptedData).Interface().(*stringfloatbool)
			if !ok {
				t.Errorf("error casting interface to desired type")
			}

			if reflect.DeepEqual(decryptedObj, tt.args.data) {
				t.Errorf("decrypted body = %+v\n want = %+v\n", decryptedObj, tt.args.data)
			}
		})
	}
}

func Test_encryptionService_EncryptStr(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successfully encrypt & decrypt string input",
			args: args{
				str: "test inputf5t67yyyyu857tfo",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := encryptionService{
				key: []byte{176, 55, 108, 116, 181, 15, 21, 190, 134, 27, 183, 18, 48, 179, 221, 123, 225, 172, 55, 54, 142, 158, 173, 59, 77, 239, 116, 99, 248, 15, 228, 254},
			}

			encrypted, err := e.EncryptStr(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("encryptionService.EncryptStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got, err := e.DecryptStr(encrypted)
			if (err != nil) != tt.wantErr {
				t.Errorf("encryptionService.EncryptStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.args.str {
				t.Errorf("encryptionService.EncryptStr() = %v, want %v", got, tt.args.str)
			}
		})
	}
}

func Test_encryptionService_EncryptBytes(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successfully encrypt & decrypt string input",
			args: args{
				b: []byte{12, 21, 21, 45, 52},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := encryptionService{
				key: []byte{176, 55, 108, 116, 181, 15, 21, 190, 134, 27, 183, 18, 48, 179, 221, 123, 225, 172, 55, 54, 142, 158, 173, 59, 77, 239, 116, 99, 248, 15, 228, 254},
			}

			encrypted, err := e.EncryptByt(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("encryptionService.EncryptStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got, err := e.DecryptByt(encrypted)
			if (err != nil) != tt.wantErr {
				t.Errorf("encryptionService.EncryptStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.args.b) {
				t.Errorf("encryptionService.EncryptStr() = %v, want %v", got, tt.args.b)
			}

			for i := 0; i < len(got); i++ {
				if got[i] != tt.args.b[i] {
					t.Errorf("encryptionService.EncryptStr() = %v, want %v", got, tt.args.b)
				}
			}
		})
	}
}
