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

package auth

import (
	"context"
	"fmt"

	"github.com/gravitational/trace"

	"github.com/gravitational/teleport/api/types"
	"github.com/gravitational/teleport/lib/modules"
	"github.com/gravitational/teleport/lib/spacelift"
)

type spaceliftIDTokenValidator interface {
	Validate(
		ctx context.Context, domain string, token string,
	) (*spacelift.IDTokenClaims, error)
}

func (a *Server) checkSpaceliftJoinRequest(
	ctx context.Context,
	req *types.RegisterUsingTokenRequest,
	pt types.ProvisionToken,
) (*spacelift.IDTokenClaims, error) {
	if req.IDToken == "" {
		return nil, trace.BadParameter("id_token not provided for spacelift join request")
	}
	token, ok := pt.(*types.ProvisionTokenV2)
	if !ok {
		return nil, trace.BadParameter("spacelift join method only supports ProvisionTokenV2, '%T' was provided", pt)
	}

	if modules.GetModules().BuildType() != modules.BuildEnterprise {
		return nil, fmt.Errorf(
			"spacelift joining: %w",
			ErrRequiresEnterprise,
		)
	}

	claims, err := a.spaceliftIDTokenValidator.Validate(
		ctx, token.Spec.Spacelift.Hostname, req.IDToken,
	)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	a.logger.InfoContext(ctx, "Spacelift run trying to join cluster",
		"claims", claims,
		"token", pt.GetName(),
	)

	return claims, trace.Wrap(checkSpaceliftAllowRules(token, claims))
}

func checkSpaceliftAllowRules(token *types.ProvisionTokenV2, claims *spacelift.IDTokenClaims) error {
	globCheck := func(want string, got string) (bool, error) {
		if token.Spec.Spacelift.EnableGlobMatching {
			return joinRuleGlobMatch(want, got)
		}
		if want == "" {
			return true, nil
		}
		return want == got, nil
	}

	// If a single rule passes, accept the IDToken
	for i, rule := range token.Spec.Spacelift.Allow {
		// Please consider keeping these field validators in the same order they
		// are defined within the ProvisionTokenSpecV2Spacelift proto spec.
		spaceIDMatch, err := globCheck(rule.SpaceID, claims.SpaceID)
		if err != nil {
			return trace.Wrap(err, "evaluating rule (%d) space_id glob match", i)
		}
		if !spaceIDMatch {
			continue
		}
		callerIDMatch, err := globCheck(rule.CallerID, claims.CallerID)
		if err != nil {
			return trace.Wrap(err, "evaluating rule (%d) caller_id glob match", i)
		}
		if !callerIDMatch {
			continue
		}
		if rule.CallerType != "" && claims.CallerType != rule.CallerType {
			continue
		}
		if rule.Scope != "" && claims.Scope != rule.Scope {
			continue
		}

		// All provided rules met.
		return nil
	}

	return trace.AccessDenied("id token claims did not match any allow rules")
}
