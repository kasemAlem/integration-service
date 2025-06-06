//
// Copyright 2021, Patrick Webster
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package gitlab

import (
	"fmt"
	"net/http"
	"net/url"
)

type (
	InstanceVariablesServiceInterface interface {
		ListVariables(opt *ListInstanceVariablesOptions, options ...RequestOptionFunc) ([]*InstanceVariable, *Response, error)
		GetVariable(key string, options ...RequestOptionFunc) (*InstanceVariable, *Response, error)
		CreateVariable(opt *CreateInstanceVariableOptions, options ...RequestOptionFunc) (*InstanceVariable, *Response, error)
		UpdateVariable(key string, opt *UpdateInstanceVariableOptions, options ...RequestOptionFunc) (*InstanceVariable, *Response, error)
		RemoveVariable(key string, options ...RequestOptionFunc) (*Response, error)
	}

	// InstanceVariablesService handles communication with the
	// instance level CI variables related methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/instance_level_ci_variables/
	InstanceVariablesService struct {
		client *Client
	}
)

var _ InstanceVariablesServiceInterface = (*InstanceVariablesService)(nil)

// InstanceVariable represents a GitLab instance level CI Variable.
//
// GitLab API docs:
// https://docs.gitlab.com/api/instance_level_ci_variables/
type InstanceVariable struct {
	Key          string            `json:"key"`
	Value        string            `json:"value"`
	VariableType VariableTypeValue `json:"variable_type"`
	Protected    bool              `json:"protected"`
	Masked       bool              `json:"masked"`
	Raw          bool              `json:"raw"`
	Description  string            `json:"description"`
}

func (v InstanceVariable) String() string {
	return Stringify(v)
}

// ListInstanceVariablesOptions represents the available options for listing variables
// for an instance.
//
// GitLab API docs:
// https://docs.gitlab.com/api/instance_level_ci_variables/#list-all-instance-variables
type ListInstanceVariablesOptions ListOptions

// ListVariables gets a list of all variables for an instance.
//
// GitLab API docs:
// https://docs.gitlab.com/api/instance_level_ci_variables/#list-all-instance-variables
func (s *InstanceVariablesService) ListVariables(opt *ListInstanceVariablesOptions, options ...RequestOptionFunc) ([]*InstanceVariable, *Response, error) {
	u := "admin/ci/variables"

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var vs []*InstanceVariable
	resp, err := s.client.Do(req, &vs)
	if err != nil {
		return nil, resp, err
	}

	return vs, resp, nil
}

// GetVariable gets a variable.
//
// GitLab API docs:
// https://docs.gitlab.com/api/instance_level_ci_variables/#show-instance-variable-details
func (s *InstanceVariablesService) GetVariable(key string, options ...RequestOptionFunc) (*InstanceVariable, *Response, error) {
	u := fmt.Sprintf("admin/ci/variables/%s", url.PathEscape(key))

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	v := new(InstanceVariable)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, nil
}

// CreateInstanceVariableOptions represents the available CreateVariable()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/instance_level_ci_variables/#create-instance-variable
type CreateInstanceVariableOptions struct {
	Key          *string            `url:"key,omitempty" json:"key,omitempty"`
	Value        *string            `url:"value,omitempty" json:"value,omitempty"`
	Description  *string            `url:"description,omitempty" json:"description,omitempty"`
	Masked       *bool              `url:"masked,omitempty" json:"masked,omitempty"`
	Protected    *bool              `url:"protected,omitempty" json:"protected,omitempty"`
	Raw          *bool              `url:"raw,omitempty" json:"raw,omitempty"`
	VariableType *VariableTypeValue `url:"variable_type,omitempty" json:"variable_type,omitempty"`
}

// CreateVariable creates a new instance level CI variable.
//
// GitLab API docs:
// https://docs.gitlab.com/api/instance_level_ci_variables/#create-instance-variable
func (s *InstanceVariablesService) CreateVariable(opt *CreateInstanceVariableOptions, options ...RequestOptionFunc) (*InstanceVariable, *Response, error) {
	u := "admin/ci/variables"

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	v := new(InstanceVariable)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, nil
}

// UpdateInstanceVariableOptions represents the available UpdateVariable()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/instance_level_ci_variables/#update-instance-variable
type UpdateInstanceVariableOptions struct {
	Value        *string            `url:"value,omitempty" json:"value,omitempty"`
	Description  *string            `url:"description,omitempty" json:"description,omitempty"`
	Masked       *bool              `url:"masked,omitempty" json:"masked,omitempty"`
	Protected    *bool              `url:"protected,omitempty" json:"protected,omitempty"`
	Raw          *bool              `url:"raw,omitempty" json:"raw,omitempty"`
	VariableType *VariableTypeValue `url:"variable_type,omitempty" json:"variable_type,omitempty"`
}

// UpdateVariable updates the position of an existing
// instance level CI variable.
//
// GitLab API docs:
// https://docs.gitlab.com/api/instance_level_ci_variables/#update-instance-variable
func (s *InstanceVariablesService) UpdateVariable(key string, opt *UpdateInstanceVariableOptions, options ...RequestOptionFunc) (*InstanceVariable, *Response, error) {
	u := fmt.Sprintf("admin/ci/variables/%s", url.PathEscape(key))

	req, err := s.client.NewRequest(http.MethodPut, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	v := new(InstanceVariable)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, nil
}

// RemoveVariable removes an instance level CI variable.
//
// GitLab API docs:
// https://docs.gitlab.com/api/instance_level_ci_variables/#remove-instance-variable
func (s *InstanceVariablesService) RemoveVariable(key string, options ...RequestOptionFunc) (*Response, error) {
	u := fmt.Sprintf("admin/ci/variables/%s", url.PathEscape(key))

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
