package AES

import (
	"fmt"
	"reflect"
	"testing"
)

type stringfloatbool struct {
	String  string  `bson:"String" encrypted:"false"`
	Float64 float64 `bson:"Float64"`
	Bool    bool    `bson:"Bool"`
}

type stringfloatboolEnc struct {
	String  string `bson:"String" encrypted:"false"`
	Float64 []byte `bson:"Float64"`
	Bool    []byte `bson:"Bool"`
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := encryptionService{
				key: []byte{176, 55, 108, 116, 181, 15, 21, 190, 134, 27, 183, 18, 48, 179, 221, 123, 225, 172, 55, 54, 142, 158, 173, 59, 77, 239, 116, 99, 248, 15, 228, 254},
			}

			encryptedData, err := e.Encrypt(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("encryptionService.Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			encryptedObj := stringfloatboolEnc{
				String:  fmt.Sprint(encryptedData["String"]),
				Float64: encryptedData["Float64"].([]byte),
				Bool:    encryptedData["Bool"].([]byte),
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

			if *decryptedObj != tt.args.data {
				t.Errorf("decrypted body = %+v\n want = %+v\n", decryptedObj, tt.args.data)
			}
		})
	}
}
