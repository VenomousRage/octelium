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

package main

import (
	"context"

	"github.com/octelium/octelium/cluster/common/components"
	"github.com/octelium/octelium/cluster/ingress/ingress"
	"go.uber.org/zap"
)

func init() {
	components.SetComponentNamespace(components.ComponentNamespaceOctelium)
	components.SetComponentType(components.Ingress)
}

func main() {
	if err := components.InitComponent(context.Background(), nil); err != nil {
		zap.L().Fatal("init component err", zap.Error(err))
	}

	err := ingress.Run()
	if err != nil {
		zap.L().Fatal("Could not run main", zap.Error(err))
	}
}
