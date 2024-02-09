//go:build unit

package common

import (
	"reflect"
	"testing"
)

func TestIsCompressed(t *testing.T) {
	type args struct {
		data []byte
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"compressed", args{[]byte{0x1f, 0x8b}}, true},
		{"not compressed", args{[]byte{0x1f, 0x8c}}, false},
	}

	for _, tt := range tests {
		testInstance := tt
		t.Run(testInstance.name, func(t *testing.T) {
			if got := IsCompressed(testInstance.args.data); got != testInstance.want {
				t.Errorf("IsCompressed() = %v, want %v", got, testInstance.want)
			}
		})
	}
}

func TestCompressData(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    []byte
		wantErr bool
	}{
		{"Compress valid data", []byte("Hello world"), []byte{
			0x1f, 0x8b, 0x8, 0x0, 0x0, 0x0,
			0x0, 0x0, 0x2, 0xff, 0xf2, 0x48, 0xcd, 0xc9, 0xc9, 0x57, 0x28, 0xcf, 0x2f, 0xca, 0x49, 0x1, 0x4, 0x0,
			0x0, 0xff, 0xff, 0x52, 0x9e, 0xd6, 0x8b, 0xb, 0x0, 0x0, 0x0,
		}, false},
	}

	for _, tt := range tests {
		testInstance := tt
		t.Run(testInstance.name, func(t *testing.T) {
			got, err := CompressData(testInstance.data)
			if (err != nil) != testInstance.wantErr {
				t.Errorf("CompressData() error = %v, wantErr %v", err, testInstance.wantErr)

				return
			}
			if !reflect.DeepEqual(got, testInstance.want) {
				t.Errorf("CompressData() = %v, want %v", got, testInstance.want)
			}
		})
	}
}

func TestUncompressData(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    []byte
		wantErr bool
	}{
		{"uncompress correct data", []byte{
			0x1f, 0x8b, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2,
			0xff, 0xf2, 0x48, 0xcd, 0xc9, 0xc9, 0x57, 0x28, 0xcf, 0x2f, 0xca, 0x49, 0x1, 0x4, 0x0, 0x0, 0xff, 0xff,
			0x52, 0x9e, 0xd6, 0x8b, 0xb, 0x0, 0x0, 0x0,
		}, []byte("Hello world"), false},
		{"uncompress wrong data", []byte{
			0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2,
			0xff, 0xf2, 0x48, 0xcd, 0xc9, 0xc9, 0x57, 0x28, 0xcf, 0x2f, 0xca, 0x49, 0x1, 0x4, 0x0, 0x0, 0xff, 0xff,
			0x52, 0x9e, 0xd6, 0x8b, 0xb, 0x0, 0x0, 0x0,
		}, nil, true},
	}

	for _, tt := range tests {
		testInstance := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := UncompressData(testInstance.data)
			if (err != nil) != testInstance.wantErr {
				t.Errorf("UncompressData() error = %v, wantErr %v", err, testInstance.wantErr)

				return
			}
			if !reflect.DeepEqual(got, testInstance.want) {
				t.Errorf("UncompressData() = %v, want %v", got, testInstance.want)
			}
		})
	}
}
