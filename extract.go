// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

// extracts parses the given structure.xml file into corresponding
// XML types.
func extract(filename string) (*project, error) {
	input, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to Open: %v", err)
	}
	p := &project{}
	d := xml.NewDecoder(input)
	if err := d.Decode(p); err != nil {
		return nil, fmt.Errorf("unable to Decode: %v", err)
	}
	return p, nil
}

type project struct {
	Files             []file             `xml:"file,omitempty"`
	Name              string             `xml:"name,attr,omitempty"`
	ProjectNamespaces []projectNamespace `xml:"namespace"`
}

type projectNamespace struct {
	Name     string `xml:"name,attr,omitempty"`
	FullName string `xml:"full_name,attr,omitempty"`

	ProjectNamespaces []projectNamespace `xml:"namespace,omitempty"`
}

type file struct {
	Path string `xml:"path,attr,omitempty"`
	Hash string `xml:"hash,attr,omitempty"`

	Docblock         *docblock        `xml:"docblock,omitempty"`
	NamespaceAliases []namespaceAlias `xml:"namespace-alias,omitempty"`
	Tags             []tag            `xml:"tag,omitempty"`
	Class            *class           `xml:"class,omitempty"`
	Interface        *iface           `xml:"interface,omitempty"`
	Trait            *trait           `xml:"trait,omitempty"`
	Constants        []constant       `xml:"constant,omitempty"` // TODO
	Functions        []fn             `xml:"function,omitempty"` // TODO

	// TODO: includes, parse_markers
}

type docblock struct {
	Line int `xml:"line,attr,omitempty"`

	Description     string `xml:"description,omitempty"`
	LongDescription string `xml:"long-description,omitempty"` // Markdown!
	Tags            []tag  `xml:"tag,omitempty"`
}

func (d *docblock) summary() string {
	if d == nil {
		return ""
	}
	s := d.Description
	if d.LongDescription != "" {
		if d.Description != "" {
			s += "\n\n"
		}
		s += d.LongDescription
	}
	return s
}

type tag struct {
	Name        string `xml:"name,attr,omitempty"`
	Description string `xml:"description,attr,omitempty"`
	Variable    string `xml:"variable,attr,omitempty"`
	Type        string `xml:"type,attr,omitempty"`
	LinkOrRef   string `xml:"link,attr,omitempty"`
	Version     string `xml:"version,omitempty"`
	MethodName  string `xml:"method_name,omitempty"`
}

type namespaceAlias struct {
	Name string `xml:"name,attr,omitempty"`
}

type class struct {
	Final     bool   `xml:"final,attr,omitempty"`
	Abstract  bool   `xml:"abstract,attr,omitempty"`
	Namespace string `xml:"namespace,attr,omitempty"`
	Line      string `xml:"line,attr,omitempty"`

	Name       string     `xml:"name,omitempty"`
	FullName   string     `xml:"full_name,omitempty"`
	Docblock   *docblock  `xml:"docblock,omitempty"`
	Implements []string   `xml:"implements,omitempty"`
	Extends    string     `xml:"extends,omitempty"`
	Properties []property `xml:"property,omitempty"`
	Methods    []method   `xml:"method,omitempty"`
	Constants  []constant `xml:"constant,omitempty"`
}

type property struct {
	Namespace  string `xml:"namespace,attr,omitempty"`
	Line       string `xml:"line,attr,omitempty"`
	Visibility string `xml:"visibility,attr,omitempty"`

	Name          string    `xml:"name,omitempty"`
	FullName      string    `xml:"full_name,omitempty"`
	Docblock      *docblock `xml:"docblock,omitempty"`
	Default       string    `xml:"property,omitempty"`
	InheritedFrom string    `xml:"inherited_from,omitempty"`
}

type method struct {
	Final      bool   `xml:"final,attr,omitempty"`
	Abstract   bool   `xml:"abstract,attr,omitempty"`
	Static     bool   `xml:"static,attr,omitempty"`
	Namespace  string `xml:"namespace,attr,omitempty"`
	Line       string `xml:"line,attr,omitempty"`
	Visibility string `xml:"visibility,attr,omitempty"`

	Name          string     `xml:"name,omitempty"`
	FullName      string     `xml:"full_name,omitempty"`
	Value         string     `xml:"value,omitempty"` // TODO
	Docblock      *docblock  `xml:"docblock,omitempty"`
	Arguments     []argument `xml:"argument,omitempty"`
	InheritedFrom string     `xml:"inherited_from,omitempty"`
	// TODO: Value
}

type fn struct {
	Namespace string `xml:"namespace,attr,omitempty"`
	Line      string `xml:"line,attr,omitempty"`
	Package   string `xml:"package,attr,omitempty"` // TODO

	Name      string     `xml:"name,omitempty"`
	FullName  string     `xml:"full_name,omitempty"`
	Docblock  *docblock  `xml:"docblock,omitempty"`
	Arguments []argument `xml:"argument,omitempty"`
}

type argument struct {
	Line        string `xml:"line,attr,omitempty"`
	ByReference bool   `xml:"by_reference,attr"`

	Name    string `xml:"name,omitempty"`
	Type    string `xml:"type,omitempty"`
	Default string `xml:"property,omitempty"`
}

type iface struct {
	Namespace string `xml:"namespace,attr,omitempty"`
	Line      string `xml:"line,attr,omitempty"`
	Package   string `xml:"package,attr,omitempty"`

	Name      string     `xml:"name,omitempty"`
	FullName  string     `xml:"full_name,omitempty"`
	Docblock  *docblock  `xml:"docblock,omitempty"`
	Methods   []method   `xml:"method,omitempty"`
	Constants []constant `xml:"constant,omitempty"`
	Extends   string     `xml:"extends,omitempty"` // TODO
}

type constant struct {
	Namespace  string `xml:"namespace,attr,omitempty"`
	Line       string `xml:"line,attr,omitempty"`
	Visibility string `xml:"visibility,attr,omitempty"`

	Name          string    `xml:"name,omitempty"`
	FullName      string    `xml:"full_name,omitempty"`
	Value         string    `xml:"value,omitempty"`
	Docblock      *docblock `xml:"docblock,omitempty"`
	InheritedFrom string    `xml:"inherited_from,omitempty"`
}

type trait struct {
	Namespace string `xml:"namespace,attr,omitempty"`
	Line      string `xml:"line,attr,omitempty"`

	Name       string     `xml:"name,omitempty"`
	FullName   string     `xml:"full_name,omitempty"`
	Docblock   *docblock  `xml:"docblock,omitempty"`
	Properties []property `xml:"property,omitempty"`
	Methods    []method   `xml:"method,omitempty"`
}

func (d *docblock) status() string {
	if d == nil {
		return ""
	}
	for _, t := range d.Tags {
		if t.Name == "deprecated" {
			return "deprecated"
		}
	}
	return ""
}

func (d *docblock) param(name string) string {
	if d == nil {
		return ""
	}
	for _, t := range d.Tags {
		if t.Name == "param" && t.Variable == name {
			return t.Description
		}
	}
	return ""
}
