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
)

func TestGetAllTitles(t *testing.T) {
	// get OK
	lovs, err := titleService.GetAllTitles(testCtx)
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}
	if lovs == nil {
		t.Fatal("expected not nil result")
	}
	if len(lovs.Items) != 1 {
		t.Fatalf("unexpected result %+v", lovs)
	}
}

func TestGetAllUnits(t *testing.T) {
	// get OK
	lovs, err := unitService.GetAllUnits(testCtx)
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}
	if lovs == nil {
		t.Fatal("expected not nil result")
	}
	if len(lovs.Items) != 2 {
		t.Fatalf("unexpected result %+v", lovs)
	}
}

func TestGetAllGenders(t *testing.T) {
	// get OK
	lovs, err := genderService.GetAllGenders(testCtx)
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}
	if lovs == nil {
		t.Fatal("expected not nil result")
	}
	if len(lovs.Items) != 2 {
		t.Fatalf("unexpected result %+v", lovs)
	}
}

func TestGetAllSpecies(t *testing.T) {
	// get OK
	lovs, err := speciesService.GetAllSpecies(testCtx)
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}
	if lovs == nil {
		t.Fatal("expected not nil result")
	}
	if len(lovs.Items) != 2 {
		t.Fatalf("unexpected result %+v", lovs)
	}
}

func TestGetAllBreeds(t *testing.T) {
	// get OK - empty
	lovs, err := breedService.GetAllBreedsBySpecies(testCtx, 1000)
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}
	if lovs == nil {
		t.Fatal("expected not nil result")
	}
	if len(lovs.Items) != 0 {
		t.Fatalf("unexpected result %+v", lovs)
	}

	// get OK - found
	lovs, err = breedService.GetAllBreedsBySpecies(testCtx, 1)
	if err != nil {
		t.Fatalf("expected nil error, but was %+v", err)
	}
	if lovs == nil {
		t.Fatal("expected not nil result")
	}
	if len(lovs.Items) != 2 {
		t.Fatalf("unexpected result %+v", lovs)
	}
}
