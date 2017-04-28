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
	"time"

	"github.com/jkusniar/lara"
)

func TestGetPatient(t *testing.T) {
	var p *lara.GetPatient
	var err error

	// get not existing
	p, err = patientService.Get(testCtx, 10000)
	if err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 404); !ok {
		t.Fatalf("expected error code 404 but was %d, %+v", actual, err)
	}

	// get OK
	p, err = patientService.Get(testCtx, 1)
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}
	if p == nil {
		t.Fatal("expected not nil result")
	}
	if p.Name != "test-pet" || len(p.Records) != 7 || len(p.Tags) != 1 {
		t.Fatalf("unexpected result %+v", p)
	}
}

func TestCreatePatient(t *testing.T) {
	var err error
	p := &lara.CreatePatient{OwnerID: 2, NewPatient: lara.NewPatient{Patient: lara.Patient{Name: "unit-pet"}}}

	id, err := patientService.Create(testCtx, p)
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}
	if id == 0 {
		t.Fatal("incorrect ID returned")
	}
}

func TestUpdatePatient(t *testing.T) {
	p := &lara.UpdatePatient{Patient: lara.Patient{BirthDate: time.Now()}, Version: 0}
	var err error

	// update SQL error (not null constraint)
	err = patientService.Update(testCtx, 1, p)
	if err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 400); !ok {
		t.Fatalf("expected error code 400 but was %d, %+v", actual, err)
	}

	// update not existing
	p.Name = "updated-name"
	err = patientService.Update(testCtx, 100, p)
	if err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 404); !ok {
		t.Fatalf("expected error code 404 but was %d, %+v", actual, err)
	}

	// update OK
	err = patientService.Update(testCtx, 1, p)
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}

	// second update without version update = fail
	p.Name = "updated-name"
	err = patientService.Update(testCtx, 1, p)
	if err == nil {
		t.Fatal("expected error")
	}
	if ok, actual := checkErrCode(err, 409); !ok {
		t.Fatalf("expected error code 409 but was %d, %+v", actual, err)
	}

	// second update with version update = OK
	p.Version = p.Version + 1
	err = patientService.Update(testCtx, 1, p)
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}
}
