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

package mock

import (
	"context"

	"github.com/jkusniar/lara"
)

// UserService mock implementation
type UserService struct {
	AuthenticateFn      func(login, password string) (*lara.User, error)
	AuthenticateInvoked bool

	RegisterFn     func(login, password string, permissions []lara.PermissionType) error
	RegiserInvoked bool
}

// Authenticate mock implementation
func (s *UserService) Authenticate(ctx context.Context, login, password string) (*lara.User, error) {
	s.AuthenticateInvoked = true
	return s.AuthenticateFn(login, password)
}

// Register mock implementation
func (s *UserService) Register(ctx context.Context, login, password string, permissions []lara.PermissionType) error {
	s.RegiserInvoked = true
	return s.RegisterFn(login, password, permissions)
}

func (s *UserService) Grant(ctx context.Context, login string, permissions []lara.PermissionType) error {
	return nil
}
func (s *UserService) Revoke(ctx context.Context, login string, permissions []lara.PermissionType) error {
	return nil
}
