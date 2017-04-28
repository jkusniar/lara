/*
   Copyright (C) 2016-2017 Contributors as noted in the AUTHORS file

   This file is part of lara, veterinary practice support software.

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package http

import (
	"net/http"
	"strings"

	"github.com/jkusniar/lara"
	"github.com/pkg/errors"
	"github.com/pressly/chi/render"
)

type loginMsg struct {
	User string `json:"username"`
	Pass string `json:"password"`
}

// AuthToken is authentication token provider/validator
type AuthToken interface {
	Create(*lara.User) (string, error)
	Parse(token string) (*lara.User, error)
}

// authenticationHandler performs user authentication an returns a JSON Web Token.
func (s *Server) authenticationHandler(w http.ResponseWriter, r *http.Request) {
	var l loginMsg
	if err := render.DecodeJSON(r.Body, &l); err != nil {
		renderBadJSONError(w, r, err)
		return
	}

	u, err := s.UserService.Authenticate(r.Context(), l.User, l.Pass)
	if err != nil {
		renderError(w, r, err)
		return
	}

	token, err := s.Token.Create(u)
	if err != nil {
		renderError(w, r, err)
		return
	}

	render.PlainText(w, r, token)
}

// requireAuthorizedUser is authorization middleware.
// Takes care of authorization token validation. Returns HTTP 401 if token missing/invalid
// User object is passed through request context after successful token validation.
func (s *Server) requireAuthorizedUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ah := r.Header.Get("Authorization")
		if ah == "" {
			renderError(w, r,
				lara.NewCodedError(http.StatusUnauthorized,
					errors.New("authorization header missing")))
			return
		}

		if len(ah) <= 7 || strings.ToUpper(ah[0:6]) != "BEARER" {
			renderError(w, r,
				lara.NewCodedError(http.StatusUnauthorized,
					errors.New("authorization header shoud contain \"Bearer\" <JWT token>")))
			return
		}

		u, err := s.Token.Parse(ah[7:])
		if err != nil {
			renderError(w, r, lara.NewCodedError(http.StatusUnauthorized, err))
			return
		}

		ctx := lara.ContextWithUser(r.Context(), u)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

// requirePermission is authorization middleware.
// Panics, if requireAuthorizedUser is not set in middleware stack first.
// Returns HTTP 403 if user doesn't have required permission
func requirePermission(perm lara.PermissionType) func(next http.Handler) http.Handler {
	f := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			u, ok := lara.UserFromContext(r.Context())
			if !ok {
				// call requireAuthorizedUser first !
				panic("no user in context")
			}

			if !u.Permissions[perm] {
				renderError(w, r,
					lara.NewCodedError(http.StatusForbidden,
						errors.Errorf("user doesn't have permission %s", perm)))
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}

	return f
}
