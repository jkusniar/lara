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
	"log"

	"github.com/jkusniar/lara"
	"github.com/pkg/errors"
)

// Password is interface for passwords manipulation
type Password interface {
	Create(pass string) ([]byte, []byte, error) // hash, salt, err
	Check(pass string, salt, hash []byte) error
}

// UserService is lara.UserService implementation backed by postgresql
type UserService struct {
	DB   *sql.DB
	Pass Password
}

// Authenticate checks password for given login and returns User structure
// Method Desn't wrap errors, that it wasn't possible to find out,
// whether password was incorrect, or login didn't exist.
// This method logs *all* errors and transforms them to unauthorizedError for http layer
func (s *UserService) Authenticate(ctx context.Context, login, password string) (*lara.User, error) {
	const q = `SELECT pass_salt, pass_hash FROM "user" WHERE login = $1`
	var hash, salt []byte

	if err := s.DB.QueryRowContext(ctx, q, login).Scan(&salt, &hash); err != nil {
		log.Printf("ERROR: Authenticate: %+v\n", err)
		return nil, unauthorizedError
	}

	if err := s.Pass.Check(password, salt, hash); err != nil {
		log.Printf("ERROR: Authenticate: Pass.Check: %+v\n", err)
		return nil, unauthorizedError
	}

	// load permissions
	const pq = `SELECT p.name
			FROM permission p
			  JOIN user_permission up ON up.permission_id = p.id
			  JOIN "user" u ON u.id = up.user_id
			WHERE u.login = $1 ORDER BY p.name`
	rows, err := s.DB.QueryContext(ctx, pq, login)
	if err != nil {
		log.Printf("ERROR: Authenticate: Get permissions query: %+v\n", err)
		return nil, unauthorizedError
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			log.Printf("ERROR: Authenticate: Scan permission name: %+v\n", err)
			return nil, unauthorizedError
		}
		perms = append(perms, p)
	}
	if err := rows.Err(); err != nil {
		log.Printf("ERROR: Authenticate: Rows processing: %+v\n", err)
		return nil, unauthorizedError
	}

	u, err := lara.MakeUser(login, perms)
	if err != nil {
		log.Printf("ERROR: Authenticate: Error creating user: %+v\n", err)
		return nil, unauthorizedError
	}
	return u, nil
}

// Register creates new user in database. If user already exists, error is
// returned.
func (s *UserService) Register(ctx context.Context, login, password string, permissions []lara.PermissionType) error {
	if len(login) == 0 {
		return requiredFieldError("login")
	}

	if len(login) > 20 {
		return lara.NewCodedError(400, errors.New("login too long (max. 20 chars)"))
	}

	if len(password) == 0 {
		return requiredFieldError("password")
	}

	err := execInTransaction(ctx, s.DB, func(tx *sql.Tx) error {
		var cnt int
		if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM "user" WHERE login = $1`, login).
			Scan(&cnt); err != nil {
			return errors.Wrap(err, "register count by login error")
		}

		if cnt != 0 {
			return lara.NewCodedError(400,
				errors.Errorf("user with login %s already exists", login))
		}

		hash, salt, err := s.Pass.Create(password)
		if err != nil {
			return err
		}

		const q = `INSERT INTO "user"(login, pass_salt, pass_hash) VALUES($1, $2, $3) RETURNING id`
		var uid uint64
		if err := tx.QueryRowContext(ctx, q, login, salt, hash).Scan(&uid); err != nil {
			return errors.Wrap(err, "register user failed")
		}

		// add permissions
		for _, p := range permissions {
			if err := grantUserPermission(ctx, tx, uid, p); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func grantUserPermission(ctx context.Context, tx *sql.Tx, uid uint64, perm lara.PermissionType) error {
	var pid uint64
	err := tx.QueryRowContext(ctx, `SELECT id FROM permission WHERE name = $1`, perm.String()).Scan(&pid)
	switch {
	case err == sql.ErrNoRows:
		if err := tx.QueryRowContext(ctx, `INSERT INTO permission (name) VALUES($1) RETURNING id`,
			perm.String()).Scan(&pid); err != nil {
			return errors.Wrap(err, "create new permission failed")
		}
	case err != nil:
		return errors.Wrap(err, "load permission query error")
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO user_permission (user_id, permission_id) VALUES ($1,$2)`,
		uid, pid)

	return errors.Wrap(err, "create user_permission failed")
}

func getUserIdByLogin(ctx context.Context, tx *sql.Tx, login string) (uint64, error) {
	var uid uint64
	err := tx.QueryRowContext(ctx, `SELECT ID FROM "user" WHERE login = $1`, login).Scan(&uid)

	return uid, errors.Wrap(err, "get user id by login error")
}

func (s *UserService) Grant(ctx context.Context, login string, permissions []lara.PermissionType) error {
	err := execInTransaction(ctx, s.DB, func(tx *sql.Tx) error {
		uid, err := getUserIdByLogin(ctx, tx, login)
		if err != nil {
			return err
		}

		for _, p := range permissions {
			if err := grantUserPermission(ctx, tx, uid, p); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
func (s *UserService) Revoke(ctx context.Context, login string, permissions []lara.PermissionType) error {
	err := execInTransaction(ctx, s.DB, func(tx *sql.Tx) error {
		uid, err := getUserIdByLogin(ctx, tx, login)
		if err != nil {
			return err
		}

		const dq = `DELETE FROM user_permission
				WHERE user_id = $1 AND permission_id = (SELECT id
									FROM permission
									WHERE name = $2)`
		for _, p := range permissions {
			_, err = tx.ExecContext(ctx, dq, uid, p.String())
			if err != nil {
				return errors.Wrapf(err, "delete user_permission %d, %s failed",
					uid, p.String())
			}
		}

		return nil
	})

	return err
}
