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

import "testing"

func TestSearchFound(t *testing.T) {
	s, err := searchService.Search(testCtx, "NLY")

	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}

	if s == nil {
		t.Fatal("expected not nil result")
	}

	if s.Total != 2 {
		t.Fatalf("expected total count 2 but was %d", s.Total)
	}

	if s.Records == nil {
		t.Fatal("expected Records not nil")
	}

	if len(s.Records) != 2 {
		t.Fatalf("expected Records length 2 but was %d", len(s.Records))
	}
}

func TestSearchNotFound(t *testing.T) {
	s, err := searchService.Search(testCtx, "kkkkk")

	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}

	if s == nil {
		t.Fatal("expected not nil result")
	}

	if s.Total != 0 {
		t.Fatalf("expected total count 0 but was %d", s.Total)
	}

	if s.Records == nil {
		t.Fatal("expected Records not nil")
	}

	if len(s.Records) != 0 {
		t.Fatalf("expected Records length 0 but was %d", len(s.Records))
	}
}
