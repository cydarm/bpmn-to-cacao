/*
 * Copyright 2023 Cydarm Technologies Pty Ltd, https://cydarm.com/
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 		http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cydarm/bpmn-to-cacao/bpmn"
	"github.com/cydarm/bpmn-to-cacao/cacao"
	"github.com/golang/glog"
)

var outDir string
var cacaoSpecVersion string

func init() {
	flag.StringVar(&outDir, "output-dir", ".", "Specify a directory for output")
	flag.StringVar(&cacaoSpecVersion, "cacao-spec", "1.1", "Specify a CACAO spec version (1.1 or 2.0)")
}

func main() {
	flag.Set("stderrthreshold", "INFO")
	flag.Parse()
	inputFiles := flag.Args()
	// validate output directory
	dirInfo, err := os.Stat(outDir)
	if err != nil {
		glog.Fatalf("Error parsing argument %s: %s", outDir, err)
	}
	if !dirInfo.IsDir() {
		glog.Fatalf("Error %s is not a directory", outDir)
	}
	if len(inputFiles) == 0 {
		glog.Fatalf("No input files were specified")
	}
	for _, inputFile := range inputFiles {
		glog.Infof("Processing %s", inputFile)
		lstat, err := os.Lstat(inputFile)
		if err != nil {
			glog.Errorf("could not lstat %s", inputFile)
		}
		inputFileBaseName := lstat.Name()
		inputData, err := ioutil.ReadFile(inputFile)
		if err != nil {
			glog.Errorf("could not read %s", inputFile)
		}
		bpmnDefinition, err := bpmn.ReadBpmn(inputData)
		if err != nil {
			glog.Errorf("processing input file failed: %s", err)
			continue
		}
		cacaoOutput, err := cacao.ConvertToCacao(bpmnDefinition, cacaoSpecVersion)
		if err != nil {
			glog.Errorf("cacao convertion failed: %s", err)
			continue
		}
		outBytes, err := json.MarshalIndent(cacaoOutput, "", "    ")
		if err != nil {
			glog.Errorf("marshaling JSON failed: %s", err)
			continue
		}
		outputFileName := fmt.Sprintf("%s/%s.cacao.json", outDir, inputFileBaseName)
		if err := os.WriteFile(outputFileName, outBytes, 0644); err != nil {
			glog.Errorf("writing file %s failed: %s", outputFileName, err)
			continue
		}
		glog.Infof("Wrote output to %s", outputFileName)
	}
}
