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

package main

import (
	"context"
	"flag"
	"strings"

	"github.com/jkusniar/lara"
	"github.com/jkusniar/lara/crypto"
	"github.com/jkusniar/lara/postgres"
	"github.com/pkg/errors"
)

func register(user, pass, host, name string, port uint, sslMode string, args []string) error {
	if len(args) != 3 {
		flag.Usage()
	}

	db, err := postgres.Open(user, pass, host, name, port, sslMode)
	if err != nil {
		return err
	}
	defer db.Close()

	service := &postgres.UserService{DB: db, Pass: crypto.NewPassword()}
	return service.Register(context.Background(), args[1], args[2], lara.DefaultPermissions)
}

func grant(user, pass, host, name string, port uint, sslMode string, args []string) error {
	if len(args) != 3 {
		flag.Usage()
	}

	db, err := postgres.Open(user, pass, host, name, port, sslMode)
	if err != nil {
		return err
	}
	defer db.Close()

	p, err := extractPermissions(args[2])
	if err != nil {
		return err
	}

	service := &postgres.UserService{DB: db, Pass: crypto.NewPassword()}
	return service.Grant(context.Background(), args[1], p)
}

func revoke(user, pass, host, name string, port uint, sslMode string, args []string) error {
	if len(args) != 3 {
		flag.Usage()
	}

	db, err := postgres.Open(user, pass, host, name, port, sslMode)
	if err != nil {
		return err
	}
	defer db.Close()

	p, err := extractPermissions(args[2])
	if err != nil {
		return err
	}

	service := &postgres.UserService{DB: db, Pass: crypto.NewPassword()}
	return service.Revoke(context.Background(), args[1], p)
}

func extractPermissions(s string) ([]lara.PermissionType, error) {
	perms := []string{}
	if len(s) > 0 {
		perms = strings.Split(s, ",")
	}

	pm, err := lara.MakePermissions(perms)
	if err != nil {
		return nil, errors.Wrap(err, "error converting permissions")
	}

	result := make([]lara.PermissionType, 0, len(pm))
	for k, v := range pm {
		if v {
			result = append(result, k)
		}
	}

	return result, nil
}
