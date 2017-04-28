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
	"time"

	"github.com/jkusniar/lara"
	"github.com/pkg/errors"
)

type RecordService struct {
	DB *sql.DB
}

type recordDTO struct {
	versionedDTO
	creatorDTO
	modifierDTO
	Date   time.Time
	Text   sql.NullString
	Billed bool
	Total  string
}

func (r *recordDTO) toGetRecord(total string, items []lara.GetRecordItem) *lara.GetRecord {
	return &lara.GetRecord{
		Versioned: lara.Versioned{
			ID:      r.ID,
			Version: r.Version},
		CreatorModifier: lara.CreatorModifier{
			Creator:  r.Creator,
			Created:  r.Created,
			Modifier: r.Modifier.String,
			Modified: r.Modified.Time},
		Date:   r.Date,
		Text:   r.Text.String,
		Billed: r.Billed,
		Total:  total,
		Items:  items,
	}
}

func (s *RecordService) Get(ctx context.Context, id uint64) (*lara.GetRecord, error) {
	const q = `SELECT
			  id,
			  rec_date,
			  data,
			  billed,
			  version,
			  creator,
			  created,
			  modifier,
			  modified
			FROM record WHERE id = $1`

	var r recordDTO
	err := s.DB.QueryRowContext(ctx, q, id).Scan(
		&r.ID,
		&r.Date,
		&r.Text,
		&r.Billed,
		&r.Version,
		&r.Creator,
		&r.Created,
		&r.Modifier,
		&r.Modified)
	switch {
	case err == sql.ErrNoRows:
		return nil, notFoundByIDError(id)
	case err != nil:
		return nil, errors.Wrap(err, "get record by id failed")
	}

	items, err := s.getRecordItems(ctx, r.ID)
	if err != nil {
		return nil, err
	}

	sum, err := s.sumItemsForRecord(ctx, r.ID)
	if err != nil {
		return nil, err
	}

	return r.toGetRecord(sum, items), nil
}

func (s *RecordService) getRecordItems(ctx context.Context, id uint64) ([]lara.GetRecordItem, error) {
	const q = `SELECT ri.id,
			  ri.prod_id,
			  ri.prod_price,
			  ri.amount,
			  ri.item_price,
			  ri.item_type,
			  p.name as product,
			  u.name as unit,
			  p.plu as plu
			FROM record_item ri
			JOIN lov_product p ON p.id = ri.prod_id
			JOIN lov_unit u ON u.id = p.unit_id
			WHERE ri.record_id = $1
			ORDER BY ri.id`

	rows, err := s.DB.QueryContext(ctx, q, id)
	if err != nil {
		return nil, errors.Wrap(err, "get record's items query error")
	}
	defer rows.Close()

	items := []lara.GetRecordItem{}
	for rows.Next() {
		var i lara.GetRecordItem
		var plu sql.NullString
		if err := rows.Scan(&i.ID,
			&i.ProductID,
			&i.ProductPrice,
			&i.Amount,
			&i.ItemPrice,
			&i.ItemType,
			&i.Product,
			&i.Unit,
			&plu); err != nil {
			return nil, errors.Wrap(err, "scan DTO error")
		}
		i.PLU = plu.String
		items = append(items, i)
	}
	err = rows.Err()

	return items, errors.Wrap(err, "rows processing errror")
}

func (s *RecordService) sumItemsForRecord(ctx context.Context, id uint64) (string, error) {
	const sq = `SELECT SUM(ri.item_price)
		FROM record r INNER JOIN record_item ri ON ri.record_id = r.id
		WHERE r.id = $1`

	var sum sql.NullString
	err := s.DB.QueryRowContext(ctx, sq, id).Scan(&sum)
	switch {
	case err == sql.ErrNoRows:
		return "", notFoundByIDError(id)
	case err != nil:
		return "", errors.Wrap(err, "get record items sum failed")
	}

	if !sum.Valid {
		return "0.00", nil
	}

	return sum.String, nil
}

func (s *RecordService) Create(ctx context.Context, r *lara.CreateRecord) (uint64, error) {
	var recID uint64
	err := execInTransaction(ctx, s.DB, func(tx *sql.Tx) (err error) {
		recID, err = createRecordTx(ctx, tx, r)
		return
	})

	return recID, err
}

func createRecordTx(ctx context.Context, tx *sql.Tx, r *lara.CreateRecord) (uint64, error) {
	if r.PatientID == 0 {
		return 0, requiredFieldError("patientId")
	}

	if err := validateRecordItems(r.Items); err != nil {
		return 0, err
	}

	const insertRecord = `INSERT INTO record (patient_id, rec_date, data, billed, creator, created)
					VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	var recID uint64

	u, ok := lara.UserFromContext(ctx)
	if !ok {
		return 0, errors.New("no user in context")
	}

	err := tx.QueryRowContext(ctx, insertRecord,
		toNullFK(r.PatientID),
		now(),
		toNullString(r.Text),
		r.Billed,
		toNullString(u.Login),
		now()).Scan(&recID)
	if err != nil {
		return 0, errors.Wrap(err, "create record failed")
	}

	err = createRecordItems(ctx, tx, recID, r.Items)

	return recID, err
}

func createRecordItems(ctx context.Context, tx *sql.Tx, recID uint64, items []lara.RecordItem) error {
	const insert = `INSERT INTO record_item (record_id, prod_id, amount, item_price, prod_price, item_type)
					VALUES ($1, $2, $3, $4, $5, $6)`

	for _, i := range items {
		_, err := tx.ExecContext(ctx, insert,
			recID,
			toNullFK(i.ProductID),
			toNullString(i.Amount),
			toNullString(i.ItemPrice),
			toNullString(i.ProductPrice),
			i.ItemType,
		)
		if err != nil {
			return errors.Wrap(err, "insert record item failed")
		}
	}

	return nil
}

func (s *RecordService) Update(ctx context.Context, id uint64, r *lara.UpdateRecord) error {
	if err := validateRecordItems(r.Items); err != nil {
		return err
	}

	err := execInTransaction(ctx, s.DB, func(tx *sql.Tx) error {
		const lock = `SELECT id FROM record WHERE id = $1 FOR UPDATE`
		const update = `UPDATE record
				SET data   = $1,
				  modifier = $2,
				  modified = $3,
				  version  = version + 1
				WHERE id = $4 AND version = $5`
		const del = `DELETE FROM record_item WHERE record_id = $1`

		var rid uint64
		err := tx.QueryRowContext(ctx, lock, id).Scan(&rid)
		switch err {
		case nil: // continue
		case sql.ErrNoRows:
			return notFoundByIDError(id)
		default:
			return errors.Wrap(err, "error selecting record by id")
		}

		u, ok := lara.UserFromContext(ctx)
		if !ok {
			return errors.New("no user in context")
		}

		res, err := tx.ExecContext(ctx, update,
			toNullString(r.Text),
			toNullString(u.Login),
			now(),
			rid,
			r.Version)
		if err != nil {
			return errors.Wrap(err, "update record failed")
		}

		count, err := res.RowsAffected()
		if err != nil {
			return errors.Wrap(err, "update record can't check updated rows")
		}

		if count != 1 {
			return versionMismatchError(id)
		}

		// del & create items
		if _, err := tx.ExecContext(ctx, del, rid); err != nil {
			return err
		}

		return createRecordItems(ctx, tx, rid, r.Items)
	})

	return err
}

func validateRecordItems(items []lara.RecordItem) error {
	for i, itm := range items {
		if itm.ProductID == 0 {
			return requiredFieldError(fmt.Sprintf("productId on item %d", i))
		}
		if len(itm.ProductPrice) == 0 {
			return requiredFieldError(fmt.Sprintf("productPrice on item %d", i))
		}
		if len(itm.Amount) == 0 {
			return requiredFieldError(fmt.Sprintf("amount on item %d", i))
		}
		if len(itm.ItemPrice) == 0 {
			return requiredFieldError(fmt.Sprintf("itemPrice on item %d", i))
		}
	}
	return nil
}
