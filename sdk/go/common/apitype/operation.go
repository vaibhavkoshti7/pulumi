// Copyright 2022, Pulumi Corporation.  All rights reserved.

package apitype

import (
	"fmt"
)

// OperationContext describes what to do.
type OperationContext struct {
	// PreRunCommands is an optional list of arbitrary commands to run before Pulumi
	// is invoked.
	//
	// BUG: We probably don't want to support this in general.
	// ref: https://github.com/pulumi/pulumi/issues/9397
	PreRunCommands []string `json:"preRunCommands"`

	// Operation is what we plan on doing.
	Operation PulumiOperation `json:"operation"`

	// EnvironmentVariables contains environment variables to be applied during the execution.
	EnvironmentVariables map[string]string `json:"environmentVariables"`
}

// ConfigurationValue is the value of a Pulumi stack configuration key value pair.
//
// NOTE: Any type information is intentionally lost, as Pulumi will
// treat everything as a string at runtime.
type ConfigurationValue struct {
	Value    string `json:"value"`
	IsSecret bool   `json:"isSecret,omitempty"`
}

// Validate validates the operation context has the required fields set with correct values.
func (op OperationContext) Validate() error {

	switch op.Operation {
	case Update, Preview, Refresh, Destroy:
		break
	default:
		return fmt.Errorf("invalid operation: %v", op.Operation)
	}

	return nil
}
