/*
 * Teleport
 * Copyright (C) 2023  Gravitational, Inc.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package config

import (
	"github.com/gravitational/trace"
	"gopkg.in/yaml.v3"

	"github.com/gravitational/teleport/lib/tbot/bot"
	"github.com/gravitational/teleport/lib/tbot/internal/encoding"
)

const ExampleServiceType = "example"

// ExampleService is a temporary example service for testing purposes. It is
// not intended to be used and exists to demonstrate how a user configurable
// service integrates with the tbot service manager.
type ExampleService struct {
	// Name of the service for logs and the /readyz endpoint.
	Name string `yaml:"name,omitempty"`

	Message string `yaml:"message"`
}

func (s *ExampleService) Type() string {
	return ExampleServiceType
}

// GetName returns the user-given name of the service, used for validation purposes.
func (s *ExampleService) GetName() string {
	return s.Name
}

func (s *ExampleService) MarshalYAML() (any, error) {
	type raw ExampleService
	return encoding.WithTypeHeader((*raw)(s), ExampleServiceType)
}

func (s *ExampleService) UnmarshalYAML(node *yaml.Node) error {
	// Alias type to remove UnmarshalYAML to avoid recursion
	type raw ExampleService
	if err := node.Decode((*raw)(s)); err != nil {
		return trace.Wrap(err)
	}
	return nil
}

func (s *ExampleService) CheckAndSetDefaults() error {
	if s.Message == "" {
		return trace.BadParameter("message: should not be empty")
	}
	return nil
}

func (s *ExampleService) GetCredentialLifetime() bot.CredentialLifetime {
	return bot.CredentialLifetime{}
}
