package valueobjects

import "testing"

func TestNewPhone(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "valid phone",
			value:   "+79991234567",
			wantErr: false,
		},
		{
			name:    "valid phone without plus",
			value:   "79991234567",
			wantErr: false,
		},
		{
			name:    "phone with spaces",
			value:   "+7 999 123 45 67",
			wantErr: false,
		},
		{
			name:    "phone with brackets",
			value:   "+7 (999) 123-45-67",
			wantErr: false,
		},
		{
			name:    "too short phone",
			value:   "123456",
			wantErr: true,
		},
		{
			name:    "too long phone",
			value:   "1234567890123456",
			wantErr: true,
		},
		{
			name:    "empty phone",
			value:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone, err := NewPhone(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPhone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && phone == nil {
				t.Errorf("NewPhone() returned nil phone for valid input")
			}
		})
	}
}

func TestPhone_String(t *testing.T) {
	phone, _ := NewPhone("+79991234567")
	result := phone.String()
	if result != "+79991234567" {
		t.Errorf("String() got = %v, want %v", result, "+79991234567")
	}
}
