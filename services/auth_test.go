package services

import (
	"os"
	"testing"
)

func TestAuthService_PublicKey(t *testing.T) {
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzkzMjk2NzUsImlkIjoiOGU3OWUwNTUtMzcxNi00ZTAxLTgzMWUtZjBiZTk2YjEyMmUzIiwib3JpZ19pYXQiOjE2Nzg3MjQ4NzUsIndhbGxldF9hZGRyIjoiN3VpaldXZDVKNjRRVUVKRGZpZVk4OUM4QlhFVXhUV1lMRkJIVzhzNDk0dFYifQ.N5zvBjmDHqvcOAeXqzCDEVEtGjE2Z1TkY2aYdyGcZEE"

	os.Setenv("JWT_TOKEN", jwt)
	asvc := AuthService{}
	err := asvc.Start()
	if err != nil {
		t.Fatal(err)
	}

	if asvc.PublicKey() != "7uijWWd5J64QUEJDfieY89C8BXEUxTWYLFBHW8s494tV" {
		t.Logf("Expected: %s - Got: %s", "7uijWWd5J64QUEJDfieY89C8BXEUxTWYLFBHW8s494tV", asvc.PublicKey())
		t.Fail()
	}
}
