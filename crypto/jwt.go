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
	"crypto/rsa"
	"io/ioutil"
	"strings"
	"time"

	"github.com/jkusniar/lara"
	"github.com/pkg/errors"
	jwt "gopkg.in/dgrijalva/jwt-go.v3"
)

// JWTToken is authentication token provider/validator using JWT
type JWTToken struct {
	pkPath string         // openssl genrsa -out lara.rsa 2048
	pubKey *rsa.PublicKey // openssl rsa -in lara.rsa -pubout > lara.rsa.pub
	issuer string         // issuer string for generated tokens
}

// NewJWTToken creates new JWTToken instance initialized from keys on disk
func NewJWTToken(privateKeyPath, pubKeyPath, issuer string) (*JWTToken, error) {
	keyBytes, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "error reading public RSA key")
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(keyBytes)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing public RSA key")
	}

	return &JWTToken{privateKeyPath, pubKey, issuer}, nil
}

type laraUserClaims struct {
	Permissions string `json:"permissions"`
	jwt.StandardClaims
}

// Create creates authentication token for given user
func (a *JWTToken) Create(u *lara.User) (string, error) {
	// load private key from disk every time new token is created
	keyBytes, err := ioutil.ReadFile(a.pkPath)
	if err != nil {
		return "", errors.Wrap(err, "error reading private key")
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if err != nil {
		return "", errors.Wrap(err, "error parsing private key")
	}

	// convert u.Permissions to slice
	perms := make([]string, 0, len(u.Permissions))
	for p := range u.Permissions {
		perms = append(perms, p.String())
	}

	claims := &laraUserClaims{
		strings.Join(perms, ","),
		jwt.StandardClaims{
			Subject:   u.Login,
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
			Issuer:    a.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	k, err := token.SignedString(key)
	if err != nil {
		return "", errors.Wrap(err, "error signing token")
	}
	return k, nil
}

// Parse validates and parses given authentication token. If valid, lara.User
// object encoded in token is returned
func (a *JWTToken) Parse(tokenString string) (*lara.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &laraUserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// required because of JWT spec vulnerability
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.Errorf("unexpected signing method: %v",
					token.Header["alg"])
			}
			return a.pubKey, nil
		})

	if err != nil {
		return nil, errors.Wrap(err, "error parsing JWT")
	}

	claims, ok := token.Claims.(*laraUserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid authorization token")
	}

	perms := []string{}
	if len(claims.Permissions) > 0 {
		perms = strings.Split(claims.Permissions, ",")
	}

	u, err := lara.MakeUser(claims.Subject, perms)
	return u, errors.Wrap(err, "error creating user from JWT")
}
