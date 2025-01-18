package modsettingslsx

import (
	"encoding/xml"
	"fmt"
	"github.com/GiGurra/bookeeper/pkg/config"
	"os"
	"regexp"
	"strings"
)

type ModSettingsXml struct {
	XMLName xml.Name   `xml:"save"`
	Version XmlVersion `xml:"version"`
	Region  XmlRegion  `xml:"region"`
}

func NewModSettingsXmlFromFile(filePath string) *ModSettingsXml {
	bs, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Errorf("failed to read file %s: %w", filePath, err))
	}

	var root ModSettingsXml
	err = xml.Unmarshal(bs, &root)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal xml: %w", err))
	}

	return &root
}

func Load(cfg *config.BaseConfig) *ModSettingsXml {
	return NewModSettingsXmlFromFile(config.Bg3ModsettingsFilePath(cfg))
}

type XmlVersion struct {
	Major    int `xml:"major,attr"`
	Minor    int `xml:"minor,attr"`
	Revision int `xml:"revision,attr"`
	Build    int `xml:"build,attr"`
}

type XmlRegion struct {
	ID         string        `xml:"id,attr"`
	Categories XmlCategories `xml:"node"`
}

type XmlCategories struct {
	ID         string         `xml:"id,attr"`
	Children   []XmlCategory  `xml:"children>node"` // This will capture all child nodes
	Attributes []XmlAttribute `xml:"attribute"`
}

type XmlCategory struct {
	ID         string         `xml:"id,attr"`
	Children   []XmlMod       `xml:"children>node"` // This will capture all child nodes
	Attributes []XmlAttribute `xml:"attribute"`
}

type XmlMod struct {
	ID         string         `xml:"id,attr"`
	Attributes []XmlAttribute `xml:"attribute"`
}

type XmlAttribute struct {
	ID    string `xml:"id,attr"`
	Value string `xml:"value,attr"`
	Type  string `xml:"type,attr"`
}

func (m *ModSettingsXml) ToXML() string {

	bs, err := xml.MarshalIndent(m, "", "  ")
	if err != nil {
		panic(fmt.Errorf("failed to marshal modsettingx xml/lsx: %w", err))
	}

	return xml.Header + makeBg3StyleXml(string(bs))
}

func (n *XmlCategories) GetXmlModOrder() []XmlMod {
	for _, child := range n.Children {
		if child.ID == "ModOrder" {
			return child.Children
		}
	}
	return nil
}

func (n *XmlCategories) GetXmlMods() []XmlMod {
	for _, child := range n.Children {
		if child.ID == "Mods" {
			return child.Children
		}
	}
	return nil
}

func (n *XmlCategories) SetXmlMods(newMods []XmlMod) {
	for i, child := range n.Children {
		if child.ID == "Mods" {
			n.Children[i].Children = newMods
		}
	}
}

func (n *XmlCategories) SetXmlModOrder(newOrder []XmlMod) {
	found := false
	for i, child := range n.Children {
		if child.ID == "ModOrder" {
			n.Children[i].Children = newOrder
			found = true
		}
	}

	if !found {
		newChildren := make([]XmlCategory, 0, len(n.Children)+1)
		newChildren = append(newChildren, XmlCategory{
			ID:       "ModOrder",
			Children: newOrder,
		})
		newChildren = append(newChildren, n.Children...)
		n.Children = newChildren
	}
}

func (n *XmlMod) GetXmlAttributeValue(id string) string {
	for _, attr := range n.Attributes {
		if attr.ID == id {
			return attr.Value
		}
	}
	return ""
}

func (n *ModSettingsXml) WithNewModSet() []XmlMod {
	return n.Region.Categories.GetXmlModOrder()
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
