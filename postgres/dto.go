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
	"database/sql"
	"strings"
)

// OwnerNameDTO is common DTO for formatted Owner's name retrieval
type OwnerNameDTO struct {
	FirstName sql.NullString
	LastName  string
	Title     sql.NullString
}

func (o *OwnerNameDTO) String() string {
	var name string
	if o.FirstName.Valid {
		name = strings.Join([]string{o.FirstName.String, o.LastName}, " ")
	} else {
		name = o.LastName
	}

	if o.Title.Valid {
		name = strings.Join([]string{name, o.Title.String}, ", ")
	}

	return name
}

// OwnerAddressDTO is common DTO for formatted Owner's address retrieval
type OwnerAddressDTO struct {
	City    sql.NullString
	Street  sql.NullString
	HouseNo sql.NullString
}

func (o *OwnerAddressDTO) String() string {
	var address string
	if o.City.Valid {
		if o.Street.Valid {
			address = strings.Join([]string{o.Street.String, " ",
				o.HouseNo.String, ", ", o.City.String}, "")
		} else {
			address = strings.Join([]string{o.City.String, o.HouseNo.String},
				" ")
		}

	}
	return address
}
