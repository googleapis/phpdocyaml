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
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var updateGoldens bool

func TestMain(m *testing.M) {
	flag.BoolVar(&updateGoldens, "update-goldens", false, "Update the golden files")
	flag.Parse()
	os.Exit(m.Run())
}

func TestGoldens(t *testing.T) {
	structure := "testdata/structure.xml"
	gotDir := "testdata/out"
	goldenDir := "testdata/golden"
	namespace := `\Google\Cloud\Vision`

	p, err := extract(structure)
	if err != nil {
		t.Fatalf("unable to parse: %v", err)
	}

	pages, toc, err := transform(p, namespace)
	if err != nil {
		t.Fatalf("unable to transform: %v", err)
	}

	ignoreFiles := map[string]bool{"docs.metadata": true}

	if updateGoldens {
		os.RemoveAll(goldenDir)

		if err := write(goldenDir, pages, toc, namespace, "1.0.0"); err != nil {
			t.Fatalf("unable to write: %v", err)
		}

		for ignore := range ignoreFiles {
			if err := os.Remove(filepath.Join(goldenDir, ignore)); err != nil {
				t.Fatalf("Remove: %v", err)
			}
		}

		t.Logf("Successfully updated goldens in %s", goldenDir)

		return
	}

	if err := write(gotDir, pages, toc, namespace, "1.0.0"); err != nil {
		t.Fatalf("unable to write: %v", err)
	}

	gotFiles, err := ioutil.ReadDir(gotDir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}

	goldens, err := ioutil.ReadDir(goldenDir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}

	if got, want := len(gotFiles)-len(ignoreFiles), len(goldens); got != want {
		t.Fatalf("parse & write got %d files in %s, want %d ignoring %v", got, gotDir, want, ignoreFiles)
	}

	for _, golden := range goldens {
		if golden.IsDir() {
			continue
		}
		gotPath := filepath.Join(gotDir, golden.Name())
		goldenPath := filepath.Join(goldenDir, golden.Name())

		gotContent, err := ioutil.ReadFile(gotPath)
		if err != nil {
			t.Fatalf("ReadFile: %v", err)
		}

		goldenContent, err := ioutil.ReadFile(goldenPath)
		if err != nil {
			t.Fatalf("ReadFile: %v", err)
		}

		if string(gotContent) != string(goldenContent) {
			t.Errorf("got %s is different from expected %s", gotPath, goldenPath)
		}
	}
}
