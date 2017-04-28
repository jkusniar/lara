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

// ProductService is lara.ProductService implementation backed by postgresql
type ProductService struct {
	DB *sql.DB
}

// Search performs DB search according to ProductSearchRequest
// TODO make search limit parametric
func (s *ProductService) Search(ctx context.Context, p *lara.ProductSearchRequest) (*lara.ProductSearchResult, error) {
	r := lara.ProductSearchResult{Total: 0, Products: []lara.Product{}}

	const cq = `SELECT count(*) FROM lov_product WHERE name ILIKE $1 AND (valid_to IS NULL OR valid_to >= $2)`
	const dq = `SELECT
			  p.id    AS id,
			  p.name  AS name,
			  u.name  AS unit,
			  p.price AS price
			FROM lov_product p
			  JOIN lov_unit u ON u.id = p.unit_id
			WHERE p.name ILIKE $1 AND (p.valid_to IS NULL OR p.valid_to >= $2)
			ORDER BY p.name
			LIMIT $3`

	if err := s.DB.QueryRowContext(ctx, cq, "%"+p.Query+"%", p.ValidTo.Local()).Scan(&r.Total); err != nil {
		return nil, errors.Wrap(err, "search product count error")
	}

	rows, err := s.DB.QueryContext(ctx, dq, "%"+p.Query+"%", p.ValidTo.Local(), 30)
	if err != nil {
		return nil, errors.Wrap(err, "search query error")
	}
	defer rows.Close()

	for rows.Next() {
		var p lara.Product
		if err := rows.Scan(&p.ID,
			&p.Name,
			&p.Unit,
			&p.Price); err != nil {
			return nil, errors.Wrap(err, "scan DTO error")
		}
		r.Products = append(r.Products, p)
	}
	err = rows.Err()

	return &r, errors.Wrap(err, "rows processing errror")
}
