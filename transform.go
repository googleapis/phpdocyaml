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
	"fmt"
	"os"
	"sort"
	"strings"
)

// transform translates from the XML input types into YAML output types.
func transform(p *project, rootNamespace string) (map[string]*page, tableOfContents, error) {
	pages := map[string]*page{}
	// TODO: consider grouping by namespace and by deprecation status.
	// TODO: cross references.
	// TODO: visibility.
	tocRoot := &tocItem{Name: rootNamespace}

	for _, f := range p.Files {
		if strings.HasPrefix(f.Path, "tests") {
			continue
		}

		if len(f.Constants) > 0 {
			fmt.Fprintf(os.Stderr, "Found unhandled constants in %s", f.Path)
		}
		if len(f.Functions) > 0 {
			fmt.Fprintf(os.Stderr, "Found unhandled functions in %s", f.Path)
		}

		if f.Class != nil {
			classPage := &page{}
			uid := f.Class.FullName
			if !strings.HasPrefix(uid, rootNamespace) {
				return nil, nil, fmt.Errorf("found %q which does not belong to namespace %q", uid, rootNamespace)
			}
			tocRoot.addItem(&tocItem{
				Name:   strings.TrimPrefix(uid, rootNamespace),
				UID:    uid,
				Status: f.Class.Docblock.status(),
			})
			classItem := &item{
				UID:        uid,
				Name:       f.Class.Name,
				ID:         f.Class.Name,
				Summary:    f.Class.Docblock.summary(),
				Langs:      onlyPHP,
				Type:       "class",
				Status:     f.Class.Docblock.status(),
				Implements: f.Class.Implements,
			}
			classPage.addItem(classItem)
			if _, ok := pages[uid]; ok {
				return nil, nil, fmt.Errorf("found duplicate UID: %q", uid)
			}
			pages[uid] = classPage

			for _, p := range f.Class.Properties {
				if p.InheritedFrom != "" {
					classItem.InheritedMembers = append(classItem.InheritedMembers, p.FullName)
					continue
				}
				t := ""
				desc := ""
				if p.Docblock != nil {
					// TODO: not first tag
					if len(p.Docblock.Tags) > 0 {
						t = p.Docblock.Tags[0].Type
						if len(p.Docblock.Tags[0].Description) > 0 {
							fmt.Println("TODO: found a class.tag.description for", p.FullName)
						}
					}
					// TODO p.Default
					desc = p.Docblock.summary()
				}
				classItem.Properties = append(classItem.Properties, docfxProperty{
					Name:        p.Name,
					Type:        t,
					Description: desc,
				})
			}
			for _, m := range f.Class.Methods {
				if m.InheritedFrom != "" {
					classItem.InheritedMembers = append(classItem.InheritedMembers, m.FullName)
					continue
				}
				mUID := m.FullName
				mItem := &item{
					UID:        mUID,
					Name:       m.Name,
					ID:         m.Name,
					Parent:     uid,
					Summary:    m.Docblock.summary(),
					Langs:      onlyPHP,
					Type:       "method",
					Status:     m.Docblock.status(),
					Parameters: arguments(m),
				}
				classItem.addChild(child(mUID))
				classPage.addItem(mItem)
			}
			for _, c := range f.Class.Constants {
				if c.InheritedFrom != "" {
					classItem.InheritedMembers = append(classItem.InheritedMembers, c.FullName)
					continue
				}
				cUID := c.FullName
				cItem := &item{
					UID:     cUID,
					Name:    c.Name,
					ID:      c.Name,
					Parent:  uid,
					Syntax:  syntax{Content: c.Value},
					Summary: c.Docblock.summary(),
					Langs:   onlyPHP,
					Type:    "constant",
					Status:  c.Docblock.status(),
				}
				classItem.addChild(child(cUID))
				classPage.addItem(cItem)
			}
		}

		// TODO: update template to include traits. Leads to broken pages right now.
		if f.Trait != nil {
			traitPage := &page{}
			uid := f.Trait.FullName
			if !strings.HasPrefix(uid, rootNamespace) {
				return nil, nil, fmt.Errorf("found %q which does not belong to namespace %q", uid, rootNamespace)
			}
			tocRoot.addItem(&tocItem{
				Name:   strings.TrimPrefix(uid, rootNamespace),
				UID:    uid,
				Status: f.Trait.Docblock.status(),
			})
			traitItem := &item{
				UID:     uid,
				Name:    f.Trait.Name,
				ID:      f.Trait.Name,
				Summary: f.Trait.Docblock.summary(),
				Langs:   onlyPHP,
				Type:    "trait",
				Status:  f.Trait.Docblock.status(),
				// TODO: f.Trait.Properties,
			}
			traitPage.addItem(traitItem)
			if _, ok := pages[uid]; ok {
				return nil, nil, fmt.Errorf("found duplicate UID: %q", uid)
			}
			pages[uid] = traitPage

			for _, m := range f.Trait.Methods {
				if m.InheritedFrom != "" {
					traitItem.InheritedMembers = append(traitItem.InheritedMembers, m.FullName)
					continue
				}
				mUID := m.FullName
				mItem := &item{
					UID:        mUID,
					Name:       m.Name,
					ID:         m.Name,
					Parent:     uid,
					Summary:    m.Docblock.summary(),
					Langs:      onlyPHP,
					Type:       "method",
					Status:     m.Docblock.status(),
					Parameters: arguments(m),
				}
				traitItem.addChild(child(mUID))
				traitPage.addItem(mItem)
			}
		}

		if f.Interface != nil {
			interfacePage := &page{}
			uid := f.Interface.FullName
			if !strings.HasPrefix(uid, rootNamespace) {
				return nil, nil, fmt.Errorf("found %q which does not belong to namespace %q", uid, rootNamespace)
			}
			tocRoot.addItem(&tocItem{
				Name:   strings.TrimPrefix(uid, rootNamespace),
				UID:    uid,
				Status: f.Interface.Docblock.status(),
			})
			interfaceItem := &item{
				UID:     uid,
				Name:    f.Interface.Name,
				ID:      f.Interface.Name,
				Summary: f.Interface.Docblock.summary(),
				Langs:   onlyPHP,
				Type:    "interface",
				Status:  f.Interface.Docblock.status(),
			}
			interfacePage.addItem(interfaceItem)
			if _, ok := pages[uid]; ok {
				return nil, nil, fmt.Errorf("found duplicate UID: %q", uid)
			}
			pages[f.Interface.FullName] = interfacePage

			for _, m := range f.Interface.Methods {
				if m.InheritedFrom != "" {
					interfaceItem.InheritedMembers = append(interfaceItem.InheritedMembers, m.FullName)
					continue
				}
				mUID := m.FullName
				mItem := &item{
					UID:        mUID,
					Name:       m.Name,
					ID:         m.Name,
					Parent:     uid,
					Summary:    m.Docblock.summary(),
					Langs:      onlyPHP,
					Type:       "method",
					Status:     m.Docblock.status(),
					Parameters: arguments(m),
				}
				interfaceItem.addChild(child(mUID))
				interfacePage.addItem(mItem)
			}
			for _, c := range f.Interface.Constants {
				if c.InheritedFrom != "" {
					interfaceItem.InheritedMembers = append(interfaceItem.InheritedMembers, c.FullName)
					continue
				}
				cUID := c.FullName
				cItem := &item{
					UID:     cUID,
					Name:    c.Name,
					ID:      c.Name,
					Parent:  uid,
					Syntax:  syntax{Content: c.Value},
					Summary: c.Docblock.summary(),
					Langs:   onlyPHP,
					Type:    "constant",
					Status:  c.Docblock.status(),
				}
				interfaceItem.addChild(child(cUID))
				interfacePage.addItem(cItem)
			}
		}
	}
	sort.Slice(tocRoot.Items, func(i, j int) bool {
		return tocRoot.Items[i].UID < tocRoot.Items[j].UID
	})
	return pages, tableOfContents{tocRoot}, nil
}

func arguments(m method) []parameter {
	params := []parameter{}
	for _, a := range m.Arguments {
		params = append(params, parameter{
			Name:        a.Name,
			Type:        a.Type,
			Description: m.Docblock.param(a.Name),
			// TODO: by reference, default
		})
	}
	return params
}

// tableOfContents represents a TOC.
type tableOfContents []*tocItem

// tocItem is an item in a TOC.
type tocItem struct {
	UID    string     `yaml:"uid,omitempty"`
	Name   string     `yaml:"name,omitempty"`
	Items  []*tocItem `yaml:"items,omitempty"`
	Href   string     `yaml:"href,omitempty"`
	Status string     `yaml:"status,omitempty"`
}

func (t *tocItem) addItem(i *tocItem) {
	t.Items = append(t.Items, i)
}

// page represents a single DocFX page.
//
// There is one page per package.
type page struct {
	Items      []*item `yaml:"items"`
	References []*item `yaml:"references,omitempty"`
}

// child represents an item child.
type child string

// syntax represents syntax.
type syntax struct {
	Content string `yaml:"content,omitempty"`
}

type example struct {
	Content string `yaml:"content,omitempty"`
	Name    string `yaml:"name,omitempty"`
}

type docfxProperty struct {
	Type        string `yaml:"type,omitempty"`
	Name        string `yaml:"name,omitempty"`
	Description string `yaml:"description,omitempty"`
}

type parameter struct {
	Type        string `yaml:"type,omitempty"`
	Name        string `yaml:"name,omitempty"`
	Description string `yaml:"description,omitempty"`
}

// item represents a DocFX item.
type item struct {
	UID              string          `yaml:"uid"`
	Name             string          `yaml:"name,omitempty"`
	ID               string          `yaml:"id,omitempty"`
	Summary          string          `yaml:"summary,omitempty"`
	Parent           string          `yaml:"parent,omitempty"`
	Type             string          `yaml:"type,omitempty"`
	Langs            []string        `yaml:"langs,omitempty"`
	Syntax           syntax          `yaml:"syntax,omitempty"`
	Examples         []example       `yaml:"codeexamples,omitempty"`
	Children         []child         `yaml:"children,omitempty"`
	AltLink          string          `yaml:"alt_link,omitempty"`
	Status           string          `yaml:"status,omitempty"`
	Implements       []string        `yaml:"implements,omitempty"`
	InheritedMembers []string        `yaml:"inheritedMembers,omitempty"`
	Properties       []docfxProperty `yaml:"properties,omitempty"`
	Parameters       []parameter     `yaml:"parameters,omitempty"`
}

func (p *page) addItem(i *item) {
	p.Items = append(p.Items, i)
}

func (i *item) addChild(c child) {
	i.Children = append(i.Children, c)
}

var onlyPHP = []string{"php"}
