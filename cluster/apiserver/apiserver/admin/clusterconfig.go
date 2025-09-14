/*
 * Copyright Octelium Labs, LLC. All rights reserved.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License version 3,
 * as published by the Free Software Foundation of the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package admin

import (
	"context"

	"github.com/asaskevich/govalidator"
	"github.com/octelium/octelium/apis/main/corev1"
	"github.com/octelium/octelium/apis/main/metav1"
	apisrvcommon "github.com/octelium/octelium/cluster/apiserver/apiserver/common"
	"github.com/octelium/octelium/cluster/apiserver/apiserver/serr"
	"github.com/octelium/octelium/cluster/common/apivalidation"
	"github.com/octelium/octelium/cluster/common/grpcutils"
	"github.com/pkg/errors"
)

func (s *Server) GetClusterConfig(ctx context.Context, req *corev1.GetClusterConfigRequest) (*corev1.ClusterConfig, error) {
	cc, err := s.octeliumC.CoreV1Utils().GetClusterConfig(ctx)
	if err != nil {
		return nil, serr.InternalWithErr(err)
	}

	return cc, nil
}

func (s *Server) UpdateClusterConfig(ctx context.Context, req *corev1.ClusterConfig) (*corev1.ClusterConfig, error) {

	if err := s.validateClusterConfig(ctx, req); err != nil {
		return nil, err
	}

	cfg, err := s.octeliumC.CoreV1Utils().GetClusterConfig(ctx)
	if err != nil {
		return nil, serr.InternalWithErr(err)
	}

	apisrvcommon.MetadataUpdate(cfg.Metadata, req.Metadata)
	cfg.Spec = req.Spec

	ccOut, err := s.octeliumC.CoreC().UpdateClusterConfig(ctx, cfg)
	if err != nil {
		return nil, serr.InternalWithErr(err)
	}

	return ccOut, nil
}

func (s *Server) validateClusterConfig(ctx context.Context, req *corev1.ClusterConfig) error {

	if err := apivalidation.ValidateCommon(req, &apivalidation.ValidateCommonOpts{
		ValidateMetadataOpts: apivalidation.ValidateMetadataOpts{},
	}); err != nil {
		return err
	}

	if req.Spec == nil {
		return grpcutils.InvalidArg("Nil spec")
	}

	if err := s.validateClusterConfigSpec(ctx, req); err != nil {
		return grpcutils.InvalidArgWithErr(err)
	}

	return nil
}

func (s *Server) validateClusterConfigSpec(ctx context.Context, c *corev1.ClusterConfig) error {

	if err := validateSession(c); err != nil {
		return err
	}

	if err := validateGateway(c); err != nil {
		return err
	}

	if err := validateDNS(c); err != nil {
		return err
	}

	if err := s.validateCCAuthorization(ctx, c); err != nil {
		return err
	}

	return nil
}

func validateSession(c *corev1.ClusterConfig) error {
	if c.Spec.Session == nil {
		return nil
	}

	if c.Spec.Session.Human != nil {
		if err := validateDuration(c.Spec.Session.Human.ClientDuration); err != nil {
			return errors.Errorf("Invalid HUMAN client session duration")
		}
		if err := validateDuration(c.Spec.Session.Human.ClientlessDuration); err != nil {
			return errors.Errorf("Invalid HUMAN clientless session duration")
		}

		if err := validateDuration(c.Spec.Session.Human.RefreshTokenDuration); err != nil {
			return err
		}

		if err := validateDuration(c.Spec.Session.Human.AccessTokenDuration); err != nil {
			return err
		}
	}

	if c.Spec.Session.Workload != nil {
		if err := validateDuration(c.Spec.Session.Workload.ClientDuration); err != nil {
			return errors.Errorf("Invalid WORKLOAD client session duration")
		}
		if err := validateDuration(c.Spec.Session.Workload.ClientlessDuration); err != nil {
			return errors.Errorf("Invalid WORKLOAD clientless session duration")
		}

		if err := validateDuration(c.Spec.Session.Workload.RefreshTokenDuration); err != nil {
			return err
		}

		if err := validateDuration(c.Spec.Session.Workload.AccessTokenDuration); err != nil {
			return err
		}
	}

	return nil
}

func validateDuration(d *metav1.Duration) error {
	return apisrvcommon.ValidateDuration(d)
}

func validateGateway(c *corev1.ClusterConfig) error {
	if c.Spec.Gateway == nil {
		return nil
	}
	if err := validateDuration(c.Spec.Gateway.WireguardKeyRotationDuration); err != nil {
		return errors.Errorf("Invalid WireGuard key rotation duration")
	}

	return nil
}

func (s *Server) validateCCAuthorization(ctx context.Context, c *corev1.ClusterConfig) error {
	if c.Spec.Authorization == nil {
		return nil
	}

	if err := s.validatePolicyOwner(ctx, c.Spec.Authorization); err != nil {
		return err
	}

	return nil
}

func validateDNS(c *corev1.ClusterConfig) error {
	if c.Spec.Dns == nil {
		return nil
	}

	if c.Spec.Dns.FallbackZone != nil {
		if c.Spec.Dns.FallbackZone.Servers != nil {
			if err := validateDuration(c.Spec.Dns.FallbackZone.CacheDuration); err != nil {
				return errors.Errorf("Invalid WORKLOAD client session duration")
			}

			if len(c.Spec.Dns.FallbackZone.Servers) > 32 {
				return grpcutils.InvalidArg("Too many DNS servers")
			}

			for _, srv := range c.Spec.Dns.FallbackZone.Servers {
				switch {
				case govalidator.IsDNSName(srv), govalidator.IsIP(srv), govalidator.IsURL(srv):
				default:
					return grpcutils.InvalidArg("Invalid DNS server: %s", srv)
				}
			}
		}
	}

	return nil
}
