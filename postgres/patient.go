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
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// PatientService is implementation of lara.PatientService using postgresql database
type PatientService struct {
	DB *sql.DB
}

type patientDTO struct {
	versionedDTO
	creatorDTO
	modifierDTO
	Name      string
	BirthDate pq.NullTime
	SpeciesID sql.NullInt64
	BreedID   sql.NullInt64
	GenderID  sql.NullInt64
	Note      sql.NullString
	Dead      bool
	Species   sql.NullString
	Breed     sql.NullString
	Gender    sql.NullString
}

func (p *patientDTO) toGetPatient(records []patientsRecordDTO, tags []lara.PatientsTag) *lara.GetPatient {
	result := lara.GetPatient{
		Versioned: lara.Versioned{
			ID:      p.ID,
			Version: p.Version},
		CreatorModifier: lara.CreatorModifier{
			Creator:  p.Creator,
			Created:  p.Created,
			Modifier: p.Modifier.String,
			Modified: p.Modified.Time},
		Patient: lara.Patient{
			Name:      p.Name,
			BirthDate: p.BirthDate.Time,
			SpeciesID: uint64(p.SpeciesID.Int64),
			BreedID:   uint64(p.BreedID.Int64),
			GenderID:  uint64(p.GenderID.Int64),
			Note:      p.Note.String,
		},
		Dead:    p.Dead,
		Species: p.Species.String,
		Breed:   p.Breed.String,
		Gender:  p.Gender.String,
		Records: []lara.PatientsRecord{},
		Tags:    tags}

	for _, r := range records {
		result.Records = append(result.Records, *r.toPatientsRecord())
	}

	return &result
}

type patientsRecordDTO struct {
	ID   uint64
	Date time.Time
	Text sql.NullString
}

func (p *patientsRecordDTO) toPatientsRecord() *lara.PatientsRecord {
	result := lara.PatientsRecord{ID: p.ID, Date: p.Date}

	if len(p.Text.String) > 30 {
		result.Text = p.Text.String[:30]
	} else {
		result.Text = p.Text.String
	}

	return &result
}

// Get is implementation of PatientService.Get using postgresql database.
func (s *PatientService) Get(ctx context.Context, id uint64) (*lara.GetPatient, error) {
	const q = `SELECT
			  p.id,
			  p.name,
			  p.birth_date,
			  p.species_id,
			  p.breed_id,
			  p.gender_id,
			  p.note,
			  p.dead,
			  p.version,
			  p.creator,
			  p.created,
			  p.modifier,
			  p.modified,
			  g.name as gender,
			  s.name as species,
			  b.name as breed
			FROM patient p
			 LEFT JOIN lov_gender g ON g.id = p.gender_id
			 LEFT JOIN lov_species s ON s.id = p.species_id
			 LEFT JOIN lov_breed b ON b.id = p.breed_id
			WHERE p.id = $1`

	var p patientDTO
	err := s.DB.QueryRowContext(ctx, q, id).Scan(
		&p.ID,
		&p.Name,
		&p.BirthDate,
		&p.SpeciesID,
		&p.BreedID,
		&p.GenderID,
		&p.Note,
		&p.Dead,
		&p.Version,
		&p.Creator,
		&p.Created,
		&p.Modifier,
		&p.Modified,
		&p.Gender,
		&p.Species,
		&p.Breed)
	switch {
	case err == sql.ErrNoRows:
		return nil, notFoundByIDError(id)
	case err != nil:
		return nil, errors.Wrap(err, "get patient by id failed")
	}

	// load records
	records, err := s.getPatientsRecords(ctx, p.ID)
	if err != nil {
		return nil, err
	}

	// load tags
	tags, err := s.getPatientsTags(ctx, p.ID)
	if err != nil {
		return nil, err
	}

	return p.toGetPatient(records, tags), nil
}

func (s *PatientService) getPatientsRecords(ctx context.Context, id uint64) ([]patientsRecordDTO, error) {
	const q = `SELECT id,
			  rec_date,
			  data
			FROM record
			WHERE patient_id = $1
			ORDER BY rec_date DESC`

	rows, err := s.DB.QueryContext(ctx, q, id)
	if err != nil {
		return nil, errors.Wrap(err, "get patient's records query error")
	}
	defer rows.Close()

	records := []patientsRecordDTO{}
	for rows.Next() {
		var r patientsRecordDTO
		if err := rows.Scan(&r.ID,
			&r.Date,
			&r.Text); err != nil {
			return nil, errors.Wrap(err, "scan DTO error")
		}
		records = append(records, r)
	}
	err = rows.Err()

	return records, errors.Wrap(err, "rows processing errror")
}

func (s *PatientService) getPatientsTags(ctx context.Context, id uint64) ([]lara.PatientsTag, error) {
	const q = `SELECT t.id,
			  tt.name,
			  t.value
			FROM tag t
			JOIN tag_type tt on tt.id = t.tag_type_id
			WHERE patient_id = $1
			ORDER BY tt.name, t.created`

	rows, err := s.DB.QueryContext(ctx, q, id)
	if err != nil {
		return nil, errors.Wrap(err, "get patient's tags query error")
	}
	defer rows.Close()

	tags := []lara.PatientsTag{}
	for rows.Next() {
		var t lara.PatientsTag
		if err := rows.Scan(&t.ID,
			&t.Type,
			&t.Value); err != nil {
			return nil, errors.Wrap(err, "scan DTO error")
		}
		tags = append(tags, t)
	}
	err = rows.Err()

	return tags, errors.Wrap(err, "rows processing errror")
}

// Create is implementation of PatientService.Create using postgresql database.
func (s *PatientService) Create(ctx context.Context, p *lara.CreatePatient) (uint64, error) {
	var pID uint64
	err := execInTransaction(ctx, s.DB, func(tx *sql.Tx) (err error) {
		pID, err = createPatientTx(ctx, tx, p)
		return
	})

	return pID, err
}

func createPatientTx(ctx context.Context, tx *sql.Tx, p *lara.CreatePatient) (uint64, error) {
	if p.OwnerID == 0 {
		return 0, requiredFieldError("ownerId")
	}

	if len(p.Name) == 0 {
		return 0, requiredFieldError("name")
	}

	const insert = `INSERT INTO patient (owner_id, name, birth_date, species_id, breed_id, gender_id, note, dead, creator, created)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	u, ok := lara.UserFromContext(ctx)
	if !ok {
		return 0, errors.New("no user in context")
	}

	var pID uint64
	err := tx.QueryRowContext(ctx, insert,
		toNullFK(p.OwnerID),
		toNullString(p.Name),
		toNullTime(p.BirthDate),
		toNullFK(p.SpeciesID),
		toNullFK(p.BreedID),
		toNullFK(p.GenderID),
		toNullString(p.Note),
		false,
		toNullString(u.Login),
		now()).Scan(&pID)
	if err != nil {
		return 0, errors.Wrap(err, "create patient failed")
	}

	_, err = createRecordTx(ctx, tx, &lara.CreateRecord{PatientID: pID, NewRecord: p.Record})

	return pID, err
}

// Update is implementation of PatientService.Update using postgresql database.
func (s *PatientService) Update(ctx context.Context, id uint64, p *lara.UpdatePatient) error {
	if len(p.Name) == 0 {
		return requiredFieldError("name")
	}

	err := execInTransaction(ctx, s.DB, func(tx *sql.Tx) error {
		const lck = `SELECT id FROM patient WHERE id = $1 FOR UPDATE`
		const upd = `UPDATE patient
				SET name     = $1,
				  birth_date = $2,
				  species_id = $3,
				  breed_id   = $4,
				  gender_id  = $5,
				  note       = $6,
				  dead       = $7,
				  modifier   = $8,
				  modified   = $9,
				  version    = version + 1
				WHERE id = $10 AND version = $11`

		var pid uint64
		err := tx.QueryRowContext(ctx, lck, id).Scan(&pid)
		switch err {
		case nil: // continue
		case sql.ErrNoRows:
			return notFoundByIDError(id)
		default:
			return errors.Wrap(err, "error selecting patient by id")
		}

		u, ok := lara.UserFromContext(ctx)
		if !ok {
			return errors.New("no user in context")
		}

		r, err := tx.ExecContext(ctx, upd,
			toNullString(p.Name),
			toNullTime(p.BirthDate),
			toNullFK(p.SpeciesID),
			toNullFK(p.BreedID),
			toNullFK(p.GenderID),
			toNullString(p.Note),
			p.Dead,
			toNullString(u.Login),
			now(),
			pid,
			p.Version)
		if err != nil {
			return errors.Wrap(err, "update patient failed")
		}

		count, err := r.RowsAffected()
		if err != nil {
			return errors.Wrap(err, "update patient can't check updated rows")
		}

		if count != 1 {
			return versionMismatchError(id)
		}

		return nil
	})

	return err
}
