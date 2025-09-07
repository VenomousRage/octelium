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

package authserver

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/octelium/octelium/apis/main/metav1"
	"github.com/octelium/octelium/apis/rsc/rcachev1"
	"github.com/octelium/octelium/pkg/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var rgxStateID = regexp.MustCompile(`^[a-zA-Z0-9]{36}$`)

type loginState struct {
	ID          string
	CallbackURL string
	IsApp       bool

	UID       string
	RequestID string
	LoginURL  string
}

func (s *server) saveLoginState(ctx context.Context, state *loginState) error {
	stateBytes, err := json.Marshal(state)
	if err != nil {
		return err
	}

	if _, err := s.octeliumC.CacheC().SetCache(ctx, &rcachev1.SetCacheRequest{
		Key:  []byte(getAuthKey(state.ID)),
		Data: stateBytes,
		Duration: &metav1.Duration{
			Type: &metav1.Duration_Minutes{
				Minutes: 8,
			},
		},
	}); err != nil {
		return err
	}

	return nil
}

func (s *server) getLoginStateFromCallback(r *http.Request, doDelete bool) (*loginState, error) {
	var stateID string

	switch r.Method {
	case http.MethodGet:
		if errMsg := r.FormValue("error"); errMsg != "" {
			return nil, errors.Errorf("%s", r.FormValue("error_description"))
		}

		if stateID = r.FormValue("state"); stateID == "" {
			return nil, errors.Errorf("No state found")
		}

	case http.MethodPost:
		if stateID = r.PostFormValue("RelayState"); stateID == "" {
			return nil, errors.Errorf("No state found")
		}

	default:
		return nil, errors.Errorf("Invalid HTTP method")
	}

	if !rgxStateID.MatchString(stateID) {
		return nil, errors.Errorf("Invalid state ID")
	}

	stateIDCookie, err := getStateIDFromCookie(r)
	if err != nil {
		return nil, errors.Errorf("Could not get state ID from cookie")
	}

	if !utils.SecureStringEqual(stateID, stateIDCookie) {
		return nil, errors.Errorf("state ID and cookie don't match")
	}

	return s.getLoginStateFromStateID(r.Context(), stateID, doDelete)
}

func getStateIDFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("octelium_login_state")
	if err != nil {
		return "", err
	}

	if !rgxStateID.MatchString(cookie.Value) {
		return "", errors.Errorf("Invalid state ID")
	}

	return cookie.Value, nil
}

func (s *server) getLoginStateFromStateID(ctx context.Context, stateID string, doDelete bool) (*loginState, error) {

	/*
		redisCmd := func() *redis.StringCmd {
			if doDelete {
				return s.redisC.GetDel(ctx, getAuthKey(stateID))
			} else {
				return s.redisC.Get(ctx, getAuthKey(stateID))
			}
		}()

		stateStr, err := redisCmd.Result()
		if err != nil {
			if err == redis.Nil {
				return nil, errors.Errorf("Invalid or expired state. Please try again.")
			}
			zap.S().Errorf("Could not get key %s: %+v", getAuthKey(stateID), err)
			return nil, errors.Errorf("Invalid or expired state. Please try again.")
		}
	*/

	res, err := s.octeliumC.CacheC().GetCache(ctx, &rcachev1.GetCacheRequest{
		Key:    []byte(getAuthKey(stateID)),
		Delete: doDelete,
	})
	if err != nil {
		return nil, err
	}

	var userState loginState
	if err := json.Unmarshal(res.Data, &userState); err != nil {
		zap.S().Errorf("Could not unmarshal json: %+v", err)
		return nil, errors.Errorf("Invalid or expired state. Please try again.")
	}

	return &userState, nil
}
