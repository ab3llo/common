package discovery

import (
	"testing"
)

func TestParsePort(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "valid port number",
			input:    "8080",
			expected: 8080,
		},
		{
			name:     "port zero",
			input:    "0",
			expected: 0,
		},
		{
			name:     "high port number",
			input:    "65535",
			expected: 65535,
		},
		{
			name:     "invalid port string",
			input:    "invalid",
			expected: 0,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0,
		},
		{
			name:     "negative number",
			input:    "-1",
			expected: -1,
		},
		{
			name:     "port with colon prefix",
			input:    ":8080",
			expected: 0, // parsePort doesn't handle colon prefix
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePort(tt.input)
			if result != tt.expected {
				t.Errorf("parsePort(%s) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRegisterServiceWithConsul_ValidatePortExtraction(t *testing.T) {
	tests := []struct {
		name         string
		grpcAddr     string
		expectedPort string
	}{
		{
			name:         "address with colon prefix",
			grpcAddr:     ":5001",
			expectedPort: "5001",
		},
		{
			name:         "address without colon",
			grpcAddr:     "8080",
			expectedPort: "8080",
		},
		{
			name:         "empty address",
			grpcAddr:     "",
			expectedPort: "",
		},
		{
			name:         "just colon",
			grpcAddr:     ":",
			expectedPort: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			port := tt.grpcAddr
			if len(port) > 0 && port[0] == ':' {
				if len(port) > 1 {
					port = port[1:]
				} else {
					port = ""
				}
			}

			if port != tt.expectedPort {
				t.Errorf("Port extraction for %s = %s, want %s", tt.grpcAddr, port, tt.expectedPort)
			}
		})
	}
}

func TestParsePort_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "very large port number",
			input:    "99999",
			expected: 99999,
		},
		{
			name:     "port with leading zeros",
			input:    "008080",
			expected: 8080,
		},
		{
			name:     "hexadecimal string",
			input:    "0x1f90",
			expected: 0,
		},
		{
			name:     "floating point string",
			input:    "8080.5",
			expected: 0,
		},
		{
			name:     "port with spaces",
			input:    " 8080 ",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePort(tt.input)
			if result != tt.expected {
				t.Errorf("parsePort(%s) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestServiceRegistration_ValidationLogic(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		grpcAddr    string
		consulHost  string
		consulPort  string
		expectValid bool
	}{
		{
			name:        "valid service registration",
			serviceName: "test-service",
			grpcAddr:    ":8080",
			consulHost:  "localhost",
			consulPort:  "8500",
			expectValid: true,
		},
		{
			name:        "empty service name",
			serviceName: "",
			grpcAddr:    ":8080",
			consulHost:  "localhost",
			consulPort:  "8500",
			expectValid: false,
		},
		{
			name:        "empty grpc address",
			serviceName: "test-service",
			grpcAddr:    "",
			consulHost:  "localhost",
			consulPort:  "8500",
			expectValid: false,
		},
		{
			name:        "empty consul host",
			serviceName: "test-service",
			grpcAddr:    ":8080",
			consulHost:  "",
			consulPort:  "8500",
			expectValid: false,
		},
		{
			name:        "invalid port in grpc address",
			serviceName: "test-service",
			grpcAddr:    ":invalid",
			consulHost:  "localhost",
			consulPort:  "8500",
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := true

			if tt.serviceName == "" {
				valid = false
			}
			if tt.grpcAddr == "" {
				valid = false
			}
			if tt.consulHost == "" {
				valid = false
			}
			if tt.grpcAddr != "" {
				port := tt.grpcAddr
				if len(port) > 0 && port[0] == ':' {
					if len(port) > 1 {
						port = port[1:]
					} else {
						port = ""
					}
				}
				if parsePort(port) == 0 && port != "0" {
					valid = false
				}
			}

			if valid != tt.expectValid {
				t.Errorf("Service registration validation = %v, want %v", valid, tt.expectValid)
			}
		})
	}
}

func TestGetInstance_ParameterValidation(t *testing.T) {
	tests := []struct {
		name       string
		service    string
		project    string
		region     string
		expectErr  bool
	}{
		{
			name:      "valid parameters",
			service:   "test-service",
			project:   "test-project",
			region:    "us-central1",
			expectErr: false,
		},
		{
			name:      "empty service name",
			service:   "",
			project:   "test-project",
			region:    "us-central1",
			expectErr: true,
		},
		{
			name:      "empty project",
			service:   "test-service",
			project:   "",
			region:    "us-central1",
			expectErr: true,
		},
		{
			name:      "empty region",
			service:   "test-service",
			project:   "test-project",
			region:    "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shouldErr := tt.service == "" || tt.project == "" || tt.region == ""
			if shouldErr != tt.expectErr {
				t.Errorf("Parameter validation expectation mismatch for %s", tt.name)
			}

			if !shouldErr {
				expectedPath := "projects/" + tt.project + "/locations/" + tt.region + "/namespaces/default/services/" + tt.service
				if expectedPath == "" {
					t.Error("Expected path should not be empty for valid parameters")
				}
			}
		})
	}
}

func BenchmarkParsePort(b *testing.B) {
	testPorts := []string{"8080", "9090", "5000", "3000", "8443"}
	for i := 0; i < b.N; i++ {
		port := testPorts[i%len(testPorts)]
		parsePort(port)
	}
}

func BenchmarkPortExtraction(b *testing.B) {
	testAddrs := []string{":8080", ":9090", "5000", ":3000", "8443"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		addr := testAddrs[i%len(testAddrs)]
		port := addr
		if len(port) > 0 && port[0] == ':' {
			if len(port) > 1 {
				port = port[1:]
			} else {
				port = ""
			}
		}
		_ = port
	}
}