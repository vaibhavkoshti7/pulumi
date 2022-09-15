// Copyright 2016-2022, Pulumi Corporation.  All rights reserved.

package apitype

import (
	"encoding/json"
	"fmt"
	"time"
)

// PulumiOperation describes what operation to perform on the
// stack as defined in the Job spec.
type PulumiOperation string

// The possible operations we can deploy.
const (
	Update  PulumiOperation = "update"
	Preview PulumiOperation = "preview"
	Destroy PulumiOperation = "destroy"
	Refresh PulumiOperation = "refresh"
)

// GetDeploymentResponse is the response from the API when getting
// a single Deployment.
type GetDeploymentResponse struct {
	// ID is the internal PK of the Deployment.
	ID string `json:"id"`

	// Created defines when the Deployment was created.
	Created string `json:"created"`

	// Created defines when the corresponding WorkflowRun was modified.
	Modified string `json:"modified"`

	// Status is the current status of the workflow runID.
	Status string `json:"status"`

	// Version is the ordinal ID for the stack
	Version int `json:"version"`

	// Username who created the deployment.
	Username string `json:"username"`

	// Name of the user who created the deployment.
	Name string `json:"name"`

	// Jobs make up all the Jobs of the Deployment.
	Jobs []DeploymentJob `json:"jobs"`
}

// ListDeploymentResponse is the response from the API when listing Deployments.
type ListDeploymentResponse struct {
	// ID is the internal PK of the Deployment.
	ID string `json:"id"`

	// Created defines when the Deployment was created.
	Created string `json:"created"`

	// Created defines when the corresponding WorkflowRun was modified.
	Modified string `json:"modified"`

	// Status is the current status of the workflow runID.
	Status string `json:"status"`

	// Version is the ordinal ID for the stack
	Version int `json:"version"`

	// Username who created the deployment.
	UserName string `json:"username"`

	// Name of the user who created the deployment.
	Name string `json:"name"`
}

type DeploymentJob struct {
	Status      JobStatus `json:"status"`
	Started     time.Time `json:"started"`
	LastUpdated time.Time `json:"lastUpdated"`
	Steps       []StepRun `json:"steps"`
}

type JobManifestResponse struct {
	Executor *ExecutorContext `json:"executorContext"`

	// NOTE: We don't always need the source context. For example running `pulumi refresh`
	// or `pulumi destroy` does not need source.
	Source *SourceContext `json:"sourceContext"`

	// Operation defines the options that the executor will use to run the Pulumi commands.
	Operation OperationContextResponse `json:"operationContext"`
}

// OperationContextResponse describes what to do.
type OperationContextResponse struct {
	// PreRunCommands is an optional list of arbitrary commands to run before Pulumi
	// is invoked.
	PreRunCommands []string `json:"preRunCommands"`

	// StackIdentity is the fully-qualified stack to operate on. (org/project/stack)
	StackIdentity string `json:"stackIdentity"`

	// Operation is what we plan on doing.
	Operation PulumiOperation `json:"operation"`

	// EnvironmentVariables contains environment variables to be applied during the execution.
	EnvironmentVariables map[string]string `json:"environmentVariables"`
}

// CreateDeploymentRequest defines the request payload that is expected when
// creating a new deployment.
type CreateDeploymentRequest struct {
	// Executor defines options that the executor is going to use to run the job.
	Executor *ExecutorContext `json:"executorContext"`

	// Source defines how the source code to the Pulumi program will be gathered.
	Source *SourceContext `json:"sourceContext,omitempty"`

	// Operation defines the options that the executor will use to run the Pulumi commands.
	Operation OperationContext `json:"operationContext"`
}

func (cdr *CreateDeploymentRequest) ToJSON() (string, error) {
	data, err := json.Marshal(cdr)
	return string(data), err
}

// Validate ensures that the CreateDeploymentRequest is valid and all required fields are set.
func (cdr *CreateDeploymentRequest) Validate() error {
	// Validate the source context.
	if cdr.Source != nil {
		if git := cdr.Source.Git; git != nil {
			if err := git.Validate(); err != nil {
				return err
			}
		}
	}

	// Validate the operation context.
	if err := cdr.Operation.Validate(); err != nil {
		return err
	}

	// Validate the operation type has a source context.
	// Pulumi operations of type "update" & "refresh" require a source context.
	switch cdr.Operation.Operation {
	case Update, Refresh:
		if cdr.Source == nil {
			return fmt.Errorf("source context is required for operation type %q", cdr.Operation.Operation)
		}
	}

	return nil
}

// CreateDeploymentResponse defines the response given when a new Deployment is created.
type CreateDeploymentResponse struct {
	// ID represents the generated Deployment ID.
	ID string `json:"id"`
}

// ParseCreateDeploymentResponse parses the CreateDeploymentResponse JSON to a struct.
func ParseCreateDeploymentResponse(input string) (*CreateDeploymentResponse, error) {
	var deploymentResponse *CreateDeploymentResponse
	if err := json.Unmarshal([]byte(input), &deploymentResponse); err != nil {
		return nil, fmt.Errorf("unmarshalling deploymentResponse manifest: %w", err)
	}
	return deploymentResponse, nil
}

// GetDeploymentRequest contains the current status of a deployment.
type GetDeploymentRequest struct {
	// ID is the internal PK of the Deployment.
	ID string `json:"id"`
	// Created defines when the Deployment was created.
	Created  string `json:"created"`
	Modified string `json:"modified"`
	Status   string `json:"status"`
}

// CreateDeploymentEnvironmentVariableRequest are all the settings needed to
// create a new DeploymentEnvironmentVariable row.
type CreateDeploymentEnvironmentVariableRequest struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Encrypt bool   `json:"encrypt"`
}

// DeleteDeploymentEnvironmentVariableRequest deletes the environment variable
// with the given name from the deployment settings.
type DeleteDeploymentEnvironmentVariableRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CreateDeploymentSettingsRequest are all the settings needed to create a new
// DeploymentSettings row.
type CreateDeploymentSettingsRequest struct {
	GitHubRepositoryID *int64  `json:"gitHubRepositoryId,omitempty"`
	BranchName         *string `json:"branchName,omitempty"`
	WorkingDirectory   *string `json:"workingDirectory,omitempty"`
	PreRunCommands     *string `json:"preRunCommands,omitempty"`

	DeployCommits       bool `json:"deployCommits"`
	PreviewPullRequests bool `json:"previewPullRequests"`
}

// GitHubRepo identifies a GitHub repository.
type GitHubRepo struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

// ListReposResponse returns the repos we have access to.
type ListReposResponse struct {
	Repos []GitHubRepo `json:"repos"`
}

// CreateDeploymentSettingsResponse includes the newly created
// DeploymentSettings ID.
type CreateDeploymentSettingsResponse struct {
	DeploymentSettingsID string `json:"deploymentSettingsId"`
}

// Trigger represents when a Deployment would be created.
type Trigger struct {
	Condition string `json:"condition"`
	Action    string `json:"action"`
}

// GetDeploymentSettingsResponse returns the DeploymentSettings details.
type GetDeploymentSettingsResponse struct {
	ID         string  `json:"id"`
	Ref        string  `json:"ref"`
	WorkingDir *string `json:"workingDirectory,omitempty"`
	Repository string  `json:"repository"`
	Triggers   []Trigger
}

// ListDeploymentSettingsResponse returns a list of DeploymentSettings.
type ListDeploymentSettingsResponse struct {
	Result []GetDeploymentSettingsResponse `json:"result"`
}

// EnvironmentVariable is a key-value pair; if the variable is secret, value will be omitted.
type EnvironmentVariable struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Value  *string `json:"value,omitempty"`
	Secret bool    `json:"secret"`
}

// ListDeploymentEnvironmentVariablesResponse returns environment variables
// associated with these DeploymentSettings.
type ListDeploymentEnvironmentVariablesResponse struct {
	EnvironmentVariables []EnvironmentVariable `json:"environmentVariables"`
}

type DeploymentStepLogLine struct {
	Timestamp time.Time `json:"timestamp"`
	Line      string    `json:"line"`
}

type DeploymentLogsStep struct {
	Name       string                  `json:"name,omitempty"`
	NextOffset *int                    `json:"nextOffset,omitempty"`
	Lines      []DeploymentStepLogLine `json:"lines"`
}

type DeploymentLogLine struct {
	Header    string    `json:"header,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Line      string    `json:"line,omitempty"`
}

type DeploymentLogs struct {
	Lines     []DeploymentLogLine `json:"lines,omitempty"`
	NextToken string              `json:"nextToken,omitempty"`
}
