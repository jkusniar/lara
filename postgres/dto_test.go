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

func TestOwnerNameDTO(t *testing.T) {
	dto := OwnerNameDTO{LastName: "b"}
	if dto.String() != "b" {
		t.Fatalf("unexpected result '%s'", dto.String())
	}

	dto = OwnerNameDTO{FirstName: sql.NullString{String: "a", Valid: true},
		LastName: "b"}
	if dto.String() != "a b" {
		t.Fatalf("unexpected result '%s'", dto.String())
	}

	dto = OwnerNameDTO{FirstName: sql.NullString{String: "a", Valid: true},
		LastName: "b",
		Title:    sql.NullString{String: "c", Valid: true}}
	if dto.String() != "a b, c" {
		t.Fatalf("unexpected result '%s'", dto.String())
	}

	dto = OwnerNameDTO{LastName: "b",
		Title: sql.NullString{String: "c", Valid: true}}
	if dto.String() != "b, c" {
		t.Fatalf("unexpected result '%s'", dto.String())
	}
}

func TestOwnerAddressDTO(t *testing.T) {
	dto := OwnerAddressDTO{}
	if dto.String() != "" {
		t.Fatalf("unexpected result %s", dto.String())
	}

	dto = OwnerAddressDTO{City: sql.NullString{String: "a", Valid: true}}
	if dto.String() != "a " {
		t.Fatalf("unexpected result '%s'", dto.String())
	}

	dto = OwnerAddressDTO{City: sql.NullString{String: "a", Valid: true},
		HouseNo: sql.NullString{String: "1", Valid: true}}
	if dto.String() != "a 1" {
		t.Fatalf("unexpected result '%s'", dto.String())
	}

	dto = OwnerAddressDTO{City: sql.NullString{String: "a", Valid: true},
		Street:  sql.NullString{String: "b", Valid: true},
		HouseNo: sql.NullString{String: "1", Valid: true}}
	if dto.String() != "b 1, a" {
		t.Fatalf("unexpected result '%s'", dto.String())
	}

	dto = OwnerAddressDTO{City: sql.NullString{String: "a", Valid: true},
		Street: sql.NullString{String: "b", Valid: true}}
	if dto.String() != "b , a" {
		t.Fatalf("unexpected result '%s'", dto.String())
	}
}
