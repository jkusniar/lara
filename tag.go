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

package lara

import (
	"context"
	"fmt"
)

// -----------------------------------------------------------------------------
// PATIENT'S TAG MANAGEMENT SERVICE

// TagType is enum defining system permissions
//go:generate stringer -type=TagType -output tag_string.go
//requires golang.org/x/tools/cmd/stringer installed locally
//if new object permission added to enum, run "go generate"
type TagType int

// TagType enum
const (
	LyssaVirus  TagType = iota // Canine Rabies Tags, tag format: YYYY-SK-9999
	Tattoo                     // Pet tattoo, tag format: regular string
	PetPassport                // EU Pet Passport, tag format: SK999999999
	RFID                       // ISO 11784/11785 FDX-B RFID chip, tag format: 999999999999999
)

// FromString creates TagType from string
func (i *TagType) FromString(tt string) error {
	// TODO: more intelligent implementation
	switch tt {
	case "LyssaVirus":
		*i = LyssaVirus
	case "Tattoo":
		*i = Tattoo
	case "PetPassport":
		*i = PetPassport
	case "RFID":
		*i = RFID
	default:
		return fmt.Errorf("bad TagType: '%s'", tt)
	}

	return nil
}

type GetTag struct {
	Versioned
	CreatorModifier
	Type  string `json:"type"`
	Value string `json:"value"`
	Data  []byte `json:"data"`
}

type CreateTag struct {
	PatientID uint64 `json:"patientId"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	Data      []byte `json:"data"`
}

type UpdateTag struct {
	Version uint64 `json:"version"`
	Value   string `json:"value"`
	Data    []byte `json:"data"`
}

type PatientByTag struct {
	TagType      string `json:"tagType"`
	Name         string `json:"name"`
	Species      string `json:"species"`
	Breed        string `json:"breed"`
	Gender       string `json:"gender"`
	OwnerID      uint64 `json:"ownerId"`      // DB primary key
	OwnerName    string `json:"ownerName"`    // formatted owner name (first, last, title)
	OwnerAddress string `json:"ownerAddress"` // owner's address
}

type TagService interface {
	Get(ctx context.Context, id uint64) (*GetTag, error)
	Update(ctx context.Context, id uint64, o *UpdateTag) error
	Create(ctx context.Context, o *CreateTag) (uint64, error)
	GetPatientByTag(ctx context.Context, tagValue string) (*PatientByTag, error)
}
