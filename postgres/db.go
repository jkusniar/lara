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
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// Open creates new connection pool to database.
// Clients of this call should defer Close on returned object.
func Open(user, pass, host, name string, port uint, sslMode string) (*sql.DB, error) {
	db, err := sql.Open("postgres",
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			user, pass, host, port, name, sslMode))

	if err != nil {
		return nil, errors.Wrap(err, "db connect failed")
	}

	err = db.Ping()

	return db, errors.Wrap(err, "db ping failed")
}

// execute in transaction wrapper
func execInTransaction(ctx context.Context, db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "cannot begin transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback() // rollback failed error ignored
			return
		}
		err = errors.Wrap(tx.Commit(), "transaction commit error")
	}()
	return txFunc(tx)
}

// predefined errors
func requiredFieldError(field string) error {
	return lara.NewCodedError(400,
		errors.Errorf("%s is required", field))
}

func notFoundByIDError(id uint64) error {
	return lara.NewCodedError(404,
		errors.Errorf("object with id %d not found", id))
}

func versionMismatchError(id uint64) error {
	return lara.NewCodedError(409,
		errors.Errorf("object with id %d modified by another user. Reload and edit again.", id))
}

var unauthorizedError = lara.NewCodedError(401,
	errors.New("username or password invalid"))

// fields for reuse in DTOs
type versionedDTO struct {
	ID      uint64
	Version uint64
}

type creatorDTO struct {
	Creator string
	Created time.Time
}

type modifierDTO struct {
	Modifier sql.NullString
	Modified pq.NullTime
}

func toNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

// On nullable FK columns value 0 is equal to null in DB
func toNullFK(i uint64) sql.NullInt64 {
	return sql.NullInt64{Int64: int64(i), Valid: i != 0}
}

func toNullTime(t time.Time) pq.NullTime {
	return pq.NullTime{Time: t, Valid: t.String() != "0001-01-01 00:00:00 +0000 UTC"}
}

func now() pq.NullTime {
	return pq.NullTime{Time: time.Now(), Valid: true}
}
