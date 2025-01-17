package modzip

import (
	"encoding/json"
	"testing"
)

func TestInspectModZip(t *testing.T) {
	zipPath := "testdata/UnlockLevelCurve_Patch_XP_x0.5-377-2-0-0-24-1721672677.zip"
	modData := InspectModZip(zipPath)
	if len(modData.MD5) == 0 {
		t.Errorf("MD5 is empty")
	}
	if len(modData.Mods) == 0 {
		t.Errorf("Mods is empty")
	}

	// print the whole thing as json
	bs, err := json.MarshalIndent(modData, "", "  ")
	if err != nil {
		t.Errorf("Failed to marshal modData: %v", err)
	}

	t.Logf("modData: %s", string(bs))
}
