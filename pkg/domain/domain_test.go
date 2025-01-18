package domain

import (
	"encoding/xml"
	"github.com/GiGurra/bookeeper/pkg/modsettingslsx"
	"github.com/google/go-cmp/cmp"
	"log"
	"testing"
)

func TestXmlCategories_GetMods(t *testing.T) {
	var modSettingsLSX modsettingslsx.ModSettingsXml
	err := xml.Unmarshal([]byte(srcXmlData), &modSettingsLSX)
	if err != nil {
		log.Fatal(err)
	}

	// Access Mods
	mods := ListActiveMods(&modSettingsLSX)

	if len(mods) != 2 {
		t.Fatalf("expected 2 mods, got %d", len(mods))
	}

	expect := []Mod{
		{
			Folder:    "GustavDev",
			MD5:       "41a80562831251b58df743c05a7af21b",
			Name:      "GustavDev",
			UUID:      "28ac9ce2-2aba-8cda-b3b5-6e922f71b6b8",
			Version64: "144396877804629717",
		},
		{
			Folder:    "BasketEquipmentSFW",
			MD5:       "",
			Name:      "BasketEquipmentSFW",
			UUID:      "b200f917-43ec-45d9-9dff-ac6191d62388",
			Version64: "144115196665790673",
		},
	}

	if diff := cmp.Diff(expect, mods); diff != "" {
		t.Fatalf("mismatch (-want +got):\n%s", diff)
	}

	// reverse order
	SetActiveModsInBg3Cfg(&modSettingsLSX, []Mod{mods[1], mods[0]})

	actual := ListActiveMods(&modSettingsLSX)
	expect = []Mod{mods[1], mods[0]}

	if diff := cmp.Diff(expect, actual); diff != "" {
		t.Fatalf("mismatch (-want +got):\n%s", diff)
	}

}

func TestXmlRoot_GetMods_largeDataSet(t *testing.T) {
	modSettingsXml := modsettingslsx.NewModSettingsXmlFromFile("testdata/modsettings.lsx")
	mods := ListActiveMods(modSettingsXml)

	expect := []Mod{
		{
			Folder:        "GustavDev",
			MD5:           "33c1655f7ae3559b579ff3b9d2c45413",
			Name:          "GustavDev",
			PublishHandle: "0",
			UUID:          "28ac9ce2-2aba-8cda-b3b5-6e922f71b6b8",
			Version64:     "145100779997082619",
		},
		{
			Folder:        "UnlockSpells_f109f659-4ac2-7fc1-ade4-4b778c86cbd8",
			MD5:           "d3e21dc4eb538ff73d7fdec72b3441ef",
			Name:          "5e Spells: WIP",
			PublishHandle: "4358882",
			UUID:          "f109f659-4ac2-7fc1-ade4-4b778c86cbd8",
			Version64:     "36028799166447637",
		}, // maybe check more of them later
	}

	if diff := cmp.Diff(expect, mods[:2]); diff != "" {
		t.Fatalf("mismatch (-want +got):\n%s", diff)
	}
}

var srcXmlData = `<?xml version="1.0" encoding="UTF-8"?>
<save>
  <version major="4" minor="0" revision="8" build="2" />
  <region id="ModuleSettings">
    <node id="root">
      <children>
        <node id="ModOrder">
          <children>
            <node id="Module">
              <attribute id="UUID" value="28ac9ce2-2aba-8cda-b3b5-6e922f71b6b8" type="FixedString" />
            </node>
            <node id="Module">
              <attribute id="UUID" value="b200f917-43ec-45d9-9dff-ac6191d62388" type="FixedString" />
            </node>
          </children>
        </node>
        <node id="Mods">
          <children>
            <node id="ModuleShortDesc">
              <attribute id="Folder" value="GustavDev" type="LSString" />
              <attribute id="MD5" value="41a80562831251b58df743c05a7af21b" type="LSString" />
              <attribute id="Name" value="GustavDev" type="LSString" />
              <attribute id="UUID" value="28ac9ce2-2aba-8cda-b3b5-6e922f71b6b8" type="FixedString" />
              <attribute id="Version64" value="144396877804629717" type="int64" />
            </node>
            <node id="ModuleShortDesc">
              <attribute id="Folder" value="BasketEquipmentSFW" type="LSString" />
              <attribute id="MD5" value="" type="LSString" />
              <attribute id="Name" value="BasketEquipmentSFW" type="LSString" />
              <attribute id="UUID" value="b200f917-43ec-45d9-9dff-ac6191d62388" type="FixedString" />
              <attribute id="Version64" value="144115196665790673" type="int64" />
            </node>
          </children>
        </node>
      </children>
    </node>
  </region>
</save>`
