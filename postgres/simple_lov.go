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
	"context"
	"database/sql"
	"fmt"

	"github.com/jkusniar/lara"
	"github.com/pkg/errors"
)

// SimpleLovService is implementation of lara.TitleService, lara.UnitService, lara.GenderService,
// lara.SpeciesService, lara.BreedService using postgresql database.
type SimpleLovService struct {
	DB *sql.DB
}

//go:generate stringer -type=listOfValuesType -output lov_string.go
//requires golang.org/x/tools/cmd/stringer installed locally
//if new object type added to enum, run "go generate"
type listOfValuesType int

const (
	title listOfValuesType = iota
	unit
	gender
	species
	breed
)

type lovTable struct {
	Name        string
	GetAllQuery string
}

var lovTableNames = map[listOfValuesType]lovTable{
	title:   {Name: "lov_title"},
	unit:    {Name: "lov_unit"},
	gender:  {Name: "lov_gender"},
	species: {Name: "lov_species"},
	breed:   {Name: "lov_breed", GetAllQuery: `SELECT id, name FROM lov_breed WHERE lov_species_id = $1`},
}

func lovGetAll(ctx context.Context, db *sql.DB, lov listOfValuesType, params ...interface{}) (*lara.LOVItemList, error) {
	q := lovTableNames[lov].GetAllQuery
	if len(q) == 0 {
		q = fmt.Sprintf(`SELECT id, name FROM %s`, lovTableNames[lov].Name)
	}

	rows, err := db.QueryContext(ctx, q, params...)
	if err != nil {
		return nil, errors.Wrapf(err, "get all from %s query error", lov)
	}
	defer rows.Close()

	var result = lara.LOVItemList{Items: []lara.LOVItem{}}
	for rows.Next() {
		var l lara.LOVItem
		if err := rows.Scan(&l.ID, &l.Name); err != nil {
			return nil, errors.Wrap(err, "scan DTO error")
		}
		result.Items = append(result.Items, l)
	}
	err = rows.Err()

	return &result, errors.Wrap(err, "rows processing errror")
}

// GetAllTitles is implementation of TitleService.GetAllTitles using postgresql database.
func (s *SimpleLovService) GetAllTitles(ctx context.Context) (*lara.LOVItemList, error) {
	return lovGetAll(ctx, s.DB, title)
}

// GetAllUnits is implementation of UnitService.GetAllUnits using postgresql database.
func (s *SimpleLovService) GetAllUnits(ctx context.Context) (*lara.LOVItemList, error) {
	return lovGetAll(ctx, s.DB, unit)
}

// GetAllSpecies is implementation of SpeciesService.GetAllSpecies using postgresql database.
func (s *SimpleLovService) GetAllSpecies(ctx context.Context) (*lara.LOVItemList, error) {
	return lovGetAll(ctx, s.DB, species)
}

// GetAllGenders is implementation of GenderService.GetAllGenders using postgresql database.
func (s *SimpleLovService) GetAllGenders(ctx context.Context) (*lara.LOVItemList, error) {
	return lovGetAll(ctx, s.DB, gender)
}

// GetAllBreedsBySpecies is implementation of BreedService.GetAllBreedsBySpecies using postgresql database.
func (s *SimpleLovService) GetAllBreedsBySpecies(ctx context.Context, speciesID uint64) (*lara.LOVItemList, error) {
	return lovGetAll(ctx, s.DB, breed, speciesID)
}
