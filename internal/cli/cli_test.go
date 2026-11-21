package cli

import "testing"

func TestSingleFile(t *testing.T) {
	for _, tc := range []struct {
		name     string
		args     []string
		wantPath string
		wantCode int
		wantOK   bool
	}{
		{"plain path", []string{"scores.4ss"}, "scores.4ss", 0, true},
		{"dash path after terminator", []string{"--", "-scores.4ss"}, "-scores.4ss", 0, true},
		{"help", []string{"--help"}, "", 0, false},
		{"missing path", nil, "", 2, false},
		{"too many paths", []string{"one", "two"}, "", 2, false},
		{"unknown flag", []string{"--unknown"}, "", 2, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path, code, ok := SingleFile(tc.args, "tool", "file")
			if path != tc.wantPath || code != tc.wantCode || ok != tc.wantOK {
				t.Fatalf("SingleFile(%v) = (%q, %d, %v), want (%q, %d, %v)", tc.args, path, code, ok, tc.wantPath, tc.wantCode, tc.wantOK)
			}
		})
	}
}
