package modzip

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestInspectModZipUsingTempFolderExtract(t *testing.T) {
	zipPath := "testdata/UnlockLevelCurve_Patch_XP_x0.5-377-2-0-0-24-1721672677.zip"
	modData := InspectModZipUsingTempFolderExtract(zipPath)
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

	expectedData := `{
          "Mods": [
            {
              "Author": "Charis",
              "Name": "UnlockLevelCurve_Patch_XP_x0.5",
              "Folder": "UnlockLevelCurve_Patch_XP_x0.5",
              "Version": "72057594037927960",
              "Description": "Halves the XP requirement",
              "UUID": "e53ae4b5-a922-47ef-b69d-d55c5745a65b",
              "Created": "2024-07-22T20:17:59.1988778+02:00",
              "Dependencies": [],
              "Group": "dafe83c2-97b3-4b05-9aef-2e1cc2e1de98"
            }
          ],
          "MD5": "87758161b02ba6eb90ac2a6c92cd746f"
        }`

	modDataExpect := ModData{}
	err = json.Unmarshal([]byte(expectedData), &modDataExpect)
	if err != nil {
		t.Errorf("Failed to unmarshal expected data: %v", err)
	}

	if diff := cmp.Diff(modDataExpect, modData); diff != "" {
		t.Errorf("modData mismatch (-want +got):\n%s", diff)
	}
}

func TestInspectModZip(t *testing.T) {
	zipPath := "testdata/UnlockLevelCurve_Patch_XP_x0.5-377-2-0-0-24-1721672677.zip"
	modData, pakFile := InspectModZip(zipPath)
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

	expectedData := `{
          "Mods": [
            {
              "Author": "Charis",
              "Name": "UnlockLevelCurve_Patch_XP_x0.5",
              "Folder": "UnlockLevelCurve_Patch_XP_x0.5",
              "Version": "72057594037927960",
              "Description": "Halves the XP requirement",
              "UUID": "e53ae4b5-a922-47ef-b69d-d55c5745a65b",
              "Created": "2024-07-22T20:17:59.1988778+02:00",
              "Dependencies": [],
              "Group": "dafe83c2-97b3-4b05-9aef-2e1cc2e1de98"
            }
          ],
          "MD5": "87758161b02ba6eb90ac2a6c92cd746f"
        }`

	modDataExpect := ModData{}
	err = json.Unmarshal([]byte(expectedData), &modDataExpect)
	if err != nil {
		t.Errorf("Failed to unmarshal expected data: %v", err)
	}

	if diff := cmp.Diff(modDataExpect, modData); diff != "" {
		t.Errorf("modData mismatch (-want +got):\n%s", diff)
	}

	if pakFile != "UnlockLevelCurve_Patch_XP_x0.5.pak" {
		t.Errorf("pakFile mismatch: %s", pakFile)
	}
}
