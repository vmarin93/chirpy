package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
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

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)
	expiredToken, _ := MakeJWT(userID, "secret", time.Nanosecond)
	tests := map[string]struct {
		tokenString string
		tokenSecret string
		wantErr     bool
		want        uuid.UUID
	}{
		"valid jwt":     {tokenString: validToken, tokenSecret: "secret", wantErr: false, want: userID},
		"invalid jwt":   {tokenString: "invalid", tokenSecret: "secret", wantErr: true, want: uuid.Nil},
		"wrong secret":  {tokenString: validToken, tokenSecret: "badSecret", wantErr: true, want: uuid.Nil},
		"expired token": {tokenString: expiredToken, tokenSecret: "secret", wantErr: true, want: uuid.Nil},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ValidateJWT(tc.tokenString, tc.tokenSecret)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tc.wantErr)
			}
			if !tc.wantErr && got != tc.want {
				t.Errorf("ValidateJWT() expects = %v, got %v", tc.want, got)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := map[string]struct {
		header  http.Header
		wantErr bool
		want    string
	}{
		"Valid Token": {
			header: http.Header{
				"Authorization": []string{"Bearer validToken"}},
			wantErr: false, want: "validToken"},
		"Missing header": {header: http.Header{}, wantErr: true, want: ""},
		"Malformed header": {
			header: http.Header{
				"Authorization": []string{"Invalid validToken"}},
			wantErr: true, want: ""},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := GetBearerToken(tc.header)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tc.wantErr)
			}
			if !tc.wantErr && got != tc.want {
				t.Errorf("GetBearerToken() expects = %v, got %v", tc.want, got)
			}
		})
	}
}
