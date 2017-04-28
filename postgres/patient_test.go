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
	"time"
)

func TestPatientsRecordDTO(t *testing.T) {
	now := time.Now()

	r := patientsRecordDTO{
		ID:   42,
		Date: now,
		Text: sql.NullString{String: "dog", Valid: true},
	}

	rr := r.toPatientsRecord()

	if rr == nil {
		t.Fatal("expected not nil")
	}

	// all field mapping
	if rr.ID != r.ID || rr.Date != r.Date || rr.Text != r.Text.String {
		t.Fatalf("data mapping error. Mapped struct is %#v", rr)
	}

	// long text trim
	r.Text = sql.NullString{String: "123456789012345678901234567890XXX", Valid: true}
	rr = r.toPatientsRecord()

	const trimmed = "123456789012345678901234567890"
	if rr.Text != trimmed {
		t.Fatalf("expected text %s, but was %s", trimmed, rr.Text)
	}
}
