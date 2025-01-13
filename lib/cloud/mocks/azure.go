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

package mocks

import (
	"context"
	"time"

	"github.com/gravitational/trace"
	"github.com/jonboulle/clockwork"
	"k8s.io/client-go/rest"

	"github.com/gravitational/teleport/lib/cloud/azure"
)

// AKSClusterEntry is an entry in the AKSMock.Clusters list.
type AKSClusterEntry struct {
	azure.ClusterCredentialsConfig
	Config *rest.Config
	TTL    time.Duration
}

// AKSMock implements the azure.AKSClient interface for tests.
type AKSMock struct {
	azure.AKSClient
	Clusters []AKSClusterEntry
	Notify   chan struct{}
	Clock    clockwork.Clock
}

func (a *AKSMock) ClusterCredentials(ctx context.Context, cfg azure.ClusterCredentialsConfig) (*rest.Config, time.Time, error) {
	defer func() {
		a.Notify <- struct{}{}
	}()
	for _, cluster := range a.Clusters {
		if cluster.ClusterCredentialsConfig.ResourceGroup == cfg.ResourceGroup &&
			cluster.ClusterCredentialsConfig.ResourceName == cfg.ResourceName &&
			cluster.ClusterCredentialsConfig.TenantID == cfg.TenantID {
			return cluster.Config, a.Clock.Now().Add(cluster.TTL), nil
		}
	}
	return nil, time.Now(), trace.NotFound("cluster not found")
}
