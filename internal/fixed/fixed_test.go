package fixed

import (
	"testing"
)

func TestUint16_16_String(t *testing.T) {
	tests := []struct {
		name string
		x    Uint16_16
		want string
	}{
		{"1.0", Uint16_16(0x00010000), "1.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.x.String(); got != tt.want {
				t.Errorf("Uint16_16.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint16_16_UnmarshalBinary(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		x       Uint16_16
		args    args
		wantErr bool
		want    string
	}{
		{
			"1.0",
			Uint16_16(0),
			args{[]byte{0x00, 0x01, 0x00, 0x00}},
			false,
			"1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.x.UnmarshalBinary(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Uint16_16.UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.x.String() != tt.want {
				t.Errorf("Uint16_16.UnmarshalBinary() string = %v, want %v", tt.x.String(), tt.want)
			}
		})
	}
}

func TestUint8_8_String(t *testing.T) {
	tests := []struct {
		name string
		x    Uint8_8
		want string
	}{
		{"1.0", Uint8_8(0x0100), "1.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.x.String(); got != tt.want {
				t.Errorf("Uint8_8.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUint8_8_UnmarshalBinary(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		x       Uint8_8
		args    args
		wantErr bool
		want    string
	}{
		{
			"1.0",
			Uint8_8(0),
			args{[]byte{0x01, 0x00}},
			false,
			"1.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.x.UnmarshalBinary(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("Uint8_8.UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.x.String() != tt.want {
				t.Errorf("Uint16_16.UnmarshalBinary() string = %v, want %v", tt.x.String(), tt.want)
			}
		})
	}
}
