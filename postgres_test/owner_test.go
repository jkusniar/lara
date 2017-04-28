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

func TestGetOwner(t *testing.T) {
	var o *lara.GetOwner
	var err error

	// get not existing
	o, err = ownerService.Get(testCtx, 100)
	if err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 404); !ok {
		t.Fatalf("expected error code 404 but was %d, %+v", actual, err)
	}

	// get OK
	o, err = ownerService.Get(testCtx, 5)
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}
	if o == nil {
		t.Fatal("expected not nil result")
	}
	// test owner's fields
	if o.FirstName != "Test" || o.LastName != "GetOwner" ||
		o.TitleID != 1 || o.CityID != 1 || o.StreetID != 1 || o.HouseNo != "1B" ||
		o.Phone1 != "000111" ||
		o.Phone2 != "000222" ||
		o.Email != "test@test.com" ||
		o.Note != "test note" ||
		o.IC != "0001" ||
		o.DIC != "0002" ||
		o.ICDPH != "0003" ||
		o.ID != 5 ||
		o.Version != 3 ||
		o.Creator != "testuser" || o.Modifier != "testmodifier" ||
		timeEmpty(o.Created) || timeEmpty(o.Modified) ||
		o.Title != "Ing." || o.City != "test city" || o.Street != "test street" ||
		len(o.Patients) != 1 {
		t.Fatalf("unexpected result %+v", o)
	}
	// test patient's fields
	if o.Patients[0].ID != 2 || o.Patients[0].Name != "get-owner-pet" ||
		o.Patients[0].Breed != "german shepard" || o.Patients[0].Species != "dog" ||
		o.Patients[0].Gender != "male" || o.Patients[0].Dead != false {
		t.Fatalf("unexpected result %+v", o.Patients[0])
	}
}

func TestUpdateOwner(t *testing.T) {
	o := &lara.UpdateOwner{Owner: lara.Owner{FirstName: "UpdatedFirst", Note: "A note", TitleID: 1}, Version: 42}
	var err error

	// update SQL error (not null constraint)
	err = ownerService.Update(testCtx, 2, o)
	if err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 400); !ok {
		t.Fatalf("expected error code 400 but was %d, %+v", actual, err)
	}

	// update not existing
	o.LastName = "UpdatedLast"
	err = ownerService.Update(testCtx, 100, o)
	if err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 404); !ok {
		t.Fatalf("expected error code 404 but was %d, %+v", actual, err)
	}

	// update OK
	err = ownerService.Update(testCtx, 2, o)
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}

	// second update without version update = fail
	o.LastName = "UpdateLast2"
	err = ownerService.Update(testCtx, 2, o)
	if err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 409); !ok {
		t.Fatalf("expected error code 409 but was %d, %+v", actual, err)
	}

	// second update with version update = OK
	o.Version = o.Version + 1
	err = ownerService.Update(testCtx, 2, o)
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}
}

func TestCreateOwner(t *testing.T) {
	var id uint64
	var err error

	// not null contraint
	id, err = ownerService.Create(testCtx, &lara.CreateOwner{Owner: lara.Owner{FirstName: "CreatedFirst"}})
	if err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 400); !ok {
		t.Fatalf("expected error code 400 but was %d, %+v", actual, err)
	}

	// ok
	id, err = ownerService.Create(testCtx, &lara.CreateOwner{Owner: lara.Owner{FirstName: "CreatedFirst",
		LastName: "CreatedLast"}, Patient: lara.NewPatient{Patient: lara.Patient{Name: "unit-pet"}}})
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}
	if id == 0 {
		t.Fatal("expected id != 0")
	}
}
