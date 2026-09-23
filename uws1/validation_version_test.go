package uws1

import "testing"

func TestSupportsUWSVersionAtLeastUsesSemVerPrecedence(t *testing.T) {
	tests := []struct {
		version string
		want    bool
	}{
		{version: "1.9.2", want: true},
		{version: "1.9.2-rc.1", want: false},
		{version: "1.9.3-rc.1", want: true},
		{version: "1.9.2-beta.2", want: false},
		{version: "1.9.2-beta.11", want: false},
		{version: "1.9.2-beta.11+build.7", want: false},
		{version: "1.9.2-beta.01", want: false},
	}

	for _, test := range tests {
		t.Run(test.version, func(t *testing.T) {
			if got := supportsUWSVersionAtLeast(test.version, 1, 9, 2); got != test.want {
				t.Fatalf("supportsUWSVersionAtLeast(%q, 1, 9, 2) = %v, want %v", test.version, got, test.want)
			}
		})
	}
}
