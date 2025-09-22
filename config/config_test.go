package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	testCases := []struct {
		name        string
		envVars     map[string]string
		configFile  string
		expectedCfg Config
		expectError bool
	}{
		{
			name:    "default values when no config file or env vars",
			envVars: map[string]string{},
			expectedCfg: Config{
				ServiceName: "gateway",
				Env:         "dev",
				DSN:         "postgresql://postgres:postgres@127.0.0.1:54322/postgres",
				PORT:        "8080",
				GCPProject:  "",
				GCPRegion:   "",
				ClerkAPIKey: "",
				ConsulHost:  "consul",
				ConsulPort:  "8500",
			},
			expectError: false,
		},
		{
			name: "environment variables override defaults",
			envVars: map[string]string{
				"SERVICE_NAME":         "test-service",
				"ENV":                  "production",
				"PORT":                 "3000",
				"CLERK_API_KEY":        "test-api-key",
				"CONSUL_HOST":          "localhost",
				"CONSUL_PORT":          "8501",
				"GOOGLE_CLOUD_PROJECT": "test-project",
				"GOOGLE_CLOUD_REGION":  "us-central1",
			},
			expectedCfg: Config{
				ServiceName: "test-service",
				Env:         "production",
				DSN:         "postgresql://postgres:postgres@127.0.0.1:54322/postgres",
				PORT:        "3000",
				GCPProject:  "test-project",
				GCPRegion:   "us-central1",
				ClerkAPIKey: "test-api-key",
				ConsulHost:  "localhost",
				ConsulPort:  "8501",
			},
			expectError: false,
		},
		{
			name: "config file overrides defaults",
			configFile: `SERVICE_NAME=file-service
ENV=staging
PORT=4000
CLERK_API_KEY=file-api-key
`,
			expectedCfg: Config{
				ServiceName: "file-service",
				Env:         "staging",
				DSN:         "postgresql://postgres:postgres@127.0.0.1:54322/postgres",
				PORT:        "4000",
				GCPProject:  "",
				GCPRegion:   "",
				ClerkAPIKey: "file-api-key",
				ConsulHost:  "consul",
				ConsulPort:  "8500",
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "config-test")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			for key, value := range tc.envVars {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			if tc.configFile != "" {
				configPath := filepath.Join(tempDir, ".env.local")
				err := os.WriteFile(configPath, []byte(tc.configFile), 0644)
				if err != nil {
					t.Fatalf("Failed to write config file: %v", err)
				}
			}

			config, err := LoadConfig(tempDir)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tc.expectError {
				if config.ServiceName != tc.expectedCfg.ServiceName {
					t.Errorf("ServiceName = %v, want %v", config.ServiceName, tc.expectedCfg.ServiceName)
				}
				if config.Env != tc.expectedCfg.Env {
					t.Errorf("Env = %v, want %v", config.Env, tc.expectedCfg.Env)
				}
				if config.DSN != tc.expectedCfg.DSN {
					t.Errorf("DSN = %v, want %v", config.DSN, tc.expectedCfg.DSN)
				}
				if config.PORT != tc.expectedCfg.PORT {
					t.Errorf("PORT = %v, want %v", config.PORT, tc.expectedCfg.PORT)
				}
				if config.GCPProject != tc.expectedCfg.GCPProject {
					t.Errorf("GCPProject = %v, want %v", config.GCPProject, tc.expectedCfg.GCPProject)
				}
				if config.GCPRegion != tc.expectedCfg.GCPRegion {
					t.Errorf("GCPRegion = %v, want %v", config.GCPRegion, tc.expectedCfg.GCPRegion)
				}
				if config.ClerkAPIKey != tc.expectedCfg.ClerkAPIKey {
					t.Errorf("ClerkAPIKey = %v, want %v", config.ClerkAPIKey, tc.expectedCfg.ClerkAPIKey)
				}
				if config.ConsulHost != tc.expectedCfg.ConsulHost {
					t.Errorf("ConsulHost = %v, want %v", config.ConsulHost, tc.expectedCfg.ConsulHost)
				}
				if config.ConsulPort != tc.expectedCfg.ConsulPort {
					t.Errorf("ConsulPort = %v, want %v", config.ConsulPort, tc.expectedCfg.ConsulPort)
				}
			}
		})
	}
}

func TestConfig_Struct(t *testing.T) {
	config := Config{
		ServiceName: "test-service",
		Env:         "test",
		DSN:         "test-dsn",
		PORT:        "8080",
		GCPProject:  "test-project",
		GCPRegion:   "us-central1",
		ClerkAPIKey: "test-key",
		ConsulHost:  "localhost",
		ConsulPort:  "8500",
	}

	if config.ServiceName != "test-service" {
		t.Errorf("ServiceName = %v, want %v", config.ServiceName, "test-service")
	}
	if config.Env != "test" {
		t.Errorf("Env = %v, want %v", config.Env, "test")
	}
	if config.DSN != "test-dsn" {
		t.Errorf("DSN = %v, want %v", config.DSN, "test-dsn")
	}
}
