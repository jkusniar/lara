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

package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"io"

	"github.com/pkg/errors"
	"golang.org/x/crypto/scrypt"
)

const (
	defaultSaltBytes = 32
	defaultHashBytes = 64
	scryptN          = 16384
	scryptR          = 8
	scryptP          = 1
)

// Password contains methods for manipulating with passwords.
type Password struct {
	// SaltBytes is salt length in bytes
	SaltBytes int
	// HashBytes is password hash length in bytes
	HashBytes int
	n, r, p   int // scrypt params
}

// NewPassword returns new Password struct with default values.
func NewPassword() *Password {
	return &Password{defaultSaltBytes, defaultHashBytes,
		scryptN, scryptR, scryptP}
}

// Check verifies whether supplied password and salt can be encrypted to supplied
// hash. If not, error is returned.
func (p *Password) Check(pass string, salt, hash []byte) error {
	computed, err := p.encrypt(pass, salt)
	if err != nil {
		return err
	}

	if subtle.ConstantTimeCompare(computed, hash) != 1 {
		return errors.New("password validation failed")
	}

	return nil
}

// Create returns password's hash and salt generated from supplied password.
func (p *Password) Create(pass string) ([]byte, []byte, error) {
	salt := make([]byte, p.SaltBytes)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, nil, errors.Wrap(err, "error generating salt")
	}

	hash, err := p.encrypt(pass, salt)
	if err != nil {
		return nil, nil, err
	}

	return hash, salt, nil
}

// generate hash from password+salt. Using scrypt key derivation function
func (p *Password) encrypt(password string, salt []byte) ([]byte, error) {
	hash, err := scrypt.Key([]byte(password), salt, p.n, p.r, p.p, p.HashBytes)
	if err != nil {
		return nil, errors.Wrap(err, "error computing password hash")
	}

	return hash, nil
}
