// Copyright 2022, Pulumi Corporation.  All rights reserved.

package apitype

type ExecutorContext struct {
	// WorkingDirectory defines the path where the work should be done when executing.
	WorkingDirectory string `json:"workingDirectory"`

	// Defines the image that the pulumi operations should run in.
	ExecutorImage string `json:"executorImage,omitempty"`
}
