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

	"github.com/jkusniar/lara"
	"github.com/pkg/errors"
)

type AddressService struct {
	DB *sql.DB
}

func (s *AddressService) searchCityOrStreet(ctx context.Context,
	countQuery, dataQuery string, params ...interface{}) (*lara.CityStreetList, error) {
	var result = lara.CityStreetList{Items: []lara.CityStreet{}}
	// Count query. Ignore last param - search limit
	if err := s.DB.QueryRowContext(ctx, countQuery, params[:len(params)-1]...).Scan(&result.Total); err != nil {
		return nil, errors.Wrap(err, "search city/street count error")
	}

	rows, err := s.DB.QueryContext(ctx, dataQuery, params...)
	if err != nil {
		return nil, errors.Wrap(err, "search city/street query error")
	}
	defer rows.Close()

	for rows.Next() {
		var dto lara.CityStreet
		if err := rows.Scan(&dto.ID, &dto.Name, &dto.ZIP); err != nil {
			return nil, errors.Wrap(err, "scan DTO error")
		}
		result.Items = append(result.Items, dto)
	}
	err = rows.Err()

	return &result, errors.Wrap(err, "rows processing errror")
}

// TODO make search limit parametric
func (s *AddressService) SearchCity(ctx context.Context, query string) (*lara.CityStreetList, error) {
	const cq = `SELECT count(*) FROM lov_city WHERE  city ILIKE $1`
	const dq = `SELECT
			  id,
			  city,
			  psc
			FROM lov_city
			WHERE city ILIKE $1
			ORDER BY city
			LIMIT $2`

	return s.searchCityOrStreet(ctx, cq, dq, "%"+query+"%", 30)
}

func (s *AddressService) SearchStreetForCity(ctx context.Context, cityID uint64, query string) (*lara.CityStreetList, error) {
	const cq = `SELECT count(*) FROM lov_street WHERE city_id = $1 AND street ILIKE $2`
	const dq = `SELECT
			  id,
			  street,
			  psc
			FROM lov_street
			WHERE city_id = $1 AND street ILIKE $2
			ORDER BY street
			LIMIT $3`
	return s.searchCityOrStreet(ctx, cq, dq, cityID, "%"+query+"%", 30)
}
