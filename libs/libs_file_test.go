package libs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileExists(t *testing.T) {
	tests := map[string]bool{
		"file.go": true,
		"aaa.bbb": false,
		"/":       false,
	}
	for fpath, mustExist := range tests {
		exist := FileExists(fpath)
		require.Equalf(t, mustExist, exist, "invalid \"%s\": %v", fpath, exist)
	}
}
