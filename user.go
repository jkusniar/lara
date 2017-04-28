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

package lara

import (
	"context"
	"fmt"
)

// -----------------------------------------------------------------------------
// USER MANAGEMENT SERVICE

// User is application user
type User struct {
	Login       string
	Permissions map[PermissionType]bool // map instead of slice, for fast searching
}

// MakeUser creates User from login and list of permissions
func MakeUser(login string, permissions []string) (u *User, err error) {
	u = &User{Login: login}
	u.Permissions, err = MakePermissions(permissions)
	return
}

// MakePermissions converts string slice to PermissionType set
func MakePermissions(perms []string) (map[PermissionType]bool, error) {
	result := make(map[PermissionType]bool)
	for _, s := range perms {
		var p PermissionType
		if err := p.FromString(s); err != nil {
			return nil, err
		}
		result[p] = true
	}

	return result, nil
}

// PermissionType is enum defining system permissions
//go:generate stringer -type=PermissionType -output permission_string.go
//requires golang.org/x/tools/cmd/stringer installed locally
//if new object permission added to enum, run "go generate"
type PermissionType int

// Permissions enum
const (
	ViewRecord PermissionType = iota
	EditRecord
	ViewReports
	EditProducts
)

// FromString creates PermissionType from string
func (i *PermissionType) FromString(perm string) error {
	// TODO: more intelligent implementation
	switch perm {
	case "ViewRecord":
		*i = ViewRecord
	case "EditRecord":
		*i = EditRecord
	case "ViewReports":
		*i = ViewReports
	case "EditProducts":
		*i = EditProducts
	default:
		return fmt.Errorf("bad PermissionType: '%s'", perm)
	}

	return nil
}

var DefaultPermissions = []PermissionType{ViewRecord}

// UserService manages application's users
type UserService interface {
	Authenticate(ctx context.Context, login, password string) (*User, error)
	Register(ctx context.Context, login, password string, permissions []PermissionType) error
	Grant(ctx context.Context, login string, permissions []PermissionType) error
	Revoke(ctx context.Context, login string, permissions []PermissionType) error
}

// -----------------------------------------------------------------------------
// User in Context

type keyType int

const userKey keyType = 0

// ContextWithUser returns a new Context that carries value u.
func ContextWithUser(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// UserFromContext returns the User value stored in ctx, if any.
func UserFromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userKey).(*User)
	return u, ok
}
