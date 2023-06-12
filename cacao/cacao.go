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

package cacao

import (
	"crypto"
	_ "crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/cydarm/bpmn-to-cacao/bpmn"
	"github.com/golang/glog"
	"github.com/google/uuid"
)

const CACAO_NAMESPACE_UUID_STRING string = "aa7caf3a-d55a-4e9a-b34e-056215fba56a"
const CACAO_SPEC_VERSION_11 string = "1.1"
const CACAO_SPEC_VERSION_20 string = "2.0"

// CACAO step types
const CACAO_STEP_TYPE_START string = "start"
const CACAO_STEP_TYPE_END string = "end"
const CACAO_STEP_TYPE_11_STEP string = "step"
const CACAO_STEP_TYPE_ACTION string = "action"
const CACAO_STEP_TYPE_11_SINGLE string = "single"
const CACAO_STEP_TYPE_PLAYBOOK_ACTION string = "playbook-action"
const CACAO_STEP_TYPE_PARALLEL string = "parallel"
const CACAO_STEP_TYPE_IF_COND string = "if-condition"
const CACAO_STEP_TYPE_SWITCH_COND string = "switch-condition"
const CACAO_STEP_TYPE_WHILE_COND string = "while-condition"

// CACAO command types
const CACAO_COMMAND_TYPE_MANUAL string = "manual"
const CACAO_COMMAND_TYPE_BASH string = "bash"
const CACAO_COMMAND_TYPE_HTTP string = "http-api"
const CACAO_COMMAND_TYPE_SSH string = "ssh"
const CACAO_COMMAND_TYPE_CALDERA string = "caldera-cmd"
const CACAO_COMMAND_TYPE_ELASTIC string = "elastic"
const CACAO_COMMAND_TYPE_JUPYTER string = "juptyer"
const CACAO_COMMAND_TYPE_KESTREL string = "kestrel"
const CACAO_COMMAND_TYPE_OPENC2 string = "openc2-json"
const CACAO_COMMAND_TYPE_SIGMA string = "sigma"
const CACAO_COMMAND_TYPE_YARA string = "yara"

// CacaoPlaybook represents a CACAO playbook
type CacaoPlaybook struct {
	Type               string                      `json:"type"`
	SpecVersion        string                      `json:"spec_version"`
	ID                 string                      `json:"id"`
	Name               string                      `json:"name"`
	Description        string                      `json:"description,omitempty"`
	PlaybookTypes      []string                    `json:"playbook_types,omitempty"`
	CreatedBy          string                      `json:"created_by,omitempty"`
	Created            *time.Time                  `json:"created"`
	Modified           *time.Time                  `json:"modified"`
	Revoked            bool                        `json:"revoked"`
	ValidFrom          *time.Time                  `json:"valid_from,omitempty"`
	ValidUntil         *time.Time                  `json:"valid_until,omitempty"`
	DerivedFrom        string                      `json:"derived-from,omitempty"`
	Priority           int                         `json:"priority"`
	Severity           int                         `json:"severity"`
	Impact             int                         `json:"impact"`
	Labels             []string                    `json:"labels,omitempty"`
	ExternalReferences []ExternalReference         `json:"external_references,omitempty"`
	Markings           []string                    `json:"markings,omitempty"`
	PlaybookVariables  map[string]PlaybookVariable `json:"playbook_variables,omitempty"`
	WorkflowStart      string                      `json:"workflow_start"`
	WorkflowException  string                      `json:"workflow_exception,omitempty"`
	Workflow           map[string]Step             `json:"workflow"`
}

// ExternalReference represents an external reference embedded in a playbook
type ExternalReference struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Source      string `json:"source"`
	URL         string `json:"url"`
	Hash        string `json:"hash"`
	ExternalID  string `json:"external_id"`
}

// PlaybookVariable represents a variable that can be used in the playbook
type PlaybookVariable struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Value       string `json:"value"`
	Constant    bool   `json:"constant"`
}

// Step represents a step in the workflow
type Step struct {
	Type         string              `json:"type"`
	Name         string              `json:"name,omitempty"`
	OnCompletion string              `json:"on_completion,omitempty"`
	Condition    string              `json:"condition,omitempty"`
	OnTrue       string              `json:"on_true,omitempty"`
	OnFalse      string              `json:"on_false,omitempty"`
	Switch       string              `json:"switch,omitempty"`
	Cases        map[string][]string `json:"cases,omitempty"`
	NextSteps    []string            `json:"next_steps,omitempty"`
	Commands     []Command           `json:"commands,omitempty"`
	InArgs       []string            `json:"in_args,omitempty"`
}

// Command represents a command that can be executed
type Command struct {
	Type        string `json:"type"`
	Command     string `json:"command"`
	Description string `json:"description"`
}

// ProcessTasks processes the tasks in the BPMN and creates the appropriate steps
func ProcessTask(task bpmn.BpmnTask, commandType string, specVersion string, stepMap, nextStepMap map[string]string, cacaoPlaybook *CacaoPlaybook) {
	taskUuid := uuid.NewHash(crypto.SHA256.New(), uuid.MustParse(CACAO_NAMESPACE_UUID_STRING), []byte(task.Id), 5)
	stepType := CACAO_STEP_TYPE_ACTION // default to action - TODO: add support for other types
	if specVersion == CACAO_SPEC_VERSION_11 {
		stepType = CACAO_STEP_TYPE_11_STEP
	}
	stepId := fmt.Sprintf("%s--%s", stepType, taskUuid)
	onCompletion := stepMap[nextStepMap[fmt.Sprintf("%s:0", task.Id)]]
	if onCompletion == "" {
		// create another end task and link it
		endStepType := CACAO_STEP_TYPE_END
		if specVersion == CACAO_SPEC_VERSION_11 {
			endStepType = CACAO_STEP_TYPE_11_STEP
		}
		endEventUuid := uuid.New()
		stepId := fmt.Sprintf("%s--%s", endStepType, endEventUuid)
		cacaoPlaybook.Workflow[stepId] = Step{
			Type: CACAO_STEP_TYPE_END,
			Name: "End",
		}
		onCompletion = stepId
	}
	internalStepType := CACAO_STEP_TYPE_ACTION
	if specVersion == CACAO_SPEC_VERSION_11 {
		internalStepType = CACAO_STEP_TYPE_11_SINGLE
	}
	if cacaoPlaybook.WorkflowStart == stepId {
		internalStepType = CACAO_STEP_TYPE_START
	}
	cacaoPlaybook.Workflow[stepId] = Step{
		Type:         internalStepType,
		Name:         task.Name,
		OnCompletion: onCompletion,
		Commands: []Command{
			{
				Type:        commandType,
				Command:     task.Name,
				Description: task.Documentation,
			},
		},
	}
}

// ProcessGateway processes a gateway and creates the appropriate steps
func ProcessGateway(gateway bpmn.BpmnGateway, specVersion string, parallel bool, stepMap, nextStepMap map[string]string, cacaoPlaybook *CacaoPlaybook) {
	gatewayUuid := uuid.NewHash(crypto.SHA256.New(), uuid.MustParse(CACAO_NAMESPACE_UUID_STRING), []byte(gateway.Id), 5)
	parallelStepType := CACAO_STEP_TYPE_PARALLEL
	ifStepType := CACAO_STEP_TYPE_IF_COND
	switchStepType := CACAO_STEP_TYPE_SWITCH_COND
	endStepType := CACAO_STEP_TYPE_END
	if specVersion == CACAO_SPEC_VERSION_11 {
		parallelStepType = CACAO_STEP_TYPE_11_STEP
		ifStepType = CACAO_STEP_TYPE_11_STEP
		switchStepType = CACAO_STEP_TYPE_11_STEP
		endStepType = CACAO_STEP_TYPE_11_STEP
	}
	if parallel {
		stepId := fmt.Sprintf("%s--%s", parallelStepType, gatewayUuid)
		step := Step{
			Type: CACAO_STEP_TYPE_PARALLEL,
		}
		for i := 0; i < len(gateway.Outgoing); i++ {
			step.NextSteps = append(step.NextSteps, stepMap[nextStepMap[fmt.Sprintf("%s:%d", gateway.Id, i)]])
		}
		cacaoPlaybook.Workflow[stepId] = step
		return
	}
	// mangle the name to make it a valid variable name
	condition := strings.ReplaceAll(gateway.Name, " ", "_")
	condition = strings.ToLower(condition)
	condition = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_' {
			return r
		}
		return -1
	}, condition)
	if cacaoPlaybook.PlaybookVariables == nil {
		cacaoPlaybook.PlaybookVariables = make(map[string]PlaybookVariable)
	}
	if condition == "" {
		condition = gateway.Id
	}
	gatewayName := gateway.Name
	if gatewayName == "" {
		gatewayName = gateway.Id
	}
	cacaoPlaybook.PlaybookVariables[condition] = PlaybookVariable{
		Type:        "integer",
		Description: gatewayName,
		Value:       "0",
		Constant:    false,
	}
	if len(gateway.Outgoing) == 2 {
		stepId := fmt.Sprintf("%s--%s", ifStepType, gatewayUuid)
		onTrue := stepMap[nextStepMap[fmt.Sprintf("%s:%s", gateway.Id, "YES")]]
		if onTrue == "" {
			// create another end task and link it
			endEventUuid := uuid.New()
			stepId := fmt.Sprintf("%s--%s", endStepType, endEventUuid)
			cacaoPlaybook.Workflow[stepId] = Step{
				Type: CACAO_STEP_TYPE_END,
				Name: "End",
			}
			onTrue = stepId
		}
		onFalse := stepMap[nextStepMap[fmt.Sprintf("%s:%s", gateway.Id, "NO")]]
		if onFalse == "" {
			// create another end task and link it
			endEventUuid := uuid.New()
			stepId := fmt.Sprintf("%s--%s", endStepType, endEventUuid)
			cacaoPlaybook.Workflow[stepId] = Step{
				Type: CACAO_STEP_TYPE_END,
				Name: "End",
			}
			onFalse = stepId
		}
		cacaoPlaybook.Workflow[stepId] = Step{
			Type:      CACAO_STEP_TYPE_IF_COND,
			Condition: fmt.Sprintf("%s == 1", condition),
			InArgs:    []string{condition},
			Name:      gatewayName,
			OnTrue:    onTrue,
			OnFalse:   onFalse,
		}
	} else if len(gateway.Outgoing) > 2 {
		stepId := fmt.Sprintf("%s--%s", switchStepType, gatewayUuid)
		step := Step{
			Type:   CACAO_STEP_TYPE_SWITCH_COND,
			InArgs: []string{condition},
			Name:   gatewayName,
			Cases:  make(map[string][]string),
			Switch: condition,
		}
		for i := 0; i < len(gateway.Outgoing); i++ {
			// find map key
			for key, val := range nextStepMap {
				if strings.HasPrefix(key, gateway.Id) {
					nameString := strings.TrimPrefix(key, gateway.Id+":")
					step.Cases[nameString] = []string{stepMap[val]}
				}
			}
		}
		cacaoPlaybook.Workflow[stepId] = step
	} else {
		glog.Errorf("exclusive gateway %s has unexpected number of outgoing flows: %d", gateway.Id, len(gateway.Outgoing))
	}
}

// ConvertToCacao converts a BPMN definition to a CACAO playbook
func ConvertToCacao(bpmnDefinition *bpmn.BpmnDefinitions, specVersion string) (*CacaoPlaybook, error) {
	if len(bpmnDefinition.Processes) != 1 {
		return nil, errors.New(fmt.Sprintf("unexpected number of process definitions: %d", len(bpmnDefinition.Processes)))
	}
	bpmnProcess := bpmnDefinition.Processes[0]
	playbookUuid := uuid.NewHash(crypto.SHA256.New(), uuid.MustParse(CACAO_NAMESPACE_UUID_STRING), []byte(bpmnProcess.Id), 5)
	// map the BPMN ID of each step to the CACAO ID
	stepMap := make(map[string]string)
	startStepType := CACAO_STEP_TYPE_START
	endStepType := CACAO_STEP_TYPE_END
	actionStepType := CACAO_STEP_TYPE_ACTION
	// playbookActionStepType := CACAO_STEP_TYPE_PLAYBOOK_ACTION
	ifStepType := CACAO_STEP_TYPE_IF_COND
	parallelStepType := CACAO_STEP_TYPE_PARALLEL
	switchStepType := CACAO_STEP_TYPE_SWITCH_COND
	// whileStepType := CACAO_STEP_TYPE_WHILE_COND
	if specVersion == CACAO_SPEC_VERSION_11 {
		startStepType = CACAO_STEP_TYPE_11_STEP
		endStepType = CACAO_STEP_TYPE_11_STEP
		actionStepType = CACAO_STEP_TYPE_11_STEP
		// playbookActionStepType = CACAO_STEP_TYPE_11_STEP
		ifStepType = CACAO_STEP_TYPE_11_STEP
		parallelStepType = CACAO_STEP_TYPE_11_STEP
		switchStepType = CACAO_STEP_TYPE_11_STEP
		// whileStepType = CACAO_STEP_TYPE_11_STEP
	}
	// process possible start events
	startStepId := ""
	if bpmnProcess.StartEvent != nil {
		startEventUuid := uuid.NewHash(crypto.SHA256.New(), uuid.MustParse(CACAO_NAMESPACE_UUID_STRING), []byte(bpmnProcess.StartEvent.Id), 5)
		startStepId = fmt.Sprintf("%s--%s", startStepType, startEventUuid)
		stepMap[bpmnProcess.StartEvent.Id] = startStepId
	}
	for _, task := range bpmnProcess.IntermediateCatchEvent {
		taskUuid := uuid.NewHash(crypto.SHA256.New(), uuid.MustParse(CACAO_NAMESPACE_UUID_STRING), []byte(task.Id), 5)
		if startStepId != "" {
			// if start step is set, treat this as an action
			stepMap[task.Id] = fmt.Sprintf("%s--%s", actionStepType, taskUuid)
		} else {
			// if there is no start step, make this the start step
			startStepId = fmt.Sprintf("%s--%s", startStepType, taskUuid)
		}
	}
	// process end event
	for _, endEvent := range bpmnProcess.EndEvent {
		endEventUuid := uuid.NewHash(crypto.SHA256.New(), uuid.MustParse(CACAO_NAMESPACE_UUID_STRING), []byte(endEvent.Id), 5)
		stepMap[endEvent.Id] = fmt.Sprintf("%s--%s", endStepType, endEventUuid)
	}
	// process tasks
	for _, serviceTask := range bpmnProcess.ServiceTask {
		serviceTaskUuid := uuid.NewHash(crypto.SHA256.New(), uuid.MustParse(CACAO_NAMESPACE_UUID_STRING), []byte(serviceTask.Id), 5)
		stepMap[serviceTask.Id] = fmt.Sprintf("%s--%s", actionStepType, serviceTaskUuid)
	}
	for _, userTask := range bpmnProcess.UserTask {
		userTaskUuid := uuid.NewHash(crypto.SHA256.New(), uuid.MustParse(CACAO_NAMESPACE_UUID_STRING), []byte(userTask.Id), 5)
		stepMap[userTask.Id] = fmt.Sprintf("%s--%s", actionStepType, userTaskUuid)
	}
	for _, manualTask := range bpmnProcess.ManualTask {
		manualTaskUuid := uuid.NewHash(crypto.SHA256.New(), uuid.MustParse(CACAO_NAMESPACE_UUID_STRING), []byte(manualTask.Id), 5)
		stepMap[manualTask.Id] = fmt.Sprintf("%s--%s", actionStepType, manualTaskUuid)
	}
	for _, userTask := range bpmnProcess.ScriptTask {
		userTaskUuid := uuid.NewHash(crypto.SHA256.New(), uuid.MustParse(CACAO_NAMESPACE_UUID_STRING), []byte(userTask.Id), 5)
		stepMap[userTask.Id] = fmt.Sprintf("%s--%s", actionStepType, userTaskUuid)
	}
	for _, userTask := range bpmnProcess.SendTask {
		userTaskUuid := uuid.NewHash(crypto.SHA256.New(), uuid.MustParse(CACAO_NAMESPACE_UUID_STRING), []byte(userTask.Id), 5)
		stepMap[userTask.Id] = fmt.Sprintf("%s--%s", actionStepType, userTaskUuid)
	}
	for _, task := range bpmnProcess.Task {
		taskUuid := uuid.NewHash(crypto.SHA256.New(), uuid.MustParse(CACAO_NAMESPACE_UUID_STRING), []byte(task.Id), 5)
		stepMap[task.Id] = fmt.Sprintf("%s--%s", actionStepType, taskUuid)
	}
	for _, task := range bpmnProcess.IntermediateThrowEvent {
		taskUuid := uuid.NewHash(crypto.SHA256.New(), uuid.MustParse(CACAO_NAMESPACE_UUID_STRING), []byte(task.Id), 5)
		stepMap[task.Id] = fmt.Sprintf("%s--%s", actionStepType, taskUuid)
	}
	for _, endEvent := range bpmnProcess.EndEvent {
		endEventUuid := uuid.NewHash(crypto.SHA256.New(), uuid.MustParse(CACAO_NAMESPACE_UUID_STRING), []byte(endEvent.Id), 5)
		stepMap[endEvent.Id] = fmt.Sprintf("%s--%s", endStepType, endEventUuid)
	}
	for _, exclusiveGateway := range bpmnProcess.ExclusiveGateway {
		exclusiveGatewayUuid := uuid.NewHash(crypto.SHA256.New(), uuid.MustParse(CACAO_NAMESPACE_UUID_STRING), []byte(exclusiveGateway.Id), 5)
		if len(exclusiveGateway.Outgoing) == 2 {
			stepMap[exclusiveGateway.Id] = fmt.Sprintf("%s--%s", ifStepType, exclusiveGatewayUuid)
		} else if len(exclusiveGateway.Outgoing) > 2 {
			stepMap[exclusiveGateway.Id] = fmt.Sprintf("%s--%s", switchStepType, exclusiveGatewayUuid)
		} else {
			glog.Errorf("exclusive gateway %s has unexpected number of outgoing flows: %d", exclusiveGateway.Id, len(exclusiveGateway.Outgoing))
		}
	}
	for _, parallelGateway := range bpmnProcess.ParallelGateway {
		parallelGatewayUuid := uuid.NewHash(crypto.SHA256.New(), uuid.MustParse(CACAO_NAMESPACE_UUID_STRING), []byte(parallelGateway.Id), 5)
		stepMap[parallelGateway.Id] = fmt.Sprintf("%s--%s", parallelStepType, parallelGatewayUuid)
	}
	for _, inclusiveGateway := range bpmnProcess.InclusiveGateway {
		// TODO: add an if step for each outgoing flow
		parallelGatewayUuid := uuid.NewHash(crypto.SHA256.New(), uuid.MustParse(CACAO_NAMESPACE_UUID_STRING), []byte(inclusiveGateway.Id), 5)
		stepMap[inclusiveGateway.Id] = fmt.Sprintf("%s--%s", parallelStepType, parallelGatewayUuid)
	}
	// map the transitions, using BMPN ID and name (if present), to BPMN target,
	// eg.
	//     Activity_1g87yhd: -> Activity_0wagh2h
	//     Gateway_1hblfsj:Yes -> Activity_0vuc752
	//     Gateway_1g3qmkj:FILEHASH -> Event_0d4dl33
	nextStepMap := make(map[string]string)
	for _, sequenceFlow := range bpmnProcess.SequenceFlow {
		var nextStepMapKey string
		if sequenceFlow.Name != "" {
			nextStepMapKey = fmt.Sprintf("%s:%s", sequenceFlow.SourceRef, strings.ToUpper(sequenceFlow.Name))
		} else {
			for i := 0; true; i++ {
				// probe until we find an unused index
				nextStepMapKey = fmt.Sprintf("%s:%d", sequenceFlow.SourceRef, i)
				if _, found := nextStepMap[nextStepMapKey]; !found {
					break
				}
			}
		}
		nextStepMap[nextStepMapKey] = sequenceFlow.TargetRef
	}

	// create the playbook
	now := time.Now()
	cacaoPlaybook := &CacaoPlaybook{
		Type:          "playbook",
		SpecVersion:   specVersion,
		ID:            fmt.Sprintf("playbook--%s", playbookUuid),
		Name:          bpmnProcess.Name,
		Created:       &now,
		Modified:      &now,
		WorkflowStart: startStepId,
		Workflow:      make(map[string]Step),
	}

	// create start steps
	if bpmnProcess.StartEvent != nil {
		startId := stepMap[bpmnProcess.StartEvent.Id]
		cacaoPlaybook.Workflow[startId] = Step{
			Type:         CACAO_STEP_TYPE_START,
			Name:         bpmnProcess.StartEvent.Name,
			OnCompletion: stepMap[nextStepMap[fmt.Sprintf("%s:0", bpmnProcess.StartEvent.Id)]],
		}
	}
	for _, task := range bpmnProcess.IntermediateCatchEvent {
		ProcessTask(task, CACAO_COMMAND_TYPE_MANUAL, specVersion, stepMap, nextStepMap, cacaoPlaybook)
	}
	// create end steps
	for _, endEvent := range bpmnProcess.EndEvent {
		endId := stepMap[endEvent.Id]
		cacaoPlaybook.Workflow[endId] = Step{
			Type: CACAO_STEP_TYPE_END,
			Name: "End",
		}
	}
	// create the action steps
	for _, task := range bpmnProcess.ServiceTask {
		ProcessTask(task, CACAO_COMMAND_TYPE_HTTP, specVersion, stepMap, nextStepMap, cacaoPlaybook)
	}
	for _, task := range bpmnProcess.UserTask {
		ProcessTask(task, CACAO_COMMAND_TYPE_MANUAL, specVersion, stepMap, nextStepMap, cacaoPlaybook)
	}
	for _, task := range bpmnProcess.ManualTask {
		ProcessTask(task, CACAO_COMMAND_TYPE_MANUAL, specVersion, stepMap, nextStepMap, cacaoPlaybook)
	}
	for _, task := range bpmnProcess.ScriptTask {
		ProcessTask(task, CACAO_COMMAND_TYPE_BASH, specVersion, stepMap, nextStepMap, cacaoPlaybook)
	}
	for _, task := range bpmnProcess.SendTask {
		ProcessTask(task, CACAO_COMMAND_TYPE_BASH, specVersion, stepMap, nextStepMap, cacaoPlaybook)
	}
	for _, task := range bpmnProcess.Task {
		ProcessTask(task, CACAO_COMMAND_TYPE_MANUAL, specVersion, stepMap, nextStepMap, cacaoPlaybook)
	}
	for _, task := range bpmnProcess.IntermediateThrowEvent {
		ProcessTask(task, CACAO_COMMAND_TYPE_MANUAL, specVersion, stepMap, nextStepMap, cacaoPlaybook)
	}
	// create the branch steps
	for _, gateway := range bpmnProcess.ExclusiveGateway {
		ProcessGateway(gateway, specVersion, false, stepMap, nextStepMap, cacaoPlaybook)
	}
	for _, gateway := range bpmnProcess.ParallelGateway {
		ProcessGateway(gateway, specVersion, true, stepMap, nextStepMap, cacaoPlaybook)
	}
	for _, gateway := range bpmnProcess.InclusiveGateway {
		ProcessGateway(gateway, specVersion, false, stepMap, nextStepMap, cacaoPlaybook)
	}
	return cacaoPlaybook, nil
}
