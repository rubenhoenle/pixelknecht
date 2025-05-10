package config

import (
	"os"
	"testing"
)

const (
	envName = "TEST_ENV_1"
)

func Test_readEnvWithFallback(t *testing.T) {
	type args struct {
		envKey        string
		fallbackValue string
	}
	tests := []struct {
		name     string
		args     args
		envVal   string
		expected string
	}{
		{
			name: "fallback",
			args: args{
				envKey:        envName,
				fallbackValue: "fallback-val-1",
			},
			envVal:   "env-value-1",
			expected: "env-value-1",
		},
		{
			name: "fallback",
			args: args{
				envKey:        envName,
				fallbackValue: "fallback-val-1",
			},
			expected: "fallback-val-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			if tt.envVal != "" {
				// set env var to test value
				os.Setenv(tt.args.envKey, tt.envVal)
			}

			actual := readEnvWithFallback(tt.args.envKey, tt.args.fallbackValue)
			if actual != tt.expected {
				t.Errorf("readEnvWithFallback(), expected = %s, actual = %s", tt.expected, actual)
			}
		})
	}
}
