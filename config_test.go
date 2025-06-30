package common

import (
	"os"
	"testing"
)

func TestGetEnvOrDefault_ReturnsEnv(t *testing.T) {
	os.Setenv("TEST_KEY", "test_value")
	t.Cleanup(func() { os.Unsetenv("TEST_KEY") })
	v := GetEnvOrDefault("TEST_KEY", "default")
	if v != "test_value" {
		t.Errorf("expected 'test_value', got '%s'", v)
	}
}

func TestGetEnvOrDefault_ReturnsDefault(t *testing.T) {
	os.Unsetenv("TEST_KEY_NOT_SET")
	v := GetEnvOrDefault("TEST_KEY_NOT_SET", "default")
	if v != "default" {
		t.Errorf("expected 'default', got '%s'", v)
	}
}
