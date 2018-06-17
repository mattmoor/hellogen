/*
Copyright 2018 Matt Moore

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package generator

import (
	"bytes"
	"io"
	"os"

	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"

	"github.com/ghodss/yaml"
)

const (
	YAMLFileType = "yaml"
)

// YAMLGen implements a do-nothing Generator.
//
// It can be used to implement static content files.
type YAMLGen struct {
	// OptionalName, if present, will be used for the generator's name, and
	// the filename (with ".yaml" appended).
	OptionalName string

	Objects []interface{}
}

func (d YAMLGen) Name() string                                        { return d.OptionalName }
func (d YAMLGen) Filter(*Context, *types.Type) bool                   { return true }
func (d YAMLGen) Namers(*Context) namer.NameSystems                   { return nil }
func (d YAMLGen) Imports(*Context) []string                           { return []string{} }
func (d YAMLGen) PackageVars(*Context) []string                       { return []string{} }
func (d YAMLGen) PackageConsts(*Context) []string                     { return []string{} }
func (d YAMLGen) GenerateType(*Context, *types.Type, io.Writer) error { return nil }
func (d YAMLGen) Filename() string                                    { return d.OptionalName + ".yaml" }
func (d YAMLGen) FileType() string                                    { return YAMLFileType }
func (d YAMLGen) Finalize(*Context, io.Writer) error                  { return nil }

func (d YAMLGen) Init(c *Context, w io.Writer) error {
	// Ideally this would use the pattern of:
	//   enc := yaml.NewEncoder(w)
	//   for _, obj := range d.Objects {
	//      obj.Encode(obj)
	//   }
	// https://github.com/go-yaml/yaml/blob/v2.2.1/yaml.go#L124
	// However, to do this with typed K8s objects requires the use
	// of a different yaml package that doesn't support this.
	for _, obj := range d.Objects {
		if _, err := w.Write([]byte("\n---\n")); err != nil {
			return err
		}
		b, err := yaml.Marshal(obj)
		if err != nil {
			return err
		}
		if _, err := w.Write(b); err != nil {
			return err
		}
	}
	return nil
}

func NewYAMLFile() FileType {
	return &yamlFileType{}
}

type yamlFileType struct{}

func (ft yamlFileType) AssembleFile(f *File, pathname string) error {
	destFile, err := os.Create(pathname)
	if err != nil {
		return err
	}
	defer destFile.Close()

	b := &bytes.Buffer{}
	et := NewErrorTracker(b)
	et.Write(f.Body.Bytes())
	if et.Error() != nil {
		return et.Error()
	}
	_, err = destFile.Write(b.Bytes())
	return err
}

func (ft yamlFileType) VerifyFile(f *File, pathname string) error {
	return nil
}
