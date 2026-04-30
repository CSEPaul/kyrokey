package libs

import "testing"

func TestConstants(t *testing.T) {
	if LogPath != "./logs/ZAP.log" {
		t.Errorf("LogPath expected %q, got %q", "./logs/ZAP.log", LogPath)
	}

	if KeyChainDB != "./conf/keychain_entries.db" {
		t.Errorf("KeyChainDB expected %q, got %q", "./conf/keychain_entries.db", KeyChainDB)
	}

	if OutputDir != "./output/" {
		t.Errorf("OutputDir expected %q, got %q", "./output/", OutputDir)
	}
}
