// Copyright 2022, Pulumi Corporation.  All rights reserved.

package apitype

// GitAuthConfig specifies git authentication configuration options.
// There are 3 different authentication options:
//   - Personal access token
//   - SSH private key (and its optional password)
//   - Basic auth username and password
//
// Only 1 authentication mode is valid.
type GitAuthConfig struct {
	PersonalAccessToken *string    `json:"accessToken,omitempty"`
	SSHAuth             *SSHAuth   `json:"sshAuth,omitempty"`
	BasicAuth           *BasicAuth `json:"basicAuth,omitempty"`
}

// SSHAuth configures ssh-based auth for git authentication.
// SSHPrivateKey is required but password is optional.
type SSHAuth struct {
	SSHPrivateKey string  `json:"sshPrivateKey"`
	Password      *string `json:"password,omitempty"`
}

// BasicAuth configures git authentication through basic auth —
// i.e. username and password. Both UserName and Password are required.
type BasicAuth struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}
