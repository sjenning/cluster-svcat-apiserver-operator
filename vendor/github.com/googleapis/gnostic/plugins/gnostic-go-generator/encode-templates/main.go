// Copyright 2017 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// encode-templates is a support tool for building text templates
// directly into an application binary.
//
// It reads files from a "templates" directory, and generates
// a Go source file containing base64-encoded representations of those
// files. This allows these files to be directly compiled into the
// executable.
package main

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"
)

const TEMPLATES_GO = `// Copyright 2017 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

func templates() map[string]string {
	return map[string]string{ {{range .TemplateStrings}}
        "{{.Name}}": "{{.Encoding}}",{{end}}
    }
}`

type NamedTemplateString struct {
	Name     string // the name of the file to be generated by the template
	Encoding string // base64-encoded string
}

func main() {
	fileinfos, err := ioutil.ReadDir("templates")
	if err != nil {
		panic(err)
	}
	templateStrings := make([]*NamedTemplateString, 0)
	for _, fileinfo := range fileinfos {
		name := fileinfo.Name()
		if filepath.Ext(name) == ".tmpl" {
			data, _ := ioutil.ReadFile("templates/" + name)
			encoding := base64.StdEncoding.EncodeToString(data)
			templateStrings = append(templateStrings,
				&NamedTemplateString{
					Name:     strings.TrimSuffix(name, ".tmpl"),
					Encoding: encoding})
		}
	}
	t, err := template.New("").Parse(TEMPLATES_GO)
	if err != nil {
		panic(err)
	}
	f := new(bytes.Buffer)
	err = t.Execute(f, struct{ TemplateStrings []*NamedTemplateString }{templateStrings})
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("templates.go", f.Bytes(), 0644)
}
