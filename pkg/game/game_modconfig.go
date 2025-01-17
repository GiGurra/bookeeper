package game

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
)

type ModSettingsLSX struct {
	XMLName xml.Name `xml:"save"`
	Version Version  `xml:"version"`
	Region  Region   `xml:"region"`
}

type Region struct {
	ID   string   `xml:"id,attr"`
	Root RootNode `xml:"node"`
}

type RootNode struct {
	ID         string         `xml:"id,attr"`
	Children   []CategoryNode `xml:"children>node"` // This will capture all child nodes
	Attributes []Attribute    `xml:"attribute"`
}

func (m *ModSettingsLSX) ToXML() string {

	bs, err := xml.MarshalIndent(m, "", "  ")
	if err != nil {
		panic(fmt.Errorf("failed to marshal modsettingx xml/lsx: %w", err))
	}

	return xml.Header + makeBg3StyleXml(string(bs))
}

type Version struct {
	Major    int `xml:"major,attr"`
	Minor    int `xml:"minor,attr"`
	Revision int `xml:"revision,attr"`
	Build    int `xml:"build,attr"`
}

type Mod struct {
	ID         string      `xml:"id,attr"`
	Attributes []Attribute `xml:"attribute"`
}

type CategoryNode struct {
	ID         string      `xml:"id,attr"`
	Children   []Mod       `xml:"children>node"` // This will capture all child nodes
	Attributes []Attribute `xml:"attribute"`
}

type Attribute struct {
	ID    string `xml:"id,attr"`
	Value string `xml:"value,attr"`
	Type  string `xml:"type,attr"`
}

func (n *RootNode) GetModOrder() []Mod {
	for _, child := range n.Children {
		if child.ID == "ModOrder" {
			return child.Children
		}
	}
	return nil
}

func (n *RootNode) GetMods() []Mod {
	for _, child := range n.Children {
		if child.ID == "Mods" {
			return child.Children
		}
	}
	return nil
}

func (n *Mod) GetAttributeValue(id string) string {
	for _, attr := range n.Attributes {
		if attr.ID == id {
			return attr.Value
		}
	}
	return ""
}

func makeBg3StyleXml(xmlData string) string {

	// TODO: Remove if it turns out we don't need this in the BG3 xml

	endOfLineToRelace := regexp.MustCompile(`></[A-Za-z]+>$`)

	makeBgish := func(xmlData string) string {
		src := strings.Split(xmlData, "\n")
		result := make([]string, 0, len(src))
		for _, line := range src {
			match := endOfLineToRelace.FindStringIndex(line)
			if match != nil {
				result = append(result, line[:match[0]]+" />")
			} else {
				result = append(result, line)
			}
		}

		return strings.Join(result, "\n")
	}

	return makeBgish(xmlData)
}
