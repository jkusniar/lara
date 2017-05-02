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

// TagService is lara.TagService implementation backed by postgresql
type TagService struct {
	DB *sql.DB
}

type tagDTO struct {
	versionedDTO
	creatorDTO
	modifierDTO
	Type  string
	Value string
	Data  []byte
}

func (t *tagDTO) toGetTag() *lara.GetTag {
	return &lara.GetTag{
		Versioned: lara.Versioned{
			ID:      t.ID,
			Version: t.Version},
		CreatorModifier: lara.CreatorModifier{
			Creator:  t.Creator,
			Created:  t.Created,
			Modifier: t.Modifier.String,
			Modified: t.Modified.Time,
		},
		Type:  t.Type,
		Value: t.Value,
		Data:  t.Data,
	}
}

// Get is implementation of TagService.Get using postgresql database.
func (s *TagService) Get(ctx context.Context, id uint64) (*lara.GetTag, error) {
	const q = `SELECT
			  t.id,
			  t.version,
			  t.creator,
			  t.created,
			  t.modifier,
			  t.modified,
			  tt.name as type,
			  t.value,
			  t.data
			FROM tag t
			 JOIN tag_type tt ON tt.id = t.tag_type_id
			WHERE t.id = $1`

	var t tagDTO
	err := s.DB.QueryRowContext(ctx, q, id).Scan(
		&t.ID,
		&t.Version,
		&t.Creator,
		&t.Created,
		&t.Modifier,
		&t.Modified,
		&t.Type,
		&t.Value,
		&t.Data)
	switch {
	case err == sql.ErrNoRows:
		return nil, notFoundByIDError(id)
	case err != nil:
		return nil, errors.Wrap(err, "get tag by id failed")
	}

	return t.toGetTag(), nil
}

// Update is implementation of TagService.Update using postgresql database.
func (s *TagService) Update(ctx context.Context, id uint64, t *lara.UpdateTag) error {
	if len(t.Value) == 0 {
		return requiredFieldError("value")
	}

	err := execInTransaction(ctx, s.DB, func(tx *sql.Tx) error {
		const lck = `SELECT id FROM tag WHERE id = $1 FOR UPDATE`
		const upd = `UPDATE tag
				SET value  = $1,
				  data     = $2,
				  modifier = $3,
				  modified = $4,
				  version  = version + 1
				WHERE id = $5 AND version = $6`

		var tid uint64
		err := tx.QueryRowContext(ctx, lck, id).Scan(&tid)
		switch err {
		case nil: // continue
		case sql.ErrNoRows:
			return notFoundByIDError(id)
		default:
			return errors.Wrap(err, "error selecting tag by id")
		}

		u, ok := lara.UserFromContext(ctx)
		if !ok {
			return errors.New("no user in context")
		}

		r, err := tx.ExecContext(ctx, upd,
			toNullString(t.Value),
			t.Data,
			toNullString(u.Login),
			now(),
			id,
			t.Version)
		if err != nil {
			return errors.Wrap(err, "update tag failed")
		}

		count, err := r.RowsAffected()
		if err != nil {
			return errors.Wrap(err, "update tag can't check updated rows")
		}

		if count != 1 {
			return versionMismatchError(id)
		}

		return nil
	})

	return err
}

// Create is implementation of TagService.Create using postgresql database.
func (s *TagService) Create(ctx context.Context, t *lara.CreateTag) (uint64, error) {
	if len(t.Value) == 0 {
		return 0, requiredFieldError("value")
	}

	if len(t.Type) == 0 {
		return 0, requiredFieldError("type")
	}

	if t.PatientID == 0 {
		return 0, requiredFieldError("patientId")
	}

	var id uint64
	err := execInTransaction(ctx, s.DB, func(tx *sql.Tx) error {
		u, ok := lara.UserFromContext(ctx)
		if !ok {
			return errors.New("no user in context")
		}

		const insert = `INSERT INTO tag (patient_id, tag_type_id, value, data, creator, created)
					VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

		ttID, err := getOrCreateTagType(ctx, tx, t.Type)
		if err != nil {
			return err
		}

		err = tx.QueryRowContext(ctx, insert,
			toNullFK(t.PatientID),
			ttID,
			toNullString(t.Value),
			t.Data,
			toNullString(u.Login),
			now()).Scan(&id)

		return errors.Wrap(err, "create tag failed")
	})

	return id, err
}

func getOrCreateTagType(ctx context.Context, tx *sql.Tx, tagType string) (uint64, error) {
	var ttID uint64

	var tt lara.TagType
	if err := tt.FromString(tagType); err != nil {
		return ttID, lara.NewCodedError(400, errors.Errorf("invalid TagType %s", tagType))
	}

	err := tx.QueryRowContext(ctx, `SELECT id FROM tag_type WHERE name = $1`, tt.String()).Scan(&ttID)
	if err == sql.ErrNoRows {
		err = tx.QueryRowContext(ctx, `INSERT INTO tag_type (name) VALUES($1) RETURNING id`,
			tt.String()).Scan(&ttID)
	}

	return ttID, errors.Wrap(err, "get or create TagType failed")
}

type patientByTagDTO struct {
	TagType string
	Name    string
	Species sql.NullString
	Breed   sql.NullString
	Gender  sql.NullString
	OwnerID uint64
	OwnerNameDTO
	OwnerAddressDTO
}

func (p *patientByTagDTO) toPatientByTag() *lara.PatientByTag {
	return &lara.PatientByTag{TagType: p.TagType,
		Name: p.Name, Species: p.Species.String, Breed: p.Breed.String, Gender: p.Gender.String,
		OwnerID: p.OwnerID, OwnerName: p.OwnerNameDTO.String(), OwnerAddress: p.OwnerAddressDTO.String()}
}

// GetPatientByTag is implementation of TagService.GetPatientByTag using postgresql database.
func (s *TagService) GetPatientByTag(ctx context.Context, tagValue string) (*lara.PatientByTag, error) {
	const q = `SELECT
			  tt.name  AS tagType,
			  p.name   AS name,
			  sp.name  AS species,
			  b.name   AS breed,
			  g.name   AS gender,
			  o.id     AS ownerId,
			  o.first_name,
			  o.last_name,
			  l.name   AS title,
			  c.city   AS city,
			  s.street AS street,
			  o.house_no
			FROM tag t
			  JOIN tag_type tt ON tt.id = t.tag_type_id
			  JOIN patient p ON p.id = t.patient_id
			  JOIN owner o ON o.id = p.owner_id
			  LEFT JOIN lov_title l ON l.id = o.title_id
			  LEFT JOIN lov_city c ON c.id = o.city_id
			  LEFT JOIN lov_street s ON s.id = o.street_id
			  LEFT JOIN lov_gender g ON g.id = p.gender_id
			  LEFT JOIN lov_species sp ON sp.id = p.species_id
			  LEFT JOIN lov_breed b ON b.id = p.breed_id
			WHERE t.value = $1`

	var p patientByTagDTO
	err := s.DB.QueryRowContext(ctx, q, tagValue).Scan(
		&p.TagType, &p.Name, &p.Species, &p.Breed, &p.Gender, &p.OwnerID,
		&p.FirstName, &p.LastName, &p.Title, &p.City,
		&p.Street, &p.HouseNo)
	switch {
	case err == sql.ErrNoRows:
		return nil, lara.NewCodedError(404, errors.Errorf("no patient with tag %s", tagValue))
	case err != nil:
		return nil, errors.Wrap(err, "get patient by tag failed")
	}

	return p.toPatientByTag(), nil
}
