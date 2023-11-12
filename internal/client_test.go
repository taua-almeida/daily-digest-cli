package internal

import (
	"os"
	"testing"
)

func TestGetGitToken(t *testing.T) {
	// Set up the environment variable for the test
	os.Setenv("GITHUB_TOKEN", "testtoken")
	defer os.Unsetenv("GITHUB_TOKEN") // Clean up after the test

	// Test for a successful case
	token, err := GetGitToken("GITHUB_TOKEN")
	if err != nil {
		t.Errorf("GetGitToken returned an error: %v", err)
	}
	if token != "testtoken" {
		t.Errorf("Expected 'testtoken', got '%s'", token)
	}

	// Test for the error case
	_, err = GetGitToken("NON_EXISTENT_ENV")
	if err == nil {
		t.Error("Expected an error for non-existent environment variable, got none")
	}
}
