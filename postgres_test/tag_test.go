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

func TestCreateTag(t *testing.T) {
	var err error
	x := &lara.CreateTag{}

	// value missing
	if _, err = tagService.Create(testCtx, x); err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 400); !ok {
		t.Fatalf("expected error code 400 but was %d, %+v", actual, err)
	}

	// type missing
	x.Value = "1982-SK-0728"
	if _, err = tagService.Create(testCtx, x); err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 400); !ok {
		t.Fatalf("expected error code 400 but was %d, %+v", actual, err)
	}

	// type conversion error
	x.Type = "Fail"
	if _, err = tagService.Create(testCtx, x); err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 400); !ok {
		t.Fatalf("expected error code 400 but was %d, %+v", actual, err)
	}

	// patientId missing
	x.Type = "Tattoo"
	if _, err = tagService.Create(testCtx, x); err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 400); !ok {
		t.Fatalf("expected error code 400 but was %d, %+v", actual, err)
	}

	// OK
	x.PatientID = 3
	id, err := tagService.Create(testCtx, x)
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}

	if id == 0 {
		t.Fatal("incorrect ID returned")
	}
}

func TestUpdateTag(t *testing.T) {
	x := &lara.UpdateTag{Version: 1, Data: []byte{11, 11, 11, 11}}
	var err error

	// value missing
	if err = tagService.Update(testCtx, 2, x); err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 400); !ok {
		t.Fatalf("expected error code 400 but was %d, %+v", actual, err)
	}

	// bad ID
	x.Value = "updated-tag-value"
	if err = tagService.Update(testCtx, 100, x); err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 404); !ok {
		t.Fatalf("expected error code 404 but was %d, %+v", actual, err)
	}

	// version mismatch
	if err = tagService.Update(testCtx, 2, x); err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 409); !ok {
		t.Fatalf("expected error code 409 but was %d, %+v", actual, err)
	}

	// OK
	x.Version = 2
	err = tagService.Update(testCtx, 2, x)
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}
}

func TestGetTag(t *testing.T) {
	var x *lara.GetTag
	var err error

	// get not existing
	x, err = tagService.Get(testCtx, 100)
	if err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 404); !ok {
		t.Fatalf("expected error code 404 but was %d, %+v", actual, err)
	}

	// get OK
	x, err = tagService.Get(testCtx, 1)
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}
	if x == nil {
		t.Fatal("expected not nil result")
	}

	if x.ID != 1 || x.Value != "2017-SK-0007" || len(x.Data) != 0 || x.Type != "LyssaVirus" {
		t.Fatalf("unexpected result %+v", x)
	}
}

func TestGetPatientByTag(t *testing.T) {
	var x *lara.PatientByTag
	var err error

	// get by empty tag
	x, err = tagService.GetPatientByTag(testCtx, "")
	if err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 404); !ok {
		t.Fatalf("expected error code 404 but was %d, %+v", actual, err)
	}

	// get not existing
	x, err = tagService.GetPatientByTag(testCtx, "NOT-EXISTING")
	if err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 404); !ok {
		t.Fatalf("expected error code 404 but was %d, %+v", actual, err)
	}

	// get OK
	x, err = tagService.GetPatientByTag(testCtx, "2017-SK-0007")
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}
	if x == nil {
		t.Fatal("expected not nil result")
	}
	if x.OwnerID != 1 || x.TagType != "LyssaVirus" {
		t.Fatalf("unexpected result %+v", x)
	}
}
