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

package bpmn

import (
	"encoding/xml"
)

// BpmnDefinitions is the root element of a BPMN 2.0 XML document.
// See http://www.omg.org/spec/BPMN/2.0/
type BpmnDefinitions struct {
	XMLName         xml.Name      `xml:"http://www.omg.org/spec/BPMN/20100524/MODEL definitions"`
	Bpmn            string        `xml:"xmlns:bpmn,attr"`
	Bpmndi          string        `xml:"xmlns:bpmndi,attr"`
	Dc              string        `xml:"xmlns:dc,attr"`
	Di              string        `xml:"xmlns:di,attr"`
	Bioc            string        `xml:"xmlns:bioc,attr"`
	Camunda         string        `xml:"xmlns:camunda,attr"`
	Id              string        `xml:"id,attr"`
	TargetNamespace string        `xml:"targetNamespace,attr"`
	Exporter        string        `xml:"exporter,attr"`
	ExporterVersion string        `xml:"exporterVersion,attr"`
	Processes       []BpmnProcess `xml:"process"`
}

// BpmnProcess is a BPMN 2.0 process.
type BpmnProcess struct {
	Id                     string             `xml:"id,attr"`
	Name                   string             `xml:"name,attr"`
	IsExecutable           bool               `xml:"isExecutable,attr"`
	CamundaVersionTag      string             `xml:"versionTag,http://camunda.org/schema/1.0/bpmn"`
	StartEvent             *BpmnStartEvent    `xml:"startEvent"`
	ServiceTask            []BpmnTask         `xml:"serviceTask"`
	UserTask               []BpmnTask         `xml:"userTask"`
	ManualTask             []BpmnTask         `xml:"manualTask"`
	ScriptTask             []BpmnTask         `xml:"scriptTask"`
	SendTask               []BpmnTask         `xml:"sendTask"`
	Task                   []BpmnTask         `xml:"task"`
	IntermediateThrowEvent []BpmnTask         `xml:"intermediateThrowEvent"`
	IntermediateCatchEvent []BpmnTask         `xml:"intermediateCatchEvent"`
	ExclusiveGateway       []BpmnGateway      `xml:"exclusiveGateway"`
	InclusiveGateway       []BpmnGateway      `xml:"inclusiveGateway"`
	ParallelGateway        []BpmnGateway      `xml:"parallelGateway"`
	EndEvent               []BpmnEndEvent     `xml:"endEvent"`
	SequenceFlow           []BpmnSequenceFlow `xml:"sequenceFlow"`
}

// BpmnStartEvent is a BPMN 2.0 start event.
type BpmnStartEvent struct {
	Id       string `xml:"id,attr"`
	Name     string `xml:"name,attr"`
	Outgoing string `xml:"outgoing"`
}

// BpmnTask is a BPMN 2.0 task.
type BpmnTask struct {
	Id            string `xml:"id,attr"`
	Name          string `xml:"name,attr"`
	Documentation string `xml:"documentation"`
	Incoming      string `xml:"incoming"`
	Outgoing      string `xml:"outgoing"`
}

// BpmnGateway is a BPMN 2.0 gateway.
type BpmnGateway struct {
	Id       string   `xml:"id,attr"`
	Name     string   `xml:"name,attr"`
	Incoming string   `xml:"incoming"`
	Outgoing []string `xml:"outgoing"`
}

// BpmnEndEvent is a BPMN 2.0 end event.
type BpmnEndEvent struct {
	Id                 string                    `xml:"id,attr"`
	Name               string                    `xml:"name,attr"`
	Incoming           []string                  `xml:"incoming"`
	SignalEventDefinit BpmnSignalEventDefinition `xml:"signalEventDefinition"`
}

// BpmnSignalEventDefinition is a BPMN 2.0 signal event definition.
type BpmnSignalEventDefinition struct {
	Id string `xml:"id,attr"`
}

// BpmnSequenceFlow is a BPMN 2.0 sequence flow.
type BpmnSequenceFlow struct {
	Id        string `xml:"id,attr"`
	SourceRef string `xml:"sourceRef,attr"`
	TargetRef string `xml:"targetRef,attr"`
	Name      string `xml:"name,attr"`
}

// ReadBpmn reads a BPMN 2.0 XML document.
func ReadBpmn(inputData []byte) (*BpmnDefinitions, error) {
	bpmnDefinitions := new(BpmnDefinitions)
	if err := xml.Unmarshal(inputData, bpmnDefinitions); err != nil {
		return nil, err
	}
	return bpmnDefinitions, nil
}
