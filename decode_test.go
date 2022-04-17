package encryption

import (
	"reflect"
	"testing"
)

func TestConvert(t *testing.T) {
	type args struct {
		s string
		t reflect.Value
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "convert int8",
			args: args{
				s: "8",
				t: reflect.ValueOf(int8(0)),
			},
			want:    int8(8),
			wantErr: false,
		},
		{
			name: "convert uint8",
			args: args{
				s: "8",
				t: reflect.ValueOf(uint8(0)),
			},
			want:    uint8(8),
			wantErr: false,
		},
		{
			name: "convert []byte",
			args: args{
				s: "[5 4 3 2 1]",
				t: reflect.ValueOf([]byte{}),
			},
			want:    []byte{5, 4, 3, 2, 1},
			wantErr: false,
		},
		{
			name: "convert int16",
			args: args{
				s: "16",
				t: reflect.ValueOf(int16(0)),
			},
			want:    int16(16),
			wantErr: false,
		},
		{
			name: "convert uint16",
			args: args{
				s: "16",
				t: reflect.ValueOf(uint16(0)),
			},
			want:    uint16(16),
			wantErr: false,
		},
		{
			name: "convert int32",
			args: args{
				s: "32",
				t: reflect.ValueOf(int32(0)),
			},
			want:    int32(32),
			wantErr: false,
		},
		{
			name: "convert uint32",
			args: args{
				s: "32",
				t: reflect.ValueOf(uint32(0)),
			},
			want:    uint32(32),
			wantErr: false,
		},
		{
			name: "convert int64",
			args: args{
				s: "64",
				t: reflect.ValueOf(int64(0)),
			},
			want:    int64(64),
			wantErr: false,
		},
		{
			name: "convert uint64",
			args: args{
				s: "64",
				t: reflect.ValueOf(uint64(0)),
			},
			want:    uint64(64),
			wantErr: false,
		},
		{
			name: "convert int",
			args: args{
				s: "64",
				t: reflect.ValueOf(int(0)),
			},
			want:    int(64),
			wantErr: false,
		},
		{
			name: "convert uint",
			args: args{
				s: "64",
				t: reflect.ValueOf(uint(0)),
			},
			want:    uint(64),
			wantErr: false,
		},
		{
			name: "convert uintptr",
			args: args{
				s: "64",
				t: reflect.ValueOf(uintptr(0)),
			},
			want:    uintptr(64),
			wantErr: false,
		},
		{
			name: "convert float32",
			args: args{
				s: "32.2",
				t: reflect.ValueOf(float32(0)),
			},
			want:    float32(32.2),
			wantErr: false,
		},
		{
			name: "convert float64",
			args: args{
				s: "64.4",
				t: reflect.ValueOf(float64(0)),
			},
			want:    float64(64.4),
			wantErr: false,
		},
		{
			name: "convert string",
			args: args{
				s: "string",
				t: reflect.ValueOf(string("")),
			},
			want:    string("string"),
			wantErr: false,
		},
		{
			name: "convert bool",
			args: args{
				s: "true",
				t: reflect.ValueOf(bool(true)),
			},
			want:    bool(true),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.args.s, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}
