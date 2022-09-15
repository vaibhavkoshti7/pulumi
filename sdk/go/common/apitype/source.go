// Copyright 2022, Pulumi Corporation.  All rights reserved.

package apitype

import (
	"errors"
	"strings"
)

// SourceContext describes some source code, and how to obtain it.
type SourceContext struct {
	Git *SourceContextGit `json:"git,omitempty"`
}

type SourceContextGit struct {
	RepoURL string `json:"repoURL"`

	Branch string `json:"branch"`

	// (optional) RepoDir is the directory to work from in the project's source repository
	// where Pulumi.yaml is located. It is used in case Pulumi.yaml is not
	// in the project source root.
	RepoDir string `json:"repoDir,omitempty"`

	// (optional) Commit is the hash of the commit to deploy. If used, HEAD will be in detached mode. This
	// is mutually exclusive with the Branch setting. Either value needs to be specified.
	Commit string `json:"commit,omitempty"`

	// (optional) GitAuth allows configuring git authentication options
	// There are 3 different authentication options:
	//   * SSH private key (and its optional password)
	//   * Personal access token
	//   * Basic auth username and password
	// Only one authentication mode will be considered if more than one option is specified,
	// with ssh private key/password preferred first, then personal access token, and finally
	// basic auth credentials.
	GitAuth *GitAuthConfig `json:"gitAuth,omitempty"`
}

// Validate validates the source context is valid.
func (scg *SourceContextGit) Validate() error {

	// The repoURL is required.
	if len(strings.TrimSpace(scg.RepoURL)) == 0 {
		return errors.New("repoDir cannot be empty")
	}

	// Either the commit or the branch must be specified.
	if scg.Commit != "" && scg.Branch != "" {
		return errors.New("commit and branch cannot both be specified")
	}

	// Either branch or commit is required.
	if scg.Commit == "" && scg.Branch == "" {
		return errors.New("at least commit or branch are required")
	}

	// Git Auth is optional, but if specified, it must be valid.
	if gitAuth := scg.GitAuth; gitAuth != nil {

		foundAuthType := false

		if pat := gitAuth.PersonalAccessToken; pat != nil {
			// Set foundAuthType to true if the personal access method is specified.
			foundAuthType = true
			if len(strings.TrimSpace(*pat)) == 0 {
				return errors.New("personal access token cannot be empty")
			}
		}

		if sshAuth := gitAuth.SSHAuth; sshAuth != nil {
			// Check if another auth type has already been specified.
			if foundAuthType {
				return errors.New("only one authentication type is allowed")
			}
			// Set the foundAuthType to true if the ssh auth method is specified.
			foundAuthType = true
			if len(strings.TrimSpace(sshAuth.SSHPrivateKey)) == 0 {
				return errors.New("ssh private key cannot be empty")
			}
		}

		if basicAuth := gitAuth.BasicAuth; basicAuth != nil {
			// Check if another auth type has already been specified.
			if foundAuthType {
				return errors.New("only one authentication type is allowed")
			}
			if len(strings.TrimSpace(basicAuth.UserName)) == 0 {
				return errors.New("basic auth username cannot be empty")
			}
			if len(strings.TrimSpace(basicAuth.Password)) == 0 {
				return errors.New("basic auth password cannot be empty")
			}
		}
	}

	return nil
}
