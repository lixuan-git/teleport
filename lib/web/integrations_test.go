/*
 * Teleport
 * Copyright (C) 2024  Gravitational, Inc.
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

package web

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gravitational/trace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gravitational/teleport/api/client/proto"
	discoveryconfigv1 "github.com/gravitational/teleport/api/gen/proto/go/teleport/discoveryconfig/v1"
	integrationv1 "github.com/gravitational/teleport/api/gen/proto/go/teleport/integration/v1"
	usertasksv1 "github.com/gravitational/teleport/api/gen/proto/go/teleport/usertasks/v1"
	"github.com/gravitational/teleport/api/types"
	"github.com/gravitational/teleport/api/types/discoveryconfig"
	"github.com/gravitational/teleport/api/types/header"
	"github.com/gravitational/teleport/api/types/usertasks"
	"github.com/gravitational/teleport/lib/auth/integration/credentials"
	"github.com/gravitational/teleport/lib/modules"
	"github.com/gravitational/teleport/lib/modules/modulestest"
	"github.com/gravitational/teleport/lib/services"
	libui "github.com/gravitational/teleport/lib/ui"
	"github.com/gravitational/teleport/lib/utils/log/logtest"
	"github.com/gravitational/teleport/lib/web/ui"
)

func TestIntegrationsCreateWithAudience(t *testing.T) {
	t.Parallel()
	wPack := newWebPack(t, 1 /* proxies */)
	proxy := wPack.proxies[0]
	authPack := proxy.authPack(t, "user", []types.Role{services.NewPresetEditorRole()})
	ctx := context.Background()

	const integrationName = "test-integration"
	cases := []struct {
		name     string
		audience string
	}{
		{
			name:     "without audiences",
			audience: types.IntegrationAWSOIDCAudienceUnspecified,
		},
		{
			name:     "with audiences",
			audience: types.IntegrationAWSOIDCAudienceAWSIdentityCenter,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			createData := ui.Integration{
				Name:    integrationName,
				SubKind: "aws-oidc",
				AWSOIDC: &ui.IntegrationAWSOIDCSpec{
					RoleARN:  "arn:aws:iam::026090554232:role/testrole",
					Audience: test.audience,
				},
			}
			createEndpoint := authPack.clt.Endpoint("webapi", "sites", wPack.server.ClusterName(), "integrations")
			createResp, err := authPack.clt.PostJSON(ctx, createEndpoint, createData)
			require.NoError(t, err)
			require.Equal(t, 200, createResp.Code())

			// check origin label stored in backend
			intgrationResource, err := wPack.server.Auth().GetIntegration(ctx, integrationName)
			require.NoError(t, err)
			require.Equal(t, test.audience, intgrationResource.GetAWSOIDCIntegrationSpec().Audience)

			// check origin label returned in the web api
			getEndpoint := authPack.clt.Endpoint("webapi", "sites", wPack.server.ClusterName(), "integrations", integrationName)
			getResp, err := authPack.clt.Get(ctx, getEndpoint, nil)
			require.NoError(t, err)
			require.Equal(t, 200, getResp.Code())

			var resp ui.Integration
			err = json.Unmarshal(getResp.Bytes(), &resp)
			require.NoError(t, err)
			require.Equal(t, createData, resp)

			err = wPack.server.Auth().DeleteIntegration(ctx, integrationName)
			require.NoError(t, err)
		})
	}
}

func TestIntegrationsCRUDRolesAnywhere(t *testing.T) {
	t.Parallel()
	wPack := newWebPack(t, 1 /* proxies */)
	proxy := wPack.proxies[0]
	authPack := proxy.authPack(t, "user", []types.Role{services.NewPresetEditorRole()})
	ctx := context.Background()

	// Create Integration
	const integrationName = "test-integration"
	trustAnchorARN := "arn:aws:rolesanywhere:eu-west-2:123456789012:trust-anchor/12345678-1234-1234-1234-123456789012"
	createData := ui.Integration{
		Name:    integrationName,
		SubKind: "aws-ra",
		AWSRA: &ui.IntegrationAWSRASpec{
			TrustAnchorARN: trustAnchorARN,
			ProfileSyncConfig: ui.AWSRAProfileSync{
				Enabled: false,
			},
		},
	}
	createEndpoint := authPack.clt.Endpoint("webapi", "sites", wPack.server.ClusterName(), "integrations")
	createResp, err := authPack.clt.PostJSON(ctx, createEndpoint, createData)
	require.NoError(t, err)
	require.Equal(t, 200, createResp.Code())

	intgrationResource, err := wPack.server.Auth().GetIntegration(ctx, integrationName)
	require.NoError(t, err)
	require.Equal(t, trustAnchorARN, intgrationResource.GetAWSRolesAnywhereIntegrationSpec().TrustAnchorARN)

	// Create Integration fails when sync is enabled but config is not set
	createDataWithoutSyncFields := ui.Integration{
		Name:    "another-integration",
		SubKind: "aws-ra",
		AWSRA: &ui.IntegrationAWSRASpec{
			TrustAnchorARN: trustAnchorARN,
			ProfileSyncConfig: ui.AWSRAProfileSync{
				Enabled:    true,
				ProfileARN: "",
				RoleARN:    "arn:aws:iam::123456789012:role/testrole",
			},
		},
	}
	createResp, err = authPack.clt.PostJSON(ctx, createEndpoint, createDataWithoutSyncFields)
	require.ErrorContains(t, err, "missing awsra.profileSync.profileArn field")
	require.Equal(t, 400, createResp.Code())

	// Get single integration
	getEndpoint := authPack.clt.Endpoint("webapi", "sites", wPack.server.ClusterName(), "integrations", integrationName)
	getResp, err := authPack.clt.Get(ctx, getEndpoint, nil)
	require.NoError(t, err)
	require.Equal(t, 200, getResp.Code())

	var resp ui.Integration
	err = json.Unmarshal(getResp.Bytes(), &resp)
	require.NoError(t, err)
	require.Equal(t, createData, resp)

	// Update integration
	updatedTrustAnchor := "arn:aws:rolesanywhere:eu-west-2:123456789012:trust-anchor/00000000-0000-0000-0000-123456789012"
	syncProfileARN := "arn:aws:rolesanywhere:eu-west-2:123456789012:profile/00000000-0000-0000-0000-123456789012"
	syncRoleARN := "arn:aws:iam::123456789012:role/testrole"
	updateIntegration := ui.UpdateIntegrationRequest{
		AWSRA: &ui.IntegrationAWSRASpec{
			TrustAnchorARN: updatedTrustAnchor,
			ProfileSyncConfig: ui.AWSRAProfileSync{
				Enabled:            true,
				ProfileARN:         syncProfileARN,
				RoleARN:            syncRoleARN,
				ProfileNameFilters: []string{"ExposedProfile-*"},
			},
		},
	}
	updateEndpoint := authPack.clt.Endpoint("webapi", "sites", wPack.server.ClusterName(), "integrations", integrationName)
	updateResp, err := authPack.clt.PutJSON(ctx, updateEndpoint, updateIntegration)
	require.NoError(t, err)
	require.Equal(t, 200, updateResp.Code())

	// List integrations
	listEndpoint := authPack.clt.Endpoint("webapi", "sites", wPack.server.ClusterName(), "integrations")
	listResp, err := authPack.clt.Get(ctx, listEndpoint, nil)
	require.NoError(t, err)
	require.Equal(t, 200, listResp.Code())

	var listRespObject ui.IntegrationsListResponse
	err = json.Unmarshal(listResp.Bytes(), &listRespObject)
	require.NoError(t, err)
	require.Len(t, listRespObject.Items, 1)
	integrationObject := listRespObject.Items[0]
	require.Equal(t, updatedTrustAnchor, integrationObject.AWSRA.TrustAnchorARN)
	require.True(t, integrationObject.AWSRA.ProfileSyncConfig.Enabled)
	require.Equal(t, syncProfileARN, integrationObject.AWSRA.ProfileSyncConfig.ProfileARN)
	require.Equal(t, syncRoleARN, integrationObject.AWSRA.ProfileSyncConfig.RoleARN)
	require.Len(t, integrationObject.AWSRA.ProfileSyncConfig.ProfileNameFilters, 1)
	require.Equal(t, "ExposedProfile-*", integrationObject.AWSRA.ProfileSyncConfig.ProfileNameFilters[0])

	// Delete Integration
	err = wPack.server.Auth().DeleteIntegration(ctx, integrationName)
	require.NoError(t, err)
}

type mockUserTasksLister struct {
	defaultPageSize int64
	userTasks       []*usertasksv1.UserTask
}

func (m *mockUserTasksLister) ListUserTasks(ctx context.Context, pageSize int64, nextToken string, filters *usertasksv1.ListUserTasksFilters) ([]*usertasksv1.UserTask, string, error) {
	var ret []*usertasksv1.UserTask
	if pageSize == 0 {
		pageSize = m.defaultPageSize
	}

	if len(m.userTasks) == 0 {
		return ret, "", nil
	}

	var sliceStart int
	if nextToken != "" {
		nextTokenInt, err := strconv.Atoi(nextToken)
		if err != nil {
			return nil, "", trace.Wrap(err)
		}
		sliceStart = nextTokenInt
	}
	userTasksSlice := m.userTasks[sliceStart:]

	for i, userTask := range userTasksSlice {
		if userTask.GetSpec().GetState() == "OPEN" {
			ret = append(ret, userTask)
			if len(ret) == int(pageSize) {
				nextTokenInt := sliceStart + i + 1
				return ret, strconv.Itoa(nextTokenInt), nil
			}
		}
	}

	return ret, "", nil
}

func TestCollectAWSOIDCAutoDiscoverStats(t *testing.T) {
	ctx := context.Background()
	logger := logtest.NewLogger()

	integrationName := "my-integration"
	integration, err := types.NewIntegrationAWSOIDC(
		types.Metadata{Name: integrationName},
		&types.AWSOIDCIntegrationSpecV1{
			RoleARN: "arn:role",
		},
	)
	require.NoError(t, err)

	deployedServiceCommand := buildCommandDeployedDatabaseService(t, true, types.Labels{"vpc": []string{"vpc1", "vpc2"}})
	deployedDatabaseServicesClient := &mockDeployedDatabaseServices{
		integration: "my-integration",
		servicesPerRegion: map[string][]*integrationv1.DeployedDatabaseService{
			"us-west-2": dummyDeployedDatabaseServices(1, deployedServiceCommand),
		},
	}

	t.Run("without discovery configs, returns just the integration", func(t *testing.T) {
		clt := &mockRelevantAWSRegionsClient{
			databaseServices: &proto.ListResourcesResponse{
				Resources: []*proto.PaginatedResource{},
			},
			databases:        make([]types.Database, 0),
			discoveryConfigs: make([]*discoveryconfig.DiscoveryConfig, 0),
		}

		req := collectIntegrationStatsRequest{
			logger:                logger,
			integration:           integration,
			discoveryConfigLister: clt,
			databaseGetter:        clt,
			awsOIDCClient:         deployedDatabaseServicesClient,
			userTasksClient:       &mockUserTasksLister{},
		}
		gotSummary, err := collectIntegrationStats(ctx, req)
		require.NoError(t, err)
		expectedSummary := &ui.IntegrationWithSummary{
			Integration: &ui.Integration{
				Name:    integrationName,
				SubKind: "aws-oidc",
				AWSOIDC: &ui.IntegrationAWSOIDCSpec{RoleARN: "arn:role"},
			},
		}
		require.Equal(t, expectedSummary, gotSummary)
	})

	t.Run("returns the number of unresolved user tasks", func(t *testing.T) {
		clt := &mockRelevantAWSRegionsClient{
			databaseServices: &proto.ListResourcesResponse{
				Resources: []*proto.PaginatedResource{},
			},
			databases:        make([]types.Database, 0),
			discoveryConfigs: make([]*discoveryconfig.DiscoveryConfig, 0),
		}

		var userTasksList []*usertasksv1.UserTask
		ec2UserTasks := 10
		for range ec2UserTasks {
			userTasksList = append(userTasksList, &usertasksv1.UserTask{Spec: &usertasksv1.UserTaskSpec{State: usertasks.TaskStateOpen, TaskType: usertasks.TaskTypeDiscoverEC2}})
		}
		rdsUserTasks := 20
		for range rdsUserTasks {
			userTasksList = append(userTasksList, &usertasksv1.UserTask{Spec: &usertasksv1.UserTaskSpec{State: usertasks.TaskStateOpen, TaskType: usertasks.TaskTypeDiscoverRDS}})
		}
		for range 100 {
			userTasksList = append(userTasksList, &usertasksv1.UserTask{Spec: &usertasksv1.UserTaskSpec{State: usertasks.TaskStateResolved, TaskType: usertasks.TaskTypeDiscoverEC2}})
		}

		userTasksClient := &mockUserTasksLister{
			defaultPageSize: 3,
			userTasks:       userTasksList,
		}

		req := collectIntegrationStatsRequest{
			logger:                logger,
			integration:           integration,
			discoveryConfigLister: clt,
			databaseGetter:        clt,
			awsOIDCClient:         deployedDatabaseServicesClient,
			userTasksClient:       userTasksClient,
		}
		gotSummary, err := collectIntegrationStats(ctx, req)
		require.NoError(t, err)
		expectedSummary := &ui.IntegrationWithSummary{
			Integration: &ui.Integration{
				Name:    integrationName,
				SubKind: "aws-oidc",
				AWSOIDC: &ui.IntegrationAWSOIDCSpec{RoleARN: "arn:role"},
			},
			UnresolvedUserTasks: ec2UserTasks + rdsUserTasks,
			AWSEC2: ui.ResourceTypeSummary{
				UnresolvedUserTasks: ec2UserTasks,
			},
			AWSRDS: ui.ResourceTypeSummary{
				UnresolvedUserTasks: rdsUserTasks,
			},
		}
		require.Equal(t, expectedSummary, gotSummary)
	})

	t.Run("collects multiple discovery configs", func(t *testing.T) {
		syncTime := time.Now()
		dcForEC2 := &discoveryconfig.DiscoveryConfig{
			Spec: discoveryconfig.Spec{AWS: []types.AWSMatcher{{
				Integration: integrationName,
				Types:       []string{"ec2"},
				Regions:     []string{"us-east-1"},
			}}},
			Status: discoveryconfig.Status{
				LastSyncTime:        syncTime,
				DiscoveredResources: 2,
				IntegrationDiscoveredResources: map[string]*discoveryconfigv1.IntegrationDiscoveredSummary{
					integrationName: {
						AwsEc2: &discoveryconfigv1.ResourcesDiscoveredSummary{Found: 2, Enrolled: 1, Failed: 1},
					},
				},
			},
		}
		dcForRDS := &discoveryconfig.DiscoveryConfig{
			Spec: discoveryconfig.Spec{AWS: []types.AWSMatcher{{
				Integration: integrationName,
				Types:       []string{"rds"},
				Regions:     []string{"us-east-1", "us-east-2", "us-west-2"},
			}}},
			Status: discoveryconfig.Status{
				LastSyncTime:        syncTime,
				DiscoveredResources: 2,
				IntegrationDiscoveredResources: map[string]*discoveryconfigv1.IntegrationDiscoveredSummary{
					integrationName: {
						AwsRds: &discoveryconfigv1.ResourcesDiscoveredSummary{Found: 2, Enrolled: 1, Failed: 1},
					},
				},
			},
		}
		dcForEKS := &discoveryconfig.DiscoveryConfig{
			Spec: discoveryconfig.Spec{AWS: []types.AWSMatcher{{
				Integration: integrationName,
				Types:       []string{"eks"},
				Regions:     []string{"us-east-1"},
			}}},
			Status: discoveryconfig.Status{
				LastSyncTime:        syncTime,
				DiscoveredResources: 2,
				IntegrationDiscoveredResources: map[string]*discoveryconfigv1.IntegrationDiscoveredSummary{
					integrationName: {
						AwsEks: &discoveryconfigv1.ResourcesDiscoveredSummary{Found: 4, Enrolled: 0, Failed: 0},
					},
				},
			},
		}
		clt := &mockRelevantAWSRegionsClient{
			discoveryConfigs: []*discoveryconfig.DiscoveryConfig{
				dcForEC2,
				dcForRDS,
				dcForEKS,
			},
			databaseServices: &proto.ListResourcesResponse{},
			databases:        make([]types.Database, 0),
		}

		req := collectIntegrationStatsRequest{
			logger:                logger,
			integration:           integration,
			discoveryConfigLister: clt,
			databaseGetter:        clt,
			awsOIDCClient:         deployedDatabaseServicesClient,
			userTasksClient:       &mockUserTasksLister{},
		}
		gotSummary, err := collectIntegrationStats(ctx, req)
		require.NoError(t, err)
		expectedSummary := &ui.IntegrationWithSummary{
			Integration: &ui.Integration{
				Name:    integrationName,
				SubKind: "aws-oidc",
				AWSOIDC: &ui.IntegrationAWSOIDCSpec{RoleARN: "arn:role"},
			},
			AWSEC2: ui.ResourceTypeSummary{
				RulesCount:                 1,
				ResourcesFound:             2,
				ResourcesEnrollmentFailed:  1,
				ResourcesEnrollmentSuccess: 1,
				DiscoverLastSync:           &syncTime,
			},
			AWSRDS: ui.ResourceTypeSummary{
				RulesCount:                 3,
				ResourcesFound:             2,
				ResourcesEnrollmentFailed:  1,
				ResourcesEnrollmentSuccess: 1,
				ECSDatabaseServiceCount:    1,
				DiscoverLastSync:           &syncTime,
			},
			AWSEKS: ui.ResourceTypeSummary{
				RulesCount:                 1,
				ResourcesFound:             4,
				ResourcesEnrollmentFailed:  0,
				ResourcesEnrollmentSuccess: 0,
				DiscoverLastSync:           &syncTime,
			},
		}
		require.Equal(t, expectedSummary, gotSummary)
	})
	t.Run("returns 0 ECS DatabaseServices if listing deployed database services returns AccessDenied", func(t *testing.T) {
		syncTime := time.Now()
		dcForRDS := &discoveryconfig.DiscoveryConfig{
			Spec: discoveryconfig.Spec{AWS: []types.AWSMatcher{{
				Integration: integrationName,
				Types:       []string{"rds"},
				Regions:     []string{"us-east-1", "us-east-2", "us-west-2"},
			}}},
			Status: discoveryconfig.Status{
				LastSyncTime:        syncTime,
				DiscoveredResources: 2,
				IntegrationDiscoveredResources: map[string]*discoveryconfigv1.IntegrationDiscoveredSummary{
					integrationName: {
						AwsRds: &discoveryconfigv1.ResourcesDiscoveredSummary{Found: 2, Enrolled: 1, Failed: 1},
					},
				},
			},
		}
		clt := &mockRelevantAWSRegionsClient{
			discoveryConfigs: []*discoveryconfig.DiscoveryConfig{
				dcForRDS,
			},
			databaseServices: &proto.ListResourcesResponse{},
			databases:        make([]types.Database, 0),
		}

		deployedDatabaseServicesClient := &mockDeployedDatabaseServices{
			listErr: trace.AccessDenied("AccessDenied to ECS:ListServices"),
		}
		req := collectIntegrationStatsRequest{
			logger:                logger,
			integration:           integration,
			discoveryConfigLister: clt,
			databaseGetter:        clt,
			awsOIDCClient:         deployedDatabaseServicesClient,
			userTasksClient:       &mockUserTasksLister{},
		}
		gotSummary, err := collectIntegrationStats(ctx, req)
		require.NoError(t, err)
		expectedSummary := &ui.IntegrationWithSummary{
			Integration: &ui.Integration{
				Name:    integrationName,
				SubKind: "aws-oidc",
				AWSOIDC: &ui.IntegrationAWSOIDCSpec{RoleARN: "arn:role"},
			},
			AWSRDS: ui.ResourceTypeSummary{
				RulesCount:                 3,
				ResourcesFound:             2,
				ResourcesEnrollmentFailed:  1,
				ResourcesEnrollmentSuccess: 1,
				ECSDatabaseServiceCount:    0,
				DiscoverLastSync:           &syncTime,
			},
		}
		require.Equal(t, expectedSummary, gotSummary)
	})
}

func TestCollectAutoDiscoveryRules(t *testing.T) {
	ctx := context.Background()
	integrationName := "my-integration"

	t.Run("without discovery configs, returns no rules", func(t *testing.T) {
		clt := &mockRelevantAWSRegionsClient{
			discoveryConfigs: make([]*discoveryconfig.DiscoveryConfig, 0),
		}

		gotRules, err := collectAutoDiscoveryRules(ctx, integrationName, "", "", nil, clt)
		require.NoError(t, err)
		expectedRules := ui.IntegrationDiscoveryRules{}
		require.Equal(t, expectedRules, gotRules)
	})

	t.Run("collects multiple discovery configs", func(t *testing.T) {
		syncTime := time.Now()
		dcForEC2 := &discoveryconfig.DiscoveryConfig{
			ResourceHeader: header.ResourceHeader{Metadata: header.Metadata{
				Name: uuid.NewString(),
			}},
			Spec: discoveryconfig.Spec{AWS: []types.AWSMatcher{{
				Integration: integrationName,
				Types:       []string{"ec2"},
				Regions:     []string{"us-east-1"},
				Tags:        types.Labels{"*": []string{"*"}},
			}}},
			Status: discoveryconfig.Status{
				LastSyncTime: syncTime,
			},
		}
		dcForRDS := &discoveryconfig.DiscoveryConfig{
			ResourceHeader: header.ResourceHeader{Metadata: header.Metadata{
				Name: uuid.NewString(),
			}},
			Spec: discoveryconfig.Spec{AWS: []types.AWSMatcher{{
				Integration: integrationName,
				Types:       []string{"rds"},
				Regions:     []string{"us-east-1", "us-east-2"},
				Tags: types.Labels{
					"env": []string{"dev", "prod"},
				},
			}}},
			Status: discoveryconfig.Status{
				LastSyncTime: syncTime,
			},
		}
		dcForEKS := &discoveryconfig.DiscoveryConfig{
			ResourceHeader: header.ResourceHeader{Metadata: header.Metadata{
				Name: uuid.NewString(),
			}},
			Spec: discoveryconfig.Spec{AWS: []types.AWSMatcher{{
				Integration: integrationName,
				Types:       []string{"eks"},
				Regions:     []string{"us-east-1"},
				Tags:        types.Labels{"*": []string{"*"}},
			}}},
			Status: discoveryconfig.Status{
				LastSyncTime: syncTime,
			},
		}
		dcForEKSWithoutStatus := &discoveryconfig.DiscoveryConfig{
			ResourceHeader: header.ResourceHeader{Metadata: header.Metadata{
				Name: uuid.NewString(),
			}},
			Spec: discoveryconfig.Spec{AWS: []types.AWSMatcher{{
				Integration: integrationName,
				Types:       []string{"eks"},
				Regions:     []string{"eu-west-1"},
				Tags:        types.Labels{"*": []string{"*"}},
			}}},
		}
		clt := &mockRelevantAWSRegionsClient{
			discoveryConfigs: []*discoveryconfig.DiscoveryConfig{
				dcForEC2,
				dcForRDS,
				dcForEKS,
				dcForEKSWithoutStatus,
			},
		}

		got, err := collectAutoDiscoveryRules(ctx, integrationName, "", "", nil, clt)
		require.NoError(t, err)
		expectedRules := []ui.IntegrationDiscoveryRule{
			{
				ResourceType: "ec2",
				Region:       "us-east-1",
				LabelMatcher: []libui.Label{
					{Name: "*", Value: "*"},
				},
				DiscoveryConfig: dcForEC2.GetName(),
				LastSync:        &syncTime,
			},
			{
				ResourceType: "eks",
				Region:       "us-east-1",
				LabelMatcher: []libui.Label{
					{Name: "*", Value: "*"},
				},
				DiscoveryConfig: dcForEKS.GetName(),
				LastSync:        &syncTime,
			},
			{
				ResourceType: "eks",
				Region:       "eu-west-1",
				LabelMatcher: []libui.Label{
					{Name: "*", Value: "*"},
				},
				DiscoveryConfig: dcForEKSWithoutStatus.GetName(),
			},
			{
				ResourceType: "rds",
				Region:       "us-east-1",
				LabelMatcher: []libui.Label{
					{Name: "env", Value: "dev"},
					{Name: "env", Value: "prod"},
				},
				DiscoveryConfig: dcForRDS.GetName(),
				LastSync:        &syncTime,
			},
			{
				ResourceType: "rds",
				Region:       "us-east-2",
				LabelMatcher: []libui.Label{
					{Name: "env", Value: "dev"},
					{Name: "env", Value: "prod"},
				},
				DiscoveryConfig: dcForRDS.GetName(),
				LastSync:        &syncTime,
			},
		}
		require.Empty(t, got.NextKey)
		require.ElementsMatch(t, expectedRules, got.Rules)
	})

	t.Run("filters resource type", func(t *testing.T) {
		syncTime := time.Now()
		dcForEC2 := &discoveryconfig.DiscoveryConfig{
			ResourceHeader: header.ResourceHeader{Metadata: header.Metadata{
				Name: uuid.NewString(),
			}},
			Spec: discoveryconfig.Spec{AWS: []types.AWSMatcher{{
				Integration: integrationName,
				Types:       []string{"ec2"},
				Regions:     []string{"us-east-1"},
				Tags:        types.Labels{"*": []string{"*"}},
			}}},
			Status: discoveryconfig.Status{
				LastSyncTime: syncTime,
			},
		}
		dcForRDS := &discoveryconfig.DiscoveryConfig{
			ResourceHeader: header.ResourceHeader{Metadata: header.Metadata{
				Name: uuid.NewString(),
			}},
			Spec: discoveryconfig.Spec{AWS: []types.AWSMatcher{{
				Integration: integrationName,
				Types:       []string{"rds"},
				Regions:     []string{"us-east-1", "us-east-2"},
				Tags: types.Labels{
					"env": []string{"dev", "prod"},
				},
			}}},
			Status: discoveryconfig.Status{
				LastSyncTime: syncTime,
			},
		}
		clt := &mockRelevantAWSRegionsClient{
			discoveryConfigs: []*discoveryconfig.DiscoveryConfig{
				dcForEC2,
				dcForRDS,
			},
		}

		got, err := collectAutoDiscoveryRules(ctx, integrationName, "", "ec2", nil, clt)
		require.NoError(t, err)
		expectedRules := []ui.IntegrationDiscoveryRule{
			{
				ResourceType: "ec2",
				Region:       "us-east-1",
				LabelMatcher: []libui.Label{
					{Name: "*", Value: "*"},
				},
				DiscoveryConfig: dcForEC2.GetName(),
				LastSync:        &syncTime,
			},
		}
		require.Empty(t, got.NextKey)
		require.ElementsMatch(t, expectedRules, got.Rules)
	})

	t.Run("filters by region", func(t *testing.T) {
		syncTime := time.Now()
		dcForRDS := &discoveryconfig.DiscoveryConfig{
			ResourceHeader: header.ResourceHeader{Metadata: header.Metadata{
				Name: uuid.NewString(),
			}},
			Spec: discoveryconfig.Spec{AWS: []types.AWSMatcher{{
				Integration: integrationName,
				Types:       []string{"rds"},
				Regions:     []string{"us-east-1", "us-east-2", "us-west-2"},
				Tags: types.Labels{
					"env": []string{"dev", "prod"},
				},
			}}},
			Status: discoveryconfig.Status{
				LastSyncTime: syncTime,
			},
		}
		clt := &mockRelevantAWSRegionsClient{
			discoveryConfigs: []*discoveryconfig.DiscoveryConfig{
				dcForRDS,
			},
		}

		got, err := collectAutoDiscoveryRules(ctx, integrationName, "", "", []string{"us-east-1", "us-east-2"}, clt)
		require.NoError(t, err)
		expectedRules := []ui.IntegrationDiscoveryRule{
			{
				ResourceType: "rds",
				Region:       "us-east-1",
				LabelMatcher: []libui.Label{
					{Name: "env", Value: "dev"},
					{Name: "env", Value: "prod"},
				},
				DiscoveryConfig: dcForRDS.GetName(),
				LastSync:        &syncTime,
			},
			{
				ResourceType: "rds",
				Region:       "us-east-2",
				LabelMatcher: []libui.Label{
					{Name: "env", Value: "dev"},
					{Name: "env", Value: "prod"},
				},
				DiscoveryConfig: dcForRDS.GetName(),
				LastSync:        &syncTime,
			},
		}
		require.Empty(t, got.NextKey)
		require.ElementsMatch(t, expectedRules, got.Rules)
	})

	t.Run("pagination", func(t *testing.T) {
		syncTime := time.Now()
		totalRules := 1000

		discoveryConfigs := make([]*discoveryconfig.DiscoveryConfig, 0, totalRules)
		for range totalRules {
			discoveryConfigs = append(discoveryConfigs,
				&discoveryconfig.DiscoveryConfig{
					ResourceHeader: header.ResourceHeader{Metadata: header.Metadata{
						Name: uuid.NewString(),
					}},
					Spec: discoveryconfig.Spec{AWS: []types.AWSMatcher{{
						Integration: integrationName,
						Types:       []string{"ec2"},
						Regions:     []string{"us-east-1"},
						Tags:        types.Labels{"*": []string{"*"}},
					}}},
					Status: discoveryconfig.Status{
						LastSyncTime: syncTime,
					},
				},
			)
		}
		clt := &mockRelevantAWSRegionsClient{
			discoveryConfigs: discoveryConfigs,
		}

		nextKey := ""
		rulesCounter := 0
		for {
			got, err := collectAutoDiscoveryRules(ctx, integrationName, nextKey, "", nil, clt)
			require.NoError(t, err)
			rulesCounter += len(got.Rules)
			nextKey = got.NextKey
			if nextKey == "" {
				break
			}
		}
		require.Equal(t, totalRules, rulesCounter)
	})
}

// TestGitHubIntegration tests CRUD on GitHub integration subkind and CA export.
// GitHub integration requires modules.BuildEnterprise.
// The test cases in this test are performed sequentially and each test case
// depends on the previous state.
func TestGitHubIntegration(t *testing.T) {
	modulestest.SetTestModules(t, modulestest.Modules{TestBuildType: modules.BuildEnterprise})

	wPack := newWebPack(t, 1 /* proxies */)
	proxy := wPack.proxies[0]
	authPack := proxy.authPack(t, "user", []types.Role{services.NewPresetEditorRole()})
	ctx := context.Background()
	orgName := "my-org"
	integrationName := "github-" + orgName

	t.Run("create", func(t *testing.T) {
		endpoint := authPack.clt.Endpoint("webapi", "sites", wPack.server.ClusterName(), "integrations")
		uiIntegration := ui.Integration{
			Name:    integrationName,
			SubKind: types.IntegrationSubKindGitHub,
			GitHub: &ui.IntegrationGitHub{
				Organization: orgName,
			},
		}
		t.Run("missing oauth", func(t *testing.T) {
			_, err := authPack.clt.PostJSON(ctx, endpoint, ui.CreateIntegrationRequest{
				Integration: uiIntegration,
			})
			require.Error(t, err)

		})
		t.Run("success", func(t *testing.T) {
			createResp, err := authPack.clt.PostJSON(ctx, endpoint, ui.CreateIntegrationRequest{
				Integration: uiIntegration,
				OAuth: &ui.IntegrationOAuthCredentials{
					ID:     "oauth-id",
					Secret: "oauth-secret",
				},
			})
			require.NoError(t, err)
			require.Equal(t, 200, createResp.Code())
		})
	})

	t.Run("get", func(t *testing.T) {
		endpoint := authPack.clt.Endpoint("webapi", "sites", wPack.server.ClusterName(), "integrations", integrationName)
		getResp, err := authPack.clt.Get(ctx, endpoint, nil)
		require.NoError(t, err)
		require.Equal(t, 200, getResp.Code())

		var resp ui.Integration
		require.NoError(t, json.Unmarshal(getResp.Bytes(), &resp))
		require.Equal(t, ui.Integration{
			Name:    integrationName,
			SubKind: types.IntegrationSubKindGitHub,
			GitHub: &ui.IntegrationGitHub{
				Organization: orgName,
			},
		}, resp)
	})

	t.Run("export ca", func(t *testing.T) {
		endpoint := authPack.clt.Endpoint("webapi", "sites", wPack.server.ClusterName(), "integrations", integrationName, "ca")
		caResp, err := authPack.clt.Get(ctx, endpoint, nil)
		require.NoError(t, err)
		require.Equal(t, 200, caResp.Code())

		var resp ui.CAKeySet
		require.NoError(t, json.Unmarshal(caResp.Bytes(), &resp))
		require.NotEmpty(t, resp.SSH)
		assert.NotEmpty(t, resp.SSH[0].PublicKey)
		assert.NotEmpty(t, resp.SSH[0].Fingerprint)
	})

	t.Run("update", func(t *testing.T) {
		endpoint := authPack.clt.Endpoint("webapi", "sites", wPack.server.ClusterName(), "integrations", integrationName)
		t.Run("bad request", func(t *testing.T) {
			_, err := authPack.clt.PutJSON(ctx, endpoint, ui.UpdateIntegrationRequest{
				OAuth: &ui.IntegrationOAuthCredentials{
					ID: "oauth-id",
				},
			})
			require.Error(t, err)
		})

		t.Run("success", func(t *testing.T) {
			_, err := authPack.clt.PutJSON(ctx, endpoint, ui.UpdateIntegrationRequest{
				OAuth: &ui.IntegrationOAuthCredentials{
					ID:     "new-oauth-id",
					Secret: "new-oauth-secret",
				},
			})
			require.NoError(t, err)

			// Credentials are only accessible by Auth at the moment.
			ig, err := wPack.server.Auth().GetIntegration(ctx, integrationName)
			require.NoError(t, err)
			require.NotNil(t, ig.GetCredentials())
			cred, err := credentials.GetByPurpose(ctx, ig.GetCredentials().GetStaticCredentialsRef(), credentials.PurposeGitHubOAuth, wPack.server.Auth())
			require.NoError(t, err)
			updatedID, updatedSecret := cred.GetOAuthClientSecret()
			assert.Equal(t, "new-oauth-id", updatedID)
			assert.Equal(t, "new-oauth-secret", updatedSecret)
		})
	})

	t.Run("delete", func(t *testing.T) {
		githubServer, err := types.NewGitHubServer(types.GitHubServerMetadata{
			Integration:  integrationName,
			Organization: orgName,
		})
		require.NoError(t, err)
		_, err = proxy.auth.AuthServer.AuthServer.CreateGitServer(ctx, githubServer)
		require.NoError(t, err)

		endpoint := authPack.clt.Endpoint("webapi", "sites", wPack.server.ClusterName(), "integrations", integrationName)
		t.Run("failed because existing git server ", func(t *testing.T) {
			_, err := authPack.clt.Delete(ctx, endpoint)
			require.Error(t, err)
		})

		t.Run("success with associated resources param", func(t *testing.T) {
			_, err := authPack.clt.Delete(ctx, endpoint+"?associatedresources=true")
			require.NoError(t, err)

			_, err = authPack.clt.Get(ctx, endpoint, nil)
			require.Error(t, err)
			require.True(t, trace.IsNotFound(err))
			_, err = proxy.auth.AuthServer.AuthServer.GetGitServer(ctx, githubServer.GetName())
			require.True(t, trace.IsNotFound(err))
		})

		t.Run("not found", func(t *testing.T) {
			_, err := authPack.clt.Delete(ctx, endpoint)
			require.Error(t, err)
			require.True(t, trace.IsNotFound(err))
		})
	})
}
