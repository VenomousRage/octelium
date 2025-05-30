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
	"embed"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/octelium/octelium/apis/main/corev1"
	"github.com/octelium/octelium/pkg/apiutils/ucorev1"
	"go.uber.org/zap"
)

//go:embed web
var fsWeb embed.FS

func (s *server) handleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cookie, err := r.Cookie("octelium_rt")
	if err != nil {
		s.setLogoutCookies(w)
		s.renderIndex(w)
		return
	}

	sess, err := s.getSessionFromRefreshToken(ctx, cookie.Value)
	if err != nil {
		s.setLogoutCookies(w)
		s.renderIndex(w)
		return
	}

	if sess.Status.Type != corev1.Session_Status_CLIENTLESS {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !sess.Status.IsBrowser {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !ucorev1.ToSession(sess).ShouldRefresh() {
		zap.S().Debugf("No need to re-authenticate Session: %s", sess.Metadata.Name)
		if vReq := r.URL.Query().Get("octelium_req"); vReq == "" {
			if redirect := r.URL.Query().Get("redirect"); redirect != "" && s.isURLSameClusterOrigin(redirect) {
				http.Redirect(w, r, redirect, http.StatusSeeOther)
				return
			}

			http.Redirect(w, r, s.getPortalURL(), http.StatusSeeOther)
			return
		}

		// There should not be a referer if there is a octelium_req parameter
		if referer := r.Header.Get("referer"); referer != "" && !s.isURLSameClusterOrigin(referer) {
			http.Redirect(w, r, s.rootURL, http.StatusSeeOther)
			return
		}
		// We still go for a login to create a token for the client
	}

	zap.S().Debugf("Session: %s needs a login", sess.Metadata.Name)

	if sess.Status.InitialAuthentication == nil || sess.Status.InitialAuthentication.Info == nil ||
		sess.Status.InitialAuthentication.Info.GetIdentityProvider() == nil {
		s.setLogoutCookies(w)
		s.renderIndex(w)
		return
	}

	switch sess.Status.InitialAuthentication.Info.GetIdentityProvider().Type {
	case corev1.Session_Status_Authentication_Info_IdentityProvider_GITHUB,
		corev1.Session_Status_Authentication_Info_IdentityProvider_OIDC,
		corev1.Session_Status_Authentication_Info_IdentityProvider_SAML:
	default:
		s.setLogoutCookies(w)
		s.renderIndex(w)
		return
	}

	provider, err := s.getWebProviderFromUID(
		sess.Status.InitialAuthentication.Info.GetIdentityProvider().IdentityProviderRef.Uid)
	if err != nil {
		zap.S().Debugf("Could not get provider with name: %s. Probably removed by Cluster admins",
			sess.Status.InitialAuthentication.Info.GetIdentityProvider().IdentityProviderRef.Uid)
		s.setLogoutCookies(w)
		s.renderIndex(w)
		return
	}

	if provider.Provider().Spec.IsDisabled {
		s.setLogoutCookies(w)
		s.renderIndex(w)
		return
	}

	loginState, err := s.doGenerateLoginState(ctx, provider, r.URL.Query().Encode(), w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, loginState.LoginURL, http.StatusSeeOther)
}

func (s *server) isURLSameClusterOrigin(arg string) bool {

	if len(arg) == 0 || len(arg) > 1500 {
		return false
	}

	redirectURL, err := url.Parse(arg)
	if err != nil {
		return false
	}

	if !strings.HasSuffix(redirectURL.Hostname(), s.domain) {
		return false
	}

	return true
}

func (s *server) handleIndex(w http.ResponseWriter, r *http.Request) {

	/*
		isValid, _, err := s.isWebSessionValid(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if isValid {
			s.redirectToPortal(w, r)
			return
		}

	*/

	s.redirectToLogin(w, r)
}

func (s *server) redirectToLogin(w http.ResponseWriter, r *http.Request) {
	murl, _ := url.Parse(fmt.Sprintf("%s/login", s.rootURL))
	murl.RawQuery = r.URL.RawQuery
	http.Redirect(w, r, murl.String(), http.StatusSeeOther)
}

func (s *server) redirectToCallbackSuccess(w http.ResponseWriter, r *http.Request, redirectURL string) {
	murl, _ := url.Parse(fmt.Sprintf("%s/callback/success", s.rootURL))
	q := murl.Query()
	q.Set("redirect", redirectURL)
	murl.RawQuery = q.Encode()

	http.Redirect(w, r, murl.String(), http.StatusSeeOther)
}

func (s *server) redirectToPortal(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, s.getPortalURL(), http.StatusSeeOther)
}

func (s *server) getPortalURL() string {
	return fmt.Sprintf("https://portal.%s", s.domain)
}

/*
func (s *server) isWebSessionValid(r *http.Request) (bool, *corev1.Session, error) {
	ctx := r.Context()

	cookie, err := r.Cookie("octelium_rt")
	if err != nil {
		return false, nil, nil
	}

	claims, err := s.jwkCtl.VerifyRefreshToken(cookie.Value)
	if err != nil {
		return false, nil, nil
	}

	sess, err := s.octeliumC.CoreC().GetSession(ctx, &rmetav1.GetOptions{Uid: claims.SessionUID})
	if err != nil {
		if grpcerr.IsNotFound(err) {
			return false, nil, nil
		}
		return false, nil, err
	}

	isValid := ucorev1.ToSession(sess).HasValidAccessTokenByTokenID(claims.TokenID) &&
		!ucorev1.ToSession(sess).ShouldRefresh()

	return isValid, sess, nil
}
*/

func (s *server) handleStatic(w http.ResponseWriter, r *http.Request) {
	blob, err := fs.ReadFile(fsWeb, filepath.Join("web", r.URL.Path))
	if err != nil {
		zap.L().Debug("Could not read index.html file from web fs", zap.Error(err))
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(r.URL.Path)))
	w.Write(blob)
}

func (s *server) handleDenied(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("%s/denied", s.rootURL), http.StatusSeeOther)
}
