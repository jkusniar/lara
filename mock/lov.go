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

package mock

import (
	"context"

	"github.com/jkusniar/lara"
)

// SimpleLovService mock implementation
type SimpleLovService struct {
	GetAllTitlesFn      func() (*lara.LOVItemList, error)
	GetAllTitlesInvoked bool

	GetAllUnitsFn      func() (*lara.LOVItemList, error)
	GetAllUnitsInvoked bool

	GetAllGendersFn      func() (*lara.LOVItemList, error)
	GetAllGendersInvoked bool

	GetAllSpeciesFn      func() (*lara.LOVItemList, error)
	GetAllSpeciesInvoked bool

	GetAllBreedsFn      func(speciesId uint64) (*lara.LOVItemList, error)
	GetAllBreedsInvoked bool
}

func (s *SimpleLovService) GetAllSpecies(ctx context.Context) (*lara.LOVItemList, error) {
	s.GetAllSpeciesInvoked = true
	return s.GetAllSpeciesFn()
}

func (s *SimpleLovService) GetAllGenders(ctx context.Context) (*lara.LOVItemList, error) {
	s.GetAllGendersInvoked = true
	return s.GetAllGendersFn()
}

func (s *SimpleLovService) GetAllUnits(ctx context.Context) (*lara.LOVItemList, error) {
	s.GetAllUnitsInvoked = true
	return s.GetAllUnitsFn()
}

func (s *SimpleLovService) GetAllTitles(ctx context.Context) (*lara.LOVItemList, error) {
	s.GetAllTitlesInvoked = true
	return s.GetAllTitlesFn()
}

func (s *SimpleLovService) GetAllBreedsBySpecies(ctx context.Context, speciesID uint64) (*lara.LOVItemList, error) {
	s.GetAllBreedsInvoked = true
	return s.GetAllBreedsFn(speciesID)
}
