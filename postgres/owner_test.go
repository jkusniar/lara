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

package postgres

import (
	"database/sql"
	"testing"
)

func TestOwnerDTO(t *testing.T) {
	o := ownerDTO{
		versionedDTO: versionedDTO{ID: 10, Version: 42},
		FirstName:    sql.NullString{String: "first", Valid: true},
		LastName:     "last",
		TitleID:      sql.NullInt64{Int64: 100, Valid: true},
		CityID:       sql.NullInt64{Int64: 101, Valid: true},
		StreetID:     sql.NullInt64{Int64: 102, Valid: true},
		HouseNo:      sql.NullString{String: "a", Valid: true},
		Phone1:       sql.NullString{String: "b", Valid: true},
		Phone2:       sql.NullString{String: "c", Valid: true},
		Email:        sql.NullString{String: "d", Valid: true},
		Note:         sql.NullString{String: "e", Valid: true},
		IC:           sql.NullString{String: "f", Valid: true},
		DIC:          sql.NullString{String: "g", Valid: true},
		ICDPH:        sql.NullString{String: "h", Valid: true},
	}

	lo := o.toGetOwner([]ownersPatientDTO{})

	if lo == nil {
		t.Fatal("expected not nil")
	}

	// all field mapping
	if lo.ID != 10 || lo.FirstName != "first" || lo.LastName != "last" ||
		lo.TitleID != 100 || lo.CityID != 101 || lo.StreetID != 102 ||
		lo.HouseNo != "a" || lo.Phone1 != "b" || lo.Phone2 != "c" ||
		lo.Email != "d" || lo.Note != "e" || lo.IC != "f" || lo.DIC != "g" ||
		lo.ICDPH != "h" || lo.Version != 42 || lo.Patients == nil || len(lo.Patients) != 0 {
		t.Fatalf("data mapping error. Mapped struct is %#v", lo)
	}

	// nil patients
	lo = o.toGetOwner(nil)
	if lo.Patients == nil || len(lo.Patients) != 0 {
		t.Fatalf("nil patients mapping error. Mapped struct is %#v", lo)
	}

	// patients mapping
	lo = o.toGetOwner([]ownersPatientDTO{{ID: 10, Name: "pet"}})
	if lo.Patients == nil || len(lo.Patients) != 1 {
		t.Fatalf("patients mapping error. Mapped struct is %#v", lo)
	}
}

func TestOwnersPatientsDTO(t *testing.T) {
	o := ownersPatientDTO{
		ID:      42,
		Name:    "name",
		Species: sql.NullString{String: "dog", Valid: true},
		Breed:   sql.NullString{String: "boxer", Valid: true},
		Sex:     sql.NullString{String: "male", Valid: true},
		Dead:    true,
	}

	lo := o.toOwnersPatient()

	if lo == nil {
		t.Fatal("expected not nil")
	}

	// all field mapping
	if lo.ID != 42 || lo.Name != "name" || lo.Species != "dog" ||
		lo.Breed != "boxer" || lo.Gender != "male" || !lo.Dead {
		t.Fatalf("data mapping error. Mapped struct is %#v", lo)
	}

	// is dead
	o1 := ownersPatientDTO{Dead: false}
	if o1.toOwnersPatient().Dead {
		t.Fatal("expected not dead")
	}
}
