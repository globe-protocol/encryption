package encryption

import (
	"reflect"
	"testing"
)

func TestEncode(t *testing.T) {
	type args struct {
		value reflect.Value
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "[]string",
			args: args{
				value: reflect.ValueOf([]string{"test", "test met iets anders"}),
			},
			want: "test°test met iets anders",
		},
		{
			name: "string",
			args: args{
				value: reflect.ValueOf("test met iets anders"),
			},
			want: "test met iets anders",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Encode(tt.args.value); got != tt.want {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
