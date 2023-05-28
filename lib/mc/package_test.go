package mc

import (
	"reflect"
	"testing"
)

func TestToVarInt(t *testing.T) {
	tests := []struct {
		name string
		arg  int
		want []byte
	}{
		{
			name: "0",
			arg:  0,
			want: []byte{0},
		},
		{
			name: "1",
			arg:  1,
			want: []byte{1},
		},
		{
			name: "127",
			arg:  127,
			want: []byte{0x7f},
		},
		{
			name: "128",
			arg:  128,
			want: []byte{0x80, 0x01},
		},
		{
			name: "255",
			arg:  255,
			want: []byte{0xff, 0x01},
		},
		{
			name: "25565",
			arg:  25565,
			want: []byte{0xdd, 0xc7, 0x01},
		},
		{
			name: "2097151",
			arg:  2097151,
			want: []byte{0xff, 0xff, 0x7f},
		},
		{
			name: "2147483647",
			arg:  2147483647,
			want: []byte{0xff, 0xff, 0xff, 0xff, 0x07},
		},
		{
			name: "-1",
			arg:  -1,
			want: []byte{0xff, 0xff, 0xff, 0xff, 0x0f},
		},
		{
			name: "-2147483648",
			arg:  -2147483648,
			want: []byte{0x80, 0x80, 0x80, 0x80, 0x08},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToVarInt(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToVarInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
