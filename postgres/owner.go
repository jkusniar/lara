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

// OwnerService is implementation of lara.OwnerService using postgresql database
type OwnerService struct {
	DB *sql.DB
}

type ownerDTO struct {
	versionedDTO
	creatorDTO
	modifierDTO
	FirstName sql.NullString
	LastName  string
	TitleID   sql.NullInt64
	CityID    sql.NullInt64
	StreetID  sql.NullInt64
	HouseNo   sql.NullString
	Phone1    sql.NullString
	Phone2    sql.NullString
	Email     sql.NullString
	Note      sql.NullString
	IC        sql.NullString
	DIC       sql.NullString
	ICDPH     sql.NullString
	Title     sql.NullString
	City      sql.NullString
	Street    sql.NullString
}

func (o *ownerDTO) toGetOwner(patients []ownersPatientDTO) *lara.GetOwner {
	result := lara.GetOwner{
		Versioned: lara.Versioned{
			ID:      o.ID,
			Version: o.Version},
		Owner: lara.Owner{
			FirstName: o.FirstName.String,
			LastName:  o.LastName,
			TitleID:   uint64(o.TitleID.Int64),
			CityID:    uint64(o.CityID.Int64),
			StreetID:  uint64(o.StreetID.Int64),
			HouseNo:   o.HouseNo.String,
			Phone1:    o.Phone1.String,
			Phone2:    o.Phone2.String,
			Email:     o.Email.String,
			Note:      o.Note.String,
			IC:        o.IC.String,
			DIC:       o.DIC.String,
			ICDPH:     o.ICDPH.String},
		CreatorModifier: lara.CreatorModifier{
			Creator:  o.Creator,
			Created:  o.Created,
			Modifier: o.Modifier.String,
			Modified: o.Modified.Time,
		},
		Title:    o.Title.String,
		City:     o.City.String,
		Street:   o.Street.String,
		Patients: []lara.OwnersPatient{}}

	for _, p := range patients {
		result.Patients = append(result.Patients, *p.toOwnersPatient())
	}

	return &result
}

type ownersPatientDTO struct {
	ID      uint64
	Name    string
	Species sql.NullString
	Breed   sql.NullString
	Sex     sql.NullString
	Dead    bool
}

func (p *ownersPatientDTO) toOwnersPatient() *lara.OwnersPatient {
	return &lara.OwnersPatient{
		ID:      p.ID,
		Name:    p.Name,
		Dead:    p.Dead,
		Gender:  p.Sex.String,
		Species: p.Species.String,
		Breed:   p.Breed.String,
	}
}

// Get is  implementation of OwnerService.Get using postgresql database.
// Returns GetOwner data by ID.
func (s *OwnerService) Get(ctx context.Context, id uint64) (*lara.GetOwner, error) {
	const q = `SELECT
			  o.id,
			  o.first_name,
			  o.last_name,
			  o.title_id,
			  o.city_id,
			  o.street_id,
			  o.house_no,
			  o.phone_1,
			  o.phone_2,
			  o.email,
			  o.note,
			  o.ic,
			  o.dic,
			  o.icdph,
			  o.version,
			  o.creator,
			  o.created,
			  o.modifier,
			  o.modified,
			  t.name AS title,
			  c.city,
			  s.street
			FROM owner o
			 LEFT JOIN lov_title t ON t.id = o.title_id
			 LEFT JOIN lov_city c ON c.id = o.city_id
			 LEFT JOIN lov_street s ON s.id = o.street_id
			WHERE o.id = $1`

	var o ownerDTO
	err := s.DB.QueryRowContext(ctx, q, id).Scan(
		&o.ID,
		&o.FirstName,
		&o.LastName,
		&o.TitleID,
		&o.CityID,
		&o.StreetID,
		&o.HouseNo,
		&o.Phone1,
		&o.Phone2,
		&o.Email,
		&o.Note,
		&o.IC,
		&o.DIC,
		&o.ICDPH,
		&o.Version,
		&o.Creator,
		&o.Created,
		&o.Modifier,
		&o.Modified,
		&o.Title,
		&o.City,
		&o.Street)
	switch {
	case err == sql.ErrNoRows:
		return nil, notFoundByIDError(id)
	case err != nil:
		return nil, errors.Wrap(err, "get owner by id failed")
	}

	patients, err := s.getOwnersPatients(ctx, o.ID)
	if err != nil {
		return nil, err
	}

	return o.toGetOwner(patients), nil
}

func (s *OwnerService) getOwnersPatients(ctx context.Context, oid uint64) ([]ownersPatientDTO, error) {
	const q = `SELECT
			  p.id,
			  p.name,
			  s.name  AS species,
			  b.name  AS breed,
			  g.name AS gender,
			  p.dead
			FROM patient p
			  LEFT JOIN lov_species s ON s.id = p.species_id
			  LEFT JOIN lov_breed b ON b.id = p.breed_id
			  LEFT JOIN lov_gender g ON g.id = p.gender_id
			WHERE p.owner_id = $1
			ORDER BY dead, p.name`

	rows, err := s.DB.QueryContext(ctx, q, oid)
	if err != nil {
		return nil, errors.Wrap(err, "get owners' patients query error")
	}
	defer rows.Close()

	patients := []ownersPatientDTO{}
	for rows.Next() {
		var p ownersPatientDTO
		if err := rows.Scan(&p.ID,
			&p.Name,
			&p.Species,
			&p.Breed,
			&p.Sex,
			&p.Dead); err != nil {
			return nil, errors.Wrap(err, "scan DTO error")
		}
		patients = append(patients, p)
	}
	err = rows.Err()

	return patients, errors.Wrap(err, "rows processing errror")
}

// Update is  implementation of OwnerService.Update using postgresql database.
func (s *OwnerService) Update(ctx context.Context, id uint64, o *lara.UpdateOwner) error {
	if len(o.LastName) == 0 {
		return requiredFieldError("lastName")
	}

	err := execInTransaction(ctx, s.DB, func(tx *sql.Tx) error {
		const lck = `SELECT id FROM owner WHERE id = $1 FOR UPDATE`
		const upd = `UPDATE owner
			SET first_name = $1,
			  last_name    = $2,
			  title_id     = $3,
			  city_id      = $4,
			  street_id    = $5,
			  house_no     = $6,
			  phone_1      = $7,
			  phone_2      = $8,
			  email        = $9,
			  note         = $10,
			  ic           = $11,
			  dic          = $12,
			  icdph        = $13,
			  modifier     = $14,
			  modified     = $15,
			  version      = version + 1
			WHERE id = $16 AND version = $17`

		var oid uint64
		err := tx.QueryRowContext(ctx, lck, id).Scan(&oid)
		switch err {
		case nil: // continue
		case sql.ErrNoRows:
			return notFoundByIDError(id)
		default:
			return errors.Wrap(err, "error selecting owner by id")
		}

		u, ok := lara.UserFromContext(ctx)
		if !ok {
			return errors.New("no user in context")
		}

		r, err := tx.ExecContext(ctx, upd,
			toNullString(o.FirstName),
			toNullString(o.LastName),
			toNullFK(o.TitleID),
			toNullFK(o.CityID),
			toNullFK(o.StreetID),
			toNullString(o.HouseNo),
			toNullString(o.Phone1),
			toNullString(o.Phone2),
			toNullString(o.Email),
			toNullString(o.Note),
			toNullString(o.IC),
			toNullString(o.DIC),
			toNullString(o.ICDPH),
			toNullString(u.Login),
			now(),
			id,
			o.Version)
		if err != nil {
			return errors.Wrap(err, "update owner failed")
		}

		count, err := r.RowsAffected()
		if err != nil {
			return errors.Wrap(err, "update owner can't check updated rows")
		}

		if count != 1 {
			return versionMismatchError(id)
		}

		return nil
	})

	return err
}

// Create is  implementation of OwnerService.Create using postgresql database.
// Creates new owner in database and returns its ID.
func (s *OwnerService) Create(ctx context.Context, o *lara.CreateOwner) (uint64, error) {
	const insert = `INSERT INTO owner (first_name, last_name, title_id, city_id, street_id, house_no,
		                    phone_1, phone_2, email, note, ic, dic, icdph, creator, created)
					VALUES
					  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
					RETURNING id`

	if len(o.LastName) == 0 {
		return 0, requiredFieldError("lastName")
	}

	var id uint64
	err := execInTransaction(ctx, s.DB, func(tx *sql.Tx) error {
		u, ok := lara.UserFromContext(ctx)
		if !ok {
			return errors.New("no user in context")
		}

		err := tx.QueryRowContext(ctx, insert,
			toNullString(o.FirstName),
			toNullString(o.LastName),
			toNullFK(o.TitleID),
			toNullFK(o.CityID),
			toNullFK(o.StreetID),
			toNullString(o.HouseNo),
			toNullString(o.Phone1),
			toNullString(o.Phone2),
			toNullString(o.Email),
			toNullString(o.Note),
			toNullString(o.IC),
			toNullString(o.DIC),
			toNullString(o.ICDPH),
			toNullString(u.Login),
			now()).Scan(&id)
		if err != nil {
			return errors.Wrap(err, "create owner failed")
		}

		_, err = createPatientTx(ctx, tx, &lara.CreatePatient{OwnerID: id, NewPatient: o.Patient})

		return err
	})

	return id, err
}
