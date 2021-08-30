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

// Command phpdocyaml transform phpDocumentor XML into DocFX YAML.
//
// phpdocyaml may change at any time in backwards incompatible ways.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// TODO: handle examples with comments like this:
//       //[snippet=gcs]

// TODO: consider generating namespace pages.

func main() {
	namespace := flag.String("namespace", "", "Required. Root namespace the docs are for. Will be the root of the TOC. Must not have a trailing \\")
	version := flag.String("version", "", "Required. The library version the docs are for")
	structure := flag.String("structure", "structure.xml", "Path to structure.xml file")
	outDir := flag.String("outdir", "out", "Where to write output")
	flag.Parse()

	if *structure == "" {
		fmt.Fprintf(os.Stderr, "Must set -structure\n\n")
		flag.Usage()
		os.Exit(1)
	}
	if *outDir == "" {
		fmt.Fprintf(os.Stderr, "Must set -outdir\n\n")
		flag.Usage()
		os.Exit(1)
	}
	if *namespace == "" {
		fmt.Fprintf(os.Stderr, "Must set -namespace\n\n")
		flag.Usage()
		os.Exit(1)
	}
	if strings.HasSuffix(*namespace, "\\") {
		fmt.Fprintf(os.Stderr, "-namespace must not end with \\\n\n")
		flag.Usage()
		os.Exit(1)
	}
	if *version == "" {
		fmt.Fprintf(os.Stderr, "Must set -version\n\n")
		flag.Usage()
		os.Exit(1)
	}

	p, err := extract(*structure)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to parse: %v", err)
		os.Exit(1)
	}

	pages, toc, err := transform(p, *namespace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to transform: %v", err)
		os.Exit(1)
	}

	if err := write(*outDir, pages, toc, *namespace, *version); err != nil {
		fmt.Fprintf(os.Stderr, "unable to write: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Success! Wrote %d pages, 1 TOC, and 1 docs.metadata.\n", len(pages))
}

func write(outDir string, pages map[string]*page, toc tableOfContents, namespace, version string) error {
	for path, p := range pages {
		if path == namespace {
			path = "index.yml"
		} else {
			// We know every path starts with namespace, or we would have errored out.
			path = strings.TrimPrefix(path, namespace)
			path = strings.ReplaceAll(path[1:], "\\", ".") // Trim leading \.
			path = filepath.Join(outDir, path+".yml")
		}
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return fmt.Errorf("os.MkdirAll: %v", err)
		}
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()
		fmt.Fprintln(f, "### YamlMime:UniversalReference")
		if err := yaml.NewEncoder(f).Encode(p); err != nil {
			return err
		}

		path = filepath.Join(outDir, "toc.yml")
		f, err = os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()
		fmt.Fprintln(f, "### YamlMime:TableOfContent")
		if err := yaml.NewEncoder(f).Encode(toc); err != nil {
			return err
		}
	}

	// Write the docuploader docs.metadata file. Not for DocFX.
	// See https://github.com/googleapis/docuploader/issues/11.
	// Example:
	/*
		update_time {
		  seconds: 1600048103
		  nanos: 183052000
		}
		name: "cloud.google.com/go"
		version: "v0.65.0"
		language: "go"
	*/

	// Replace the \s with .s, ignoring the leading \.
	namespaceDots := strings.ReplaceAll(namespace[1:], "\\", ".")

	path := filepath.Join(outDir, "docs.metadata")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	now := time.Now().UTC()
	fmt.Fprintf(f, `update_time {
	seconds: %d
	nanos: %d
}
name: %q
version: %q
language: "php"
`, now.Unix(), now.Nanosecond(), namespaceDots, version)
	return nil
}
