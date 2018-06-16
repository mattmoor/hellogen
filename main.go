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

package main

import (
	"k8s.io/gengo/args"

	"github.com/golang/glog"

	"github.com/mattmoor/hellogen/generators"
)

// Feed in the whole example directory and see what it logs.
// go run main.go \
//    -v 5 \
//    --stderrthreshold INFO \
//    -i $(echo $(go list ./examples/...) | sed 's/ /,/g')
func main() {
	arguments := args.Default()

	// Override defaults.
	arguments.OutputFileBaseName = "hello_generated"

	// Run it.
	if err := arguments.Execute(
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		generators.Packages,
	); err != nil {
		glog.Fatalf("Error: %v", err)
	}
	glog.Info("Completed successfully.")
}
