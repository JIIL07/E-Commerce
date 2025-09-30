package tests

import (
	"testing"
	"time"

	"ecommerce-backend/internal/utils"
)

func TestGenerateUUID(t *testing.T) {
	uuid1 := utils.GenerateUUID()
	uuid2 := utils.GenerateUUID()
	
	if uuid1 == uuid2 {
		t.Error("Generated UUIDs should be unique")
	}
	
	if len(uuid1) != 36 {
		t.Errorf("UUID should be 36 characters long, got %d", len(uuid1))
	}
}

func TestGenerateRandomString(t *testing.T) {
	str1 := utils.GenerateRandomString(10)
	str2 := utils.GenerateRandomString(10)
	
	if len(str1) != 10 {
		t.Errorf("Random string should be 10 characters long, got %d", len(str1))
	}
	
	if str1 == str2 {
		t.Error("Random strings should be different")
	}
}

func TestHashPassword(t *testing.T) {
	password := "testpassword123"
	hash, err := utils.HashPassword(password)
	
	if err != nil {
		t.Errorf("HashPassword failed: %v", err)
	}
	
	if hash == password {
		t.Error("Hashed password should not equal original password")
	}
	
	if len(hash) == 0 {
		t.Error("Hash should not be empty")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testpassword123"
	hash, _ := utils.HashPassword(password)
	
	if !utils.CheckPasswordHash(password, hash) {
		t.Error("CheckPasswordHash should return true for correct password")
	}
	
	if utils.CheckPasswordHash("wrongpassword", hash) {
		t.Error("CheckPasswordHash should return false for incorrect password")
	}
}

func TestGenerateOTP(t *testing.T) {
	otp := utils.GenerateOTP(6)
	
	if len(otp) != 6 {
		t.Errorf("OTP should be 6 digits long, got %d", len(otp))
	}
	
	for _, char := range otp {
		if char < '0' || char > '9' {
			t.Error("OTP should contain only digits")
		}
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"Test & Special Characters!", "test-special-characters"},
		{"Multiple   Spaces", "multiple-spaces"},
		{"UPPERCASE", "uppercase"},
	}
	
	for _, test := range tests {
		result := utils.Slugify(test.input)
		if result != test.expected {
			t.Errorf("Slugify(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1", true},
		{"0", false},
		{"yes", true},
		{"no", false},
		{"invalid", false},
	}
	
	for _, test := range tests {
		result, err := utils.ParseBool(test.input)
		if err != nil && test.expected {
			t.Errorf("ParseBool(%s) returned error %v, expected %v", test.input, err, test.expected)
		}
		if result != test.expected {
			t.Errorf("ParseBool(%s) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestIsEmailValid(t *testing.T) {
	validEmails := []string{
		"test@example.com",
		"user.name@domain.co.uk",
		"user+tag@example.org",
	}
	
	invalidEmails := []string{
		"invalid-email",
		"@example.com",
		"user@",
		"user@.com",
	}
	
	for _, email := range validEmails {
		if !utils.IsEmailValid(email) {
			t.Errorf("IsEmailValid(%s) should return true", email)
		}
	}
	
	for _, email := range invalidEmails {
		if utils.IsEmailValid(email) {
			t.Errorf("IsEmailValid(%s) should return false", email)
		}
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{1024, "1.00 KB"},
		{1048576, "1.00 MB"},
		{1073741824, "1.00 GB"},
		{500, "500 B"},
	}
	
	for _, test := range tests {
		result := utils.FormatBytes(test.bytes)
		if result != test.expected {
			t.Errorf("FormatBytes(%d) = %s, expected %s", test.bytes, result, test.expected)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	duration := 2*time.Hour + 30*time.Minute + 45*time.Second
	result := utils.FormatDuration(duration)
	
	if result != "2h30m45s" {
		t.Errorf("FormatDuration returned %s, expected 2h30m45s", result)
	}
}
