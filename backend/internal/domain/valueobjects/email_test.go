package valueobjects

import "testing"

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "valid email",
			value:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "valid email with subdomain",
			value:   "user@mail.example.com",
			wantErr: false,
		},
		{
			name:    "invalid email without @",
			value:   "testexample.com",
			wantErr: true,
		},
		{
			name:    "invalid email without domain",
			value:   "test@",
			wantErr: true,
		},
		{
			name:    "empty email",
			value:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && email.Value != tt.value {
				t.Errorf("NewEmail() got = %v, want %v", email.Value, tt.value)
			}
		})
	}
}

func TestEmail_String(t *testing.T) {
	email, _ := NewEmail("test@example.com")
	if email.String() != "test@example.com" {
		t.Errorf("String() got = %v, want %v", email.String(), "test@example.com")
	}
}
