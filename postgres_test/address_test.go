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

func TestSearchCityFound(t *testing.T) {
	c, err := addressService.SearchCity(testCtx, "ci")

	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}

	if c == nil {
		t.Fatal("expected not nil result")
	}

	if c.Total != 2 {
		t.Fatalf("expected total count 2 but was %d", c.Total)
	}

	if c.Items == nil {
		t.Fatal("expected Items not nil")
	}

	if len(c.Items) != 2 {
		t.Fatalf("expected Items length 2 but was %d", len(c.Items))
	}

	if c.Items[1].ZIP != "11100" || c.Items[1].Name != "test city 2" || c.Items[1].ID != 2 {
		t.Fatalf("unexpected city returned %+v", c.Items[1])
	}
}

func TestSearchCityNotFound(t *testing.T) {
	c, err := addressService.SearchCity(testCtx, "xxxxxxx")

	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}

	if c == nil {
		t.Fatal("expected not nil result")
	}

	if c.Total != 0 {
		t.Fatalf("expected total count 0 but was %d", c.Total)
	}

	if c.Items == nil {
		t.Fatal("expected Items not nil")
	}

	if len(c.Items) != 0 {
		t.Fatalf("expected Items length 0 but was %d", len(c.Items))
	}
}

func TestSearchStreetFound(t *testing.T) {
	s, err := addressService.SearchStreetForCity(testCtx, 1, "ee")

	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}

	if s == nil {
		t.Fatal("expected not nil result")
	}

	if s.Total != 2 {
		t.Fatalf("expected total count 2 but was %d", s.Total)
	}

	if s.Items == nil {
		t.Fatal("expected Items not nil")
	}

	if len(s.Items) != 2 {
		t.Fatalf("expected Items length 2 but was %d", len(s.Items))
	}

	if s.Items[0].ZIP != "88800" || s.Items[0].Name != "another street" || s.Items[0].ID != 2 {
		t.Fatalf("unexpected street returned %+v", s.Items[0])
	}
}

func TestSearchStreetNotFound(t *testing.T) {
	s, err := addressService.SearchStreetForCity(testCtx, 1, "xxx")

	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}

	if s == nil {
		t.Fatal("expected not nil result")
	}

	if s.Total != 0 {
		t.Fatalf("expected total count 0 but was %d", s.Total)
	}

	if s.Items == nil {
		t.Fatal("expected Items not nil")
	}

	if len(s.Items) != 0 {
		t.Fatalf("expected Items length 0 but was %d", len(s.Items))
	}
}
