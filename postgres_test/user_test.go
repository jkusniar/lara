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

package postgres_test

import (
	"testing"

	"github.com/jkusniar/lara"
)

func TestRegister(t *testing.T) {
	// login empty
	err := userService.Register(testCtx, "", "x", lara.DefaultPermissions)
	if err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 400); !ok {
		t.Fatalf("expected error code 400 but was %d, %+v", actual, err)
	}

	// login too long
	err = userService.Register(testCtx, "aaaaaaaaaabbbbbbbbbbX", "x", lara.DefaultPermissions)
	if err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 400); !ok {
		t.Fatalf("expected error code 400 but was %d, %+v", actual, err)
	}

	// password empty
	err = userService.Register(testCtx, "aaa", "", lara.DefaultPermissions)
	if err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 400); !ok {
		t.Fatalf("expected error code 400 but was %d, %+v", actual, err)
	}

	// Register existing user
	err = userService.Register(testCtx, "test", "TestPassword", lara.DefaultPermissions)
	if err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 400); !ok {
		t.Fatalf("expected error code 400 but was %d, %+v", actual, err)
	}

	// Register new user (first permission is created, second is already existing
	if err := userService.Register(testCtx, "test2", "TestPassword",
		[]lara.PermissionType{lara.ViewRecord, lara.EditRecord}); err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}
}

func TestAuthenticate(t *testing.T) {
	// non-existing user
	_, err := userService.Authenticate(testCtx, "jimi", "TestPassword")
	if err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 401); !ok {
		t.Fatalf("expected error code 401 but was %d, %+v", actual, err)
	}

	// bad password
	_, err = userService.Authenticate(testCtx, "test", "BadPass")
	if err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 401); !ok {
		t.Fatalf("expected error code 401 but was %d, %+v", actual, err)
	}

	// OK
	if u, err := userService.Authenticate(testCtx, "test", "TestPassword"); err != nil {
		t.Fatalf("expected nil, but was %+v", err)
	} else if u == nil || u.Login != "test" || len(u.Permissions) != 2 || !u.Permissions[lara.EditRecord] {
		t.Fatalf("unexpected result %+v", u)
	}

}

func TestGrant(t *testing.T) {
	// Grant already granted
	if err := userService.Grant(testCtx, "test", []lara.PermissionType{lara.ViewReports}); err == nil {
		t.Fatal("expected error")
	}

	// Grant OK
	if err := userService.Grant(testCtx, "test", []lara.PermissionType{lara.ViewRecord}); err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}
}

func TestRevoke(t *testing.T) {
	// Revoke not granted (silently passes)
	if err := userService.Revoke(testCtx, "test", []lara.PermissionType{lara.EditProducts}); err != nil {
		t.Fatalf("expected nil error for not granted , but was %+v", err)
	}

	// Revoke OK
	if err := userService.Revoke(testCtx, "test", []lara.PermissionType{lara.EditRecord}); err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}
}
