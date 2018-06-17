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

package generators

import (
	"fmt"
	"path/filepath"
	"strings"

	"k8s.io/gengo/args"
	"k8s.io/gengo/examples/set-gen/sets"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"

	_ "github.com/ghodss/yaml"
	"github.com/golang/glog"
)

// These are the comment tags that carry parameters for hello generation.
// TODO(mattmoor): Write something to parse these from comments.
const (
	// All of the tags we are sensitive too will have this root.
	baseTagName = "hello:"

	// Package tag for package-level configuration.
	pkgTagName = baseTagName + "package"

	// Function tag for function-level configuration.
	funcTagName = baseTagName + "function"

	// Type; tag for type-level configuration.
	typeTagName = baseTagName + "type"
)

// TODO: This is created only to reduce number of changes in a single PR.
// Remove it and use PublicNamer instead.
func theNamer() *namer.NameStrategy {
	return &namer.NameStrategy{
		Join: func(pre string, in []string, post string) string {
			return strings.Join(in, "_")
		},
		PrependPackageNames: 1,
	}
}

// NameSystems returns the name system used by the generators in this package.
func NameSystems() namer.NameSystems {
	return namer.NameSystems{
		"public": theNamer(),
		"raw":    namer.NewRawNamer("", nil),
	}
}

// DefaultNameSystem returns the default name system for ordering the types to be
// processed by the generators in this package.
func DefaultNameSystem() string {
	return "public"
}

func extractTags(tag string, comments []string) []string {
	tagVals := types.ExtractCommentTags("+", comments)[tag]
	if len(tagVals) == 0 {
		glog.V(5).Infof("No matching comment lines: %v", comments)
		return nil
	}

	glog.V(5).Infof("Got %d tagVals: %+v", len(tagVals), tagVals)
	return tagVals
}

func Packages(context *generator.Context, arguments *args.GeneratorArgs) generator.Packages {
	boilerplate, err := arguments.LoadGoBoilerplate()
	if err != nil {
		glog.Fatalf("Failed loading boilerplate: %v", err)
	}

	inputs := sets.NewString(context.Inputs...)
	packages := generator.Packages{}
	header := append([]byte(fmt.Sprintf("// +build !%s\n\n", arguments.GeneratedBuildTag)), boilerplate...)

	for i := range inputs {
		glog.V(5).Infof("Considering pkg %q", i)
		pkg := context.Universe[i]
		if pkg == nil {
			// If the input had no Go files, for example.
			continue
		}

		extractTags(pkgTagName, pkg.Comments)

		for _, t := range pkg.Functions {
			glog.V(5).Infof("  saw function %q (%+v)", t.Name.String(), t.Members)
			extractTags(funcTagName, t.CommentLines)
			fnt := t.Underlying // We're a DeclarationOf
			if fnt.Signature != nil {
				if fnt.Signature.Receiver != nil {
					glog.V(5).Infof("  receiver:", fnt.Signature.Receiver.Name.String())
				}
				glog.V(5).Info("  args:")
				for _, p := range fnt.Signature.Parameters {
					// This will be the types of the parameters,
					// it is unclear whether it is possible to access
					// the names of the formals.
					glog.V(5).Infof("  - %v", p.Name.String())
				}
				glog.V(5).Info("  results:")
				for _, res := range fnt.Signature.Results {
					// Same as Parameters, but named results are less typical.
					glog.V(5).Infof("  - %v", res.Name.String())
				}
			}
		}

		for _, t := range pkg.Types {
			glog.V(5).Infof("  saw type %q", t.Name.String())
			extractTags(typeTagName, t.CommentLines)
		}

		packages = append(packages,
			&generator.DefaultPackage{
				PackageName: strings.Split(filepath.Base(pkg.Path), ".")[0],
				PackagePath: pkg.Path,
				HeaderText:  header,
				GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
					return []generator.Generator{
						generator.DefaultGen{
							OptionalName: arguments.OutputFileBaseName,
							OptionalBody: []byte(`var this = "is a body"`),
						},
					}
				},
				FilterFunc: func(c *generator.Context, t *types.Type) bool {
					return t.Name.Package == pkg.Path
				},
			})

		packages = append(packages,
			&generator.DefaultPackage{
				PackageName: strings.Split(filepath.Base(pkg.Path), ".")[0],
				PackagePath: pkg.Path,
				HeaderText:  header,
				GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
					return []generator.Generator{
						// TODO(mattmoor): This is a patch on top of k8s.io/gengo (in vendor),
						// since it wasn't clear how to easily hook in new filetypes otherwise.
						generator.YAMLGen{
							OptionalName: arguments.OutputFileBaseName,
							Objects: []interface{}{
								map[string]string{
									"foo": "bar",
								},
							},
						},
					}
				},
				FilterFunc: func(c *generator.Context, t *types.Type) bool {
					return t.Name.Package == pkg.Path
				},
			})
	}
	return packages
}
