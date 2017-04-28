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

package crypto_test

import (
	"testing"

	"github.com/jkusniar/lara/crypto"
)

var (
	pass = "TestPassword"
	salt = []byte{148, 129, 33, 2, 114, 70, 13, 125, 19, 59, 236, 188, 169,
		43, 179, 22, 164, 246, 23, 140, 240, 198, 94, 81, 237, 234, 60,
		166, 83, 14, 156, 202}
	hash = []byte{90, 82, 147, 105, 110, 26, 43, 54, 8, 242, 42, 23,
		144, 99, 55, 49, 107, 110, 203, 69, 63, 233, 194, 224, 118, 122,
		116, 89, 105, 89, 147, 107, 44, 88, 205, 54, 92, 28, 169, 36, 71,
		89, 88, 177, 96, 234, 166, 46, 152, 109, 202, 110, 218, 187, 246,
		127, 20, 163, 155, 169, 116, 91, 241, 64}
)

func TestCheckPassword(t *testing.T) {
	if crypto.NewPassword().Check(pass, salt, hash) != nil {
		t.Fatal("Failed test password validation")
	}
}

func TestCheckPasswordFailed(t *testing.T) {
	if crypto.NewPassword().Check("BadPassword", salt, hash) == nil {
		t.Fatal("Test password validation should fail")
	}
}

func TestCreatePassword(t *testing.T) {
	p := crypto.NewPassword()
	h, s, err := p.Create(pass)

	if err != nil {
		t.Fatalf("Failed create password test: %+v", err)
	}

	if len(h) != p.HashBytes {
		t.Fatalf("computed hash has incorrect length %d", len(h))
	}

	if len(s) != p.SaltBytes {
		t.Fatalf("computed salt has incorrect length %d", len(s))
	}

	// print data for TestCreatePassword
	//t.Logf("%v\n", hash)
	//t.Logf("%v\n", salt)
}
