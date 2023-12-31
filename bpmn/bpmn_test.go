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

package bpmn_test

import (
	"testing"

	"github.com/cydarm/bpmn-to-cacao/bpmn"
	"github.com/stretchr/testify/assert"
)

const inputDataTestString string = `<?xml version="1.0" encoding="UTF-8"?>
<bpmn:definitions xmlns:bpmn="http://www.omg.org/spec/BPMN/20100524/MODEL" xmlns:bpmndi="http://www.omg.org/spec/BPMN/20100524/DI" xmlns:dc="http://www.omg.org/spec/DD/20100524/DC" xmlns:di="http://www.omg.org/spec/DD/20100524/DI" xmlns:bioc="http://bpmn.io/schema/bpmn/biocolor/1.0" xmlns:camunda="http://camunda.org/schema/1.0/bpmn" id="Definitions_0xcsshl" targetNamespace="http://bpmn.io/schema/bpmn" exporter="Camunda Modeler" exporterVersion="3.7.2">
  <bpmn:process id="ProcessAV-EDRAlert" name="Process AV-EDR Alert" isExecutable="true" camunda:versionTag="Shareable_Workflow">
    <bpmn:startEvent id="StartEvent_1" name="Endpoint / AV Alerts on System">
      <bpmn:outgoing>Flow_1bgfopa</bpmn:outgoing>
    </bpmn:startEvent>
    <bpmn:serviceTask id="Activity_18ru9dm" name="SOAR Processes AV/EDR Alert">
      <bpmn:incoming>Flow_1bgfopa</bpmn:incoming>
      <bpmn:outgoing>Flow_017q5eb</bpmn:outgoing>
    </bpmn:serviceTask>
    <bpmn:exclusiveGateway id="Gateway_1hblfsj" name="Does alert meet policy threshold for COA review?">
      <bpmn:incoming>Flow_017q5eb</bpmn:incoming>
      <bpmn:outgoing>Flow_1jkwvw5</bpmn:outgoing>
      <bpmn:outgoing>Flow_1g10y9a</bpmn:outgoing>
    </bpmn:exclusiveGateway>
    <bpmn:endEvent id="Event_1ttzlep" name="Identify Systems and IOCs">
      <bpmn:incoming>Flow_1f96l27</bpmn:incoming>
      <bpmn:incoming>Flow_006qjb3</bpmn:incoming>
      <bpmn:signalEventDefinition id="SignalEventDefinition_1f935tp" />
    </bpmn:endEvent>
    <bpmn:exclusiveGateway id="Gateway_147ah6j" name="Does alert meet threshold for more data collection?">
      <bpmn:incoming>Flow_1g10y9a</bpmn:incoming>
      <bpmn:outgoing>Flow_031zd3o</bpmn:outgoing>
      <bpmn:outgoing>Flow_0w9k4zf</bpmn:outgoing>
    </bpmn:exclusiveGateway>
    <bpmn:serviceTask id="Activity_1g87yhd" name="SOAR Collects Internal Data on System">
      <bpmn:incoming>Flow_031zd3o</bpmn:incoming>
      <bpmn:outgoing>Flow_110b3rh</bpmn:outgoing>
    </bpmn:serviceTask>
    <bpmn:endEvent id="Event_18clgak" name="End">
      <bpmn:incoming>Flow_0w9k4zf</bpmn:incoming>
    </bpmn:endEvent>
    <bpmn:sequenceFlow id="Flow_017q5eb" sourceRef="Activity_18ru9dm" targetRef="Gateway_1hblfsj" />
    <bpmn:sequenceFlow id="Flow_1jkwvw5" name="Yes" sourceRef="Gateway_1hblfsj" targetRef="Activity_0vuc752" />
    <bpmn:sequenceFlow id="Flow_1g10y9a" name="No" sourceRef="Gateway_1hblfsj" targetRef="Gateway_147ah6j" />
    <bpmn:sequenceFlow id="Flow_031zd3o" name="Yes" sourceRef="Gateway_147ah6j" targetRef="Activity_1g87yhd" />
    <bpmn:sequenceFlow id="Flow_0w9k4zf" name="No" sourceRef="Gateway_147ah6j" targetRef="Event_18clgak" />
    <bpmn:sequenceFlow id="Flow_1bgfopa" sourceRef="StartEvent_1" targetRef="Activity_18ru9dm" />
    <bpmn:sequenceFlow id="Flow_110b3rh" sourceRef="Activity_1g87yhd" targetRef="Activity_0wagh2h" />
    <bpmn:serviceTask id="Activity_0wagh2h" name="SOAR Marks System Requires Monitoring">
      <bpmn:incoming>Flow_110b3rh</bpmn:incoming>
      <bpmn:outgoing>Flow_1f96l27</bpmn:outgoing>
    </bpmn:serviceTask>
    <bpmn:sequenceFlow id="Flow_1f96l27" sourceRef="Activity_0wagh2h" targetRef="Event_1ttzlep" />
    <bpmn:serviceTask id="Activity_0vuc752" name="SOAR Marks Case as Ready for COA Review">
      <bpmn:incoming>Flow_1jkwvw5</bpmn:incoming>
      <bpmn:outgoing>Flow_006qjb3</bpmn:outgoing>
    </bpmn:serviceTask>
    <bpmn:sequenceFlow id="Flow_006qjb3" sourceRef="Activity_0vuc752" targetRef="Event_1ttzlep" />
  </bpmn:process>
  <bpmndi:BPMNDiagram id="BPMNDiagram_1">
    <bpmndi:BPMNPlane id="BPMNPlane_1" bpmnElement="PrcssAV-EDRAlrt">
      <bpmndi:BPMNEdge id="Flow_006qjb3_di" bpmnElement="Flow_006qjb3">
        <di:waypoint x="660" y="177" />
        <di:waypoint x="832" y="177" />
      </bpmndi:BPMNEdge>
      <bpmndi:BPMNEdge id="Flow_1f96l27_di" bpmnElement="Flow_1f96l27">
        <di:waypoint x="790" y="300" />
        <di:waypoint x="850" y="300" />
        <di:waypoint x="850" y="195" />
      </bpmndi:BPMNEdge>
      <bpmndi:BPMNEdge id="Flow_110b3rh_di" bpmnElement="Flow_110b3rh">
        <di:waypoint x="660" y="300" />
        <di:waypoint x="690" y="300" />
      </bpmndi:BPMNEdge>
      <bpmndi:BPMNEdge id="Flow_1bgfopa_di" bpmnElement="Flow_1bgfopa">
        <di:waypoint x="215" y="177" />
        <di:waypoint x="300" y="177" />
      </bpmndi:BPMNEdge>
      <bpmndi:BPMNEdge id="Flow_0w9k4zf_di" bpmnElement="Flow_0w9k4zf">
        <di:waypoint x="480" y="325" />
        <di:waypoint x="480" y="410" />
        <di:waypoint x="832" y="410" />
        <bpmndi:BPMNLabel>
          <dc:Bounds x="488" y="332" width="15" height="14" />
        </bpmndi:BPMNLabel>
      </bpmndi:BPMNEdge>
      <bpmndi:BPMNEdge id="Flow_031zd3o_di" bpmnElement="Flow_031zd3o">
        <di:waypoint x="505" y="300" />
        <di:waypoint x="560" y="300" />
        <bpmndi:BPMNLabel>
          <dc:Bounds x="524" y="282" width="18" height="14" />
        </bpmndi:BPMNLabel>
      </bpmndi:BPMNEdge>
      <bpmndi:BPMNEdge id="Flow_1g10y9a_di" bpmnElement="Flow_1g10y9a">
        <di:waypoint x="480" y="202" />
        <di:waypoint x="480" y="275" />
        <bpmndi:BPMNLabel>
          <dc:Bounds x="488" y="215" width="15" height="14" />
        </bpmndi:BPMNLabel>
      </bpmndi:BPMNEdge>
      <bpmndi:BPMNEdge id="Flow_1jkwvw5_di" bpmnElement="Flow_1jkwvw5">
        <di:waypoint x="505" y="177" />
        <di:waypoint x="560" y="177" />
        <bpmndi:BPMNLabel>
          <dc:Bounds x="524" y="159" width="18" height="14" />
        </bpmndi:BPMNLabel>
      </bpmndi:BPMNEdge>
      <bpmndi:BPMNEdge id="Flow_017q5eb_di" bpmnElement="Flow_017q5eb">
        <di:waypoint x="400" y="177" />
        <di:waypoint x="455" y="177" />
      </bpmndi:BPMNEdge>
      <bpmndi:BPMNShape id="_BPMNShape_StartEvent_2" bpmnElement="StartEvent_1" bioc:stroke="rgb(67, 160, 71)" bioc:fill="rgb(200, 230, 201)">
        <dc:Bounds x="179" y="159" width="36" height="36" />
        <bpmndi:BPMNLabel>
          <dc:Bounds x="155" y="202" width="84" height="27" />
        </bpmndi:BPMNLabel>
      </bpmndi:BPMNShape>
      <bpmndi:BPMNShape id="Activity_18ru9dm_di" bpmnElement="Activity_18ru9dm">
        <dc:Bounds x="300" y="137" width="100" height="80" />
      </bpmndi:BPMNShape>
      <bpmndi:BPMNShape id="Gateway_1hblfsj_di" bpmnElement="Gateway_1hblfsj" isMarkerVisible="true">
        <dc:Bounds x="455" y="152" width="50" height="50" />
        <bpmndi:BPMNLabel>
          <dc:Bounds x="440" y="102" width="81" height="40" />
        </bpmndi:BPMNLabel>
      </bpmndi:BPMNShape>
      <bpmndi:BPMNShape id="Event_1ttzlep_di" bpmnElement="Event_1ttzlep" bioc:stroke="rgb(229, 57, 53)" bioc:fill="rgb(255, 205, 210)">
        <dc:Bounds x="832" y="159" width="36" height="36" />
        <bpmndi:BPMNLabel>
          <dc:Bounds x="809" y="121.5" width="81" height="27" />
        </bpmndi:BPMNLabel>
      </bpmndi:BPMNShape>
      <bpmndi:BPMNShape id="Gateway_147ah6j_di" bpmnElement="Gateway_147ah6j" isMarkerVisible="true">
        <dc:Bounds x="455" y="275" width="50" height="50" />
        <bpmndi:BPMNLabel>
          <dc:Bounds x="355" y="273" width="90" height="40" />
        </bpmndi:BPMNLabel>
      </bpmndi:BPMNShape>
      <bpmndi:BPMNShape id="Activity_1g87yhd_di" bpmnElement="Activity_1g87yhd">
        <dc:Bounds x="560" y="260" width="100" height="80" />
      </bpmndi:BPMNShape>
      <bpmndi:BPMNShape id="Event_18clgak_di" bpmnElement="Event_18clgak" bioc:stroke="rgb(229, 57, 53)" bioc:fill="rgb(255, 205, 210)">
        <dc:Bounds x="832" y="392" width="36" height="36" />
        <bpmndi:BPMNLabel>
          <dc:Bounds x="841" y="435" width="20" height="14" />
        </bpmndi:BPMNLabel>
      </bpmndi:BPMNShape>
      <bpmndi:BPMNShape id="Activity_1nj4o95_di" bpmnElement="Activity_0wagh2h">
        <dc:Bounds x="690" y="260" width="100" height="80" />
      </bpmndi:BPMNShape>
      <bpmndi:BPMNShape id="Activity_1hb1254_di" bpmnElement="Activity_0vuc752">
        <dc:Bounds x="560" y="137" width="100" height="80" />
      </bpmndi:BPMNShape>
    </bpmndi:BPMNPlane>
  </bpmndi:BPMNDiagram>
</bpmn:definitions>`

func TestReadBpmn(t *testing.T) {
	bpmnDefinitions, err := bpmn.ReadBpmn([]byte(inputDataTestString))
	if err != nil {
		t.Fatalf("could not read input: %s", err)
	}
	assert.NotNil(t, bpmnDefinitions)
	assert.Equal(t, 1, len(bpmnDefinitions.Processes))
	assert.Equal(t, "Process AV-EDR Alert", bpmnDefinitions.Processes[0].Name)
	assert.Equal(t, "Endpoint / AV Alerts on System", bpmnDefinitions.Processes[0].StartEvent.Name)
	assert.Equal(t, 4, len(bpmnDefinitions.Processes[0].ServiceTask))
	assert.Equal(t, 2, len(bpmnDefinitions.Processes[0].ExclusiveGateway))
	assert.Equal(t, 2, len(bpmnDefinitions.Processes[0].EndEvent))
}
