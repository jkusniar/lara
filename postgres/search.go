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

// SearchService is lara.SearchService implementation backed by postgresql
type SearchService struct {
	DB *sql.DB
}

type searchDTO struct {
	ID uint64
	OwnerNameDTO
	OwnerAddressDTO
}

func (s *searchDTO) toRecord() *lara.SearchRecord {
	return &lara.SearchRecord{ID: s.ID, Name: s.OwnerNameDTO.String(), Address: s.OwnerAddressDTO.String()}
}

// Search performs DB search of "q" string in owner/pet records
// TODO make search limit parametric
func (s *SearchService) Search(ctx context.Context, q string) (*lara.SearchResult, error) {
	r := lara.SearchResult{Total: 0, Records: []lara.SearchRecord{}}

	const cq = `SELECT count(*) FROM owner WHERE last_name ILIKE $1`
	const dq = `SELECT
			  o.id,
			  o.first_name,
			  o.last_name,
			  t.name   AS title,
			  c.city   AS city,
			  s.street AS street,
			  o.house_no
			FROM owner o
			  LEFT JOIN lov_title t ON t.id = o.title_id
			  LEFT JOIN lov_city c ON c.id = o.city_id
			  LEFT JOIN lov_street s ON s.id = o.street_id
			WHERE o.last_name
			      ILIKE $1
			ORDER BY o.last_name
			LIMIT $2`

	if err := s.DB.QueryRowContext(ctx, cq, "%"+q+"%").Scan(&r.Total); err != nil {
		return nil, errors.Wrap(err, "search count error")
	}

	rows, err := s.DB.QueryContext(ctx, dq, "%"+q+"%", 30)
	if err != nil {
		return nil, errors.Wrap(err, "search query error")
	}
	defer rows.Close()

	for rows.Next() {
		var dto searchDTO
		if err := rows.Scan(&dto.ID,
			&dto.FirstName,
			&dto.LastName,
			&dto.Title,
			&dto.City,
			&dto.Street,
			&dto.HouseNo); err != nil {
			return nil, errors.Wrap(err, "scan DTO error")
		}
		r.Records = append(r.Records, *dto.toRecord())
	}
	err = rows.Err()

	return &r, errors.Wrap(err, "rows processing errror")
}
