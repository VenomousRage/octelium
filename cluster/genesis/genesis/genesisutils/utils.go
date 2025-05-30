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

package genesisutils

import (
	"context"

	"github.com/octelium/octelium/apis/cluster/cbootstrapv1"
	"github.com/octelium/octelium/apis/main/corev1"
	"github.com/octelium/octelium/apis/rsc/rmetav1"
	"github.com/octelium/octelium/cluster/apiserver/apiserver/admin"
	"github.com/octelium/octelium/cluster/common/octeliumc"
	"github.com/octelium/octelium/pkg/grpcerr"
	"go.uber.org/zap"
)

type InstallCtx struct {
	ClusterConfig *corev1.ClusterConfig
	Bootstrap     *cbootstrapv1.Config
	Region        *corev1.Region
}

func CreateOrUpdateService(ctx context.Context, octeliumC octeliumc.ClientInterface, svc *corev1.Service) error {

	adminSrv := admin.NewServer(&admin.Opts{
		OcteliumC:  octeliumC,
		IsEmbedded: true,
	})

	if svc.Spec == nil {
		svc.Spec = &corev1.Service_Spec{}
	}

	if svc.Status == nil {
		svc.Status = &corev1.Service_Status{}
	}

	if itm, err := octeliumC.CoreC().GetService(ctx, &rmetav1.GetOptions{
		Name: svc.Metadata.Name,
	}); err == nil {

		itm.Spec = svc.Spec
		itm.Status.ManagedService = svc.Status.ManagedService
		itm.Metadata.Labels = svc.Metadata.Labels
		itm.Metadata.SystemLabels = svc.Metadata.SystemLabels
		itm.Metadata.DisplayName = svc.Metadata.DisplayName
		itm.Metadata.Annotations = svc.Metadata.Annotations
		itm.Metadata.PicURL = svc.Metadata.PicURL

		zap.L().Debug("Updating Service",
			zap.String("name", svc.Metadata.Name))

		if _, err := octeliumC.CoreC().UpdateService(ctx, itm); err != nil {
			return err
		}

		return nil
	} else if !grpcerr.IsNotFound(err) {
		return err
	}

	zap.L().Debug("Creating Service",
		zap.String("name", svc.Metadata.Name))

	if _, err := adminSrv.DoCreateService(ctx, svc, true); err != nil {
		return err
	}

	return nil
}

func CreateOrUpdateNamespace(ctx context.Context, octeliumC octeliumc.ClientInterface, ns *corev1.Namespace) error {

	if ns.Status == nil {
		ns.Status = &corev1.Namespace_Status{}
	}
	if ns.Spec == nil {
		ns.Spec = &corev1.Namespace_Spec{}
	}

	if itm, err := octeliumC.CoreC().GetNamespace(ctx, &rmetav1.GetOptions{
		Name: ns.Metadata.Name,
	}); err == nil {
		itm.Spec = ns.Spec
		itm.Metadata.Labels = ns.Metadata.Labels
		itm.Metadata.SystemLabels = ns.Metadata.SystemLabels
		itm.Metadata.IsSystemHidden = ns.Metadata.IsSystemHidden
		itm.Metadata.IsUserHidden = ns.Metadata.IsUserHidden
		itm.Metadata.IsSystem = ns.Metadata.IsSystem
		itm.Metadata.DisplayName = ns.Metadata.DisplayName
		itm.Metadata.Annotations = ns.Metadata.Annotations
		itm.Metadata.PicURL = ns.Metadata.PicURL

		_, err = octeliumC.CoreC().UpdateNamespace(ctx, itm)
		if err != nil {
			return err
		}

		return nil
	} else if !grpcerr.IsNotFound(err) {
		return err
	}

	_, err := octeliumC.CoreC().CreateNamespace(ctx, ns)
	if err != nil {
		return err
	}

	return nil
}
