package auth

import (
	"testing"
)

func TestPasswordHashing(t *testing.T) {
	password1 := "testPass1"
	password2 := "testPass2"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)
	tests := map[string]struct {
		password string
		hash     string
		wantErr  bool
		want     bool
	}{
		"password matches hash":                 {password: password1, hash: hash1, wantErr: false, want: true},
		"password doesn't match different hash": {password: password1, hash: hash2, wantErr: false, want: false},
		"incorrect password":                    {password: "incorrect", hash: hash1, wantErr: false, want: false},
		"empty password":                        {password: "", hash: hash1, wantErr: false, want: false},
		"invalid hash":                          {password: password2, hash: "invalid", wantErr: true, want: false},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			match, err := CheckPasswordHash(tc.password, tc.hash)
			if (err != nil) != tc.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tc.wantErr)
			}
			if !tc.wantErr && match != tc.want {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tc.want, match)
			}
		})
	}
}
