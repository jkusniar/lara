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

import "fmt"

// CodedError is an error with HTTP status code hint for error handler
type CodedError struct {
	err  error
	code int
}

// Error returns error string
func (ce CodedError) Error() string {
	return ce.err.Error()
}

// Code returns HTTP error code
func (ce CodedError) Code() int {
	return ce.code
}

// Format is used by fmt package when printing error
func (ce CodedError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", ce.err)
			return
		}
		fallthrough
	case 's':
		fmt.Fprintf(s, "%s", ce.err)
	}
}

// NewCodedError creates new CodedError instance
func NewCodedError(code int, err error) error {
	return CodedError{err: err, code: code}
}
