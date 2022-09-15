// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

package apitype

import (
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

// GitHubPullRequest describes a GitHub pull request that a stack update was associated with.
type GitHubPullRequest struct {
	Title    string `json:"title"`
	Number   int    `json:"number"`
	RepoSlug string `json:"repoSlug"`

	CreatedAt int64  `json:"createdAt"`
	ClosedAt  *int64 `json:"closedAt,omitempty"`
	WasMerged bool   `json:"wasMerged"`

	SourceRef string `json:"sourceRef"`
	TargetRef string `json:"targetRef"`

	MergeCommitSHA *string `json:"mergeCommitSHA,omitempty"`
}

// GetUpdateTimelineResponse is the response type returning the timeline for a given update.
type GetUpdateTimelineResponse struct {
	// Update is the information about the completed update for the stack.
	Update UpdateInfo `json:"update"`

	// CollatedUpdateEvents is the set of update events that are "relevant" to the update.
	// It will contain the sequences of previews that were performed in the same "group",
	// e.g. a GitHub PR, that lead to the final update.
	//
	// If the sequence of previews happened spanned when the stack was updated, e.g. the first
	// preview was at version N - 1 and the last preview was at version N + 1, then the update
	// to stack version N will be included.
	//
	// Will be empty if the update couldn't be associated with a particular sequence of previews.
	CollatedUpdateEvents []UpdateInfo `json:"collatedUpdateEvents"`

	// CollatedPullRequest is the source pull request that was merged, which triggered the stack's
	// update.
	CollatedPullRequest *GitHubPullRequest `json:"collatedPullRequest,omitempty"`

	// Previews contains all previews for the given stack update that weren't part of the collated
	// set. For example, from other GitHub pull requests that haven't been merged, previews ran on
	// developer machines, etc.
	Previews []UpdateInfo `json:"previews"`
}

// A WorkflowDefinition defines a workflow.
type WorkflowDefinition struct {
	Name     string                   `json:"name" yaml:"name"`
	Env      map[string]WorkflowValue `json:"env,omitempty" yaml:"env,omitempty"`
	Elements []WorkflowElement        `json:"elements" yaml:"elements"`
}

// A WorkflowElement defines a job or gate that makes up part of a workflow definition.
type WorkflowElement struct {
	Name string `json:"name" yaml:"name"`
	Kind string `json:"kind" yaml:"kind"`

	Job *JobDefinition `json:"-" yaml:"-"`
}

type workflowElement struct {
	Name string `json:"name" yaml:"name"`
	Kind string `json:"kind" yaml:"kind"`
}

type workflowJobElement struct {
	*JobDefinition

	Name string `json:"name" yaml:"name"`
	Kind string `json:"kind" yaml:"kind"`
}

const (
	WorkflowElementKindJob = "job"
)

func (v WorkflowElement) MarshalJSON() ([]byte, error) {
	switch v.Kind {
	case WorkflowElementKindJob:
		return json.Marshal(workflowJobElement{
			JobDefinition: v.Job,
			Name:          v.Name,
			Kind:          v.Kind,
		})
	default:
		return nil, fmt.Errorf("unknown workflow element kind %q", v.Kind)
	}
}

func (v *WorkflowElement) UnmarshalJSON(bytes []byte) error {
	var element workflowElement
	if err := json.Unmarshal(bytes, &element); err != nil {
		return err
	}
	switch element.Kind {
	case WorkflowElementKindJob:
		var job workflowJobElement
		if err := json.Unmarshal(bytes, &job); err != nil {
			return err
		}
		v.Name, v.Kind, v.Job = element.Name, element.Kind, job.JobDefinition
		return nil
	default:
		return fmt.Errorf("unknown workflow element kind %q", element.Kind)
	}
}

func (v WorkflowElement) MarshalYAML() (interface{}, error) {
	switch v.Kind {
	case WorkflowElementKindJob:
		return yaml.Marshal(workflowJobElement{
			JobDefinition: v.Job,
			Name:          v.Name,
			Kind:          v.Kind,
		})
	default:
		return nil, fmt.Errorf("unknown workflow element kind %q", v.Kind)
	}
}

func (v *WorkflowElement) UnmarshalYAML(node *yaml.Node) error {
	var element workflowElement
	if err := node.Decode(&element); err != nil {
		return err
	}
	switch element.Kind {
	case WorkflowElementKindJob:
		var job workflowJobElement
		if err := node.Decode(&job); err != nil {
			return err
		}
		v.Name, v.Kind, v.Job = element.Name, element.Kind, job.JobDefinition
		return nil
	default:
		return fmt.Errorf("unknown workflow element kind %q", element.Kind)
	}
}

// A JobDefinition defines a job that makes up part of a workflow definition.
type JobDefinition struct {
	OS           string                   `json:"os,omitempty" yaml:"os,omitempty"`
	Architecture string                   `json:"architecture,omitempty" yaml:"architecture,omitempty"`
	Image        *DockerImage             `json:"image,omitempty" yaml:"image,omitempty"`
	Timeout      WorkflowTimeout          `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Env          map[string]WorkflowValue `json:"env,omitempty" yaml:"env,omitempty"`
	Steps        []StepDefinition         `json:"steps" yaml:"steps"`
}

// A StepDefinition defines a step that makes up part of a job definition.
type StepDefinition struct {
	Name    string          `json:"name" yaml:"name"`
	Timeout WorkflowTimeout `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Run     string          `json:"run" yaml:"run"`
}

// A DockerImage describes a Docker image reference + optional credentials for use with aa job definition.
type DockerImage struct {
	Reference   string                  `json:"reference" yaml:"reference"`
	Credentials *DockerImageCredentials `json:"credentials,omitempty" yaml:"credentials,omitempty"`
}

// DockerImageCredentials describes the credentials needed to access a Docker repository.
type DockerImageCredentials struct {
	Username string        `json:"username" yaml:"username"`
	Password WorkflowValue `json:"password" yaml:"password"`
}

// A WorkflowTimeout describes a job or step timeout as part of a workflow definition.
type WorkflowTimeout time.Duration

func (v WorkflowTimeout) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(v).String())
}

func (v *WorkflowTimeout) UnmarshalJSON(bytes []byte) error {
	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*v = WorkflowTimeout(d)
	return nil
}

func (v WorkflowTimeout) MarshalYAML() (interface{}, error) {
	return time.Duration(v).String(), nil
}

func (v *WorkflowTimeout) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err != nil {
		return err
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*v = WorkflowTimeout(d)
	return nil
}

// A WorkflowValue describes a possibly-secret value.
type WorkflowValue struct {
	Value  string // Plaintext if Secret is false; ciphertext otherwise.
	Secret bool
}

type secretWorkflowValue struct {
	Secret string `json:"secret" yaml:"secret"`
}

func (v WorkflowValue) MarshalJSON() ([]byte, error) {
	if v.Secret {
		return json.Marshal(secretWorkflowValue{Secret: v.Value})
	}
	return json.Marshal(v.Value)
}

func (v *WorkflowValue) UnmarshalJSON(bytes []byte) error {
	var secret secretWorkflowValue
	if err := json.Unmarshal(bytes, &secret); err == nil {
		v.Value, v.Secret = secret.Secret, true
		return nil
	}

	var plaintext string
	if err := json.Unmarshal(bytes, &plaintext); err != nil {
		return err
	}
	v.Value, v.Secret = plaintext, false
	return nil
}

func (v WorkflowValue) MarshalYAML() (interface{}, error) {
	if v.Secret {
		return secretWorkflowValue{Secret: v.Value}, nil
	}
	return v.Value, nil
}

func (v *WorkflowValue) UnmarshalYAML(node *yaml.Node) error {
	var secret secretWorkflowValue
	if err := node.Decode(&secret); err == nil {
		v.Value, v.Secret = secret.Secret, true
		return nil
	}

	var plaintext string
	if err := node.Decode(&plaintext); err != nil {
		return err
	}
	v.Value, v.Secret = plaintext, false
	return nil
}

// WorkflowRunStatus describes the status of a workflow run.
type WorkflowRunStatus string

const (
	// WorkflowRunStatusRunning indicates that a workflow run is running.
	WorkflowRunStatusRunning = "running"
	// WorkflowRunStatusFailed indicates that a workflow run has failed.
	WorkflowRunStatusFailed = "failed"
	// WorkflowRunStatusSucceeded indicates that a workflow run has succeeded.
	WorkflowRunStatusSucceeded = "succeeded"
)

// WorkflowRun contains information about a workflow run.
type WorkflowRun struct {
	ID     string `json:"id"`
	OrgID  string `json:"orgId"`
	UserID string `json:"userId"`

	Status        WorkflowRunStatus `json:"status"`
	StartedAt     time.Time         `json:"startedAt"`
	LastUpdatedAt time.Time         `json:"lastUpdatedAt"`

	JobTimeout time.Time `json:"jobTimeout"`

	Jobs []JobRun `json:"jobs"`
}

// JobStatus describes the status of a job run.
type JobStatus string

const (
	// JobStatusNotStarted indicates that a job has not yet started.
	JobStatusNotStarted = "not-started"
	// JobStatusAccepted indicates that a job has been accepted for execution, but is not yet running.
	JobStatusAccepted = "accepted"
	// JobStatusNotStarted indicates that a job is running.
	JobStatusRunning = "running"
	// JobStatusNotStarted indicates that a job has failed.
	JobStatusFailed = "failed"
	// JobStatusNotStarted indicates that a job has succeeded.
	JobStatusSucceeded = "succeeded"
)

// JobRun contains information about a job run.
type JobRun struct {
	Status      JobStatus     `json:"status"`
	Started     time.Time     `json:"started"`
	LastUpdated time.Time     `json:"lastUpdated"`
	Timeout     time.Duration `json:"timeout"`
	Steps       []StepRun     `json:"steps"`
	WorkerInfo  interface{}   `json:"worker,omitempty"`
}

// StepStatus describes the status of a step in a workflow job run.
type StepStatus string

const (
	// StepStatusNotStarted indicates that a step has not yet started.
	StepStatusNotStarted = "not-started"
	// StepStatusNotStarted indicates that a step is running.
	StepStatusRunning = "running"
	// StepStatusNotStarted indicates that a step has failed.
	StepStatusFailed = "failed"
	// StepStatusNotStarted indicates that a step has succeeded.
	StepStatusSucceeded = "succeeded"
)

// StepRun contains information about a step run.
type StepRun struct {
	Name        string     `json:"name"`
	Status      StepStatus `json:"status"`
	Started     time.Time  `json:"started"`
	LastUpdated time.Time  `json:"lastUpdated"`
}

type UpdateJobWorkerRequest struct {
	Worker interface{} `json:"worker"`
}

type UpdateStepStatusRequest struct {
	Status StepStatus `json:"status"`
}

type StartWorkflowRunResponse struct {
	RunID string `json:"runID"`
}

type StepLogLine struct {
	Timestamp time.Time `json:"t"`
	Line      string    `json:"l"`
}

type AppendStepLogsRequest struct {
	Offset int           `json:"offset"`
	Lines  []StepLogLine `json:"lines"`
}

type GetStepLogsResponse struct {
	Lines []StepLogLine `json:"lines"`
}
