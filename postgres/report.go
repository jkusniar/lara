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
	"time"

	"github.com/jkusniar/lara"
	"github.com/pkg/errors"
)

// ReportService is lara.ReportService implementation backed by postgresql
type ReportService struct {
	DB  *sql.DB
	Loc *time.Location
}

// GetIncomeStatistics counts records and income for specified time period
//
// TODO: better query - all SUMs in one query - over months
// SELECT extract(year FROM r.rec_date) AS yr,
//         extract(month FROM r.rec_date) AS mon,
//         sum(case when r.billed = true then ri.item_price else 0.0 end) AS "Billed",
//         sum(case when r.billed = false then ri.item_price else 0.0 end) AS "Not Billed",
//         sum(ri.item_price) AS "Total"
//   FROM record r
//   JOIN record_item ri on ri.record_id = r.id
//   GROUP BY yr, mon ORDER by yr, mon;
//
func (s *ReportService) GetIncomeStatistics(ctx context.Context,
	r *lara.ReportRequest) (*lara.IncomeStatistics, error) {
	resp := lara.IncomeStatistics{Income: "0.00", IncomeBilled: "0.00", IncomeNotBilled: "0.00"}

	from := r.ValidFrom.In(s.Loc)
	to := r.ValidTo.In(s.Loc)

	if err := s.DB.QueryRowContext(ctx,
		`SELECT count(*) FROM record WHERE rec_date >= $1 AND rec_date <= $2`,
		from, to).Scan(&resp.Records); err != nil {
		return nil, errors.Wrap(err, "income statistics count records error")
	}

	if err := s.scanIncome(ctx, &resp.Income,
		`SELECT SUM(ri.item_price) FROM record_item ri
			INNER JOIN record r ON r.id = ri.record_id
			WHERE r.rec_date >= $1 AND r.rec_date <= $2`,
		from, to); err != nil {
		return nil, err
	}

	if err := s.scanIncome(ctx, &resp.IncomeBilled,
		`SELECT SUM(ri.item_price) FROM record_item ri
			INNER JOIN record r ON r.id = ri.record_id
			WHERE r.rec_date >= $1 AND r.rec_date <= $2
			AND r.billed = '1'`,
		from, to); err != nil {
		return nil, err
	}

	if err := s.scanIncome(ctx, &resp.IncomeNotBilled,
		`SELECT SUM(ri.item_price) FROM record_item ri
			INNER JOIN record r ON r.id = ri.record_id
			WHERE r.rec_date >= $1 AND r.rec_date <= $2
			AND r.billed = '0'`,
		from, to); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *ReportService) scanIncome(ctx context.Context, income *string, query string, from, to time.Time) error {
	var sum sql.NullString
	if err := s.DB.QueryRowContext(ctx, query, from, to).Scan(&sum); err != nil {
		return errors.Wrap(err, "income statistics SUM select error")
	}

	if sum.Valid {
		*income = sum.String
	}

	return nil
}
