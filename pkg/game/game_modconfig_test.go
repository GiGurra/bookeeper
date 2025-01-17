package game

import (
	"encoding/xml"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"log"
	"strings"
	"testing"
)

func TestNode_GetMods(t *testing.T) {
	var modSettingsLSX ModSettingsLSX
	err := xml.Unmarshal([]byte(srcXmlData), &modSettingsLSX)
	if err != nil {
		log.Fatal(err)
	}

	// Access ModOrder modules
	modOrder := modSettingsLSX.Region.Root.GetModOrder()
	for _, module := range modOrder {
		uuid := module.GetAttributeValue("UUID")
		fmt.Printf("Module UUID: %s\n", uuid)
	}

	// Access Mods
	mods := modSettingsLSX.Region.Root.GetMods()
	for _, mod := range mods {
		folder := mod.GetAttributeValue("Folder")
		name := mod.GetAttributeValue("Name")
		uuid := mod.GetAttributeValue("UUID")
		fmt.Printf("Mod: %s (%s) - UUID: %s\n", name, folder, uuid)
	}

	// marshal back the xml
	expected := srcXmlData
	actual := modSettingsLSX.ToXML()

	if diff := cmp.Diff(expected, actual); diff != "" {
		linesExpected := strings.Split(expected, "\n")
		linesActual := strings.Split(actual, "\n")
		for i := range linesExpected {
			lineExpected := linesExpected[i]
			lineActual := linesActual[i]
			if diff := cmp.Diff(lineExpected, lineActual); diff != "" {
				fmt.Printf("expected: %s\n", lineExpected)
				fmt.Printf("  actual: %s\n", lineActual)
			}
		}
		t.Fatalf("mismatch (-want +got):\n%s", diff)
	}

	fmt.Printf("XML: %s\n", actual)

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
