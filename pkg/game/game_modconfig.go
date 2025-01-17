package game

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type ModOrder struct {
	ID       string     `xml:"id,attr"`
	Children ModuleList `xml:"children"`
}

type ModuleList struct {
	Modules []Module `xml:"node"`
}

type Module struct {
	ID        string    `xml:"id,attr"`
	Attribute Attribute `xml:"attribute"`
}

type Mods struct {
	ID       string         `xml:"id,attr"`
	Children ModuleDescList `xml:"children"`
}

type ModuleDescList struct {
	Modules []ModuleShortDesc `xml:"node"`
}

type ModuleShortDesc struct {
	ID         string      `xml:"id,attr"`
	Attributes []Attribute `xml:"attribute"`
}

type ModSettingsLSX struct {
	XMLName xml.Name `xml:"save"`
	Version Version  `xml:"version"`
	Region  Region   `xml:"region"`
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

type Region struct {
	ID   string   `xml:"id,attr"`
	Root RootNode `xml:"node"`
}

type Mod struct {
	ID string `xml:"id,attr"`
	// Attributes for ModuleShortDesc
	Attributes []Attribute `xml:"attribute"`
}

type CategoryNode struct {
	ID       string `xml:"id,attr"`
	Children []Mod  `xml:"children>node"` // This will capture all child nodes
	// Attributes for ModuleShortDesc
	Attributes []Attribute `xml:"attribute"`
}

type RootNode struct {
	ID       string         `xml:"id,attr"`
	Children []CategoryNode `xml:"children>node"` // This will capture all child nodes
	// Attributes for ModuleShortDesc
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

	makeMoreCompact := func(xmlData string) string {
		// Change to short xml, ugly way. Replace ></KEY> with />
		src := strings.Split(xmlData, "\n")
		result := make([]string, 0, len(src))
		for _, line := range src {

			start := strings.Index(line, "></")
			if start == -1 {
				result = append(result, line)
				continue
			}
			// find first occurence of > after start
			end := strings.Index(line[start+1:], ">")
			if end == -1 {
				result = append(result, line)
				continue
			}
			// replace with />
			result = append(result, line[:start]+"/>"+line[start+end+2:])
		}

		return strings.Join(result, "\n")
	}

	addSpaceAfterLastAttribute := func(xmlData string) string {
		src := strings.Split(xmlData, "\n")
		result := make([]string, 0, len(src))
		for _, line := range src {
			if strings.HasSuffix(line, "\"/>") {
				result = append(result, line[:len(line)-3]+"\" />")
			} else {
				result = append(result, line)
			}

		}
		return strings.Join(result, "\n")
	}

	return addSpaceAfterLastAttribute(makeMoreCompact(xmlData))
}
