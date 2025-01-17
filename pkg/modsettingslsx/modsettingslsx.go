package modsettingslsx

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
)

type XmlRoot struct {
	XMLName xml.Name   `xml:"save"`
	Version XmlVersion `xml:"version"`
	Region  XmlRegion  `xml:"region"`
}

type XmlVersion struct {
	Major    int `xml:"major,attr"`
	Minor    int `xml:"minor,attr"`
	Revision int `xml:"revision,attr"`
	Build    int `xml:"build,attr"`
}

type XmlRegion struct {
	ID   string        `xml:"id,attr"`
	Root XmlCategories `xml:"node"`
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

func (m *XmlRoot) ToXML() string {

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

func (n *XmlMod) GetXmlAttributeValue(id string) string {
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
